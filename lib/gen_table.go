package pgxgen

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// generate var def for {struct-name}Table
func genTable(s *Struct) (string, error) {
	out := AutoCommentLn(fmt.Sprintf("%sTable describes the table corresponding with type %s", s.Name, s.Name))
	out += fmt.Sprintf("var %sTable = %sTableType{\n", s.Name, s.Name)
	// include ordered list of field-encoder funcs:
	encoders, err := genFieldEncoderArray(s)
	if err != nil {
		return "", err
	}
	out += encoders
	// include ordered list of field-scanner funcs:
	decoders, err := genFieldScannerArray(s)
	if err != nil {
		return "", err
	}
	out += decoders
	// include ordered list of column names:
	out += genColNameArray(s)
	// include ordered list of column types:
	out += genColTypeArray(s)
	// include ordered list of column aliases:
	out += genAliasArray(s)
	// include ordered list of column format codes:
	out += genFormatArray(s)
	// include ordered list of column oids:
	out += genOidArray(s)

	return out + "}\n\n", nil
}

// include ordered list of field-encoder funcs (see genTable)
func genFieldEncoderArray(s *Struct) (string, error) {
	out := fmt.Sprintf("UnboundEncoders: [%d]func(*%s) pgx.Encoder {\n", len(s.Columns), s.Name)
	for _, c := range s.Columns {
		f := c.StructField
		op := c.EncodeOp
		if op == Op(0) && c.Type != "json" {
			return "", fmt.Errorf("no column encoder available for field: %s.%s (coltype=%s, fieldtype=%s)", s.Name, f.Name, c.Type, f.Type)
		}
		dtName := DataTypeNames[c.Type]
		if dtName == "" {
			return "", fmt.Errorf("no column oid available for field: %s.%s (coltype=%s, fieldtype=%s)", s.Name, f.Name, c.Type, f.Type)
		}

		out += fmt.Sprintf("// Encode v.%s as %s\n", f.Name, c.Type)
		out += fmt.Sprintf("func(v *%s) pgx.Encoder {\n", s.Name)
		deref := ""
		if f.Type[0] == '*' {
			deref = "*"
		}
		switch {
		default:
			if c.Type != "json" {
				var castPrefix, castSuffix string
				if op.MaskCast() != Op(0) {
					castPrefix, castSuffix = op.FormatCast()+"(", ")"
				}
				out += fmt.Sprintf("return pgtypes.%sEncoder(%s%sv.%s%s)\n", dtName, castPrefix, deref, f.Name, castSuffix)
			} else {
				switch f.Type {
				case "string", "*string":
					out += fmt.Sprintf("return pgtypes.JSONEncoderString(%sv.%s)\n", deref, f.Name)
				case "[]byte", "*[]byte":
					out += fmt.Sprintf("return pgtypes.JSONEncoderBytes(%sv.%s)\n", deref, f.Name)
				default:
					out += fmt.Sprintf("return pgtypes.JSONEncoder(v.%s)\n", f.Name)
				}
			}
		case op.CustomEncode():
			out += fmt.Sprintf("return v.%s\n", f.Name)
		case op.HstoreMapEncode():
			out += fmt.Sprintf("return pgtypes.HstoreMapEncoder(%sv.%s)\n", deref, f.Name)
		case op.UuidStringEncode():
			out += fmt.Sprintf("return pgtypes.UUIDEncoderString(%sv.%s)\n", deref, f.Name)
		}
		out += "},\n"
	}
	return out + "},\n", nil
}

// include ordered list of field-scanner funcs (see genTable)
func genFieldScannerArray(s *Struct) (string, error) {
	count := len(s.Columns)
	out := fmt.Sprintf("UnboundScanners: [%d]func(*%s) pgx.Scanner{\n", count, s.Name)
	for _, c := range s.Columns {
		decoder, err := genFieldScanner(s, &c)
		if err != nil {
			return "", err
		}
		out += decoder
	}
	return out + "},\n", nil
}

// generate field-scanner func (see genFieldScannerArray)
func genFieldScanner(s *Struct, c *Column) (string, error) {
	f := c.StructField
	coltype := c.Type
	op := c.DecodeOp
	if op == Op(0) && c.Type != "json" {
		return "", fmt.Errorf("no column decoder available for field: %s.%s (coltype=%s, fieldtype=%s)", s.Name, f.Name, c.Type, f.Type)
	}
	dtName := DataTypeNames[coltype]
	if dtName == "" {
		return "", fmt.Errorf("no column oid available for field: %s.%s (coltype=%s, fieldtype=%s)", s.Name, f.Name, c.Type, f.Type)
	}
	// TODO(wd): check overflow, when necessary
	out := fmt.Sprintf("// Decode column %s::%s into v.%s\n", c.Name, coltype, f.Name)
	out += fmt.Sprintf("func(v *%s) pgx.Scanner {\n", s.Name)
	takeAddr := ""
	if len(f.Type) != 0 && f.Type[0] != '*' {
		takeAddr = "&"
	}
	switch {
	default:
		if c.Type != "json" {
			if op.MaskCast() == Op(0) {
				out += fmt.Sprintf("return pgtypes.%sScanner(%sv.%s)\n", dtName, takeAddr, f.Name)
			} else {
				cast := op.FormatCast()
				if cast == "" {
					return "", fmt.Errorf("no column oid available for field: %s.%s (coltype=%s, fieldtype=%s)", s.Name, f.Name, c.Type, f.Type)
				}
				out += fmt.Sprintf("return pgtypes.Into%s(%sv.%s)\n", strings.Title(cast), takeAddr, f.Name)
			}
		} else {
			switch f.Type {
			case "string", "*string":
				out += fmt.Sprintf("return pgtypes.JSONScannerString(%sv.%s)\n", takeAddr, f.Name)
			case "[]byte", "*[]byte":
				out += fmt.Sprintf("return pgtypes.JSONScannerBytes(%sv.%s)\n", takeAddr, f.Name)
			default:
				out += fmt.Sprintf("return pgtypes.JSONScanner(%sv.%s)\n", takeAddr, f.Name)
			}
		}
	case op.CustomScan():
		out += fmt.Sprintf("return %sv.%s\n", takeAddr, f.Name)
	case op.HstoreMapDecode():
		out += fmt.Sprintf("return pgtypes.HstoreMapScanner(%sv.%s)\n", takeAddr, f.Name)
	case op.UuidDecode():
		if op.UuidStringDecode() {
			out += fmt.Sprintf("return pgtypes.UUIDScannerString(%sv.%s)\n", takeAddr, f.Name)
		} else {
			out += fmt.Sprintf("return pgtypes.UUIDScanner(%sv.%s)\n", takeAddr, f.Name)
		}
	}
	out += "},\n"

	return out, nil
}

// include ordered list of column names (see genTable)
func genColNameArray(s *Struct) string {
	out := fmt.Sprintf("Names: [%d]string{\n", len(s.Columns))
	for _, c := range s.Columns {
		out += fmt.Sprintf("\"%s\",\n", c.Name)
	}
	return out + "},\n"
}

// include ordered list of column types (see genTable)
func genColTypeArray(s *Struct) string {
	out := fmt.Sprintf("Types: [%d]string{\n", len(s.Columns))
	for _, c := range s.Columns {
		out += fmt.Sprintf("\"%s\",\n", c.Type)
	}
	return out + "},\n"
}

// include ordered list of column aliases (see genTable)
func genAliasArray(s *Struct) string {
	out := AutoCommentLn("Aliases contains an ordered list of column names aliased as hex-encoded indexes, for faster look-ups during decoding")
	out += fmt.Sprintf("Aliases: [%d]string{\n", len(s.Columns))
	for i, c := range s.Columns {
		// 256 columns are supported, for now:
		hexIdx := hex.EncodeToString([]byte{byte(i)})
		if len(hexIdx) == 1 {
			hexIdx = "0" + hexIdx
		}
		shortName := "__" + hexIdx
		out += fmt.Sprintf("\"%s as %s::%s\",\n", c.Name, shortName, c.Type)
	}
	return out + "},\n"
}

// include ordered list of column oids (see genTable)
func genOidArray(s *Struct) string {
	out := fmt.Sprintf("Oids: [%d]pgx.Oid{\n", len(s.Columns))
	for _, c := range s.Columns {
		out += fmt.Sprintf("pgtypes.%sOid,\n", DataTypeNames[c.Type])
	}
	return out + "},\n"
}

// include ordered list of column format codes (text=0, binary=1) --  (see genTable)
func genFormatArray(s *Struct) string {
	out := fmt.Sprintf("Formats: [%d]int{", len(s.Columns))
	for i, c := range s.Columns {
		if BinaryDataTypes[DataTypeNames[c.Type]] {
			out += "1"
		} else {
			out += "0"
		}
		if i != len(s.Columns)-1 {
			out += ", "
		}
	}
	return out + "},\n"
}
