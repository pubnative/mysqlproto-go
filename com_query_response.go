package mysqlproto

import (
	"errors"
	"io"
)

type ResultSet struct {
	stream io.Reader
}

type ResultSetRow struct {
	packet []byte
}

func (r ResultSetRow) ReadString(offset uint64) (string, uint64) {
	count, intOffset, _ := lenDecInt(r.packet[offset:])
	until := offset + intOffset + count
	result := string(r.packet[offset+intOffset : until])
	return result, until
}

func (r ResultSet) Row() (ResultSetRow, bool, error) {
	packet, err := ReadPacket(r.stream)
	if err != nil {
		return ResultSetRow{}, false, err
	}

	if packet.Payload[0] == 0xfe { // EOF
		return ResultSetRow{}, true, nil
	}

	return ResultSetRow{packet.Payload}, false, nil
}

func ComQueryResponse(stream io.Reader) (ResultSet, error) {
	packet, err := ReadPacket(stream)
	if err != nil {
		return ResultSet{}, err
	}

	if packet.Payload[0] == 0xff {
		return ResultSet{}, errors.New(string(packet.Payload))
	}

	columns, _, _ := lenDecInt(packet.Payload)
	skip := int(columns) + 1 // skip column definition + first EOF
	for i := 0; i < skip; i++ {
		packet, err := ReadPacket(stream)
		if err != nil {
			return ResultSet{}, err
		}

		if packet.Payload[0] == 0xff {
			return ResultSet{}, errors.New(string(packet.Payload))
		}
	}

	return ResultSet{stream}, nil
}
