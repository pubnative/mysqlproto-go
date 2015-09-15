package mysqlproto

import (
	"bytes"
	"io"
)

const PACKET_OK = 0x00
const PACKET_ERR = 0xff
const PACKET_EOF = 0xfe

const PACKET_BUFFER_SIZE = 512

type Packet struct {
	SequenceID byte
	Payload    []byte
}

type Stream struct {
	stream io.ReadWriteCloser
	buffer []byte
	read   int
	left   int
}

func NewStream(stream io.ReadWriteCloser) *Stream {
	return &Stream{stream, nil, 0, 0}
}

func (s *Stream) Write(data []byte) (int, error) {
	return s.stream.Write(data)
}

func (s *Stream) Close() error {
	return s.stream.Close()
}

func (s *Stream) NextPacket() (Packet, error) {
	scale := func(size int) {
		if size < PACKET_BUFFER_SIZE {
			size = PACKET_BUFFER_SIZE
		}
		buf := make([]byte, size)
		copy(buf, s.buffer[s.read:s.read+s.left])
		s.buffer = buf
		s.read = 0
	}

	if len(s.buffer)-s.read < 3 { // size of the packet
		scale(PACKET_BUFFER_SIZE)
	}

	if s.left < 3 {
		read, err := io.ReadAtLeast(s.stream, s.buffer[s.read+s.left:], 3-s.left)
		if err != nil {
			return Packet{}, err
		}
		s.left += read
	}

	length := int(uint32(s.buffer[s.read]) | uint32(s.buffer[s.read+1])<<8 | uint32(s.buffer[s.read+2])<<16)
	total := length + 4
	if total > len(s.buffer)-s.read {
		scale(total)
	}

	if total > s.left {
		read, err := io.ReadAtLeast(s.stream, s.buffer[s.read+s.left:], total-s.left)
		if err != nil {
			return Packet{}, err
		}
		s.left += read
	}

	packet := Packet{
		SequenceID: s.buffer[s.read+3],
		Payload:    s.buffer[s.read+4 : s.read+total],
	}

	s.left -= total
	s.read += total

	return packet, nil
}

// For testing

type buffer struct {
	*bytes.Buffer
}

func newBuffer(data []byte) *buffer {
	return &buffer{bytes.NewBuffer(data)}
}

func (b *buffer) Close() error {
	return nil
}
