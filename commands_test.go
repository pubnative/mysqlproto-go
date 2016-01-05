package mysqlproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandPacketWithoutPayload(t *testing.T) {
	pkt := CommandPacket(COM_QUIT, nil)
	assert.Equal(t, pkt, []byte{0x01, 0x0, 0x0, 0x0, COM_QUIT})
}

func TestCommandPacketWithPayload(t *testing.T) {
	query := []byte("SELECT * FROM people")
	pkt := CommandPacket(COM_QUERY, query)
	assert.Equal(t, pkt, append([]byte{0x15, 0x0, 0x0, 0x0, COM_QUERY}, query...))
}
