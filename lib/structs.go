package pgxgen

import (
	"github.com/wdamron/astx"
)

type Struct struct {
	Struct astx.Struct
}

func (s *Struct) StructName() string {
	return s.Struct.Name
}

func (s *Struct) GetColumns() []Column {
	cols := []Column{}
	for i, f := range s.Struct.Fields {
		colname := GetFieldColumnName(f)
		coltype := GetFieldColumnType(f)
		if colname == "" || coltype == "" {
			continue
		}
		if Decoders[coltype] == nil {
			continue
		}
		op := Decoders[coltype][f.Type]
		if op == Op(0) {
			continue
		}
		col := Column{
			Name:        colname,
			Type:        coltype,
			StructField: &s.Struct.Fields[i],
			DecodeOp:    op,
		}
		cols = append(cols, col)
	}

	return cols
}
