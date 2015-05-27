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
		if !IsColumn(f) {
			continue
		}
		col := NewColumn(&s.Struct.Fields[i])
		cols = append(cols, *col)
	}
	if len(cols) != 0 {
		s.Columns = cols
	}
	return s
}
