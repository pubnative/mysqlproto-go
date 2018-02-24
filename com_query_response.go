package mysqlproto

import (
	"errors"
	"fmt"
)

type ResultSet struct {
	Columns []Column

	conn Conn
}

// https://dev.mysql.com/doc/internals/en/com-query-response.html#column-definition
type Column struct {
	Catalog      string
	Schema       string
	Table        string
	OrgTable     string
	Name         string
	OrgName      string
	CharacterSet uint16
	ColumnLength uint64
	ColumnType   Type
	Flags        uint16
	Decimals     byte
}

func (r ResultSet) Row() ([]byte, error) {
	packet, err := r.conn.NextPacket()
	if err != nil {
		return nil, err
	}

	if packet.Payload[0] == EOF_PACKET {
		return nil, nil
	}

	return packet.Payload, nil
}

// https://dev.mysql.com/doc/internals/en/com-query-response.html
func ComQueryResponse(conn Conn) (ResultSet, error) {
	read := func() ([]byte, error) {
		packet, err := conn.NextPacket()
		if err != nil {
			return nil, err
		}
		if len(packet.Payload) == 0 {
			return nil, errors.New("mysqlproto: empty payload")
		}
		if packet.Payload[0] == ERR_PACKET {
			return nil, parseError(packet.Payload, conn.CapabilityFlags)
		}
		return packet.Payload, nil
	}

	payload, err := read()
	if err != nil {
		return ResultSet{}, err
	}

	colCount, _, _ := lenDecInt(payload)
	columns := make([]Column, int(colCount))
	for i := 0; i < int(colCount); i++ {
		payload, err := read()
		if err != nil {
			return ResultSet{}, err
		}

		column := Column{}
		bytes, offset, _ := ReadRowValue(payload, 0)
		column.Catalog = string(bytes)

		bytes, offset, _ = ReadRowValue(payload, offset)
		column.Schema = string(bytes)

		bytes, offset, _ = ReadRowValue(payload, offset)
		column.Table = string(bytes)

		bytes, offset, _ = ReadRowValue(payload, offset)
		column.OrgTable = string(bytes)

		bytes, offset, _ = ReadRowValue(payload, offset)
		column.Name = string(bytes)

		bytes, offset, _ = ReadRowValue(payload, offset)
		column.OrgName = string(bytes)

		bytes, _, _ = ReadRowValue(payload, offset)
		if len(bytes) < 10 {
			return ResultSet{}, fmt.Errorf("mysqlproto: invalid column payload: %x", bytes)
		}

		column.CharacterSet = uint16(bytes[0]) | uint16(bytes[1])<<8
		column.ColumnLength = uint64(bytes[2]) | uint64(bytes[3])<<8 | uint64(bytes[4])<<16 | uint64(bytes[5])<<32
		column.ColumnType = Type(bytes[6])
		column.Flags = uint16(bytes[7]) | uint16(bytes[9])<<8
		column.Decimals = bytes[9]

		columns[i] = column
	}

	payload, err = read()
	if err != nil {
		return ResultSet{}, err
	}
	if payload[0] != EOF_PACKET {
		return ResultSet{}, parseError(payload, conn.CapabilityFlags)
	}

	rs := ResultSet{
		Columns: columns,
		conn:    conn,
	}
	return rs, nil
}
