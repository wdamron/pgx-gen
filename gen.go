package pgxgen

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/wdamron/astx"
)

const DRIVER = "github.com/wdamron/pgx"
const PGTYPES_PKG = "github.com/wdamron/pgx-gen/pgtypes"
const UUID_PKG = "github.com/satori/go.uuid"

// type File holds information extracted from a Go source file
type File struct {
	Pkg, Driver string
	Structs     []Struct
	File        *astx.File
}

// NewFile extracts information from f into a new File.
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

// Gen generates and formats code for f, returning bytes or nil if an error has
// occurred
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

	uuidImported := false
	uuidPkg := fmt.Sprintf(`"%s"`, UUID_PKG)
	for _, imp := range f.File.Imports {
		if imp.Path == uuidPkg {
			uuidImported = true
			break
		}
	}

	// stdImports/otherImports contain mappings from paths to names for imports:
	var stdImports, otherImports = make(map[string]string, 0), make(map[string]string, 0)
	var body string
	for _, s := range f.Structs {
		cols := s.Columns
		count := len(cols)
		if count == 0 {
			continue
		}

		// ensure std packages are imported when columns are present:
		stdImports["errors"] = ""
		stdImports["encoding/hex"] = ""
		// ensure driver is imported when columns are present:
		otherImports[f.Driver] = ""
		// ensure pgtypes is imported when columns are present:
		otherImports[PGTYPES_PKG] = ""

		for _, c := range cols {
			switch c.Type {
			case "uuid":
				ftype := c.StructField.Type
				if (ftype == "uuid.UUID" || ftype == "*uuid.UUID") && !uuidImported {
					return nil, fmt.Errorf("package %s must be imported to use uuid column data types\n(see %s)", uuidPkg, f.File.AbsPath)
				}
			// case "json":
			// 	switch c.StructField.Type {
			// 	// include json package when encoding/decoding Go types:
			// 	default:
			// 		stdImports["encoding/json"] = ""
			// 	// json is not encoded/decoded for string and []byte types:
			// 	case "string", "*string", "[]byte", "*[]byte":
			// 	}
			default:
			}
		}

		// generate type def for {struct-name}TableType struct:
		body += genTableType(&s)

		// generate var def for {struct-name}Table:
		table, err := genTable(&s)
		if err != nil {
			return nil, err
		}
		body += table

		// generate method def for {struct-name}TableType.Index:
		body += genIndexMethod(&s)

		// generate method def for {struct-name}TableType.Indexes:
		body += genIndexesMethod(&s)

		// generate method def for ({struct-name})TableType.Alias:
		body += genAliasMethod(&s)

		// generate method def for ({struct-name})TableType.AliasAll:
		body += genAliasAllMethod(&s)

		// generate method def for {struct-name}.DecodeRow:
		body += genRowDecoder(&s)

		// generate type def for {struct-name}FieldEncoders:
		body += genFieldEncodersType(&s)

		// generate method def for ({struct-name})TableType.Encoders:
		body += genEncodersGetter(&s)

		// generate method def for {struct-name}FieldEncoders.Bind:
		body += genEncodersBind(&s)

		// generate type def for {struct-name}FieldScanners:
		body += genFieldScannersType(&s)

		// generate method def for ({struct-name})TableType.Scanners:
		body += genScannersGetter(&s)

		// generate method def for {struct-name}FieldScanners.Bind:
		body += genScannersBind(&s)
	}

	out += genImports(f, stdImports, otherImports)
	out += body

	return []byte(out), nil
}

func AutoComment(s string) string {
	s = strings.TrimSpace(s)
	out := "//"
	if len(out)+1+len(s) <= 80 {
		return out + " " + s
	}
	lineSz := len(out)
	for _, w := range strings.Fields(s) {
		if lineSz+1+len(w) >= 80 {
			out += "\n//"
			lineSz = 2
		}
		out += " " + w
		lineSz += 1 + len(w)

	}
	if out == "//" {
		return ""
	}
	return out
}

func AutoCommentLn(s string) string {
	return AutoComment(s) + "\n"
}

func AutoCommentf(format string, args ...interface{}) string {
	if format == "" {
		return ""
	}
	s := fmt.Sprintf(format, args...)
	if format[len(format)-1] == '\n' {
		return AutoCommentLn(s)
	}
	return AutoComment(s)
}
