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
	assert.Equal(t, packet[4:8], []byte{0x6f, 0x00, 0x00, 0x00})
	assert.Equal(t, packet[8], byte(0x21))
	assert.Equal(t, packet[9:32], make([]byte, 23))
	assert.Equal(t, string(packet[32:36]), "root")
	assert.Equal(t, packet[36], byte(0x00))
	assert.Equal(t, int(packet[37]), 20)
	assert.Equal(t, hex.EncodeToString(packet[38:58]), "d7dbc13284c74850c777657fc3d7eb80b7185a25")
	assert.Equal(t, string(packet[58:65]), "mysqldb")
	assert.Equal(t, packet[65], byte(0x00))
	assert.Equal(t, string(packet[66:87]), "mysql_native_password")
	assert.Equal(t, packet[87], byte(0x00))
	assert.Equal(t, int(packet[88]), 22)
	assert.Equal(t, int(packet[89]), 14)
	assert.Equal(t, string(packet[90:104]), "client_version")
	assert.Equal(t, int(packet[104]), 6)
	assert.Equal(t, string(packet[105:]), "5.6.25")
}
