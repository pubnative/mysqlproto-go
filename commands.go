package mysqlproto

// https://dev.mysql.com/doc/internals/en/command-phase.html
const (
	COM_SLEEP byte = iota
	COM_QUIT
	COM_INIT_DB
	COM_QUERY
	COM_FIELD_LIST
	COM_CREATE_DB
	COM_DROP_DB
	COM_REFRESH
	COM_SHUTDOWN
	COM_STATISTICS
	COM_PROCESS_INFO
	COM_CONNECT
	COM_PROCESS_KILL
	COM_DEBUG
	COM_PING
	COM_TIME
	COM_DELAYED_INSERT
	COM_CHANGE_USER
	COM_BINLOG_DUMP
	COM_TABLE_DUMP
	COM_CONNECT_OUT
	COM_REGISTER_SLAVE
	COM_STMT_EXECUTE
	COM_STMT_SEND_LONG_DATA
	COM_STMT_CLOSE
	COM_STMT_RESET
	COM_SET_OPTION
	COM_STMT_FETCH
	COM_DAEMON
	COM_BINLOG_DUMP_GTID
	COM_RESET_CONNECTION
)

func CommandPacket(command byte, payload []byte) []byte {
	size := len(payload) + 1 // command byte

	packet := make([]byte, size+4)
	packet[0] = byte(size)
	packet[1] = byte(size >> 8)
	packet[2] = byte(size >> 16)
	packet[3] = byte(0x00) // sequence ID always 0x00
	packet[4] = command
	copy(packet[5:], payload)

	return packet
}
