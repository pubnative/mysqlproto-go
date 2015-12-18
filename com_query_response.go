package mysqlproto

type ResultSet struct {
	conn Conn
}

func (r ResultSet) Row() ([]byte, error) {
	packet, err := r.conn.NextPacket()
	if err != nil {
		return nil, err
	}

	if packet.Payload[0] == PACKET_EOF {
		return nil, nil
	}

	return packet.Payload, nil
}

func ComQueryResponse(conn Conn) (ResultSet, error) {
	packet, err := conn.NextPacket()
	if err != nil {
		return ResultSet{}, err
	}

	if packet.Payload[0] == PACKET_ERR {
		return ResultSet{}, parseError(packet.Payload, conn.CapabilityFlags)
	}

	columns, _, _ := lenDecInt(packet.Payload)
	skip := int(columns) + 1 // skip column definition + first EOF
	for i := 0; i < skip; i++ {
		packet, err := conn.NextPacket()
		if err != nil {
			return ResultSet{}, err
		}

		if packet.Payload[0] == PACKET_ERR {
			return ResultSet{}, parseError(packet.Payload, conn.CapabilityFlags)
		}
	}

	return ResultSet{conn}, nil
}
