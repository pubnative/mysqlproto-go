package mysqlproto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandshakeResponse41(t *testing.T) {
	capabilityFlags := uint32(0x807ff7ff)
	characterSet := byte(0x21)
	username := "root"
	password := "user123"
	authPluginData := "abcdefghijklmnopqrst"
	database := "mysqldb"
	authPluginName := "mysql_native_password"
	connectAttrs := map[string]string{"client_version": "5.6.25"}

	packet := HandshakeResponse41(
		capabilityFlags,
		characterSet,
		username,
		password,
		authPluginData,
		database,
		authPluginName,
		connectAttrs,
	)

	assert.Equal(t, packet[:4], []byte{0xff, 0xf7, 0x7f, 0x80})
	assert.Equal(t, packet[4:8], []byte{0x78, 0x00, 0x00, 0x00})
	assert.Equal(t, packet[8], byte(0x21))
	assert.Equal(t, packet[9:41], make([]byte, 32))
	assert.Equal(t, string(packet[41:45]), "root")
	assert.Equal(t, packet[45], byte(0x00))
	assert.Equal(t, int(packet[46]), 20)
	assert.Equal(t, hex.EncodeToString(packet[47:67]), "d7dbc13284c74850c777657fc3d7eb80b7185a25")
	assert.Equal(t, string(packet[67:74]), "mysqldb")
	assert.Equal(t, packet[74], byte(0x00))
	assert.Equal(t, string(packet[75:96]), "mysql_native_password")
	assert.Equal(t, packet[96], byte(0x00))
	assert.Equal(t, int(packet[97]), 22)
	assert.Equal(t, int(packet[98]), 14)
	assert.Equal(t, string(packet[99:113]), "client_version")
	assert.Equal(t, int(packet[113]), 6)
	assert.Equal(t, string(packet[114:]), "5.6.25")
}
