package mysqlproto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandshakeResponse41(t *testing.T) {
	{
		capabilityFlags := uint32(0x807ff7df)
		characterSet := byte(0x21)
		username := "root"
		password := "user123"
		authPluginData := []byte("abcdefghijklmnopqrst")
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

		assert.Equal(t, packet[:3], []byte{0x6f, 0x00, 0x00})
		assert.Equal(t, packet[3], byte(0x01))
		assert.Equal(t, packet[4:8], []byte{0xdf, 0xf7, 0x5f, 0x80})
		assert.Equal(t, packet[8:12], []byte{0x6f, 0x00, 0x00, 0x00})
		assert.Equal(t, packet[12], byte(0x21))
		assert.Equal(t, packet[13:36], make([]byte, 23))
		assert.Equal(t, string(packet[36:40]), "root")
		assert.Equal(t, packet[40], byte(0x00))
		assert.Equal(t, int(packet[41]), 20)
		assert.Equal(t, hex.EncodeToString(packet[42:62]), "d7dbc13284c74850c777657fc3d7eb80b7185a25")
		assert.Equal(t, string(packet[62:69]), "mysqldb")
		assert.Equal(t, packet[69], byte(0x00))
		assert.Equal(t, string(packet[70:91]), "mysql_native_password")
		assert.Equal(t, packet[91], byte(0x00))
		assert.Equal(t, int(packet[92]), 22)
		assert.Equal(t, int(packet[93]), 14)
		assert.Equal(t, string(packet[94:108]), "client_version")
		assert.Equal(t, int(packet[108]), 6)
		assert.Equal(t, string(packet[109:]), "5.6.25")
	}

	// connect without database
	{
		capabilityFlags := uint32(0x807ff7df)
		characterSet := byte(0x21)
		username := "root"
		password := "user123"
		authPluginData := []byte("abcdefghijklmnopqrst")
		database := ""
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

		assert.Equal(t, packet[:3], []byte{0x67, 0x00, 0x00})
		assert.Equal(t, packet[3], byte(0x01))
		assert.Equal(t, packet[4:8], []byte{0xd7, 0xf7, 0x5f, 0x80})
		assert.Equal(t, packet[8:12], []byte{0x67, 0x00, 0x00, 0x00})
		assert.Equal(t, packet[12], byte(0x21))
		assert.Equal(t, packet[13:36], make([]byte, 23))
		assert.Equal(t, string(packet[36:40]), "root")
		assert.Equal(t, packet[40], byte(0x00))
		assert.Equal(t, int(packet[41]), 20)
		assert.Equal(t, hex.EncodeToString(packet[42:62]), "d7dbc13284c74850c777657fc3d7eb80b7185a25")
		assert.Equal(t, string(packet[62:83]), "mysql_native_password")
		assert.Equal(t, packet[83], byte(0x00))
		assert.Equal(t, int(packet[84]), 22)
		assert.Equal(t, int(packet[85]), 14)
		assert.Equal(t, string(packet[86:100]), "client_version")
		assert.Equal(t, int(packet[100]), 6)
		assert.Equal(t, string(packet[101:]), "5.6.25")
	}
}

func TestNativePassword(t *testing.T) {
	data := []byte("abcdefghijklmnopqrst")
	pass := "user123"

	hash := nativePassword(pass, data)
	assert.Len(t, hash, 20)
	assert.Equal(t, hex.EncodeToString(hash), "d7dbc13284c74850c777657fc3d7eb80b7185a25")
}
