package pgtypes

import (
	"fmt"
	"math"

	"github.com/wdamron/pgx"
)

type g_int16Scanner struct {
	v *int16
}

func IntoInt16(v *int16) pgx.Scanner {
	return g_int16Scanner{v}
}

func (s g_int16Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int16"))
		return vr.Err()
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return vr.Err()
	}

	switch vr.Type().DataType {
	case BoolOid:
		*s.v = int16(vr.ReadByte())
		return vr.Err()
	case Int2Oid:
		*s.v = vr.ReadInt16()
		return vr.Err()
	case Int4Oid:
		v := vr.ReadInt32()
		if v < math.MinInt16 || v > math.MaxInt16 {
			vr.Fatal(fmt.Errorf("%T %d out of range for int16", v, v))
			return vr.Err()
		}
		*s.v = int16(v)
		return vr.Err()
	case Int8Oid:
		v := vr.ReadInt64()
		if v < math.MinInt16 || v > math.MaxInt16 {
			vr.Fatal(fmt.Errorf("%T %d out of range for int16", v, v))
			return vr.Err()
		}
		*s.v = int16(v)
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into int16", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_uint16Scanner struct {
	v *uint16
}

func IntoUint16(v *uint16) pgx.Scanner {
	return g_uint16Scanner{v}
}

func (s g_uint16Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int16"))
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
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint16", v, v))
			return vr.Err()
		}
		*s.v = uint16(v)
		return vr.Err()
	case Int2Oid:
		v := vr.ReadInt16()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint16", v, v))
			return vr.Err()
		}
		*s.v = uint16(v)
		return vr.Err()
	case Int4Oid:
		v := vr.ReadInt32()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint16", v, v))
			return vr.Err()
		}
		if v > math.MaxUint16 {
			vr.Fatal(fmt.Errorf("%T %d is larger than max uint16", v, v))
			return vr.Err()
		}
		*s.v = uint16(v)
		return vr.Err()
	case Int8Oid:
		v := vr.ReadInt64()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint16", v, v))
			return vr.Err()
		}
		if v > math.MaxUint16 {
			vr.Fatal(fmt.Errorf("%T %d is larger than max uint16", v, v))
			return vr.Err()
		}
		*s.v = uint16(v)
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into uint16", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_int32Scanner struct {
	v *int32
}

func IntoInt32(v *int32) pgx.Scanner {
	return g_int32Scanner{v}
}

func (s g_int32Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int32"))
		return vr.Err()
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return vr.Err()
	}

	switch vr.Type().DataType {
	case BoolOid:
		*s.v = int32(vr.ReadByte())
		return vr.Err()
	case Int2Oid:
		*s.v = int32(vr.ReadInt16())
		return vr.Err()
	case Int4Oid:
		*s.v = int32(vr.ReadInt32())
		return vr.Err()
	case Int8Oid:
		v := vr.ReadInt64()
		if v < math.MinInt32 || v > math.MaxInt32 {
			vr.Fatal(fmt.Errorf("%T %d out of range for int32", v, v))
			return vr.Err()
		}
		*s.v = int32(v)
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into int32", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_uint32Scanner struct {
	v *uint32
}

func IntoUint32(v *uint32) pgx.Scanner {
	return g_uint32Scanner{v}
}

func (s g_uint32Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into uint32"))
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
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint32", v, v))
			return vr.Err()
		}
		*s.v = uint32(v)
		return vr.Err()
	case Int2Oid:
		v := vr.ReadInt16()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint32", v, v))
			return vr.Err()
		}
		*s.v = uint32(v)
		return vr.Err()
	case Int4Oid:
		v := vr.ReadInt32()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint32", v, v))
			return vr.Err()
		}
		*s.v = uint32(v)
		return vr.Err()
	case Int8Oid:
		v := vr.ReadInt64()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint32", v, v))
			return vr.Err()
		}
		if v > math.MaxUint32 {
			vr.Fatal(fmt.Errorf("%T %d is larger than max uint32", v, v))
			return vr.Err()
		}
		*s.v = uint32(v)
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into uint32", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_int64Scanner struct {
	v *int64
}

func IntoInt64(v *int64) pgx.Scanner {
	return g_int64Scanner{v}
}

func (s g_int64Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int64"))
		return vr.Err()
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return vr.Err()
	}

	switch vr.Type().DataType {
	case BoolOid:
		*s.v = int64(vr.ReadByte())
		return vr.Err()
	case Int2Oid:
		*s.v = int64(vr.ReadInt16())
		return vr.Err()
	case Int4Oid:
		*s.v = int64(vr.ReadInt32())
		return vr.Err()
	case Int8Oid:
		*s.v = int64(vr.ReadInt64())
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into int64", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_uint64Scanner struct {
	v *uint64
}

func IntoUint64(v *uint64) pgx.Scanner {
	return g_uint64Scanner{v}
}

func (s g_uint64Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into uint64"))
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
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint64", v, v))
			return vr.Err()
		}
		*s.v = uint64(v)
		return vr.Err()
	case Int2Oid:
		v := vr.ReadInt16()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint64", v, v))
			return vr.Err()
		}
		*s.v = uint64(v)
		return vr.Err()
	case Int4Oid:
		v := vr.ReadInt32()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint64", v, v))
			return vr.Err()
		}
		*s.v = uint64(v)
		return vr.Err()
	case Int8Oid:
		v := vr.ReadInt64()
		if v < 0 {
			vr.Fatal(fmt.Errorf("Cannot decode negative value into uint64", v, v))
			return vr.Err()
		}
		*s.v = uint64(v)
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into uint64", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_float32Scanner struct {
	v *float32
}

func IntoFloat32(v *float32) pgx.Scanner {
	return g_float32Scanner{v}
}

func (s g_float32Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into float32"))
		return vr.Err()
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return vr.Err()
	}

	switch vr.Type().DataType {
	case Float4Oid:
		*s.v = math.Float32frombits(uint32(vr.ReadInt32()))
		return vr.Err()
	case Float8Oid:
		v := math.Float64frombits(uint64(vr.ReadInt64()))
		if math.Abs(v) < math.SmallestNonzeroFloat32 || v > math.MaxFloat32 {
			vr.Fatal(fmt.Errorf("Value %d out of range for float32", v))
			return vr.Err()
		}
		*s.v = float32(v)
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into float32", vr.Type().DataType)))
		return vr.Err()
	}
}

type g_float64Scanner struct {
	v *float64
}

func IntoFloat64(v *float64) pgx.Scanner {
	return g_float64Scanner{v}
}

func (s g_float64Scanner) Scan(vr *pgx.ValueReader) error {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into float64"))
		return vr.Err()
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return vr.Err()
	}

	switch vr.Type().DataType {
	case Float4Oid:
		*s.v = float64(math.Float32frombits(uint32(vr.ReadInt32())))
		return vr.Err()
	case Float8Oid:
		*s.v = math.Float64frombits(uint64(vr.ReadInt64()))
		return vr.Err()
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into float64", vr.Type().DataType)))
		return vr.Err()
	}
}
