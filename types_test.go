package mysqlproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeString(t *testing.T) {
	testCases := []struct {
		typ Type
		hex byte
		str string
	}{
		{typ: TypeDecimal, hex: 0x00, str: "DECIMAL"},
		{typ: TypeTiny, hex: 0x01, str: "TINY"},
		{typ: TypeShort, hex: 0x02, str: "SHORT"},
		{typ: TypeLong, hex: 0x03, str: "LONG"},
		{typ: TypeFloat, hex: 0x04, str: "FLOAT"},
		{typ: TypeDouble, hex: 0x05, str: "DOUBLE"},
		{typ: TypeNULL, hex: 0x06, str: "NULL"},
		{typ: TypeTimestamp, hex: 0x07, str: "TIMESTAMP"},
		{typ: TypeLongLong, hex: 0x08, str: "LONGLONG"},
		{typ: TypeInt24, hex: 0x09, str: "INT24"},
		{typ: TypeDate, hex: 0x0a, str: "DATE"},
		{typ: TypeTime, hex: 0x0b, str: "TIME"},
		{typ: TypeDateTime, hex: 0x0c, str: "DATETIME"},
		{typ: TypeYear, hex: 0x0d, str: "YEAR"},
		{typ: TypeNewDate, hex: 0x0e, str: "NEWDATE"},
		{typ: TypeVarchar, hex: 0x0f, str: "VARCHAR"},
		{typ: TypeBit, hex: 0x10, str: "BIT"},
		{typ: TypeTimestamp2, hex: 0x11, str: "TIMESTAMP2"},
		{typ: TypeDateTime2, hex: 0x12, str: "DATETIME2"},
		{typ: TypeTime2, hex: 0x13, str: "TIME2"},
		{typ: TypeNewDecimal, hex: 0xf6, str: "NEWDECIMAL"},
		{typ: TypeEnum, hex: 0xf7, str: "ENUM"},
		{typ: TypeSet, hex: 0xf8, str: "SET"},
		{typ: TypeTinyBLOB, hex: 0xf9, str: "TINY_BLOB"},
		{typ: TypeMediumBLOB, hex: 0xfa, str: "MEDIUM_BLOB"},
		{typ: TypeLongBLOB, hex: 0xfb, str: "LONG_BLOB"},
		{typ: TypeBLOB, hex: 0xfc, str: "BLOB"},
		{typ: TypeVarString, hex: 0xfd, str: "VAR_STRING"},
		{typ: TypeString, hex: 0xfe, str: "STRING"},
		{typ: TypeGEOMETRY, hex: 0xff, str: "GEOMETRY"},
	}

	for i, tc := range testCases {
		assert.Equal(t, byte(tc.typ), tc.hex, "test case %v", i)
		assert.Equal(t, tc.typ.String(), tc.str, "test case %v", i)
	}
}
