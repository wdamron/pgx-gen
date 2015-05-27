package pgxgen

type Op uint64

type OpMap map[string]Op

// The low 56 bits of an Op are reserved for flags. Up to 56 flags may be defined.
const (
	OpAssign Op = 1 << iota
	OpPtrAssign
	OpPass
	OpDerefPass
	OpCustomScan
	OpCustomEncode
	OpCheckOverflow
	OpHstoreMapEncode
	OpHstoreMapDecode
	OpUuidEncode
	OpUuidDecode
)

// The high 8 bits of an op are reserved for the Op's cast type (if any).
// Casts are mutually exclusive, and up to 256 different cast types may be defined.
const (
	OpCastString Op = iota << 56
	OpCastBytes
	OpCastByte
	OpCastInt
	OpCastUint
	OpCastInt8
	OpCastUint8
	OpCastInt16
	OpCastUint16
	OpCastInt32
	OpCastUint32
	OpCastInt64
	OpCastUint64
	OpCastFloat32
	OpCastFloat64
)

func (op Op) MaskFlags() Op {
	return op & 0x00FFFFFFFFFFFFFF
}

func (op Op) MaskCast() Op {
	return op & 0xFF00000000000000
}

func (op Op) Assign() bool {
	return op&OpAssign != 0
}

func (op Op) PtrAssign() bool {
	return op&OpPtrAssign != 0
}

func (op Op) Pass() bool {
	return op&OpPass != 0
}

func (op Op) DerefPass() bool {
	return op&OpDerefPass != 0
}

func (op Op) CustomScan() bool {
	return op&OpCustomScan != 0
}

func (op Op) CustomEncode() bool {
	return op&OpCustomEncode != 0
}

func (op Op) CheckOverflow() bool {
	return op&OpCheckOverflow != 0
}

func (op Op) HstoreMapEncode() bool {
	return op&OpHstoreMapEncode != 0
}

func (op Op) HstoreMapDecode() bool {
	return op&OpHstoreMapDecode != 0
}

func (op Op) UuidEncode() bool {
	return op&OpUuidEncode != 0
}

func (op Op) UuidDecode() bool {
	return op&OpUuidDecode != 0
}

func (op Op) FormatCast() string {
	switch op.MaskCast() {
	case OpCastString:
		return "string"
	case OpCastBytes:
		return "[]byte"
	case OpCastByte:
		return "byte"
	case OpCastInt:
		return "int"
	case OpCastUint:
		return "uint"
	case OpCastInt8:
		return "int8"
	case OpCastUint8:
		return "uint8"
	case OpCastInt16:
		return "int16"
	case OpCastUint16:
		return "uint16"
	case OpCastInt32:
		return "int32"
	case OpCastUint32:
		return "uint32"
	case OpCastInt64:
		return "int64"
	case OpCastUint64:
		return "uint64"
	case OpCastFloat32:
		return "float32"
	case OpCastFloat64:
		return "float64"
	}
	return ""
}
