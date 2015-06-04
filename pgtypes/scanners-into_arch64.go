// +build amd64 arm64 ppc64 ppc64le

package pgtypes

import (
	"fmt"

	"github.com/wdamron/pgx"
)

type g_intScanner struct {
	v *int
}

func IntoInt(v *int) pgx.Scanner {
	return g_intScanner{v}
}

func (s g_intScanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int"))
		return vr.Err()
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return vr.Err()
	}

	switch vr.Type().DataType {
	case BoolOid:
		*s.v = int(vr.ReadByte())
		return vr.Err()
	case Int2Oid:
		*s.v = int(vr.ReadInt16())
		return vr.Err()
	case Int4Oid:
		*s.v = int(vr.ReadInt32())
		return vr.Err()
	case Int8Oid:
		*s.v = int(vr.ReadInt64())
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into int", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_uintScanner struct {
	v *uint
}

func IntoUint(v *uint) pgx.Scanner {
	return g_uintScanner{v}
}

func (s g_uintScanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int"))
		return vr.Err()
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return vr.Err()
	}

	switch vr.Type().DataType {
	case BoolOid:
		v := int8(vr.ReadByte())
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint", v, v))
			return vr.Err()
		}
		*s.v = uint(v)
		return vr.Err()
	case Int2Oid:
		v := vr.ReadInt16()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint", v, v))
			return vr.Err()
		}
		*s.v = uint(v)
		return vr.Err()
	case Int4Oid:
		v := vr.ReadInt32()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint", v, v))
			return vr.Err()
		}
		*s.v = uint(v)
		return vr.Err()
	case Int8Oid:
		v := vr.ReadInt64()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint", v, v))
			return vr.Err()
		}
		*s.v = uint(v)
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into uint", vr.Type().DataType)))
		return vr.Err()
	}
}
