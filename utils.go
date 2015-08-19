package mysqlproto

import (
	"io"
)

func readNullStr(stream io.Reader) ([]byte, error) {
	data := make([]byte, 1)
	idx := 0
	br := false
	for {
		_, err := stream.Read(data[idx:])
		if err != nil {
			return data, err
		}

		if data[idx] == 0x00 {
			if br {
				break
			}
		} else {
			br = true
		}

		data = append(data, 0)
		idx += 1
	}

	return data[:len(data)-1], nil // remove null-character
}
