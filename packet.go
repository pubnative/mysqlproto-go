package mysqlproto

import (
	"errors"
	"strconv"
)

const PACKET_OK = 0x00
const PACKET_ERR = 0xff
const PACKET_EOF = 0xfe

var ErrERRPacketPayload = errors.New("Invalid ERR_PACKET payload.")

type Packet struct {
	SequenceID byte
	Payload    []byte
}

type ERRPacket struct {
	Header         byte // always 0xff
	ErrorCode      uint16
	SQLStateMarker string
	SQLState       string
	ErrorMessage   string
}

func ParseERRPacket(data []byte) (ERRPacket, error) {
	if len(data) == 0 || data[0] != PACKET_ERR {
		return ERRPacket{}, ErrERRPacketPayload
	}

	pkt := ERRPacket{
		Header:         data[0],
		ErrorCode:      uint16(data[1]) | uint16(data[2])<<8,
		SQLStateMarker: string(data[3]),
		SQLState:       string(data[4:9]),
		ErrorMessage:   string(data[9:]),
	}

	return pkt, nil
}

// https://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html
func (p ERRPacket) Error() string {
	return "mysqlproto: Error: " + strconv.Itoa(int(p.ErrorCode)) +
		" SQLSTATE: " + p.SQLState +
		" Message: " + p.ErrorMessage
}

func parseError(data []byte) error {
	pkt, err := ParseERRPacket(data)
	if err != nil {
		return err
	}
	return pkt
}
