package pgtypes

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"github.com/wdamron/pgx"
)

func decodeBytes(vr *pgx.ValueReader) []byte {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into []byte"))
		return nil
	}
	return vr.ReadBytes(vr.Len())
}

func decodeString(vr *pgx.ValueReader) string {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into string"))
		return ""
	}
	return vr.ReadString(vr.Len())
}

func decode1dArrayHeader(vr *pgx.ValueReader) (length int32, err error) {
	numDims := vr.ReadInt32()
	if numDims > 1 {
		return 0, pgx.ProtocolError(fmt.Sprintf("Expected array to have 0 or 1 dimension, but it had %v", numDims))
	}

	vr.ReadInt32() // 0 if no nulls / 1 if there is one or more nulls -- but we don't care
	vr.ReadInt32() // element oid

	if numDims == 0 {
		return 0, nil
	}

	length = vr.ReadInt32()

	idxFirstElem := vr.ReadInt32()
	if idxFirstElem != 1 {
		return 0, pgx.ProtocolError(fmt.Sprintf("Expected array's first element to start a index 1, but it is %d", idxFirstElem))
	}

	return length, nil
}

func DecodeBool(vr *pgx.ValueReader) bool {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into bool"))
		return false
	}

	if vr.Type().DataType != BoolOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into bool", vr.Type().DataType)))
		return false
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return false
	}

	if vr.Len() != 1 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an bool: %d", vr.Len())))
		return false
	}

	b := vr.ReadByte()
	return b != 0
}

func DecodeInt2(vr *pgx.ValueReader) int16 {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int16"))
		return 0
	}

	if vr.Type().DataType != Int2Oid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into int16", vr.Type().DataType)))
		return 0
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return 0
	}

	if vr.Len() != 2 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an int2: %d", vr.Len())))
		return 0
	}

	return vr.ReadInt16()
}

func DecodeInt4(vr *pgx.ValueReader) int32 {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int32"))
		return 0
	}

	if vr.Type().DataType != Int4Oid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into int32", vr.Type().DataType)))
		return 0
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return 0
	}

	if vr.Len() != 4 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an int4: %d", vr.Len())))
		return 0
	}

	return vr.ReadInt32()
}

func DecodeInt8(vr *pgx.ValueReader) int64 {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into int64"))
		return 0
	}

	if vr.Type().DataType != Int8Oid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into int8", vr.Type().DataType)))
		return 0
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return 0
	}

	if vr.Len() != 8 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an int8: %d", vr.Len())))
		return 0
	}

	return vr.ReadInt64()
}

func DecodeFloat4(vr *pgx.ValueReader) float32 {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into float32"))
		return 0
	}

	if vr.Type().DataType != Float4Oid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into float32", vr.Type().DataType)))
		return 0
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return 0
	}

	if vr.Len() != 4 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an float4: %d", vr.Len())))
		return 0
	}

	i := vr.ReadInt32()
	return math.Float32frombits(uint32(i))
}

func DecodeFloat8(vr *pgx.ValueReader) float64 {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into float64"))
		return 0
	}

	if vr.Type().DataType != Float8Oid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into float64", vr.Type().DataType)))
		return 0
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return 0
	}

	if vr.Len() != 8 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an float8: %d", vr.Len())))
		return 0
	}

	i := vr.ReadInt64()
	return math.Float64frombits(uint64(i))
}

func DecodeBytea(vr *pgx.ValueReader) []byte {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != ByteaOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []byte", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	return vr.ReadBytes(vr.Len())
}

func DecodeText(vr *pgx.ValueReader) string {
	return decodeString(vr)
}

func DecodeVarchar(vr *pgx.ValueReader) string {
	return decodeString(vr)
}

func DecodeDate(vr *pgx.ValueReader) time.Time {
	var zeroTime time.Time

	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into time.Time"))
		return zeroTime
	}

	if vr.Type().DataType != DateOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into time.Time", vr.Type().DataType)))
		return zeroTime
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return zeroTime
	}

	if vr.Len() != 4 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an date: %d", vr.Len())))
	}
	dayOffset := vr.ReadInt32()
	return time.Date(2000, 1, int(1+dayOffset), 0, 0, 0, 0, time.Local)
}

func DecodeTimestamp(vr *pgx.ValueReader) time.Time {
	var zeroTime time.Time

	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into timestamp"))
		return zeroTime
	}

	if vr.Type().DataType != TimestampOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into time.Time", vr.Type().DataType)))
		return zeroTime
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return zeroTime
	}

	if vr.Len() != 8 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an timestamp: %d", vr.Len())))
	}

	microsecSinceY2K := vr.ReadInt64()
	microsecSinceUnixEpoch := microsecFromUnixEpochToY2K + microsecSinceY2K
	return time.Unix(microsecSinceUnixEpoch/1000000, (microsecSinceUnixEpoch%1000000)*1000)
}

func DecodeTimestampTz(vr *pgx.ValueReader) time.Time {
	var zeroTime time.Time

	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into time.Time"))
		return zeroTime
	}

	if vr.Type().DataType != TimestampTzOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into time.Time", vr.Type().DataType)))
		return zeroTime
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return zeroTime
	}

	if vr.Len() != 8 {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an timestamptz: %d", vr.Len())))
		return zeroTime
	}

	microsecSinceY2K := vr.ReadInt64()
	microsecSinceUnixEpoch := microsecFromUnixEpochToY2K + microsecSinceY2K
	return time.Unix(microsecSinceUnixEpoch/1000000, (microsecSinceUnixEpoch%1000000)*1000)
}

func DecodeOid(vr *pgx.ValueReader) pgx.Oid {
	if vr.Len() == -1 {
		vr.Fatal(pgx.ProtocolError("Cannot decode null into Oid"))
		return 0
	}

	if vr.Type().DataType != OidOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into pgx.Oid", vr.Type().DataType)))
		return 0
	}

	// Oid needs to decode text format because it is used in loadPgTypes
	switch vr.Type().FormatCode {
	case TextFormatCode:
		s := vr.ReadString(vr.Len())
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received invalid Oid: %v", s)))
		}
		return pgx.Oid(n)
	case BinaryFormatCode:
		if vr.Len() != 4 {
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an Oid: %d", vr.Len())))
			return 0
		}
		return pgx.Oid(vr.ReadInt32())
	default:
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return pgx.Oid(0)
	}
}

func DecodeJSONString(vr *pgx.ValueReader) string {
	return decodeString(vr)
}

func DecodeJSONBytes(vr *pgx.ValueReader) []byte {
	return decodeBytes(vr)
}

func DecodeUUID(vr *pgx.ValueReader) uuid.UUID {
	var u uuid.UUID
	switch vr.Len() {
	case -1:
		vr.Fatal(pgx.ProtocolError("Cannot decode null into uuid.UUID"))
		return u
	case 16:
		u, err := uuid.FromBytes(vr.ReadBytes(16))
		if err != nil {
			vr.Fatal(err)
		}
		return u
	default:
		vr.Fatal(pgx.ProtocolError("invalid byte-length for uuid (should be 16)"))
		return u
	}
}

func DecodeUUIDString(vr *pgx.ValueReader) string {
	u := DecodeUUID(vr)
	if vr.Err() != nil {
		return ""
	}
	return u.String()
}

func DecodeHstore(vr *pgx.ValueReader) pgx.Hstore {
	size := int(vr.ReadInt32())
	h := make(pgx.Hstore, size)
	for i := 0; i < size; i++ {
		k := vr.ReadString(vr.ReadInt32())
		v := vr.ReadString(vr.ReadInt32())
		if vr.Err() != nil {
			return nil
		}
		h[k] = v
	}
	return h
}

func DecodeHstoreMap(vr *pgx.ValueReader) map[string]string {
	size := int(vr.ReadInt32())
	h := make(map[string]string, size)
	for i := 0; i < size; i++ {
		k := vr.ReadString(vr.ReadInt32())
		v := vr.ReadString(vr.ReadInt32())
		if vr.Err() != nil {
			return nil
		}
		h[k] = v
	}
	return h
}

func DecodeBoolArray(vr *pgx.ValueReader) []bool {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != BoolArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []bool", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]bool, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 1:
			if vr.ReadByte() == 1 {
				a[i] = true
			}
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an bool element: %d", elSize)))
			return nil
		}
	}

	return a
}

func DecodeInt2Array(vr *pgx.ValueReader) []int16 {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != Int2ArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []int16", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]int16, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 2:
			a[i] = vr.ReadInt16()
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an int2 element: %d", elSize)))
			return nil
		}
	}

	return a
}

func DecodeInt4Array(vr *pgx.ValueReader) []int32 {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != Int4ArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []int32", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]int32, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 4:
			a[i] = vr.ReadInt32()
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an int4 element: %d", elSize)))
			return nil
		}
	}

	return a
}

func DecodeInt8Array(vr *pgx.ValueReader) []int64 {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != Int8ArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []int64", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]int64, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 8:
			a[i] = vr.ReadInt64()
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an int8 element: %d", elSize)))
			return nil
		}
	}

	return a
}

func DecodeFloat4Array(vr *pgx.ValueReader) []float32 {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != Float4ArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []float32", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]float32, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 4:
			n := vr.ReadInt32()
			a[i] = math.Float32frombits(uint32(n))
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an float4 element: %d", elSize)))
			return nil
		}
	}

	return a
}

func DecodeFloat8Array(vr *pgx.ValueReader) []float64 {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != Float8ArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []float64", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]float64, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 8:
			n := vr.ReadInt64()
			a[i] = math.Float64frombits(uint64(n))
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an float4 element: %d", elSize)))
			return nil
		}
	}

	return a
}

func DecodeTextArray(vr *pgx.ValueReader) []string {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != TextArrayOid && vr.Type().DataType != VarcharArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []string", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]string, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		if elSize == -1 {
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		}

		a[i] = vr.ReadString(elSize)
	}

	return a
}

func DecodeVarcharArray(vr *pgx.ValueReader) []string {
	return DecodeTextArray(vr)
}

func DecodeTimestampArray(vr *pgx.ValueReader) []time.Time {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != TimestampArrayOid && vr.Type().DataType != TimestampTzArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []time.Time", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]time.Time, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 8:
			microsecSinceY2K := vr.ReadInt64()
			microsecSinceUnixEpoch := microsecFromUnixEpochToY2K + microsecSinceY2K
			a[i] = time.Unix(microsecSinceUnixEpoch/1000000, (microsecSinceUnixEpoch%1000000)*1000)
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for an time.Time element: %d", elSize)))
			return nil
		}
	}

	return a
}

func DecodeTimestampTzArray(vr *pgx.ValueReader) []time.Time {
	return DecodeTimestampArray(vr)
}

func DecodeUUIDArray(vr *pgx.ValueReader) []uuid.UUID {
	if vr.Len() == -1 {
		return nil
	}

	if vr.Type().DataType != UUIDArrayOid {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Cannot decode oid %v into []uuid.UUID", vr.Type().DataType)))
		return nil
	}

	if vr.Type().FormatCode != BinaryFormatCode {
		vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Unknown field description format code: %v", vr.Type().FormatCode)))
		return nil
	}

	numElems, err := decode1dArrayHeader(vr)
	if err != nil {
		vr.Fatal(err)
		return nil
	}

	a := make([]uuid.UUID, int(numElems))
	for i := 0; i < len(a); i++ {
		elSize := vr.ReadInt32()
		switch elSize {
		case 16:
			if err = a[i].UnmarshalBinary(vr.ReadBytes(16)); err != nil {
				vr.Fatal(err)
			}
		case -1:
			vr.Fatal(pgx.ProtocolError("Cannot decode null element"))
			return nil
		default:
			vr.Fatal(pgx.ProtocolError(fmt.Sprintf("Received an invalid size for uuid.UUID element: %d", elSize)))
			return nil
		}
	}

	return a
}
