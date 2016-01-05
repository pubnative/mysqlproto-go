package mysqlproto

func ComQueryRequest(query []byte) []byte {
	return CommandPacket(COM_QUERY, query)
}
