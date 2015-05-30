package pgtypes

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"github.com/wdamron/pgx"
)

type ScannerFunc func(*pgx.ValueReader) error

func (fn ScannerFunc) Scan(vr *pgx.ValueReader) error {
	return fn(vr)
}

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

type boolScanner struct {
	v *bool
}

func BoolScanner(v *bool) pgx.Scanner {
	return boolScanner{v}
}

func (s boolScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeBool(vr)
	return vr.Err()
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

type int2Scanner struct {
	v *int16
}

func Int2Scanner(v *int16) pgx.Scanner {
	return int2Scanner{v}
}

func (s int2Scanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeInt2(vr)
	return vr.Err()
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

type int4Scanner struct {
	v *int32
}

func Int4Scanner(v *int32) pgx.Scanner {
	return int4Scanner{v}
}

func (s int4Scanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeInt4(vr)
	return vr.Err()
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

type int8Scanner struct {
	v *int64
}

func Int8Scanner(v *int64) pgx.Scanner {
	return int8Scanner{v}
}

func (s int8Scanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeInt8(vr)
	return vr.Err()
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

type float4Scanner struct {
	v *float32
}

func Float4Scanner(v *float32) pgx.Scanner {
	return float4Scanner{v}
}

func (s float4Scanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeFloat4(vr)
	return vr.Err()
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

type float8Scanner struct {
	v *float64
}

func Float8Scanner(v *float64) pgx.Scanner {
	return float8Scanner{v}
}

func (s float8Scanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeFloat8(vr)
	return vr.Err()
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

type byteaScanner struct {
	v *[]byte
}

func ByteaScanner(v *[]byte) pgx.Scanner {
	return byteaScanner{v}
}

func (s byteaScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeBytea(vr)
	return vr.Err()
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

type textScanner struct {
	v *string
}

func TextScanner(v *string) pgx.Scanner {
	return textScanner{v}
}

func (s textScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeText(vr)
	return vr.Err()
}

func DecodeText(vr *pgx.ValueReader) string {
	return decodeString(vr)
}

func VarcharScanner(v *string) pgx.Scanner {
	return TextScanner(v)
}

func DecodeVarchar(vr *pgx.ValueReader) string {
	return decodeString(vr)
}

type dateScanner struct {
	v *time.Time
}

func DateScanner(v *time.Time) pgx.Scanner {
	return dateScanner{v}
}

func (s dateScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeDate(vr)
	return vr.Err()
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

type timestampScanner struct {
	v *time.Time
}

func TimestampScanner(v *time.Time) pgx.Scanner {
	return timestampScanner{v}
}

func (s timestampScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeTimestamp(vr)
	return vr.Err()
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

type timestampTzScanner struct {
	v *time.Time
}

func TimestampTzScanner(v *time.Time) pgx.Scanner {
	return timestampTzScanner{v}
}

func (s timestampTzScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeTimestampTz(vr)
	return vr.Err()
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

type oidScanner struct {
	v *pgx.Oid
}

func OidScanner(v *pgx.Oid) pgx.Scanner {
	return oidScanner{v}
}

func (s oidScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeOid(vr)
	return vr.Err()
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

type jsonScanner struct {
	v interface{}
}

func (s jsonScanner) Scan(vr *pgx.ValueReader) error {
	b := decodeBytes(vr)
	if vr.Err() != nil {
		return vr.Err()
	}
	return json.Unmarshal(b, s.v)
}

func JSONScanner(v interface{}) pgx.Scanner {
	return jsonScanner{v}
}

type jsonScannerString struct {
	v *string
}

func JSONScannerString(v *string) pgx.Scanner {
	return jsonScannerString{v}
}

func (s jsonScannerString) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeJSONString(vr)
	return vr.Err()
}

func DecodeJSONString(vr *pgx.ValueReader) string {
	return decodeString(vr)
}

type jsonScannerBytes struct {
	v *[]byte
}

func JSONScannerBytes(v *[]byte) pgx.Scanner {
	return jsonScannerBytes{v}
}

func (s jsonScannerBytes) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeJSONBytes(vr)
	return vr.Err()
}

func DecodeJSONBytes(vr *pgx.ValueReader) []byte {
	return decodeBytes(vr)
}

type uuidScanner struct {
	v *uuid.UUID
}

func UUIDScanner(v *uuid.UUID) pgx.Scanner {
	return uuidScanner{v}
}

func (s uuidScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeUUID(vr)
	return vr.Err()
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

type uuidScannerString struct {
	v *string
}

func UUIDScannerString(v *string) pgx.Scanner {
	return uuidScannerString{v}
}

func (s uuidScannerString) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeUUIDString(vr)
	return vr.Err()
}

func DecodeUUIDString(vr *pgx.ValueReader) string {
	u := DecodeUUID(vr)
	if vr.Err() != nil {
		return ""
	}
	return u.String()
}

type hstoreScanner struct {
	v *pgx.Hstore
}

func HstoreScanner(v *pgx.Hstore) pgx.Scanner {
	return hstoreScanner{v}
}

func (s hstoreScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeHstore(vr)
	return vr.Err()
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

type hstoreMapScanner struct {
	v *map[string]string
}

func HstoreMapScanner(v *map[string]string) pgx.Scanner {
	return hstoreMapScanner{v}
}

func (s hstoreMapScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeHstoreMap(vr)
	return vr.Err()
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

type boolArrayScanner struct {
	v *[]bool
}

func BoolArrayScanner(v *[]bool) pgx.Scanner {
	return boolArrayScanner{v}
}

func (s boolArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeBoolArray(vr)
	return vr.Err()
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

type int2ArrayScanner struct {
	v *[]int16
}

func Int2ArrayScanner(v *[]int16) pgx.Scanner {
	return int2ArrayScanner{v}
}

func (s int2ArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeInt2Array(vr)
	return vr.Err()
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

type int4ArrayScanner struct {
	v *[]int32
}

func Int4ArrayScanner(v *[]int32) pgx.Scanner {
	return int4ArrayScanner{v}
}

func (s int4ArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeInt4Array(vr)
	return vr.Err()
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

type int8ArrayScanner struct {
	v *[]int64
}

func Int8ArrayScanner(v *[]int64) pgx.Scanner {
	return int8ArrayScanner{v}
}

func (s int8ArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeInt8Array(vr)
	return vr.Err()
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

type float4ArrayScanner struct {
	v *[]float32
}

func Float4ArrayScanner(v *[]float32) pgx.Scanner {
	return float4ArrayScanner{v}
}

func (s float4ArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeFloat4Array(vr)
	return vr.Err()
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

type float8ArrayScanner struct {
	v *[]float64
}

func Float8ArrayScanner(v *[]float64) pgx.Scanner {
	return float8ArrayScanner{v}
}

func (s float8ArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeFloat8Array(vr)
	return vr.Err()
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

type textArrayScanner struct {
	v *[]string
}

func TextArrayScanner(v *[]string) pgx.Scanner {
	return textArrayScanner{v}
}

func (s textArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeTextArray(vr)
	return vr.Err()
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

func VarcharArrayScanner(v *[]string) pgx.Scanner {
	return TextArrayScanner(v)
}

func DecodeVarcharArray(vr *pgx.ValueReader) []string {
	return DecodeTextArray(vr)
}

type timestampArrayScanner struct {
	v *[]time.Time
}

func TimestampArrayScanner(v *[]time.Time) pgx.Scanner {
	return timestampArrayScanner{v}
}

func (s timestampArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeTimestampArray(vr)
	return vr.Err()
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

func TimestampTzArrayScanner(v *[]time.Time) pgx.Scanner {
	return TimestampArrayScanner(v)
}

func DecodeTimestampTzArray(vr *pgx.ValueReader) []time.Time {
	return DecodeTimestampArray(vr)
}

type uuidArrayScanner struct {
	v *[]uuid.UUID
}

func UUIDArrayScanner(v *[]uuid.UUID) pgx.Scanner {
	return uuidArrayScanner{v}
}

func (s uuidArrayScanner) Scan(vr *pgx.ValueReader) error {
	*s.v = DecodeUUIDArray(vr)
	return vr.Err()
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
