package mysqlproto

import (
	"crypto/sha1"
)

func HandshakeResponse41(
	capabilityFlags uint32,
	characterSet byte,
	username string,
	password string,
	authPluginData string,
	database string,
	authPluginName string,
	connectAttrs map[string]string,
) []byte {
	capabilityFlags |= CLIENT_PROTOCOL_41 // must be always set

	var packetSize uint32 = 0
	packetSize += 4                         // capability flags
	packetSize += 4                         // packet size
	packetSize += 1                         // character set
	packetSize += 32                        // reserved string
	packetSize += uint32(len(username)) + 1 // + null character

	var authResponse []byte
	switch authPluginName {
	case "mysql_native_password":
		authResponse = nativePassword(authPluginData, password)
	case "mysql_old_password":
		panic(`auth method "mysql_old_password" not supported`) // todo
	default:
		panic(`invalid auth method "` + authPluginName + `"`)
	}
	packetSize += uint32(len(authResponse))

	var authResponseLen []byte
	if capabilityFlags&CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA > 0 ||
		capabilityFlags&CLIENT_SECURE_CONNECTION > 0 {
		authResponseLen = lenEncInt(uint64(len(authResponse)))
		packetSize += uint32(len(authResponseLen))
	} else {
		authResponse = append(authResponse, 0x00)
		packetSize += 1
	}

	if l := len(database); l > 0 {
		capabilityFlags |= CLIENT_CONNECT_WITH_DB
		packetSize += uint32(l) + 1 // + null character
	}

	if l := len(authPluginName); l > 0 {
		capabilityFlags |= CLIENT_PLUGIN_AUTH
		packetSize += uint32(l) + 1 // + null character
	}

	var attrData []byte
	if len(connectAttrs) > 0 {
		var data []byte
		capabilityFlags |= CLIENT_CONNECT_ATTRS
		for key, value := range connectAttrs {
			data = append(data, lenEncStr(key)...)
			data = append(data, lenEncStr(value)...)
		}

		total := lenEncInt(uint64(len(data)))
		attrData = make([]byte, len(total)+len(data))

		copy(attrData[:len(total)], total)
		copy(attrData[len(total):], data)
	}

	packetSize += uint32(len(attrData))

	packet := make([]byte, 0, packetSize)

	packet = append(packet,
		byte(capabilityFlags),
		byte(capabilityFlags>>8),
		byte(capabilityFlags>>16),
		byte(capabilityFlags>>24),
	)

	packet = append(packet,
		byte(packetSize),
		byte(packetSize>>8),
		byte(packetSize>>16),
		byte(packetSize>>24),
	)

	packet = append(packet, characterSet)

	packet = append(packet, make([]byte, 32)...)

	packet = append(packet, username...)
	packet = append(packet, 0x00)

	packet = append(packet, authResponseLen...)
	packet = append(packet, authResponse...)

	packet = append(packet, database...)
	packet = append(packet, 0x00)

	packet = append(packet, authPluginName...)
	packet = append(packet, 0x00)

	packet = append(packet, attrData...)

	return packet
}

// https://dev.mysql.com/doc/internals/en/secure-password-authentication.html#packet-Authentication::Native41
// SHA1( password ) XOR SHA1( "20-bytes random data from server" <concat> SHA1( SHA1( password ) ) )
func nativePassword(authPluginData string, password string) []byte {
	if len(password) == 0 {
		return nil
	}

	hash := sha1.New()
	hash.Write([]byte(password))
	hashPass := hash.Sum(nil)

	hash = sha1.New()
	hash.Write(hashPass)
	doubleHashPass := hash.Sum(nil)

	hash = sha1.New()
	hash.Write([]byte(authPluginData))
	hash.Write(doubleHashPass)
	salt := hash.Sum(nil)

	for i, b := range hashPass {
		hashPass[i] = b ^ salt[i]
	}

	return hashPass
}