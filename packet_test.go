package mysqlproto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPacket(t *testing.T) {
	stream := bytes.NewBuffer([]byte{
		0x07, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x02,
		0x01, 0x02, 0x03,
	})

	packet, err := ReadPacket(stream)
	assert.Nil(t, err)
	assert.Equal(t, packet.SequenceID, byte(0x02))
	assert.Equal(t, packet.Payload, []byte{0x00, 0x00, 0x00, 0x02, 0x01, 0x02, 0x03})
}
