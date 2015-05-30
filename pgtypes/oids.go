package pgtypes

const (
	TextFormatCode   = 0
	BinaryFormatCode = 1
)

const (
	BoolOid, BoolArrayOid               = 16, 1000
	Int2Oid, Int2ArrayOid               = 21, 1005
	Int4Oid, Int4ArrayOid               = 23, 1007
	Int8Oid, Int8ArrayOid               = 20, 1016
	Float4Oid, Float4ArrayOid           = 700, 1021
	Float8Oid, Float8ArrayOid           = 701, 1022
	ByteaOid                            = 17
	TextOid, TextArrayOid               = 25, 1009
	VarcharOid, VarcharArrayOid         = 1043, 1015
	OidOid                              = 26
	DateOid                             = 1082
	TimestampOid, TimestampArrayOid     = 1114, 1115
	TimestampTzOid, TimestampTzArrayOid = 1184, 1185
	JSONOid                             = 114
	UUIDOid                             = 2950
	HstoreOid                           = 0    // hstore data types have a non-constant oid
	JSONArrayOid                        = 199  // json[] data types are not currently supported
	UUIDArrayOid                        = 2951 // uuid[] data types are not currently supported
	XMLOid                              = 142  // xml data types are not currently supported
)
