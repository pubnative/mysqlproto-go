package mysqlproto

import (
	"io"
)

const PACKET_OK = 0x00
const PACKET_ERR = 0xff
const PACKET_EOF = 0xfe

type Packet struct {
	SequenceID byte
	Payload    []byte
}

func (p Proto) ReadPacket(stream io.Reader) (Packet, error) {
	if _, err := io.ReadFull(stream, p.header); err != nil {
		return Packet{}, err
	}

	length := uint32(byte(p.header[0]) | p.header[1]<<8 | p.header[2]<<16)
	seqID := p.header[3]
	payload := make([]byte, length)
	if _, err := io.ReadFull(stream, payload); err != nil {
		return Packet{}, err
	}

	return Packet{seqID, payload}, nil
}
