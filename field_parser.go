package mysqlproto

func ParseString(row []byte, offset uint64) (string, uint64) {
	count, intOffset, _ := lenDecInt(row[offset:])
	until := offset + intOffset + count
	result := string(row[offset+intOffset : until])
	return result, until
}
