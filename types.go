package mysqlproto

// https://dev.mysql.com/doc/internals/en/com-query-response.html#column-type
type Type byte

const (
	TypeDecimal    Type = 0x00
	TypeTiny       Type = 0x01
	TypeShort      Type = 0x02
	TypeLong       Type = 0x03
	TypeFloat      Type = 0x04
	TypeDouble     Type = 0x05
	TypeNULL       Type = 0x06
	TypeTimestamp  Type = 0x07
	TypeLongLong   Type = 0x08
	TypeInt24      Type = 0x09
	TypeDate       Type = 0x0a
	TypeTime       Type = 0x0b
	TypeDateTime   Type = 0x0c
	TypeYear       Type = 0x0d
	TypeNewDate    Type = 0x0e
	TypeVarchar    Type = 0x0f
	TypeBit        Type = 0x10
	TypeTimestamp2 Type = 0x11
	TypeDateTime2  Type = 0x12
	TypeTime2      Type = 0x13
	TypeNewDecimal Type = 0xf6
	TypeEnum       Type = 0xf7
	TypeSet        Type = 0xf8
	TypeTinyBLOB   Type = 0xf9
	TypeMediumBLOB Type = 0xfa
	TypeLongBLOB   Type = 0xfb
	TypeBLOB       Type = 0xfc
	TypeVarString  Type = 0xfd
	TypeString     Type = 0xfe
	TypeGEOMETRY   Type = 0xff
)

func (t Type) String() string {
	switch t {
	case TypeDecimal:
		return "DECIMAL"
	case TypeTiny:
		return "TINY"
	case TypeShort:
		return "SHORT"
	case TypeLong:
		return "LONG"
	case TypeFloat:
		return "FLOAT"
	case TypeDouble:
		return "DOUBLE"
	case TypeNULL:
		return "NULL"
	case TypeTimestamp:
		return "TIMESTAMP"
	case TypeLongLong:
		return "LONGLONG"
	case TypeInt24:
		return "INT24"
	case TypeDate:
		return "DATE"
	case TypeTime:
		return "TIME"
	case TypeDateTime:
		return "DATETIME"
	case TypeYear:
		return "YEAR"
	case TypeNewDate:
		return "NEWDATE"
	case TypeVarchar:
		return "VARCHAR"
	case TypeBit:
		return "BIT"
	case TypeTimestamp2:
		return "TIMESTAMP2"
	case TypeDateTime2:
		return "DATETIME2"
	case TypeTime2:
		return "TIME2"
	case TypeNewDecimal:
		return "NEWDECIMAL"
	case TypeEnum:
		return "ENUM"
	case TypeSet:
		return "SET"
	case TypeTinyBLOB:
		return "TINY_BLOB"
	case TypeMediumBLOB:
		return "MEDIUM_BLOB"
	case TypeLongBLOB:
		return "LONG_BLOB"
	case TypeBLOB:
		return "BLOB"
	case TypeVarString:
		return "VAR_STRING"
	case TypeString:
		return "STRING"
	case TypeGEOMETRY:
		return "GEOMETRY"
	default:
		return "UNKNOWN"
	}
}
