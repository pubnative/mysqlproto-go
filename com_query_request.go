package mysqlproto

func ComQueryRequest(query []byte) []byte {
	l := len(query) + 1 // + command byte

	packet := make([]byte, l+4)
	packet[0] = byte(l)
	packet[1] = byte(l >> 8)
	packet[2] = byte(l >> 16)
	packet[3] = byte(0x00) // sequence ID always 0x00

	packet[4] = COM_QUERY
	copy(packet[5:], query)

	return packet
}
