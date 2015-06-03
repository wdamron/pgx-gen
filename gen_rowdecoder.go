package pgxgen

import (
	"fmt"
)

const rowDecoderFmt = `
%s
func (v *%s) DecodeRow(r *pgx.Rows) error {
	for _ = range r.FieldDescriptions() {
		vr, ok := r.NextColumn()
		if !ok {
			if vr != nil && vr.Err() != nil {
				return vr.Err()
			}
			break
		}
		colname := vr.Type().Name

		// Fast path (aliased columns):
		if len(colname) == 4 && colname[:2] == "__" {
			b, err := hex.DecodeString(colname[2:4])
			if err != nil {
				return err
			}
			index := int(b[0])
			if index < 0 || index > len(%sTable.UnboundScanners) - 1 {
				return errors.New("column decoder index out of range")
			}
			bound := %sTable.UnboundScanners[index](v)
			if err = bound.Scan(vr); err != nil {
				return err
			}
			continue
		}
		
		// Slow path:
		index := %sTable.Index(colname)
		if index < 0 {
			return errors.New("column decoder for " + colname + " not found in %sTable")
		}
		bound := %sTable.UnboundScanners[index](v)
		if err := bound.Scan(vr); err != nil {
			return err
		}
	}
	return nil
}

`

// generate method def for {struct-name}.DecodeRow
func genRowDecoder(s *Struct) string {
	doc := AutoCommentLn("DecodeRow decodes a single row/result from r into v.")
	doc += "//\n"
	doc += AutoComment("If an error is returned, the caller should call Rows.Close()")
	return fmt.Sprintf(rowDecoderFmt, doc, s.Name, s.Name, s.Name, s.Name, s.Name, s.Name)
}
