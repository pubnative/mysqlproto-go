package mysqlproto

import (
	"io"
)

type Conn struct {
	*Stream
	CapabilityFlags uint32
}

func Handshake(rw io.ReadWriteCloser, capabilityFlags uint32,
	username, password, database string,
	connectAttrs map[string]string) (Conn, error) {
	stream := NewStream(rw)
	handshakeV10, err := ReadHandshakeV10(stream)
	if err != nil {
		return Conn{}, err
	}

	flags := handshakeV10.CapabilityFlags & capabilityFlags

	res := HandshakeResponse41(
		flags,
		handshakeV10.CharacterSet,
		username,
		password,
		handshakeV10.AuthPluginData,
		database,
		handshakeV10.AuthPluginName,
		connectAttrs,
	)

	conn := Conn{
		stream,
		uint32(res[4]) | uint32(res[5])<<8 | uint32(res[6])<<12 | uint32(res[7])<<16,
	}

	if _, err = conn.Write(res); err != nil {
		return conn, err
	}

	packet, err := conn.NextPacket()
	if err != nil {
		return conn, err
	}

	if packet.Payload[0] == PACKET_OK {
		return conn, nil
	} else {
		return conn, parseError(packet.Payload, conn.CapabilityFlags)
	}
}
