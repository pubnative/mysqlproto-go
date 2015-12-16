package mysqlproto

import (
	"errors"
)

type ResultSet struct {
	stream *Stream
}

func (r ResultSet) Row() ([]byte, error) {
	packet, err := r.stream.NextPacket()
	if err != nil {
		return nil, err
	}

	if packet.Payload[0] == PACKET_EOF {
		return nil, nil
	}

	return packet.Payload, nil
}

func ComQueryResponse(stream *Stream) (ResultSet, error) {
	packet, err := stream.NextPacket()
	if err != nil {
		return ResultSet{}, err
	}

	if packet.Payload[0] == PACKET_ERR {
		return ResultSet{}, errors.New(string(packet.Payload))
	}

	columns, _, _ := lenDecInt(packet.Payload)
	skip := int(columns) + 1 // skip column definition + first EOF
	for i := 0; i < skip; i++ {
		packet, err := stream.NextPacket()
		if err != nil {
			return ResultSet{}, err
		}

		if packet.Payload[0] == PACKET_ERR {
			return ResultSet{}, errors.New(string(packet.Payload))
		}
	}

	return ResultSet{stream}, nil
}
