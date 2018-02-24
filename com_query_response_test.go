package mysqlproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComQueryResponseColumnReader(t *testing.T) {
	buf := newBuffer([]byte{
		// DB name "test"
		// table name "people" AS "p"

		// total records
		0x01, 0x00, 0x00, 0x01, 0x05,

		// id INT
		0x25, 0x00, 0x00, 0x02, 0x03, 0x64, 0x65, 0x66, 0x04, 0x74, 0x65, 0x73, 0x74, 0x01, 0x70, 0x06, 0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x02, 0x69, 0x64, 0x02, 0x69, 0x64, 0x0c, 0x3f, 0x00, 0x0b, 0x00, 0x00, 0x00, 0x03, 0x03, 0x42, 0x00, 0x00, 0x00,

		// firstname VARCHAR(255) AS name
		0x2e, 0x00, 0x00, 0x03, 0x03, 0x64, 0x65, 0x66, 0x04, 0x74, 0x65, 0x73, 0x74, 0x01, 0x70, 0x06, 0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x0c, 0x21, 0x00, 0xfd, 0x02, 0x00, 0x00, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00,

		// married TINYINT
		0x2f, 0x00, 0x00, 0x04, 0x03, 0x64, 0x65, 0x66, 0x04, 0x74, 0x65, 0x73, 0x74, 0x01, 0x70, 0x06, 0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x07, 0x6d, 0x61, 0x72, 0x72, 0x69, 0x65, 0x64, 0x07, 0x6d, 0x61, 0x72, 0x72, 0x69, 0x65, 0x64, 0x0c, 0x3f, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,

		// score DECIMAL(6,2)
		0x2b, 0x00, 0x00, 0x05, 0x03, 0x64, 0x65, 0x66, 0x04, 0x74, 0x65, 0x73, 0x74, 0x01, 0x70, 0x06, 0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x0c, 0x3f, 0x00, 0x08, 0x00, 0x00, 0x00, 0xf6, 0x00, 0x00, 0x02, 0x00, 0x00,

		// note TEXT
		0x29, 0x00, 0x00, 0x06, 0x03, 0x64, 0x65, 0x66, 0x04, 0x74, 0x65, 0x73, 0x74, 0x01, 0x70, 0x06, 0x70, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x0c, 0x21, 0x00, 0xfd, 0xff, 0x02, 0x00, 0xfc, 0x10, 0x00, 0x00, 0x00, 0x00,

		// EOF
		0x05, 0x00, 0x00, 0x07, 0xfe, 0x00, 0x00, 0x22, 0x00,
	})

	conn := Conn{Stream: &Stream{stream: buf}}
	rs, err := ComQueryResponse(conn)
	assert.NoError(t, err)
	assert.Len(t, rs.Columns, 5)

	id := rs.Columns[0]
	assert.Equal(t, id.Catalog, "def")
	assert.Equal(t, id.Schema, "test")
	assert.Equal(t, id.Table, "p")
	assert.Equal(t, id.OrgTable, "people")
	assert.Equal(t, id.Name, "id")
	assert.Equal(t, id.OrgName, "id")
	assert.Equal(t, id.CharacterSet, uint16(63))
	assert.Equal(t, id.ColumnLength, uint64(11))
	assert.Equal(t, id.ColumnType.String(), "LONG")
	assert.Equal(t, id.Flags, uint16(3))
	assert.Equal(t, id.Decimals, uint8(0))

	name := rs.Columns[1]
	assert.Equal(t, name.Catalog, "def")
	assert.Equal(t, name.Schema, "test")
	assert.Equal(t, name.Table, "p")
	assert.Equal(t, name.OrgTable, "people")
	assert.Equal(t, name.Name, "name")
	assert.Equal(t, name.OrgName, "firstname")
	assert.Equal(t, name.CharacterSet, uint16(33))
	assert.Equal(t, name.ColumnLength, uint64(765))
	assert.Equal(t, name.ColumnType.String(), "VAR_STRING")
	assert.Equal(t, name.Flags, uint16(0))
	assert.Equal(t, name.Decimals, uint8(0))

	married := rs.Columns[2]
	assert.Equal(t, married.Catalog, "def")
	assert.Equal(t, married.Schema, "test")
	assert.Equal(t, married.Table, "p")
	assert.Equal(t, married.OrgTable, "people")
	assert.Equal(t, married.Name, "married")
	assert.Equal(t, married.OrgName, "married")
	assert.Equal(t, married.CharacterSet, uint16(63))
	assert.Equal(t, married.ColumnLength, uint64(4))
	assert.Equal(t, married.ColumnType.String(), "TINY")
	assert.Equal(t, married.Flags, uint16(0))
	assert.Equal(t, married.Decimals, uint8(0))

	score := rs.Columns[3]
	assert.Equal(t, score.Catalog, "def")
	assert.Equal(t, score.Schema, "test")
	assert.Equal(t, score.Table, "p")
	assert.Equal(t, score.OrgTable, "people")
	assert.Equal(t, score.Name, "score")
	assert.Equal(t, score.OrgName, "score")
	assert.Equal(t, score.CharacterSet, uint16(63))
	assert.Equal(t, score.ColumnLength, uint64(8))
	assert.Equal(t, score.ColumnType.String(), "NEWDECIMAL")
	assert.Equal(t, score.Flags, uint16(512))
	assert.Equal(t, score.Decimals, uint8(2))

	note := rs.Columns[4]
	assert.Equal(t, note.Catalog, "def")
	assert.Equal(t, note.Schema, "test")
	assert.Equal(t, note.Table, "p")
	assert.Equal(t, note.OrgTable, "people")
	assert.Equal(t, note.Name, "note")
	assert.Equal(t, note.OrgName, "note")
	assert.Equal(t, note.CharacterSet, uint16(33))
	assert.Equal(t, note.ColumnLength, uint64(196605))
	assert.Equal(t, note.ColumnType.String(), "BLOB")
	assert.Equal(t, note.Flags, uint16(16))
	assert.Equal(t, note.Decimals, uint8(0))
}