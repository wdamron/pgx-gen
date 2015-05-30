package pgtypes

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/satori/go.uuid"
	"github.com/wdamron/pgx"
)

const (
	len1  = "\x00\x00\x00\x01"
	len2  = "\x00\x00\x00\x02"
	len4  = "\x00\x00\x00\x04"
	len8  = "\x00\x00\x00\x08"
	len10 = "\x00\x00\x00\x0a"
	len16 = "\x00\x00\x00\x10"
)

type boolEncoder struct {
	v bool
}

func BoolEncoder(v bool) pgx.Encoder {
	return &boolEncoder{v}
}

func (e *boolEncoder) FormatCode() int16 { return 1 }

func (e *boolEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != BoolOid {
		return fmt.Errorf("BoolEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeBool(wbuf, e.v)
}

func EncodeBool(wbuf *pgx.WriteBuf, v bool) error {
	var cast byte
	if v {
		cast = 1
	}
	wbuf.WriteBytes(append([]byte(len1), cast))
	return nil
}

type int2Encoder struct {
	v int16
}

func Int2Encoder(v int16) pgx.Encoder {
	return &int2Encoder{v}
}

func (e *int2Encoder) FormatCode() int16 { return 1 }

func (e *int2Encoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Int2Oid {
		return fmt.Errorf("Int2Encoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeInt2(wbuf, e.v)
}

func EncodeInt2(wbuf *pgx.WriteBuf, v int16) error {
	wbuf.WriteBytes(append([]byte(len2), byte(v>>8), byte(v)))
	return nil
}

type int4Encoder struct {
	v int32
}

func Int4Encoder(v int32) pgx.Encoder {
	return &int4Encoder{v}
}

func (e *int4Encoder) FormatCode() int16 { return 1 }

func (e *int4Encoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Int4Oid {
		return fmt.Errorf("Int4Encoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeInt4(wbuf, e.v)
}

func EncodeInt4(wbuf *pgx.WriteBuf, v int32) error {
	wbuf.WriteBytes(append([]byte(len4), byte(v>>24), byte(v>>16), byte(v>>8), byte(v)))
	return nil
}

type int8Encoder struct {
	v int64
}

func Int8Encoder(v int64) pgx.Encoder {
	return &int8Encoder{v}
}

func (e *int8Encoder) FormatCode() int16 { return 1 }

func (e *int8Encoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Int8Oid {
		return fmt.Errorf("Int8Encoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeInt8(wbuf, e.v)
}

func EncodeInt8(wbuf *pgx.WriteBuf, v int64) error {
	b := []byte(len8)
	b = append(b, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32))
	b = append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	wbuf.WriteBytes(b)
	return nil
}

type float4Encoder struct {
	v float32
}

func Float4Encoder(v float32) pgx.Encoder {
	return &float4Encoder{v}
}

func (e *float4Encoder) FormatCode() int16 { return 1 }

func (e *float4Encoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Float4Oid {
		return fmt.Errorf("Float4Encoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeFloat4(wbuf, e.v)
}

func EncodeFloat4(wbuf *pgx.WriteBuf, v float32) error {
	return EncodeInt4(wbuf, int32(math.Float32bits(v)))
}

type float8Encoder struct {
	v float64
}

func Float8Encoder(v float64) pgx.Encoder {
	return &float8Encoder{v}
}

func (e *float8Encoder) FormatCode() int16 { return 1 }

func (e *float8Encoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Float8Oid {
		return fmt.Errorf("Float8Encoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeFloat8(wbuf, e.v)
}

func EncodeFloat8(wbuf *pgx.WriteBuf, v float64) error {
	return EncodeInt8(wbuf, int64(math.Float64bits(v)))
}

type byteaEncoder struct {
	v []byte
}

func ByteaEncoder(v []byte) pgx.Encoder {
	return &byteaEncoder{v}
}

func (e *byteaEncoder) FormatCode() int16 { return 1 }

func (e *byteaEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != ByteaOid {
		return fmt.Errorf("ByteaEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeBytea(wbuf, e.v)
}

func EncodeBytea(wbuf *pgx.WriteBuf, v []byte) error {
	totalLen := 4 + len(v)
	b := make([]byte, totalLen)
	binary.BigEndian.PutUint32(b, uint32(len(v)))
	if len(v) != 0 {
		copy(b[4:totalLen], v)
	}
	wbuf.WriteBytes(b)
	return nil
}

type textEncoder struct {
	v string
}

func TextEncoder(v string) pgx.Encoder {
	return &textEncoder{v}
}

func (e *textEncoder) FormatCode() int16 { return 0 }

func (e *textEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	return EncodeText(wbuf, e.v)
}

func EncodeText(wbuf *pgx.WriteBuf, v string) error {
	totalLen := 4 + len(v)
	b := make([]byte, totalLen)
	binary.BigEndian.PutUint32(b, uint32(len(v)))
	if len(v) != 0 {
		copy(b[4:totalLen], v)
	}
	wbuf.WriteBytes(b)
	wbuf.WriteString(v)
	return nil
}

type textEncoderBytes struct {
	v []byte
}

func TextEncoderBytes(v []byte) pgx.Encoder {
	return &textEncoderBytes{v}
}

func (e *textEncoderBytes) FormatCode() int16 { return 0 }

func (e *textEncoderBytes) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	return EncodeTextBytes(wbuf, e.v)
}

func EncodeTextBytes(wbuf *pgx.WriteBuf, v []byte) error {
	totalLen := 4 + len(v)
	b := make([]byte, totalLen)
	binary.BigEndian.PutUint32(b, uint32(len(v)))
	if len(v) != 0 {
		copy(b[4:totalLen], v)
	}
	wbuf.WriteBytes(append(b, v...))
	return nil
}

type varcharEncoder struct {
	v string
}

func VarcharEncoder(v string) pgx.Encoder {
	return &varcharEncoder{v}
}

func (e *varcharEncoder) FormatCode() int16 { return 0 }

func (e *varcharEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != VarcharOid {
		return fmt.Errorf("VarcharEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeVarchar(wbuf, e.v)
}

func EncodeVarchar(wbuf *pgx.WriteBuf, v string) error {
	return EncodeText(wbuf, v)
}

type varcharEncoderBytes struct {
	v []byte
}

func VarcharEncoderBytes(v []byte) pgx.Encoder {
	return &varcharEncoderBytes{v}
}

func (e *varcharEncoderBytes) FormatCode() int16 { return 0 }

func (e *varcharEncoderBytes) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != VarcharOid {
		return fmt.Errorf("VarcharEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeTextBytes(wbuf, e.v)
}

type dateEncoder struct {
	v time.Time
}

func DateEncoder(v time.Time) pgx.Encoder {
	return &dateEncoder{v}
}

func (e *dateEncoder) FormatCode() int16 { return 0 }

func (e *dateEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != DateOid {
		return fmt.Errorf("DateEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeDate(wbuf, e.v)
}

func EncodeDate(wbuf *pgx.WriteBuf, v time.Time) error {
	wbuf.WriteString(len10 + v.Format("2006-01-02"))
	return nil
}

type timestampEncoder struct {
	v time.Time
}

func TimestampEncoder(v time.Time) pgx.Encoder {
	return &timestampEncoder{v}
}

func (e *timestampEncoder) FormatCode() int16 { return 1 }

func (e *timestampEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != TimestampOid {
		return fmt.Errorf("TimestampEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeTimestamp(wbuf, e.v)
}

func EncodeTimestamp(wbuf *pgx.WriteBuf, v time.Time) error {
	return EncodeTimestampTz(wbuf, v)
}

type timestampTzEncoder struct {
	v time.Time
}

func TimestampTzEncoder(v time.Time) pgx.Encoder {
	return &timestampTzEncoder{v}
}

func (e *timestampTzEncoder) FormatCode() int16 { return 1 }

func (e *timestampTzEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != TimestampTzOid {
		return fmt.Errorf("TimestampTzEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeTimestampTz(wbuf, e.v)
}

func EncodeTimestampTz(wbuf *pgx.WriteBuf, v time.Time) error {
	microsecSinceUnixEpoch := v.Unix()*1000000 + int64(v.Nanosecond())/1000
	microsecSinceY2K := microsecSinceUnixEpoch - microsecFromUnixEpochToY2K
	x := microsecSinceY2K
	b := []byte(len8)
	b = append(b, byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32))
	b = append(b, byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
	wbuf.WriteBytes(b)
	return nil
}

type oidEncoder struct {
	v pgx.Oid
}

func OidEncoder(v pgx.Oid) pgx.Encoder {
	return &oidEncoder{v}
}

func (e *oidEncoder) FormatCode() int16 { return 1 }

func (e *oidEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != OidOid {
		return fmt.Errorf("OidEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeOid(wbuf, e.v)
}

func EncodeOid(wbuf *pgx.WriteBuf, v pgx.Oid) error {
	return EncodeInt4(wbuf, int32(v))
}

func encodeArrayHeaderBytes(oid pgx.Oid, length, sizePerItem int) []byte {
	b := make([]byte, 24)
	binary.BigEndian.PutUint32(b[:4], uint32(20+length*sizePerItem))
	binary.BigEndian.PutUint32(b[4:8], 1)                // number of dimensions
	binary.BigEndian.PutUint32(b[8:12], 0)               // no nulls
	binary.BigEndian.PutUint32(b[12:16], uint32(oid))    // type of elements
	binary.BigEndian.PutUint32(b[16:20], uint32(length)) // number of elements
	binary.BigEndian.PutUint32(b[20:24], 1)              // index of first element
	return b
}

type boolArrayEncoder struct {
	v []bool
}

func BoolArrayEncoder(v []bool) pgx.Encoder {
	return &boolArrayEncoder{v}
}

func (e *boolArrayEncoder) FormatCode() int16 { return 1 }

func (e *boolArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != BoolArrayOid {
		return fmt.Errorf("BoolArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeBoolArray(wbuf, e.v)
}

func EncodeBoolArray(wbuf *pgx.WriteBuf, vs []bool) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(BoolOid, len(vs), 5))
	for _, v := range vs {
		if err := EncodeBool(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type int2ArrayEncoder struct {
	v []int16
}

func Int2ArrayEncoder(v []int16) pgx.Encoder {
	return &int2ArrayEncoder{v}
}

func (e *int2ArrayEncoder) FormatCode() int16 { return 1 }

func (e *int2ArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Int2ArrayOid {
		return fmt.Errorf("Int2ArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeInt2Array(wbuf, e.v)
}

func EncodeInt2Array(wbuf *pgx.WriteBuf, vs []int16) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(Int2Oid, len(vs), 6))
	for _, v := range vs {
		if err := EncodeInt2(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type int4ArrayEncoder struct {
	v []int32
}

func Int4ArrayEncoder(v []int32) pgx.Encoder {
	return &int4ArrayEncoder{v}
}

func (e *int4ArrayEncoder) FormatCode() int16 { return 1 }

func (e *int4ArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Int4ArrayOid {
		return fmt.Errorf("Int4ArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeInt4Array(wbuf, e.v)
}

func EncodeInt4Array(wbuf *pgx.WriteBuf, vs []int32) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(Int4Oid, len(vs), 8))
	for _, v := range vs {
		if err := EncodeInt4(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type int8ArrayEncoder struct {
	v []int64
}

func Int8ArrayEncoder(v []int64) pgx.Encoder {
	return &int8ArrayEncoder{v}
}

func (e *int8ArrayEncoder) FormatCode() int16 { return 1 }

func (e *int8ArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Int8ArrayOid {
		return fmt.Errorf("Int8ArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeInt8Array(wbuf, e.v)
}

func EncodeInt8Array(wbuf *pgx.WriteBuf, vs []int64) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(Int8Oid, len(vs), 12))
	for _, v := range vs {
		if err := EncodeInt8(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type float4ArrayEncoder struct {
	v []float32
}

func Float4ArrayEncoder(v []float32) pgx.Encoder {
	return &float4ArrayEncoder{v}
}

func (e *float4ArrayEncoder) FormatCode() int16 { return 1 }

func (e *float4ArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Float4ArrayOid {
		return fmt.Errorf("Float4ArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeFloat4Array(wbuf, e.v)
}

func EncodeFloat4Array(wbuf *pgx.WriteBuf, vs []float32) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(Int4Oid, len(vs), 8))
	for _, v := range vs {
		if err := EncodeFloat4(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type float8ArrayEncoder struct {
	v []float64
}

func Float8ArrayEncoder(v []float64) pgx.Encoder {
	return &float8ArrayEncoder{v}
}

func (e *float8ArrayEncoder) FormatCode() int16 { return 1 }

func (e *float8ArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != Float8ArrayOid {
		return fmt.Errorf("Float8ArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeFloat8Array(wbuf, e.v)
}

func EncodeFloat8Array(wbuf *pgx.WriteBuf, vs []float64) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(Int8Oid, len(vs), 12))
	for _, v := range vs {
		if err := EncodeFloat8(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type textArrayEncoder struct {
	v []string
}

func TextArrayEncoder(v []string) pgx.Encoder {
	return &textArrayEncoder{v}
}

func (e *textArrayEncoder) FormatCode() int16 { return 1 }

func (e *textArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != TextArrayOid {
		return fmt.Errorf("TextArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeTextArray(wbuf, e.v)
}

func EncodeTextArray(wbuf *pgx.WriteBuf, vs []string) error {
	return encodeTextArray(wbuf, vs, TextOid)
}

func encodeTextArray(wbuf *pgx.WriteBuf, vs []string, oid pgx.Oid) error {
	var totalStringSize int
	for _, v := range vs {
		totalStringSize += len(v)
	}
	size := 20 + len(vs)*4 + totalStringSize
	header := make([]byte, 4+20)
	binary.BigEndian.PutUint32(header[:4], uint32(size))
	binary.BigEndian.PutUint32(header[4:8], 1)                 // number of dimensions
	binary.BigEndian.PutUint32(header[8:12], 0)                // no nulls
	binary.BigEndian.PutUint32(header[12:16], uint32(oid))     // type of elements
	binary.BigEndian.PutUint32(header[16:20], uint32(len(vs))) // number of elements
	binary.BigEndian.PutUint32(header[20:24], 1)               // index of first element
	wbuf.WriteBytes(header)
	var enc func(*pgx.WriteBuf, string) error
	switch oid {
	default:
		enc = EncodeText
	case VarcharOid:
		enc = EncodeVarchar
	}
	for _, v := range vs {
		if err := enc(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type varcharArrayEncoder struct {
	v []string
}

func VarcharArrayEncoder(v []string) pgx.Encoder {
	return &varcharArrayEncoder{v}
}

func (e *varcharArrayEncoder) FormatCode() int16 { return 1 }

func (e *varcharArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != VarcharArrayOid {
		return fmt.Errorf("VarcharArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeVarcharArray(wbuf, e.v)
}

func EncodeVarcharArray(wbuf *pgx.WriteBuf, vs []string) error {
	return encodeTextArray(wbuf, vs, VarcharOid)
}

type timestampArrayEncoder struct {
	v []time.Time
}

func TimestampArrayEncoder(v []time.Time) pgx.Encoder {
	return &timestampArrayEncoder{v}
}

func (e *timestampArrayEncoder) FormatCode() int16 { return 1 }

func (e *timestampArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != TimestampArrayOid {
		return fmt.Errorf("TimestampArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeTimestampArray(wbuf, e.v)
}

func EncodeTimestampArray(wbuf *pgx.WriteBuf, vs []time.Time) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(TimestampOid, len(vs), 12))
	for _, v := range vs {
		if err := EncodeTimestamp(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type timestampTzArrayEncoder struct {
	v []time.Time
}

func TimestampTzArrayEncoder(v []time.Time) pgx.Encoder {
	return &timestampTzArrayEncoder{v}
}

func (e *timestampTzArrayEncoder) FormatCode() int16 { return 1 }

func (e *timestampTzArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != TimestampTzArrayOid {
		return fmt.Errorf("TimestampTzArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeTimestampTzArray(wbuf, e.v)
}

func EncodeTimestampTzArray(wbuf *pgx.WriteBuf, vs []time.Time) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(TimestampTzOid, len(vs), 12))
	for _, v := range vs {
		if err := EncodeTimestampTz(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type hstoreEncoder struct {
	v pgx.Hstore
}

func HstoreEncoder(v pgx.Hstore) pgx.Encoder {
	return &hstoreEncoder{v}
}

func HstoreMapEncoder(v map[string]string) pgx.Encoder {
	h := pgx.Hstore(v)
	return &hstoreEncoder{h}
}

func (e *hstoreEncoder) FormatCode() int16 { return 1 }

func (e *hstoreEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	return EncodeHstore(wbuf, e.v)
}

func EncodeHstore(wbuf *pgx.WriteBuf, kv pgx.Hstore) error {
	wbuf.WriteInt32(int32(len(kv)))
	for k, v := range kv {
		wbuf.WriteInt32(int32(len(k)))
		wbuf.WriteString(k)
		wbuf.WriteInt32(int32(len(v)))
		wbuf.WriteString(v)
	}
	return nil
}

type uuidEncoder struct {
	v uuid.UUID
}

func UUIDEncoder(v uuid.UUID) pgx.Encoder {
	return &uuidEncoder{v}
}

func UUIDEncoderString(v string) pgx.Encoder {
	u, _ := uuid.FromString(v)
	return &uuidEncoder{u}
}

func (e *uuidEncoder) FormatCode() int16 { return 1 }

func (e *uuidEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != UUIDOid {
		return fmt.Errorf("UUIDEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeUUID(wbuf, e.v)
}

func EncodeUUID(wbuf *pgx.WriteBuf, v uuid.UUID) error {
	wbuf.WriteBytes(append([]byte(len16), v[:16]...))
	return nil
}

type uuidArrayEncoder struct {
	v []uuid.UUID
}

func UUIDArrayEncoder(v []uuid.UUID) pgx.Encoder {
	return &uuidArrayEncoder{v}
}

func (e *uuidArrayEncoder) FormatCode() int16 { return 1 }

func (e *uuidArrayEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != UUIDArrayOid {
		return fmt.Errorf("UUIDArrayEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeUUIDArray(wbuf, e.v)
}

func EncodeUUIDArray(wbuf *pgx.WriteBuf, vs []uuid.UUID) error {
	wbuf.WriteBytes(encodeArrayHeaderBytes(UUIDOid, len(vs), 16))
	for _, v := range vs {
		if err := EncodeUUID(wbuf, v); err != nil {
			return err
		}
	}
	return nil
}

type jsonEncoder struct {
	v interface{}
}

func JSONEncoder(v interface{}) pgx.Encoder {
	return &jsonEncoder{v}
}

func (e *jsonEncoder) FormatCode() int16 { return 0 }

func (e *jsonEncoder) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != JSONOid {
		return fmt.Errorf("JSONEncoder.Encode cannot encode into OID: %d", oid)
	}

	return EncodeJSON(wbuf, e.v)
}

func EncodeJSON(wbuf *pgx.WriteBuf, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	wbuf.WriteInt32(int32(len(b)))
	wbuf.WriteBytes(b)
	return nil
}

type jsonEncoderString struct {
	v string
}

func JSONEncoderString(v string) pgx.Encoder {
	return &jsonEncoderString{v}
}

func (e *jsonEncoderString) FormatCode() int16 { return 0 }

func (e *jsonEncoderString) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != JSONOid {
		return fmt.Errorf("JSONEncoder.Encode cannot encode into OID: %d", oid)
	}
	wbuf.WriteInt32(int32(len(e.v)))
	wbuf.WriteString(e.v)
	return nil
}

type jsonEncoderBytes struct {
	v []byte
}

func JSONEncoderBytes(v []byte) pgx.Encoder {
	return &jsonEncoderBytes{v}
}

func (e *jsonEncoderBytes) FormatCode() int16 { return 0 }

func (e *jsonEncoderBytes) Encode(wbuf *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != JSONOid {
		return fmt.Errorf("JSONEncoder.Encode cannot encode into OID: %d", oid)
	}
	wbuf.WriteInt32(int32(len(e.v)))
	wbuf.WriteBytes(e.v)
	return nil
}
