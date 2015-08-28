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

func ReadPacket(stream io.Reader) (Packet, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(stream, header); err != nil {
		return Packet{}, err
	}

	length := uint32(byte(header[0]) | header[1]<<8 | header[2]<<16)
	seqID := header[3]
	payload := make([]byte, length)
	if _, err := io.ReadFull(stream, payload); err != nil {
		return Packet{}, err
	}

	return Packet{seqID, payload}, nil
}
