package mysqlproto

func ReadRowValue(row []byte, offset uint64) ([]byte, uint64, bool) {
	count, intOffset, null := lenDecInt(row[offset:])
	until := offset + intOffset + count
	return row[offset+intOffset : until], until, null
}

// https://dev.mysql.com/doc/internals/en/integer.html#packet-Protocol::LengthEncodedInteger
func lenEncInt(i uint64) []byte {
	if i < 251 {
		return []byte{byte(i)}
	} else if i >= 251 && i < 1<<16 {
		return []byte{0xfc, byte(i), byte(i >> 8)}
	} else if i >= 1<<16 && i < 1<<24 {
		return []byte{0xfd, byte(i), byte(i >> 8), byte(i >> 16)}
	} else {
		return []byte{0xfe, byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24),
			byte(i >> 32), byte(i >> 40), byte(i >> 48), byte(i >> 56),
		}
	}
}

func lenDecInt(b []byte) (uint64, uint64, bool) { // int, offset, is null
	switch b[0] {
	case 0xfb:
		return 0, 1, true
	case 0xfc:
		return uint64(b[1]) | uint64(b[2])<<8, 3, false
	case 0xfd:
		return uint64(b[1]) | uint64(b[2])<<8 | uint64(b[3])<<16, 4, false
	case 0xfe:
		return uint64(b[1]) | uint64(b[2])<<8 | uint64(b[3])<<16 |
			uint64(b[4])<<24 | uint64(b[5])<<32 | uint64(b[6])<<40 |
			uint64(b[7])<<48 | uint64(b[8])<<56, 9, false
	default:
		return uint64(b[0]), 1, false
	}
}

func lenEncStr(s string) []byte {
	size := lenEncInt(uint64(len(s)))
	return append(size, s...)
}
