// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeV10
package mysqlproto

import (
	"io"
)

type HandshakeV10 struct {
	ProtocolVersion byte
	ServerVersion   string
	ConnectionId    [4]byte
	AuthPluginData  string
	CapabilityFlags uint32
	CharacterSet    byte
	StatusFlags     [2]byte
	AuthPluginName  string
}

func NewHandshakeV10(stream io.Reader) (HandshakeV10, error) {
	packet := HandshakeV10{}

	data := make([]byte, 1)
	if _, err := stream.Read(data); err != nil {
		return packet, err
	}
	packet.ProtocolVersion = data[0]

	srvVer, err := readNullStr(stream)
	if err != nil {
		return packet, err
	}
	packet.ServerVersion = string(srvVer)

	if _, err := stream.Read(packet.ConnectionId[:4]); err != nil {
		return packet, err
	}

	authData := make([]byte, 8)
	if _, err := stream.Read(authData); err != nil {
		return packet, err
	}
	packet.AuthPluginData = string(authData)

	// skip filler
	if _, err := stream.Read(make([]byte, 1)); err != nil {
		return packet, err
	}

	// 1 extra byte for character set
	// to test if more data available in the packet
	data = make([]byte, 3)
	read, err := stream.Read(data)
	if err != nil {
		return packet, err
	}
	packet.CapabilityFlags = uint32(data[0]) | uint32(data[1])<<8

	if read != 3 {
		return packet, nil
	}
	packet.CharacterSet = data[2]

	if _, err := stream.Read(packet.StatusFlags[:2]); err != nil {
		return packet, err
	}

	upperFlags := make([]byte, 2)
	if _, err := stream.Read(upperFlags); err != nil {
		return packet, err
	}

	packet.CapabilityFlags = ((uint32(upperFlags[0]) | uint32(upperFlags[1])<<8) << 16) | packet.CapabilityFlags

	var authDataLen uint8 = 0
	if packet.CapabilityFlags&CLIENT_PLUGIN_AUTH > 0 {
		data = make([]byte, 1)
		if _, err := stream.Read(data); err != nil {
			return packet, err
		}
		authDataLen = uint8(data[0])
	}

	// skip reserved 10 bytes
	data = make([]byte, 10)
	if _, err := stream.Read(data); err != nil {
		return packet, err
	}

	if packet.CapabilityFlags&CLIENT_SECURE_CONNECTION > 0 {
		var read uint8 = 13
		if read < authDataLen-8 {
			read = authDataLen - 8
		}

		data = make([]byte, read)
		if _, err := stream.Read(data); err != nil {
			return packet, err
		}
		packet.AuthPluginData += string(data[:len(data)-1]) // remove null-character
	}

	if packet.CapabilityFlags&CLIENT_PLUGIN_AUTH > 0 {
		data, err = readNullStr(stream)
		if err != nil {
			return packet, err
		}
		packet.AuthPluginName = string(data)
	}

	return packet, nil
}
