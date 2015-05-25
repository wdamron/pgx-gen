package pgxgen

import (
	"github.com/wdamron/astx"
)

type Struct struct {
	astx.Struct
	Columns []Column
}

func NewStruct(as *astx.Struct) *Struct {
	s := &Struct{Struct: *as}
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
	if len(cols) != 0 {
		s.Columns = cols
	}
	return s
}
