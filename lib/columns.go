package pgxgen

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/wdamron/astx"
)

const (
	ColumnTagName = "pgx"
	ColumnNameKey = "col"
	ColumnTypeKey = "type"
)

type Column struct {
	Name, Type  string
	StructField *astx.StructField
	DecodeOp    Op
}

func GenColumnDecoderBody(c *Column) (string, error) {
	f := c.StructField
	coltype := c.Type
	if coltype == "" {
		return "", fmt.Errorf("no column type found in field: %s", f.Name)
	}
	if Decoders[coltype] == nil {
		return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
	}
	op := Decoders[coltype][f.Type]
	if op == Op(0) {
		return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
	}
	dec := DecoderNames[coltype]
	if dec == "" {
		return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
	}
	oid := OidNames[coltype]
	if oid == "" {
		return "", fmt.Errorf("no oid name found for field: %s (type=%s)", f.Name, coltype)
	}

	// TODO(wd): check overflow, when necessary
	ret := fmt.Sprintf("if vr.Type().DataType != pgx.%s {\nreturn errors.New(\"mismatched data types\")\n}\n", oid)
	ret += "x := "
	if op.MaskCast() != Op(0) {
		ret += fmt.Sprintf("%s(", op.FormatCast())
	}
	ret += fmt.Sprintf("vr.%s()", dec)
	if op.MaskCast() != Op(0) {
		ret += ")"
	}
	ret += "\nif vr.Err() != nil {\nreturn vr.Err()\n}\n"
	if op.PtrAssign() {
		ret += "*"
	}
	ret += fmt.Sprintf("v.%s = x\nreturn nil", f.Name)

	return ret, nil
}

func GenDefaultRenamerBody(cols []Column) string {
	body := "aliases = []string{\n"
	for i, col := range cols {
		// 256 columns are supported, for now:
		hexIdx := hex.EncodeToString([]byte{byte(i)})
		if len(hexIdx) == 1 {
			hexIdx = "0" + hexIdx
		}
		shortName := "__" + hexIdx
		body += fmt.Sprintf("\"%s as %s::%s\",\n", col.Name, shortName, col.Type)
	}
	body += "}\nreturn aliases, nil"
	return body
}

func GenRenamerCases(s *Struct, cols []Column) string {
	const caseTemplate = "case \"%s\": aliases = append(aliases, \"%s\")\n"
	var cases string

	for i, col := range cols {
		// 256 columns are supported, for now:
		hexIdx := hex.EncodeToString([]byte{byte(i)})
		if len(hexIdx) == 1 {
			hexIdx = "0" + hexIdx
		}
		shortName := "__" + hexIdx
		alias := fmt.Sprintf("%s as %s::%s", col.Name, shortName, col.Type)
		cases += fmt.Sprintf(caseTemplate, col.Name, alias)
	}
	cases += `default: return nil, errors.New("column " + name + " not found in type ` + s.StructName() + `")`
	return cases
}

func GenRowDecoderCases(s *Struct, cols []Column) string {
	const caseTemplate = "case \"%s\": dec = %sTable.Decoders[%d]\n"
	var cases string

	for i, col := range cols {
		cases += fmt.Sprintf(caseTemplate, col.Name, s.StructName(), i)
	}
	cases += `default: return errors.New("unknown column name: " + colname)`
	return cases
}

func GetFieldColumnName(f astx.StructField) string {
	return KVFind(f, ColumnNameKey)
}

func GetFieldColumnType(f astx.StructField) string {
	return KVFind(f, ColumnTypeKey)
}

func KVFind(f astx.StructField, key string) string {
	if f.RawTag == "" || f.RawTag == "``" {
		return ""
	}
	t := f.Tag.Get(ColumnTagName)
	if t == "" {
		return ""
	}
	split := strings.Split(t, ";")
	for _, pair := range split {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			continue
		}
		trimmed := strings.TrimSpace(kv[0])
		if trimmed != key {
			continue
		}
		return strings.TrimSpace(kv[1])
	}
	return ""
}
