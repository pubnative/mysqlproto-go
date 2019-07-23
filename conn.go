package mysqlproto

import (
	"errors"
	"io"
	"net"
	"time"
)

type Conn struct {
	*Stream
	CapabilityFlags uint32
}

var ErrNoStream = errors.New("mysqlproto: stream is not set")

func ConnectPlainHandshake(rw net.Conn, capabilityFlags uint32,
	username, password, database string,
	connectAttrs map[string]string,
	readTimeout time.Duration) (Conn, error) {
	stream := NewStream(rw, readTimeout)
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
		uint32(res[4]) | uint32(res[5])<<8 | uint32(res[6])<<16 | uint32(res[7])<<24,
	}

	if _, err = conn.Write(res); err != nil {
		return conn, err
	}

	packet, err := conn.NextPacket()
	if err != nil {
		return conn, err
	}

	if packet.Payload[0] == OK_PACKET {
		return conn, nil
	}

	return conn, parseError(packet.Payload, conn.CapabilityFlags)
}

func (c Conn) Close() error {
	if c.Stream == nil {
		return ErrNoStream
	}

	_, err := c.Write(CommandPacket(COM_QUIT, nil))
	if err != nil {
		c.Stream.Close()
		return err
	}

	pkt, err := c.NextPacket()
	if err != nil {
		if err != io.EOF {
			c.Stream.Close()
			return err
		}

		return c.Stream.Close()
	}

	if pkt.Payload[0] == OK_PACKET {
		return c.Stream.Close()
	}

	err = parseError(pkt.Payload, c.CapabilityFlags)
	if err != nil {
		c.Stream.Close()
		return err
	}

	return c.Stream.Close()
}
