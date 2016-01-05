package mysqlproto

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnCloseStreamIsNotSet(t *testing.T) {
	conn := Conn{}
	assert.Equal(t, conn.Close(), ErrNoStream)
}

func TestConnCloseWriteError(t *testing.T) {
	errSocket := errors.New("can't write into socket")
	buf := &buffer{
		writeFn: func(data []byte) (int, error) {
			if string(data) == string(CommandPacket(COM_QUIT, nil)) {
				return 0, errSocket
			}

			return 0, nil
		},
	}

	conn := Conn{Stream: &Stream{stream: buf}}
	err := conn.Close()
	assert.Equal(t, err, errSocket)
	assert.True(t, buf.closed)
}

func TestConnCloseServerReplyEOF(t *testing.T) {
	buf := newBuffer([]byte{})
	conn := Conn{Stream: &Stream{stream: buf}}
	err := conn.Close()
	assert.Nil(t, err)
	assert.True(t, buf.closed)
}

func TestConnCloseServerReplyOKPacket(t *testing.T) {
	buf := newBuffer([]byte{0x1, 0x0, 0x0, 0x1, OK_PACKET})
	conn := Conn{Stream: &Stream{stream: buf}}
	err := conn.Close()
	assert.Nil(t, err)
	assert.True(t, buf.closed)
}

func TestConnCloseServerReplyERRPacket(t *testing.T) {
	data := []byte{
		0x17, 0x00, 0x00, 0x01, 0xff, 0x48,
		0x04, 0x23, 0x48, 0x59, 0x30, 0x30,
		0x30, 0x4e, 0x6f, 0x20, 0x74, 0x61,
		0x62, 0x6c, 0x65, 0x73, 0x20, 0x75,
		0x73, 0x65, 0x64,
	}
	buf := newBuffer(data)
	conn := Conn{Stream: NewStream(buf), CapabilityFlags: CLIENT_PROTOCOL_41}
	err := conn.Close()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "mysqlproto: Error: 1096 SQLSTATE: HY000 Message: No tables used")
	assert.True(t, buf.closed)
}

func TestConnCloseServerReplyInvalidPacket(t *testing.T) {
	data := []byte{
		0x8, 0x0, 0x0, 0x0,
		0xdd, 0x48, 0x04, 0x23,
		0x48, 0x59, 0x30, 0x30,
	}
	buf := newBuffer(data)
	conn := Conn{Stream: NewStream(buf)}
	err := conn.Close()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "mysqlproto: invalid ERR_PACKET payload: dd48042348593030")
	assert.True(t, buf.closed)
}
