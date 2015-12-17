package mysqlproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseERRPacketInvalidPayload(t *testing.T) {
	data := []byte{
		0xfe, 0x48, 0x04, 0x23, 0x48, 0x59,
		0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73,
		0x20, 0x75, 0x73, 0x65, 0x64,
	}
	_, err := ParseERRPacket(data)
	assert.Equal(t, err, ErrERRPacketPayload)
}

func TestParseERRPacket(t *testing.T) {
	data := []byte{
		0xff, 0x48, 0x04, 0x23, 0x48, 0x59,
		0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73,
		0x20, 0x75, 0x73, 0x65, 0x64,
	}
	pkt, err := ParseERRPacket(data)
	assert.Nil(t, err)
	assert.Equal(t, pkt.Header, byte(0xff))
	assert.Equal(t, pkt.ErrorCode, uint16(1096))
	assert.Equal(t, pkt.SQLStateMarker, "#")
	assert.Equal(t, pkt.SQLState, "HY000")
	assert.Equal(t, pkt.ErrorMessage, "No tables used")
}

func TestERRPacketError(t *testing.T) {
	pkt := ERRPacket{
		ErrorCode:    uint16(1038),
		SQLState:     "HY001",
		ErrorMessage: "Out of sort memory",
	}
	assert.Equal(t, pkt.Error(), "mysqlproto: Error: 1038 SQLSTATE: HY001 Message: Out of sort memory")
}
