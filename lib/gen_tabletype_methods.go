package pgxgen

import (
	"encoding/hex"
	"fmt"
)

// generate method def for ({struct-name})TableType.Index
func genIndexMethod(s *Struct) string {
	out := AutoCommentLn(fmt.Sprintf("Index returns the index of the column in %sTable with the given name.", s.Name))
	out += "//\n"
	out += AutoCommentLn("If no matching column is found, the returned index will be -1.")
	out += fmt.Sprintf("func (t *%sTableType) Index(colname string) int {\n", s.Name)
	out += "switch colname {\n"

	for i, c := range s.Columns {
		out += fmt.Sprintf("case \"%s\": return %d\n", c.Name, i)
	}
	out += "}\nreturn -1\n}\n\n"
	return out
}

const indexesMethodFmt = `
%s
func (t *%sTableType) Indexes(colnames ...string) ([]int, error) {
	indexes := make([]int, len(colnames))
	for i, colname := range colnames {
		index := t.Index(colname)
		if index < 0 {
			return nil, errors.New("column " + colname + " not found in %sTable")
		}
		indexes[i] = index
	}
	return indexes, nil
}

`

// generate method def for ({struct-name})TableType.Indexes
func genIndexesMethod(s *Struct) string {
	doc := AutoCommentLn(fmt.Sprintf("Indexes returns a slice of indexes of the given columns in %sTable with the given name.", s.Name))
	doc += "//\n"
	doc += AutoComment("If any of the columns are not found, an error will be returned and the returned slice of indexes will be nil.")
	return fmt.Sprintf(indexesMethodFmt, doc, s.Name, s.Name)
}

const aliasMethodFmt = `
%s
func (t *%sTableType) Alias(colnames ...string) ([]string, error) {
	aliases := []string{}
	// If no columns are specified, alias all columns:
	if len(colnames) == 0 {
		return %sTable.Aliases[:%d], nil
	}
	indexes, err := %sTable.Indexes(colnames...)
	if err != nil {
		return nil, err
	}
	for _, index := range indexes {
		aliases = append(aliases, %sTable.Aliases[index])
	}
	return aliases, nil
}

`

// generate method def for ({struct-name})TableType.Alias
func genAliasMethod(s *Struct) string {
	doc := AutoCommentLn("Alias aliases column names as hex-encoded indexes, for faster look-ups during decoding.")
	doc += "//\n"
	doc += AutoComment("If no column names are provided, all columns will be aliased, in which case AliasAll may be a faster alternative.")
	return fmt.Sprintf(aliasMethodFmt, doc, s.Name, s.Name, len(s.Columns), s.Name, s.Name)
}

// generate method def for ({struct-name})TableType.AliasAll
func genAliasAllMethod(s *Struct) string {
	out := AutoCommentLn("AliasAll aliases column names as hex-encoded indexes, for faster look-ups during decoding")
	out += fmt.Sprintf("func (t *%sTableType) AliasAll() string {\n", s.Name)
	out += "return \""

	lastIdx := len(s.Columns) - 1
	for i, c := range s.Columns {
		// 256 columns are supported, for now:
		hexIdx := hex.EncodeToString([]byte{byte(i)})
		if len(hexIdx) == 1 {
			hexIdx = "0" + hexIdx
		}
		shortName := "__" + hexIdx
		out += fmt.Sprintf("%s as %s::%s", c.Name, shortName, c.Type)
		if i != lastIdx {
			out += ", "
		}
	}

	return out + "\"\n}\n\n"
}
