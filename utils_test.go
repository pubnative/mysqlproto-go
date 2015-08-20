package mysqlproto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadNullStr(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x00})
	str, err := readNullStr(buf)
	assert.Nil(t, err)
	assert.Len(t, str, 0)

	buf = bytes.NewBuffer([]byte{0x00, 0xaa})
	str, err = readNullStr(buf)
	assert.Nil(t, err)
	assert.Len(t, str, 0)

	buf = bytes.NewBuffer([]byte{0xaa, 0xfe, 0x01, 0x00})
	str, err = readNullStr(buf)
	assert.Nil(t, err)
	assert.Len(t, str, 3)
	assert.Equal(t, str, []byte{0xaa, 0xfe, 0x01})

	buf = bytes.NewBuffer([]byte{0xaa, 0xfe, 0x01})
	str, err = readNullStr(buf)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "EOF")
}
