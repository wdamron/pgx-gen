package exp

import (
	"go/parser"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
)

const (
	ColumnTagName = "pgx"
	ColumnNameKey = "name"
	ColumnTypeKey = "type"
)

type Program struct {
	loader.Program
}

func Load(importPath string) (*Program, error) {
	cfg := &loader.Config{
		Fset:        nil,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	cfg.Import(importPath)
	p, err := cfg.Load()
	return &Program{*p}, err
}

func (p *Program) StructTypes() []Struct {
	var structs []Struct
	for _, pkgInfo := range p.Imported {
		scope := pkgInfo.Pkg.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			s, ok := NewStruct(obj)
			if !ok {
				continue
			}
			structs = append(structs, *s)
		}
	}
	return structs
}

type Struct struct {
	Name   string
	Fields []Field
	obj    types.Object
	typ    *types.Struct
}

func NewStruct(obj types.Object) (*Struct, bool) {
	st, isStruct := obj.Type().(*types.Struct)
	if !isStruct {
		return nil, false
	}
	count := st.NumFields()
	fields := make([]Field, count)
	for i, _ := range fields {
		f := NewField(st, i)
		fields[i] = *f
	}
	s := &Struct{
		Name:   obj.Name(),
		Fields: fields,
		obj:    obj,
		typ:    st,
	}
	return s, true
}

type Field struct {
	Name     string
	Spec     map[string]string
	Type     types.Type
	PtrDepth int
	field    *types.Var
}

func NewField(s *types.Struct, i int) *Field {
	field := s.Field(i)
	spec := parseTag(s.Tag(i))
	typ := field.Type()
	ptrDepth := 0
	for {
		if p, isPtr := typ.(*types.Pointer); isPtr {
			typ = p.Elem()
			ptrDepth++
			continue
		}
		break
	}
	return &Field{
		Name:     field.Name(),
		Spec:     spec,
		Type:     typ,
		PtrDepth: ptrDepth,
		field:    field,
	}
}

func parseTag(structTag string) map[string]string {
	spec := map[string]string{}
	t := reflect.StructTag(structTag).Get(ColumnTagName)
	split := strings.Split(t, ";")
	for _, pair := range split {
		kv := strings.Split(pair, ":")
		if len(kv) == 0 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		if k == "" {
			continue
		}
		if len(kv) > 1 {
			v := strings.TrimSpace(kv[1])
			if k == ColumnTypeKey {
				spec[ColumnTypeKey] = normalizeDataType(v)
				continue
			}
			boolVal, err := strconv.ParseBool(v)
			if err != nil && v != "" {
				spec[k] = v
				continue
			}
			if !boolVal {
				spec[k] = "0"
				continue
			}
		}
		spec[k] = "1"
	}
	if len(spec) == 0 {
		return nil
	}
	return spec
}

func normalizeDataType(dataType string) string {
	if dataType == "" {
		return ""
	}
	switch dataType {
	case "custom", "bytea", "text", "date", "text[]", "varchar[]", "timestampTz[]", "hstore", "json", "uuid", "oid":
		return dataType
	case "bool", "boolean":
		return "bool"
	case "int2", "smallint":
		return "int2"
	case "int4", "int", "integer":
		return "int4"
	case "int8", "bigint":
		return "int8"
	case "float4", "float(32)", "real", "float32":
		return "real"
	case "float8", "float(53)", "float", "double", "double precision", "float64":
		return "float"
	case "varchar", "character varying":
		return "varchar"
	case "timestamp", "timestamp without time zone", "time", "time without time zone":
		return "timestamp"
	case "timestampTz", "timestamp with time zone", "time with time zone":
		return "timestampTz"
	case "bool[]", "boolean[]":
		return "boolean[]"
	case "int2[]", "smallint[]":
		return "int2[]"
	case "int4[]", "int[]", "integer[]":
		return "int4[]"
	case "int8[]", "bigint[]":
		return "int8[]"
	case "float4[]", "float(32)[]", "real[]", "float32[]":
		return "real[]"
	case "float8[]", "float(53)[]", "float[]", "double[]":
		return "float[]"
	case "timestamp[]", "time[]":
		return "timestamp[]"
	default:
		if strings.HasPrefix(dataType, "varchar") || strings.HasPrefix(dataType, "character varying") {
			return "varchar"
		}
		if strings.HasPrefix(dataType, "float(") {
			width, _ := strconv.ParseInt(dataType[len("float("):len(dataType)-1], 0, 8)
			if width >= 1 && width <= 24 {
				return "real"
			}
			return "float"
		}
	}
	return ""
}
