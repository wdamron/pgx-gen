package pgxgen

import (
	"strconv"
	"strings"

	"github.com/wdamron/astx"
)

const (
	ColumnTagName = "pgx"
	ColumnNameKey = "name"
	ColumnTypeKey = "type"
)

type Column struct {
	Name, Type  string
	StructField *astx.StructField
	Spec        map[string]string
	EncodeOp    Op
	DecodeOp    Op
}

func IsColumn(f astx.StructField) bool {
	return f.Tag.Get(ColumnTagName) != ""
}

func NewColumn(f *astx.StructField) *Column {
	spec := GetFieldColumnSpec(f)
	colname := spec[ColumnNameKey]
	coltype := spec[ColumnTypeKey]
	col := &Column{
		Name:        colname,
		Type:        coltype,
		StructField: f,
		Spec:        spec,
	}
	if Encoders[f.Type] != nil {
		col.EncodeOp = Encoders[f.Type][coltype]
	}
	if Decoders[coltype] != nil {
		col.DecodeOp = Decoders[coltype][f.Type]
	}
	return col
}

func GetFieldColumnSpec(f *astx.StructField) map[string]string {
	spec := map[string]string{}
	t := f.Tag.Get(ColumnTagName)
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
				spec[ColumnTypeKey] = NormalizeDataType(v)
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
	return spec
}
