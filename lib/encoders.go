package pgxgen

var Encoders = map[string]OpMap{
	"bool": {
		"bool": OpPass,
	},
	"*bool": {
		"bool": OpDerefPass,
	},
	"int": {
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64,
	},
	"*int": {
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64,
	},
	"uint": {
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64 | OpCheckOverflow,
	},
	"*uint": {
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64 | OpCheckOverflow,
	},
	"int16": {
		"int2": OpPass,
		"int4": OpPass | OpCastInt32,
		"int8": OpPass | OpCastInt64,
	},
	"*int16": {
		"int2": OpDerefPass,
		"int4": OpDerefPass | OpCastInt32,
		"int8": OpDerefPass | OpCastInt64,
	},
	"uint16": {
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32,
		"int8": OpPass | OpCastInt64,
	},
	"*uint16": {
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32,
		"int8": OpDerefPass | OpCastInt64,
	},
	"int32": {
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass,
		"int8": OpPass | OpCastInt64,
	},
	"*int32": {
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass,
		"int8": OpDerefPass | OpCastInt64,
	},
	"uint32": {
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64,
	},
	"*uint32": {
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64,
	},
	"int64": {
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass,
	},
	"*int64": {
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass,
	},
	"uint64": {
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64 | OpCheckOverflow,
	},
	"*uint64": {
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64 | OpCheckOverflow,
	},
	"float32": {
		"real":  OpPass,
		"float": OpPass | OpCastFloat64,
	},
	"*float32": {
		"real":  OpDerefPass,
		"float": OpDerefPass | OpCastFloat64,
	},
	"float64": {
		"real":  OpPass | OpCastFloat32 | OpCheckOverflow,
		"float": OpPass,
	},
	"*float64": {
		"real":  OpDerefPass | OpCastFloat32 | OpCheckOverflow,
		"float": OpDerefPass,
	},
	"string": {
		"bytea":   OpPass | OpCastBytes,
		"text":    OpPass,
		"varchar": OpPass,
	},
	"*string": {
		"bytea":   OpDerefPass | OpCastBytes,
		"text":    OpDerefPass,
		"varchar": OpDerefPass,
	},
	"[]byte": {
		"bytea":   OpPass,
		"text":    OpPass | OpCastString,
		"varchar": OpPass | OpCastString,
	},
	"*[]byte": {
		"bytea":   OpDerefPass,
		"text":    OpDerefPass | OpCastString,
		"varchar": OpDerefPass | OpCastString,
	},
	"time.Time": {
		"date":        OpPass,
		"timestamp":   OpPass,
		"timestampTz": OpPass,
	},
	"*time.Time": {
		"date":        OpDerefPass,
		"timestamp":   OpDerefPass,
		"timestampTz": OpDerefPass,
	},
	"[]bool": {
		"bool[]": OpPass,
	},
	"*[]bool": {
		"bool[]": OpDerefPass,
	},
	"[]int16": {
		"int2[]": OpPass,
	},
	"*[]int16": {
		"int2[]": OpPass,
	},
	"[]int32": {
		"int4[]": OpPass,
	},
	"*[]int32": {
		"int4[]": OpDerefPass,
	},
	"[]int64": {
		"int8[]": OpPass,
	},
	"*[]int64": {
		"int8[]": OpDerefPass,
	},
	"[]float32": {
		"real[]": OpPass,
	},
	"*[]float32": {
		"real[]": OpDerefPass,
	},
	"[]float64": {
		"float[]": OpPass,
	},
	"*[]float64": {
		"float[]": OpDerefPass,
	},
	"[]string": {
		"text[]":    OpPass,
		"varchar[]": OpPass,
	},
	"*[]string": {
		"text[]":    OpDerefPass,
		"varchar[]": OpDerefPass,
	},
	"[]time.Time": {
		"timestamp[]":   OpPass,
		"timestampTz[]": OpPass,
	},
	"*[]time.Time": {
		"timestamp[]":   OpPass,
		"timestampTz[]": OpPass,
	},

	// pgx built-in encoders:
	"pgx.Hstore": {
		"hstore": OpCustomEncode,
	},
	"*pgx.Hstore": {
		"hstore": OpCustomEncode,
	},
	"pgx.NullBool": {
		"bool": OpCustomEncode,
	},
	"*pgx.NullBool": {
		"bool": OpCustomEncode,
	},
	"pgx.NullInt16": {
		"int2": OpCustomEncode,
	},
	"*pgx.NullInt16": {
		"int2": OpCustomEncode,
	},
	"pgx.NullInt32": {
		"int4": OpCustomEncode,
	},
	"*pgx.NullInt32": {
		"int4": OpCustomEncode,
	},
	"pgx.NullInt64": {
		"int8": OpCustomEncode,
	},
	"*pgx.NullInt64": {
		"int8": OpCustomEncode,
	},
	"pgx.NullFloat32": {
		"real": OpCustomEncode,
	},
	"*pgx.NullFloat32": {
		"real": OpCustomEncode,
	},
	"pgx.NullFloat64": {
		"float": OpCustomEncode,
	},
	"*pgx.NullFloat64": {
		"float": OpCustomEncode,
	},
	"pgx.NullString": {
		"text":    OpCustomEncode,
		"varchar": OpCustomEncode,
	},
	"*pgx.NullString": {
		"text":    OpCustomEncode,
		"varchar": OpCustomEncode,
	},
	"pgx.NullTime": {
		"timestamp":   OpCustomEncode,
		"timestampTz": OpCustomEncode,
	},
	"*pgx.NullTime": {
		"timestamp":   OpCustomEncode,
		"timestampTz": OpCustomEncode,
	},
	"pgx.NullHstore": {
		"hstore": OpCustomEncode,
	},
	"*pgx.NullHstore": {
		"hstore": OpCustomEncode,
	},

	// wrappers around pgx built-in encoders:
	"map[string]string": {
		"hstore": OpHstoreMapEncode,
	},
	"*map[string]string": {
		"hstore": OpDerefPass | OpHstoreMapEncode,
	},
}
