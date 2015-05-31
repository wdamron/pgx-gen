package pgxgen

import (
	"fmt"
)

const fieldScannersFmt = `
%s
type %sFieldScanners []int

`

// generate type def for {struct-name}FieldScanners
func genFieldScannersType(s *Struct) string {
	doc := AutoCommentLn(fmt.Sprintf("type %sFieldScanners binds query/statement results to a value of type %s.", s.Name, s.Name))
	doc += "//\n"
	doc += AutoComment(fmt.Sprintf("Results are bound positionally, in correspondence with the field indexes stored within the %sFieldScanners slice.", s.Name))
	return fmt.Sprintf(fieldScannersFmt, doc, s.Name)
}

const scannersGetterFmt = `
%s
func (t *%sTableType) Scanners(colnames ...string) (%sFieldScanners, error) {
	indexes, err := %sTable.Indexes(colnames...)
	if err != nil {
		return nil, err
	}
	return %sFieldScanners(indexes), nil
}

`

// generate method def for {struct-name}TableType.Scanners
func genScannersGetter(s *Struct) string {
	doc := AutoCommentLn(fmt.Sprintf("Scanners creates an unbound instance of type %sFieldScanners for the columns/fields named by colnames.", s.Name))
	doc += "//\n"
	doc += AutoComment(fmt.Sprintf("Call %sFieldScanners.Bind to bind scanners from %sFieldScanners.", s.Name, s.Name))
	return fmt.Sprintf(scannersGetterFmt, doc, s.Name, s.Name, s.Name, s.Name)
}

const scannersBindFmt = `
%s
func (fs %sFieldScanners) Bind(v *%s) ([]pgx.Scanner, error) {
	bound := make([]pgx.Scanner, len(fs))
	for i, index := range fs {
		if index < 0 || index > len(%sTable.UnboundScanners) {
			return nil, errors.New("column scanner index out of range")
		}
		bound[i] = %sTable.UnboundScanners[index](v)
	}
	return bound, nil
}

`

// generate method def for {struct-name}FieldScanners.Bind
func genScannersBind(s *Struct) string {
	doc := AutoCommentLn("Bind binds query/statement result scanners for v.")
	doc += "//\n"
	doc += AutoComment(fmt.Sprintf("Scanners are bound positionally, in correspondence with the field indexes stored within the %sFieldScanners slice.", s.Name))
	return fmt.Sprintf(scannersBindFmt, doc, s.Name, s.Name, s.Name, s.Name)
}
