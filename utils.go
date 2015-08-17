package mysqlproto

import (
	"io"
)

func readNullStr(stream io.Reader) ([]byte, error) {
	data := make([]byte, 1)
	idx := 0
	for {
		_, err := stream.Read(data[idx:])
		if err != nil {
			return data, err
		}

		if data[idx] == 0x00 {
			break
		}

		data = append(data, 0)
		idx += 1
	}

	return data[:len(data)-1], nil // remove null-character
}

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

func lenEncStr(s string) []byte {
	size := lenEncInt(uint64(len(s)))
	return append(size, s...)
}
