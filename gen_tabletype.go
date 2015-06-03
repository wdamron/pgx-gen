package pgxgen

import (
	"fmt"
)

// generate type def for {struct-name}TableType struct
func genTableType(s *Struct) string {
	cols := len(s.Columns)

	out := AutoCommentf("%sTableType is the type of %sTable, which describes the table corresponding with type %s\n", s.Name, s.Name, s.Name)
	out += fmt.Sprintf("type %sTableType struct {\n", s.Name)

	out += AutoCommentf("UnboundEncoders are used by %sParamsEncoder.Bind to bind query/statement parameters from a value of type %s\n", s.Name, s.Name)
	out += fmt.Sprintf("UnboundEncoders    [%d]func(*%s) pgx.Encoder\n", cols, s.Name)

	out += AutoCommentf("UnboundScanners are used by %sParamsScanner.Bind to bind query/statement results to fields within type %s\n", s.Name, s.Name)
	out += fmt.Sprintf("UnboundScanners    [%d]func(*%s) pgx.Scanner\n", cols, s.Name)

	out += AutoCommentLn("Names contains an ordered list of column names")
	out += fmt.Sprintf("Names       [%d]string\n", cols)

	out += AutoCommentLn("Types contains an ordered list of column types")
	out += fmt.Sprintf("Types       [%d]string\n", cols)

	out += AutoCommentLn("Aliases contains an ordered list of column names aliased as hex-encoded indexes, for faster look-ups during decoding")
	out += fmt.Sprintf("Aliases     [%d]string\n", cols)

	out += AutoCommentLn("Formats contains an ordered list of column format codes (text=0, binary=1)")
	out += fmt.Sprintf("Formats     [%d]int\n", cols)

	out += AutoCommentLn("Oids contains an ordered list of column oid codes (corresponding with Postgres types)")
	out += fmt.Sprintf("Oids        [%d]pgx.Oid\n", cols)

	out += "}\n\n"

	return out

}
