package pgxgen

var Encoders = map[string]OpMap{
	"bool": OpMap{
		"bool": OpPass,
	},
	"*bool": OpMap{
		"bool": OpDerefPass,
	},
	"int": OpMap{
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64,
	},
	"*int": OpMap{
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64,
	},
	"uint": OpMap{
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64 | OpCheckOverflow,
	},
	"*uint": OpMap{
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64 | OpCheckOverflow,
	},
	"int16": OpMap{
		"int2": OpPass,
		"int4": OpPass | OpCastInt32,
		"int8": OpPass | OpCastInt64,
	},
	"*int16": OpMap{
		"int2": OpDerefPass,
		"int4": OpDerefPass | OpCastInt32,
		"int8": OpDerefPass | OpCastInt64,
	},
	"uint16": OpMap{
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32,
		"int8": OpPass | OpCastInt64,
	},
	"*uint16": OpMap{
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32,
		"int8": OpDerefPass | OpCastInt64,
	},
	"int32": OpMap{
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass,
		"int8": OpPass | OpCastInt64,
	},
	"*int32": OpMap{
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass,
		"int8": OpDerefPass | OpCastInt64,
	},
	"uint32": OpMap{
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64,
	},
	"*uint32": OpMap{
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64,
	},
	"int64": OpMap{
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass,
	},
	"*int64": OpMap{
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass,
	},
	"uint64": OpMap{
		"int2": OpPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpPass | OpCastInt64 | OpCheckOverflow,
	},
	"*uint64": OpMap{
		"int2": OpDerefPass | OpCastInt16 | OpCheckOverflow,
		"int4": OpDerefPass | OpCastInt32 | OpCheckOverflow,
		"int8": OpDerefPass | OpCastInt64 | OpCheckOverflow,
	},
	"float32": OpMap{
		"float4": OpPass,
		"float8": OpPass | OpCastFloat64,
	},
	"*float32": OpMap{
		"float4": OpDerefPass,
		"float8": OpDerefPass | OpCastFloat64,
	},
	"float64": OpMap{
		"float4": OpPass | OpCastFloat32 | OpCheckOverflow,
		"float8": OpPass,
	},
	"*float64": OpMap{
		"float4": OpDerefPass | OpCastFloat32 | OpCheckOverflow,
		"float8": OpDerefPass,
	},
	"string": OpMap{
		"bytea":   OpPass | OpCastBytes,
		"text":    OpPass,
		"varchar": OpPass,
	},
	"*string": OpMap{
		"bytea":   OpDerefPass | OpCastBytes,
		"text":    OpDerefPass,
		"varchar": OpDerefPass,
	},
	"[]byte": OpMap{
		"bytea":   OpPass,
		"text":    OpPass | OpCastString,
		"varchar": OpPass | OpCastString,
	},
	"*[]byte": OpMap{
		"bytea":   OpDerefPass,
		"text":    OpDerefPass | OpCastString,
		"varchar": OpDerefPass | OpCastString,
	},
	"time.Time": OpMap{
		"date":        OpPass,
		"timestamp":   OpPass,
		"timestampTz": OpPass,
	},
	"*time.Time": OpMap{
		"date":        OpDerefPass,
		"timestamp":   OpDerefPass,
		"timestampTz": OpDerefPass,
	},
	"[]bool": OpMap{
		"[]bool": OpPass,
	},
	"*[]bool": OpMap{
		"[]bool": OpDerefPass,
	},
	"[]int16": OpMap{
		"[]int2": OpPass,
	},
	"*[]int16": OpMap{
		"[]int2": OpPass,
	},
	"[]int32": OpMap{
		"[]int4": OpPass,
	},
	"*[]int32": OpMap{
		"[]int4": OpDerefPass,
	},
	"[]int64": OpMap{
		"[]int8": OpPass,
	},
	"*[]int64": OpMap{
		"[]int8": OpDerefPass,
	},
	"[]float32": OpMap{
		"[]float4": OpPass,
	},
	"*[]float32": OpMap{
		"[]float4": OpDerefPass,
	},
	"[]float64": OpMap{
		"[]float8": OpPass,
	},
	"*[]float64": OpMap{
		"[]float8": OpDerefPass,
	},
	"[]string": OpMap{
		"[]text":    OpPass,
		"[]varchar": OpPass,
	},
	"*[]string": OpMap{
		"[]text":    OpDerefPass,
		"[]varchar": OpDerefPass,
	},
	"[]time.Time": OpMap{
		"[]timestamp":   OpPass,
		"[]timestampTz": OpPass,
	},
	"*[]time.Time": OpMap{
		"[]timestamp":   OpPass,
		"[]timestampTz": OpPass,
	},
}
