package pgxgen

import (
	"encoding/hex"
	"fmt"
	"go/format"

	"github.com/wdamron/astx"
)

const DRIVER = "github.com/wdamron/pgx"

type File struct {
	Pkg, Driver string
	Structs     []Struct
	File        *astx.File
}

func NewFile(f *astx.File) *File {
	structs := make([]Struct, 0, len(f.Structs))
	for _, astStruct := range f.Structs {
		s := NewStruct(&astStruct)
		structs = append(structs, *s)
	}
	return &File{
		Pkg:     f.Package,
		Driver:  DRIVER,
		Structs: structs,
		File:    f,
	}
}

func (f *File) Gen() ([]byte, error) {
	out, err := f.gen()
	if err != nil {
		return nil, err
	}

	return format.Source(out)
}

func (f *File) gen() ([]byte, error) {
	// Always write the header (package name, pgxgen comment):
	out := genHeader(f)

	// If any structs in f contain column specs, write the import statement.
	// Then, for each column found, write the following:
	//   1. type def for {struct-name}ColumnDecoder func
	//   2. type def for {struct-name}TableType struct
	//   3. method def for ({struct-name})TableType.Alias
	//     * (alias selected column-names/types or all column-names/types)
	//   4. method def for ({struct-name})TableType.AliasAll
	//     * (alias all column-names/types; returns const string)
	//   5. var def for {struct-name}Table
	//     * (instance of {struct-name}TableType)
	//     * include ordered list of column-decoder funcs
	//     * include ordered list of column names
	//     * include ordered list of column aliases (e.g. x as __01::[]varchar)
	//   6. method def for {struct-name}.DecodeRow
	//     * (Decode values directly from a pgx conn into a struct!)
	wroteImports := false
	for _, s := range f.Structs {
		cols := s.Columns
		count := len(cols)
		if count == 0 {
			continue
		}
		// write import statement, if not already written:
		if !wroteImports {
			out += genImports(f)
			wroteImports = true
		}

		// 1. type def for {struct-name}ColumnDecoder func:
		out += genColDecoderType(&s)

		// 2. type def for {struct-name}TableType struct:
		out += genTableType(&s)

		// 3. method def for ({struct-name})TableType.Alias:
		out += genAliasMethod(&s)

		// 4. method def for ({struct-name})TableType.AliasAll:
		out += genAliasAllMethod(&s)

		// 5. var def for {struct-name}Table:
		table, err := genTable(&s)
		if err != nil {
			return nil, err
		}
		out += table

		// 6. method def for {struct-name}.DecodeRow:
		out += genRowDecoder(&s)

		out += genParamsEncoderType(&s)
		out += genEncoderFactory(&s)
		out += genBindVal(&s)
		out += genEncodeParamFormats(&s)
		out += genEncodeParams(&s)
	}

	return []byte(out), nil
}

const headerFmt = `
package %s

// Generated by pgxgen (see %s)

`

func genHeader(f *File) string {
	return fmt.Sprintf(headerFmt, f.Pkg, f.File.Path)
}

const importsFmt = `
import (
	"errors"
	"encoding/hex"
	
	"%s"
)

`

func genImports(f *File) string {
	return fmt.Sprintf(importsFmt, f.Driver)
}

// 1. generate type def for {struct-name}ColumnDecoder func
func genColDecoderType(s *Struct) string {
	out := fmt.Sprintf("// %sColumnDecoder funcs can be used to decode a single column value\n", s.Name)
	out += fmt.Sprintf("// from a pgx.ValueReader into type %s\n", s.Name)
	out += fmt.Sprintf("type %sColumnDecoder func(*%s, *pgx.ValueReader) error\n\n", s.Name, s.Name)
	return out
}

const tableTypeFmt = `
// %sTableType is the type of %sTable, which describes the table
// corresponding with type %s
type %sTableType struct {
	// Decoders can be used to decode a single column value from a
	// pgx.ValueReader into type %s
	Decoders    [%d]%sColumnDecoder
	ColumnNames [%d]string
	// Aliases contains an ordered list of column names aliased as hex-encoded
	// indexes, for faster look-ups during decoding
	Aliases     [%d]string
}

`

// 2. generate type def for {struct-name}TableType struct
func genTableType(s *Struct) string {
	cols := len(s.Columns)
	return fmt.Sprintf(tableTypeFmt, s.Name, s.Name, s.Name, s.Name, s.Name, cols, s.Name, cols, cols)
}

const aliasMethodFmt = `
// Alias aliases column names as hex-encoded indexes, for faster
// look-ups during decoding.
// If no column names are provided, all columns will be aliased,
// in which case AliasAll may be a faster alternative.
func (t *%sTableType) Alias(cols ...string) ([]string, error) {
	aliases := []string{}
	// If no columns are specified, alias all columns:
	if len(cols) == 0 {
		return %sTable.Aliases[:%d], nil
	}
	for _, name := range cols {
		switch name {
		%s
		}
	}
	return aliases, nil
}

`

// 3. generate method def for ({struct-name})TableType.Alias
func genAliasMethod(s *Struct) string {
	const caseFmt = "case \"%s\": aliases = append(aliases, \"%s as %s::%s\")\n"
	var cases string
	for i, c := range s.Columns {
		// 256 columns are supported, for now:
		hexIdx := hex.EncodeToString([]byte{byte(i)})
		if len(hexIdx) == 1 {
			hexIdx = "0" + hexIdx
		}
		shortName := "__" + hexIdx
		cases += fmt.Sprintf(caseFmt, c.Name, c.Name, shortName, c.Type)
	}
	cases += `default: return nil, errors.New("column " + name + " not found in type ` + s.Name + `")`

	return fmt.Sprintf(aliasMethodFmt, s.Name, s.Name, len(s.Columns), cases)
}

// 4. generate method def for ({struct-name})TableType.AliasAll
func genAliasAllMethod(s *Struct) string {
	out := "// AliasAll aliases column names as hex-encoded indexes, for faster\n"
	out += "// look-ups during decoding\n"
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

// 5. generate var def for {struct-name}Table
func genTable(s *Struct) (string, error) {
	out := fmt.Sprintf("// %sTable describes the table corresponding with type %s\n", s.Name, s.Name)
	out += fmt.Sprintf("var %sTable = %sTableType{\n", s.Name, s.Name)
	// include ordered list of column-decoder funcs:
	decoders, err := genColDecoderArray(s)
	if err != nil {
		return "", err
	}
	out += decoders
	// include ordered list of column names:
	out += genColNameArray(s)
	// include ordered list of column aliases:
	out += genAliasArray(s)

	out += genFormatArray(s)
	encoders, err := genEncoderArray(s)
	if err != nil {
		return "", err
	}
	out += encoders

	return out + "}\n\n", nil
}

// 5.a. include ordered list of column-decoder funcs
func genColDecoderArray(s *Struct) (string, error) {
	count := len(s.Columns)
	out := fmt.Sprintf("Decoders: [%d]%sColumnDecoder{\n", count, s.Name)
	for _, c := range s.Columns {
		decoder, err := genColumnDecoder(s, &c)
		if err != nil {
			return "", err
		}
		out += decoder
	}
	return out + "},\n", nil
}

func genColumnDecoder(s *Struct, c *Column) (string, error) {
	f := c.StructField
	coltype := c.Type
	if coltype == "" {
		return "", fmt.Errorf("no column type found in field: %s", f.Name)
	}
	if Decoders[coltype] == nil {
		return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
	}
	op := Decoders[coltype][f.Type]
	if op == Op(0) {
		return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
	}
	oidName := OidNames[coltype]
	if oidName == "" {
		return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
	}
	// TODO(wd): check overflow, when necessary
	out := fmt.Sprintf("// Decode column %s::%s into (*%s).%s\n", c.Name, coltype, s.Name, f.Name)
	out += fmt.Sprintf("func(v *%s, vr *pgx.ValueReader) error {\n", s.Name)
	var prefix, suffix string
	if op.MaskCast() != Op(0) {
		prefix, suffix = op.FormatCast()+"(", ")"
	}
	out += fmt.Sprintf("x := %svr.Decode%s()%s\n", prefix, oidName, suffix)
	out += "if vr.Err() != nil {\nreturn vr.Err()\n}\n"
	if op.PtrAssign() {
		out += "*"
	}
	out += fmt.Sprintf("v.%s = x\nreturn nil\n},\n", f.Name)

	return out, nil
}

// 5.b. include ordered list of column names
func genColNameArray(s *Struct) string {
	out := fmt.Sprintf("ColumnNames: [%d]string{\n", len(s.Columns))
	for _, c := range s.Columns {
		out += fmt.Sprintf("\"%s\",\n", c.Name)
	}
	return out + "},\n"
}

// 5.c. include ordered list of column aliases
func genAliasArray(s *Struct) string {
	out := "// Aliases contains an ordered list of column names aliased as hex-encoded\n"
	out += "// indexes, for faster look-ups during decoding\n"
	out += fmt.Sprintf("Aliases: [%d]string{\n", len(s.Columns))
	for i, col := range s.Columns {
		// 256 columns are supported, for now:
		hexIdx := hex.EncodeToString([]byte{byte(i)})
		if len(hexIdx) == 1 {
			hexIdx = "0" + hexIdx
		}
		shortName := "__" + hexIdx
		out += fmt.Sprintf("\"%s as %s::%s\",\n", col.Name, shortName, col.Type)
	}
	return out + "},\n"
}

// 5.d. include ordered list of column format codes (text=0, binary=1)
func genFormatArray(s *Struct) string {
	out := fmt.Sprintf("Formats: [%d]int{", len(s.Columns))
	for i, c := range s.Columns {
		if BinaryFmtOids[OidNames[c.Type]] {
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

// 5.e. include ordered list of field encoders
func genEncoderArray(s *Struct) (string, error) {
	out := fmt.Sprintf("Encoders: [%d]func(*%s, *pgx.WriteBuf) error {\n", len(s.Columns), s.Name)
	for _, c := range s.Columns {
		f := c.StructField
		if Encoders[f.Type] == nil {
			return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
		}
		op := Encoders[f.Type][c.Type]
		if op == Op(0) {
			return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
		}
		oidName := OidNames[c.Type]
		if oidName == "" {
			return "", fmt.Errorf("no column decoder available for field: %s", f.Name)
		}

		out += fmt.Sprintf("func(v *%s, wbuf *pgx.WriteBuf) error {\n", s.Name)
		var castPrefix, castSuffix string
		if op.MaskCast() != Op(0) {
			castPrefix, castSuffix = op.FormatCast()+"(", ")"
		}
		deref := ""
		if op.DerefPass() {
			deref = "*"
		}
		out += fmt.Sprintf("wbuf.Encode%s(%s%sv.%s%s)\n", oidName, castPrefix, deref, c.StructField.Name, castSuffix)
		out += "return nil\n},\n"
	}
	return out + "},\n", nil
}

const rowDecoderFmt = `
// DecodeRow decodes a SQL row into type %s.
// If an error is returned, the caller should call Rows.Close()
func (v *%s) DecodeRow(r *pgx.Rows) error {
	for _ = range r.FieldDescriptions() {
		vr, ok := r.NextColumn()
		if !ok {
			if vr != nil && vr.Err() {
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
			if index < 0 || index > %d {
				return errors.New("column-decoder index out of range")
			}
			dec := %sTable.Decoders[index]
			if err = dec(v, vr); err != nil {
				return err
			}
			continue
		}
		
		// Slow path:
		var dec %sColumnDecoder
		switch colname {
		%s
		}
		if err := dec(v, vr); err != nil {
			return err
		}
	}
	return nil
}

`

// 6. generate method def for {struct-name}.DecodeRow
func genRowDecoder(s *Struct) string {
	const caseFmt = "case \"%s\": dec = %sTable.Decoders[%d]\n"
	var cases string
	for i, col := range s.Columns {
		cases += fmt.Sprintf(caseFmt, col.Name, s.Name, i)
	}
	cases += `default: return errors.New("unknown column name: " + colname)`

	return fmt.Sprintf(rowDecoderFmt, s.Name, s.Name, len(s.Columns)-1, s.Name, s.Name, cases)
}

const paramsEncoderFmt = `
// %sParamsEncoder encodes parameter formats and values for a prepared statement
// to a buffer, and implements pgx.ParamsEncoder
type %sParamsEncoder struct {
	Formats []int
	ValEncoders []func(v *%s, wbuf *pgx.WriteBuf)
	v *%s
}
`

func genParamsEncoderType(s *Struct) string {
	return fmt.Sprintf(paramsEncoderFmt, s.Name, s.Name, s.Name, s.Name)
}

const newEncoderFmt = `
func (t *%sTableType) Encoder(cols ...string) (*%sParamsEncoder, error) {
	formats := make([]int, 0, len(cols))
	encoders := make([]func(*%s, *pgx.WriteBuf), 0, len(cols))

	for _, colname := range cols {
		switch colname {
		%s
		}
	}

	pe := &%sParamsEncoder{
		Formats: formats,
		ValEncoders: encoders,
	}
	return pe, nil
}
`

func genEncoderFactory(s *Struct) string {
	const caseFmt = `case "%s":
		formats = append(formats, %sTable.Formats[%d])
		encoders = append(encoders, %sTable.Encoders[%d])
`
	var cases string
	for i, col := range s.Columns {
		cases += fmt.Sprintf(caseFmt, col.Name, s.Name, i, s.Name, i)
	}
	cases += `default: return nil, errors.New("unknown column name: " + colname)`

	return fmt.Sprintf(newEncoderFmt, s.Name, s.Name, s.Name, cases, s.Name)
}

const bindValFmt = `
func (pe %sParamsEncoder) Bind(v *%s) pgx.ParamsEncoder {
	return &%sParamsEncoder{
		Formats: pe.Formats,
		ValEncoders: pe.ValEncoders,
		v: v,
	}
}
`

func genBindVal(s *Struct) string {
	return fmt.Sprintf(bindValFmt, s.Name, s.Name, s.Name)
}

const encodeParamFormatsFmt = `
// EncodeParamFormats encodes the formats of all params for the given statement
// into wbuf
func (pe *%sParamsEncoder) EncodeParamFormats(wbuf *pgx.WriteBuf) error {
	for _, format := range pe.Formats {
		wbuf.WriteInt16(int16(format))
	}
	return nil
}

`

func genEncodeParamFormats(s *Struct) string {
	return fmt.Sprintf(encodeParamFormatsFmt, s.Name)
}

const encodeParamsFmt = `
func (pe *%sParamsEncoder) EncodeParams(wbuf *pgx.WriteBuf) error {
	for _, enc := range pe.ValEncoders {
		enc(pe.v, wbuf)
	}
	return nil
}
`

func genEncodeParams(s *Struct) string {
	return fmt.Sprintf(encodeParamsFmt, s.Name)
}