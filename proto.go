package mysqlproto

type Proto struct {
	header []byte
}

func NewProto() Proto {
	return Proto{
		header: make([]byte, 4),
	}
}
