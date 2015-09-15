// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeV10
package mysqlproto

import (
	"bytes"
	"errors"
)

type HandshakeV10 struct {
	ProtocolVersion byte
	ServerVersion   string
	ConnectionId    [4]byte
	AuthPluginData  []byte
	CapabilityFlags uint32
	CharacterSet    byte
	StatusFlags     [2]byte
	AuthPluginName  string
}

func ReadHandshakeV10(stream *Stream) (HandshakeV10, error) {
	pkt, err := stream.NextPacket()
	if err != nil {
		return HandshakeV10{}, err
	}

	data := pkt.Payload

	if data[0] == PACKET_EOF {
		return HandshakeV10{}, errors.New(string(data))
	}

	pos := 0
	packet := HandshakeV10{
		ProtocolVersion: data[pos],
	}
	pos += 1

	null := bytes.IndexByte(data[pos:], 0x00)
	packet.ServerVersion = string(data[pos : pos+null])
	pos += null + 1 // skip null terminator

	packet.ConnectionId = [4]byte{
		data[pos],
		data[pos+1],
		data[pos+2],
		data[pos+3],
	}
	pos += 4

	authDataPos := pos
	pos += 8 // 8 bytes auth data plugin

	pos += 1 // skip filler

	packet.CapabilityFlags = uint32(data[pos]) | uint32(data[pos+1])<<8
	pos += 2

	if len(data) == pos {
		packet.AuthPluginData = data[authDataPos : authDataPos+8]
		return packet, nil
	}

	packet.CharacterSet = data[pos]
	pos += 1

	packet.StatusFlags = [2]byte{data[pos], data[pos+1]}
	pos += 2

	packet.CapabilityFlags = ((uint32(data[pos]) | uint32(data[pos+1])<<8) << 16) | packet.CapabilityFlags
	pos += 2

	var authDataLen uint8 = 0
	if packet.CapabilityFlags&CLIENT_PLUGIN_AUTH > 0 {
		authDataLen = uint8(data[pos])
	}
	pos += 1

	pos += 10 // skip reserved 10 bytes

	if packet.CapabilityFlags&CLIENT_SECURE_CONNECTION > 0 {
		var read uint8 = 13
		if read < authDataLen-8 {
			read = authDataLen - 8
		}

		packet.AuthPluginData = make([]byte, read+7) // without null-character
		copy(packet.AuthPluginData[:8], data[authDataPos:authDataPos+8])
		copy(packet.AuthPluginData[8:], data[pos:pos+int(read)-1]) // remove null-character
		pos += int(read)
	}

	if packet.CapabilityFlags&CLIENT_PLUGIN_AUTH > 0 {
		null := bytes.IndexByte(data[pos:], 0x00)
		packet.AuthPluginName = string(data[pos : pos+null])
	}

	return packet, nil
}
