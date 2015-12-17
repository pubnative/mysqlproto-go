package mysqlproto

import (
	"bytes"
	"io"
)

const PACKET_BUFFER_SIZE = 1500 // default MTU

type Stream struct {
	stream   io.ReadWriteCloser
	buffer   []byte
	read     int
	left     int
	syscalls int
}

func NewStream(stream io.ReadWriteCloser) *Stream {
	return &Stream{stream, nil, 0, 0, 0}
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
		read, err := s.readAtLeast(s.buffer[s.read+s.left:], 3-s.left)
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
		read, err := s.readAtLeast(s.buffer[s.read+s.left:], total-s.left)
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

func (s *Stream) Syscalls() int {
	return s.syscalls
}

func (s *Stream) ResetStats() {
	s.syscalls = 0
}

func (s *Stream) readAtLeast(buf []byte, min int) (n int, err error) {
	for n < min && err == nil {
		var nn int
		nn, err = s.stream.Read(buf[n:])
		s.syscalls += 1
		n += nn
	}
	return
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
