package mysqlproto

import (
	"fmt"
	"strconv"
)

// https://dev.mysql.com/doc/internals/en/generic-response-packets.html
const (
	OK_PACKET  byte = 0x00
	ERR_PACKET byte = 0xff
	EOF_PACKET byte = 0xfe
)

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

type OKPacket struct {
	Header              byte // 0x00 or 0xfe
	AffectedRows        uint64
	LastInsertID        uint64
	StatusFlags         uint16
	Warnings            uint16
	Info                string
	SessionStateChanges string
}

// https://dev.mysql.com/doc/internals/en/packet-OK_Packet.html
func ParseOKPacket(data []byte, capabilityFlags uint32) (OKPacket, error) {
	if len(data) == 0 || (data[0] != OK_PACKET && data[0] != EOF_PACKET) {
		return OKPacket{}, fmt.Errorf("mysqlproto: invalid OK_PACKET payload: %x", data)
	}

	offset := 0
	header := data[offset]
	offset += 1
	affectedRows, offsetInt, _ := lenDecInt(data[1:])
	offset += int(offsetInt)
	lastInsertID, offsetInt, _ := lenDecInt(data[offset:])
	offset += int(offsetInt)

	var statusFlags, warnings uint16
	if capabilityFlags&CLIENT_PROTOCOL_41 > 0 {
		statusFlags = uint16(data[offset]) | uint16(data[offset+1])<<8
		warnings = uint16(data[offset+2]) | uint16(data[offset+3])<<8
		offset += 4
	} else if capabilityFlags&CLIENT_TRANSACTIONS > 0 {
		statusFlags = uint16(data[offset]) | uint16(data[offset+1])<<8
		offset += 2
	}

	var info, sessionStateChanges string
	if capabilityFlags&CLIENT_SESSION_TRACK > 0 {
		size, intOffset, _ := lenDecInt(data[offset:])
		info = string(data[offset+int(intOffset) : offset+int(intOffset)+int(size)])
		offset += int(intOffset) + int(size)

		if statusFlags&SERVER_SESSION_STATE_CHANGED > 0 {
			size, intOffset, _ = lenDecInt(data[offset:])
			sessionStateChanges = string(data[offset+int(intOffset) : offset+int(intOffset)+int(size)])
			offset += int(intOffset) + int(size)
		}
	} else {
		// Documentation says that in this case info is string<EOF> type
		// but apparently it's string<lenenc>
		// https://github.com/mysql/mysql-server/blob/5.6/sql/protocol.cc#L248
		// https://github.com/mysql/mysql-server/blob/5.6/sql/protocol.cc#L585
		_, infoOffset, _ := lenDecInt(data[offset:])
		info = string(data[offset+int(infoOffset):])
	}

	pkt := OKPacket{
		Header:              header,
		AffectedRows:        affectedRows,
		LastInsertID:        lastInsertID,
		StatusFlags:         statusFlags,
		Warnings:            warnings,
		Info:                info,
		SessionStateChanges: sessionStateChanges,
	}

	return pkt, nil
}

// https://dev.mysql.com/doc/internals/en/packet-ERR_Packet.html
func ParseERRPacket(data []byte, capabilityFlags uint32) (ERRPacket, error) {
	if len(data) == 0 || data[0] != ERR_PACKET {
		return ERRPacket{}, fmt.Errorf("mysqlproto: invalid ERR_PACKET payload: %x", data)
	}

	pkt := ERRPacket{
		Header:    data[0],
		ErrorCode: uint16(data[1]) | uint16(data[2])<<8,
	}

	offset := 3
	if capabilityFlags&CLIENT_PROTOCOL_41 > 0 {
		pkt.SQLStateMarker = string(data[3])
		pkt.SQLState = string(data[4:9])
		offset = 9
	}

	pkt.ErrorMessage = string(data[offset:])

	return pkt, nil
}

// https://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html
func (p ERRPacket) Error() string {
	return "mysqlproto: Error: " + strconv.Itoa(int(p.ErrorCode)) +
		" SQLSTATE: " + p.SQLState +
		" Message: " + p.ErrorMessage
}

func parseError(data []byte, capabilityFlags uint32) error {
	pkt, err := ParseERRPacket(data, capabilityFlags)
	if err != nil {
		return err
	}
	return pkt
}
