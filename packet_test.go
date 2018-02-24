package mysqlproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOKPacketInvalidPayload(t *testing.T) {
	data := []byte{0xff}
	_, err := ParseOKPacket(data, 0)
	assert.Equal(t, err.Error(), "mysqlproto: invalid OK_PACKET payload: ff")
}

func TestParseOKPacketUpdateReply(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x28, 0x52, 0x6f, 0x77, 0x73,
		0x20, 0x6d, 0x61, 0x74, 0x63, 0x68,
		0x65, 0x64, 0x3a, 0x20, 0x31, 0x20,
		0x20, 0x43, 0x68, 0x61, 0x6e, 0x67,
		0x65, 0x64, 0x3a, 0x20, 0x31, 0x20,
		0x20, 0x57, 0x61, 0x72, 0x6e, 0x69,
		0x6e, 0x67, 0x73, 0x3a, 0x20, 0x30,
	}
	pkt, err := ParseOKPacket(data, CLIENT_PROTOCOL_41|CLIENT_SESSION_TRACK)
	assert.Nil(t, err)
	assert.Equal(t, pkt.Header, byte(0x00))
	assert.Equal(t, pkt.AffectedRows, uint64(1))
	assert.Equal(t, pkt.LastInsertID, uint64(0))
	assert.Equal(t, pkt.StatusFlags, SERVER_STATUS_AUTOCOMMIT)
	assert.Equal(t, pkt.Warnings, uint16(0))
	assert.Equal(t, pkt.Info, "Rows matched: 1  Changed: 1  Warnings: 0")
	assert.Equal(t, pkt.SessionStateChanges, "")
}

func TestParseOKPacketInsertReply(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0xfd, 0x9f, 0x86,
		0x01, 0x02, 0x00, 0x00, 0x00,
	}
	pkt, err := ParseOKPacket(data, CLIENT_PROTOCOL_41|CLIENT_SESSION_TRACK)
	assert.Nil(t, err)
	assert.Equal(t, pkt.Header, byte(0x00))
	assert.Equal(t, pkt.AffectedRows, uint64(1))
	assert.Equal(t, pkt.LastInsertID, uint64(99999))
	assert.Equal(t, pkt.StatusFlags, SERVER_STATUS_AUTOCOMMIT)
	assert.Equal(t, pkt.Warnings, uint16(0))
	assert.Equal(t, pkt.Info, "")
	assert.Equal(t, pkt.SessionStateChanges, "")
}

func TestParseERRPacketInvalidPayload(t *testing.T) {
	data := []byte{
		0xfe, 0x48, 0x04, 0x23, 0x48, 0x59,
		0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73,
		0x20, 0x75, 0x73, 0x65, 0x64,
	}
	_, err := ParseERRPacket(data, CLIENT_PROTOCOL_41)
	assert.Equal(t, err.Error(), "mysqlproto: invalid ERR_PACKET payload: fe48042348593030304e6f207461626c65732075736564")
}

func TestParseERRPacketCLIENT_PROTOCOL_41(t *testing.T) {
	data := []byte{
		0xff, 0x48, 0x04, 0x23, 0x48, 0x59,
		0x30, 0x30, 0x30, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73,
		0x20, 0x75, 0x73, 0x65, 0x64,
	}
	pkt, err := ParseERRPacket(data, CLIENT_PROTOCOL_41)
	assert.Nil(t, err)
	assert.Equal(t, pkt.Header, byte(0xff))
	assert.Equal(t, pkt.ErrorCode, uint16(1096))
	assert.Equal(t, pkt.SQLStateMarker, "#")
	assert.Equal(t, pkt.SQLState, "HY000")
	assert.Equal(t, pkt.ErrorMessage, "No tables used")
}

func TestParseERRPacket(t *testing.T) {
	data := []byte{
		0xff, 0x48, 0x04, 0x4e, 0x6f, 0x20,
		0x74, 0x61, 0x62, 0x6c, 0x65, 0x73,
		0x20, 0x75, 0x73, 0x65, 0x64,
	}
	pkt, err := ParseERRPacket(data, 0)
	assert.Nil(t, err)
	assert.Equal(t, pkt.Header, byte(0xff))
	assert.Equal(t, pkt.ErrorCode, uint16(1096))
	assert.Equal(t, pkt.SQLStateMarker, "")
	assert.Equal(t, pkt.SQLState, "")
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
