package mysqlproto

import (
	"bytes"
	"net"
	"time"
)

const PACKET_BUFFER_SIZE = 1500 // default MTU

type Stream struct {
	stream   net.Conn
	buffer   []byte
	read     int
	left     int
	syscalls int
	ReadTimeout time.Duration
}

func NewStream(stream net.Conn, readTimeout time.Duration) *Stream {
	return &Stream{stream, nil, 0, 0, 0, readTimeout}
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
		if s.ReadTimeout > 0 {
			if err = s.stream.SetReadDeadline(time.Now().Add(s.ReadTimeout)); err != nil {
				return
			}
		}

		nn, err = s.stream.Read(buf[n:])
		s.syscalls += 1
		n += nn
	}
	return
}

// For testing

type buffer struct {
	*bytes.Buffer
	closed  bool
	writeFn func([]byte) (int, error)
}

func newBuffer(data []byte) *buffer {
	return &buffer{bytes.NewBuffer(data), false, nil}
}

func (b *buffer) Close() error {
	b.closed = true
	return nil
}

func (b *buffer) Write(data []byte) (int, error) {
	if b.writeFn == nil {
		return 0, nil
	}

	return b.writeFn(data)
}
func (b *buffer) RemoteAddr() net.Addr { return MockAddr{} }
func (b *buffer) LocalAddr() net.Addr { return MockAddr{} }
func (b *buffer) SetDeadline(t time.Time) error { return nil}
func (b *buffer) SetReadDeadline(t time.Time) error { return nil}
func (b *buffer) SetWriteDeadline(t time.Time) error { return nil}

type MockAddr struct {}
func (m MockAddr) Network() string { return "" }
func (m MockAddr) String() string { return "" }
