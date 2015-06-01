package pgxgen

import (
	"fmt"
)

const fieldEncodersFmt = `
%s
type %sFieldEncoders []int

`

// generate type def for {struct-name}FieldEncoders
func genFieldEncodersType(s *Struct) string {
	doc := AutoCommentf("type %sFieldEncoders binds query/statement parameters from a value of type %s.\n", s.Name, s.Name)
	doc += "//\n"
	doc += AutoCommentf("Parameters are bound positionally, in correspondence with the field indexes stored within the %sFieldEncoders slice.", s.Name)
	return fmt.Sprintf(fieldEncodersFmt, doc, s.Name)
}

const encodersGetterFmt = `
%s
func (t *%sTableType) Encoders(colnames ...string) (%sFieldEncoders, error) {
	indexes, err := %sTable.Indexes(colnames...)
	if err != nil {
		return nil, err
	}
	return %sFieldEncoders(indexes), nil
}

`

// generate method def for {struct-name}TableType.Encoders
func genEncodersGetter(s *Struct) string {
	doc := AutoCommentf("Encoders creates an unbound instance of type %sFieldEncoders for the columns/fields named by colnames.\n", s.Name)
	doc += "//\n"
	doc += AutoCommentf("Call %sFieldEncoders.Bind to bind encoders from %sFieldEncoders.", s.Name, s.Name)
	return fmt.Sprintf(encodersGetterFmt, doc, s.Name, s.Name, s.Name, s.Name)
}

const encodersBindFmt = `
%s
func (fe %sFieldEncoders) Bind(v *%s) ([]pgx.Encoder, error) {
	bound := make([]pgx.Encoder, len(fe))
	for i, index := range fe {
		if index < 0 || index > len(%sTable.UnboundEncoders) {
			return nil, errors.New("column encoder index out of range")
		}
		bound[i] = %sTable.UnboundEncoders[index](v)
	}
	return bound, nil
}

`

// generate method def for {struct-name}FieldEncoders.Bind
func genEncodersBind(s *Struct) string {
	doc := AutoCommentLn("Bind binds query/statement parameter encoders for v.")
	doc += "//\n"
	doc += AutoCommentf("Encoders are bound positionally, in correspondence with the field indexes stored within the %sFieldEncoders slice.", s.Name)
	return fmt.Sprintf(encodersBindFmt, doc, s.Name, s.Name, s.Name, s.Name)
}
