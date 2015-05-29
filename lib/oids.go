package pgxgen

import (
	"strconv"
	"strings"
)

const (
	JSONOid = 114
	// json[] data types are not currently supported
	JSONArrayOid = 199
	UUIDOid      = 2950
	// uuid[] data types are not currently supported
	UUIDArrayOid = 2951
	// xml data types are not currently supported
	XMLOid = 142
)

var DataTypeNames = map[string]string{
	"bool":          "Bool",
	"int2":          "Int2",
	"int4":          "Int4",
	"int8":          "Int8",
	"float4":        "Float4",
	"float8":        "Float8",
	"bytea":         "Bytea",
	"text":          "Text",
	"varchar":       "Varchar",
	"date":          "Date",
	"timestampTz":   "TimestampTz",
	"timestamp":     "Timestamp",
	"bool[]":        "BoolArray",
	"int2[]":        "Int2Array",
	"int4[]":        "Int4Array",
	"int8[]":        "Int8Array",
	"real[]":        "Float4Array",
	"float[]":       "Float8Array",
	"text[]":        "TextArray",
	"varchar[]":     "VarcharArray",
	"timestamp[]":   "TimestampArray",
	"timestampTz[]": "TimestampTzArray",
	"hstore":        "Hstore",
	"json":          "Json",
	"uuid":          "Uuid",
	"oid":           "Oid",
}

var BinaryDataTypes = map[string]bool{
	"Bool":             true,
	"Bytea":            true,
	"Int2":             true,
	"Int4":             true,
	"Int8":             true,
	"Float4":           true,
	"Float8":           true,
	"TimestampTz":      true,
	"TimestampTzArray": true,
	"Timestamp":        true,
	"TimestampArray":   true,
	"BoolArray":        true,
	"Int2Array":        true,
	"Int4Array":        true,
	"Int8Array":        true,
	"Float4Array":      true,
	"Float8Array":      true,
	"TextArray":        true,
	"VarcharArray":     true,
	"Oid":              true,
	"Uuid":             true,
}

func NormalizeDataType(dataType string) string {
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
