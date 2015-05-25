package pgxgen

import (
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
	DecodeOp    Op
}

func GetFieldColumnName(f astx.StructField) string {
	return KVFind(f, ColumnNameKey)
}

func GetFieldColumnType(f astx.StructField) string {
	return KVFind(f, ColumnTypeKey)
}

func KVFind(f astx.StructField, key string) string {
	if f.Tag == "" {
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
