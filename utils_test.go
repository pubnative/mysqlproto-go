package mysqlproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLenDecInt(t *testing.T) {
	b := []byte{0xfb}
	num, offset, null := lenDecInt(b)
	assert.Equal(t, num, uint64(0))
	assert.Equal(t, offset, uint64(1))
	assert.True(t, null)

	b = []byte{0xfc, 0x01, 0x02}
	num, offset, null = lenDecInt(b)
	assert.Equal(t, num, uint64(513))
	assert.Equal(t, offset, uint64(3))
	assert.False(t, null)

	b = []byte{0xfd, 0x01, 0x02, 0x03}
	num, offset, null = lenDecInt(b)
	assert.Equal(t, num, uint64(197121))
	assert.Equal(t, offset, uint64(4))
	assert.False(t, null)

	b = []byte{0xfe, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	num, offset, null = lenDecInt(b)
	assert.Equal(t, num, uint64(578437695752307201))
	assert.Equal(t, offset, uint64(9))
	assert.False(t, null)

	b = []byte{0xfa}
	num, offset, null = lenDecInt(b)
	assert.Equal(t, num, uint64(250))
	assert.Equal(t, offset, uint64(1))
	assert.False(t, null)
}
