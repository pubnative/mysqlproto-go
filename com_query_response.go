package mysqlproto

import (
	"errors"
	"io"
)

type ResultSet struct {
	streamPacket *StreamPacket
}

func (r ResultSet) Row() ([]byte, error) {
	packet, err := r.streamPacket.NextPacket()
	if err != nil {
		return nil, err
	}

	if packet.Payload[0] == 0xfe { // EOF
		return nil, nil
	}

	return packet.Payload, nil
}

func ComQueryResponse(stream io.Reader) (ResultSet, error) {
	streamPkt := NewStreamPacket(stream)

	packet, err := streamPkt.NextPacket()
	if err != nil {
		return ResultSet{}, err
	}

	if packet.Payload[0] == 0xff {
		return ResultSet{}, errors.New(string(packet.Payload))
	}

	columns, _, _ := lenDecInt(packet.Payload)
	skip := int(columns) + 1 // skip column definition + first EOF
	for i := 0; i < skip; i++ {
		packet, err := streamPkt.NextPacket()
		if err != nil {
			return ResultSet{}, err
		}

		if packet.Payload[0] == 0xff {
			return ResultSet{}, errors.New(string(packet.Payload))
		}
	}

	return ResultSet{streamPkt}, nil
}
