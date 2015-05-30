package example

// Generated by pgxgen (see example.go)

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/satori/go.uuid"
	"github.com/wdamron/pgx"
	"github.com/wdamron/pgx-gen/pgtypes"
)

// PointTableType is the type of PointTable, which describes the table
// corresponding with type Point
type PointTableType struct {
	// UnboundEncoders are used by PointParamsEncoder.Bind to bind
	// query/statement parameters from a value of type Point
	UnboundEncoders [10]func(*Point) pgx.Encoder
	// Decoders can be used to decode a single column value from a
	// pgx.ValueReader into type Point
	Decoders [10]func(*Point, *pgx.ValueReader) error
	// Names contains an ordered list of column names
	Names [10]string
	// Types contains an ordered list of column types
	Types [10]string
	// Aliases contains an ordered list of column names aliased as hex-encoded
	// indexes, for faster look-ups during decoding
	Aliases [10]string
	// Formats contains an ordered list of column format codes (text=0, binary=1)
	Formats [10]int
	// Oids contains an ordered list of column oid codes (corresponding with
	// Postgres types)
	Oids [10]pgx.Oid
}

// PointTable describes the table corresponding with type Point
var PointTable = PointTableType{
	UnboundEncoders: [10]func(*Point) pgx.Encoder{
		// Encode v.X as varchar[]
		func(v *Point) pgx.Encoder {
			return pgtypes.VarcharArrayEncoder(v.X)
		},
		// Encode v.Y as int4
		func(v *Point) pgx.Encoder {
			return pgtypes.Int4Encoder(int32(*v.Y))
		},
		// Encode v.Z as int4
		func(v *Point) pgx.Encoder {
			return v.Z
		},
		// Encode v.H as hstore
		func(v *Point) pgx.Encoder {
			h := pgx.Hstore(*v.H)
			return pgtypes.HstoreEncoder(h)
		},
		// Encode v.H2 as hstore
		func(v *Point) pgx.Encoder {
			return v.H2
		},
		// Encode v.u as uuid
		func(v *Point) pgx.Encoder {
			return pgtypes.UUIDEncoderString(v.u)
		},
		// Encode v.u2 as uuid
		func(v *Point) pgx.Encoder {
			return pgtypes.UUIDEncoder(*v.u2)
		},
		// Encode v.j as json
		func(v *Point) pgx.Encoder {
			return pgtypes.JSONEncoderString(*v.j)
		},
		// Encode v.j2 as json
		func(v *Point) pgx.Encoder {
			return pgtypes.JSONEncoder(v.j2)
		},
		// Encode v.j3 as json
		func(v *Point) pgx.Encoder {
			return pgtypes.JSONEncoderBytes(v.j3)
		},
	},
	Decoders: [10]func(*Point, *pgx.ValueReader) error{
		// Decode column x::varchar[] into (*Point).X
		func(v *Point, vr *pgx.ValueReader) error {
			x := vr.DecodeVarcharArray()
			if vr.Err() != nil {
				return vr.Err()
			}
			v.X = x
			return nil
		},
		// Decode column y::int4 into (*Point).Y
		func(v *Point, vr *pgx.ValueReader) error {
			x := int64(vr.DecodeInt4())
			if vr.Err() != nil {
				return vr.Err()
			}
			*v.Y = x
			return nil
		},
		// Decode column z::int4 into (*Point).Z
		func(v *Point, vr *pgx.ValueReader) error {
			return v.Z.Scan(vr)
		},
		// Decode column h::hstore into (*Point).H
		func(v *Point, vr *pgx.ValueReader) error {
			h := pgx.Hstore(*v.H)
			return h.Scan(vr)
		},
		// Decode column h2::hstore into (*Point).H2
		func(v *Point, vr *pgx.ValueReader) error {
			return v.H2.Scan(vr)
		},
		// Decode column id::uuid into (*Point).u
		func(v *Point, vr *pgx.ValueReader) error {
			b := vr.ReadBytes(vr.ReadInt32())
			if vr.Err() != nil {
				return vr.Err()
			}
			if len(b) != 16 {
				return errors.New("invalid length for uuid (should be 16)")
			}
			u, err := uuid.FromBytes(b)
			if err != nil {
				return err
			}
			v.u = u.String()
			return nil
		},
		// Decode column id2::uuid into (*Point).u2
		func(v *Point, vr *pgx.ValueReader) error {
			b := vr.ReadBytes(vr.ReadInt32())
			if vr.Err() != nil {
				return vr.Err()
			}
			if len(b) != 16 {
				return errors.New("invalid length for uuid (should be 16)")
			}
			u, err := uuid.FromBytes(b)
			if err != nil {
				return err
			}
			*v.u2 = u
			return nil
		},
		// Decode column j::json into (*Point).j
		func(v *Point, vr *pgx.ValueReader) error {
			s := vr.ReadString(vr.ReadInt32())
			if vr.Err() != nil {
				return vr.Err()
			}
			*v.j = s
			return nil
		},
		// Decode column j2::json into (*Point).j2
		func(v *Point, vr *pgx.ValueReader) error {
			b := vr.ReadBytes(vr.ReadInt32())
			if vr.Err() != nil {
				return vr.Err()
			}
			return json.Unmarshal(b, &v.j2)
		},
		// Decode column j3::json into (*Point).j3
		func(v *Point, vr *pgx.ValueReader) error {
			b := vr.ReadBytes(vr.ReadInt32())
			if vr.Err() != nil {
				return vr.Err()
			}
			v.j3 = b
			return nil
		},
	},
	Names: [10]string{
		"x",
		"y",
		"z",
		"h",
		"h2",
		"id",
		"id2",
		"j",
		"j2",
		"j3",
	},
	Types: [10]string{
		"varchar[]",
		"int4",
		"int4",
		"hstore",
		"hstore",
		"uuid",
		"uuid",
		"json",
		"json",
		"json",
	},
	// Aliases contains an ordered list of column names aliased as hex-encoded
	// indexes, for faster look-ups during decoding
	Aliases: [10]string{
		"x as __00::varchar[]",
		"y as __01::int4",
		"z as __02::int4",
		"h as __03::hstore",
		"h2 as __04::hstore",
		"id as __05::uuid",
		"id2 as __06::uuid",
		"j as __07::json",
		"j2 as __08::json",
		"j3 as __09::json",
	},
	Formats: [10]int{1, 1, 1, 0, 0, 1, 1, 0, 0, 0},
	Oids: [10]pgx.Oid{
		pgx.VarcharArrayOid,
		pgx.Int4Oid,
		pgx.Int4Oid,
		pgx.Oid(0),
		pgx.Oid(0),
		pgx.Oid(2950),
		pgx.Oid(2950),
		pgx.Oid(114),
		pgx.Oid(114),
		pgx.Oid(114),
	},
}

// Index returns the index of the column in PointTable with the given name.
// If no matching column is found, the returned index will be -1.
func (t *PointTableType) Index(colname string) int {
	switch colname {
	case "x":
		return 0
	case "y":
		return 1
	case "z":
		return 2
	case "h":
		return 3
	case "h2":
		return 4
	case "id":
		return 5
	case "id2":
		return 6
	case "j":
		return 7
	case "j2":
		return 8
	case "j3":
		return 9
	}
	return -1
}

// Indexes returns a slice of indexes of the given columns in PointTable
// with the given name. If any of the columns are not found, an error will be
// returned and the returned slice of indexes will be nil.
func (t *PointTableType) Indexes(colnames ...string) ([]int, error) {
	indexes := make([]int, len(colnames))
	for i, colname := range colnames {
		index := t.Index(colname)
		if index < 0 {
			return nil, errors.New("column " + colname + " not found in PointTable")
		}
		indexes[i] = index
	}
	return indexes, nil
}

// Alias aliases column names as hex-encoded indexes, for faster look-ups during
// decoding. If no column names are provided, all columns will be aliased, in
// which case AliasAll may be a faster alternative.
func (t *PointTableType) Alias(colnames ...string) ([]string, error) {
	aliases := []string{}
	// If no columns are specified, alias all columns:
	if len(colnames) == 0 {
		return PointTable.Aliases[:10], nil
	}
	indexes, err := PointTable.Indexes(colnames...)
	if err != nil {
		return nil, err
	}
	for _, index := range indexes {
		aliases = append(aliases, PointTable.Aliases[index])
	}
	return aliases, nil
}

// AliasAll aliases column names as hex-encoded indexes, for faster look-ups
// during decoding
func (t *PointTableType) AliasAll() string {
	return "x as __00::varchar[], y as __01::int4, z as __02::int4, h as __03::hstore, h2 as __04::hstore, id as __05::uuid, id2 as __06::uuid, j as __07::json, j2 as __08::json, j3 as __09::json"
}

// DecodeRow decodes a single row/result from r into v.
// If an error is returned, the caller should call Rows.Close()
func (v *Point) DecodeRow(r *pgx.Rows) error {
	for _ = range r.FieldDescriptions() {
		vr, ok := r.NextColumn()
		if !ok {
			if vr != nil && vr.Err() != nil {
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
			if index < 0 || index > len(PointTable.Decoders)-1 {
				return errors.New("column decoder index out of range")
			}
			dec := PointTable.Decoders[index]
			if err = dec(v, vr); err != nil {
				return err
			}
			continue
		}

		// Slow path:
		index := PointTable.Index(colname)
		if index < 0 {
			return errors.New("column decoder for " + colname + " not found in PointTable")
		}
		if err := PointTable.Decoders[index](v, vr); err != nil {
			return err
		}
	}
	return nil
}

// PointParamsEncoder binds query/statement parameters from a value
// of type Point
type PointParamsEncoder struct {
	Indexes []int
}

// Encoder creates an unbound instance of type PointParamsEncoder
// for the columns/fields named by colnames
func (t *PointTableType) Encoder(colnames ...string) (*PointParamsEncoder, error) {
	indexes, err := PointTable.Indexes(colnames...)
	if err != nil {
		return nil, err
	}

	pe := &PointParamsEncoder{Indexes: indexes}
	return pe, nil
}

// Bind binds query/statement parameter encoders from v
func (pe PointParamsEncoder) Bind(v *Point) ([]pgx.Encoder, error) {
	encoders := make([]pgx.Encoder, len(pe.Indexes))
	for i, index := range pe.Indexes {
		if index < 0 || index > len(PointTable.UnboundEncoders) {
			return nil, errors.New("column encoder index out of range")
		}
		encoders[i] = PointTable.UnboundEncoders[index](v)
	}
	return encoders, nil
}