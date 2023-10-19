package binarystream

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var ErrBufferUnderflow = fmt.Errorf("buffer underflow")

type BinaryStream struct {
	r         io.Reader
	buf       []byte
	byteOrder binary.ByteOrder
}

func NewReaderFromBytes(buf []byte, byteOrder binary.ByteOrder) *BinaryStream {
	return NewReader(bytes.NewReader(buf), byteOrder)
}

func NewReader(r io.Reader, byteOrder binary.ByteOrder) *BinaryStream {
	return &BinaryStream{
		r:         r,
		byteOrder: byteOrder,
	}
}

func (s *BinaryStream) ensureData(n int) error {
	extraRequired := n - len(s.buf)
	if extraRequired <= 0 {
		return nil
	}
	data, err := io.ReadAll(io.LimitReader(s.r, int64(extraRequired)))
	if err != nil {
		return err
	}
	s.buf = append(s.buf, data...)
	if len(s.buf) < n {
		return ErrBufferUnderflow
		panic("buffer underflow") // shouldn't happen
	}
	return nil
}

func (s *BinaryStream) ReadNullTerminatedString() (string, error) {
	str, err := s.PeekNullTerminatedString()
	if err != nil {
		return "", err
	}
	n := len(str) + 1
	s.buf = s.buf[n:]
	return str, nil
}

func (s *BinaryStream) PeekNullTerminatedString() (string, error) {
	var b []byte
	i := 0
	for {
		err := s.ensureData(i + 1)
		if err != nil {
			return "", err
		}
		if s.buf[i] == 0x00 {
			break
		}
		b = append(b, s.buf[i])
	}
	return string(b), nil
}

func (s *BinaryStream) ReadUint8PrefixedString() (string, error) {
	n, err := s.ReadUint8()
	if err != nil {
		return "", err
	}
	b, err := s.ReadBytes(int(n))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *BinaryStream) ReadUint16PrefixedString() (string, error) {
	n, err := s.PeekUint16()
	if err != nil {
		return "", err
	}
	err = s.ensureData(2 + int(n))
	if err != nil {
		return "", err
	}
	s.buf = s.buf[2:]
	b, err := s.ReadBytes(int(n))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *BinaryStream) ReadUint32PrefixedString() (string, error) {
	n, err := s.PeekUint32()
	if err != nil {
		return "", err
	}
	err = s.ensureData(4 + int(n))
	if err != nil {
		return "", err
	}
	s.buf = s.buf[4:]
	b, err := s.ReadBytes(int(n))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *BinaryStream) ReadUint64PrefixedString() (string, error) {
	n, err := s.PeekUint64()
	if err != nil {
		return "", err
	}
	err = s.ensureData(8 + int(n))
	if err != nil {
		return "", err
	}
	s.buf = s.buf[8:]
	b, err := s.ReadBytes(int(n))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *BinaryStream) ReadBytes(n int) ([]byte, error) {
	err := s.ensureData(n)
	if err != nil {
		return nil, err
	}
	b := s.buf[:n][:]
	s.buf = s.buf[n:]
	return b, nil
}

func (s *BinaryStream) ReadUint64() (uint64, error) {
	x, err := s.PeekUint64()
	if err != nil {
		return 0, err
	}
	s.buf = s.buf[8:]
	return x, nil
}

func (s *BinaryStream) PeekUint64() (uint64, error) {
	err := s.ensureData(8)
	if err != nil {
		return 0, err
	}
	var val uint64
	err = binary.Read(bytes.NewReader(s.buf), s.byteOrder, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *BinaryStream) ReadUint32() (uint32, error) {
	x, err := s.PeekUint32()
	if err != nil {
		return 0, err
	}
	s.buf = s.buf[4:]
	return x, nil
}

func (s *BinaryStream) PeekUint32() (uint32, error) {
	err := s.ensureData(4)
	if err != nil {
		return 0, err
	}
	var val uint32
	err = binary.Read(bytes.NewReader(s.buf), s.byteOrder, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *BinaryStream) ReadUint16() (uint16, error) {
	x, err := s.PeekUint16()
	if err != nil {
		return 0, err
	}
	s.buf = s.buf[2:]
	return x, nil
}

func (s *BinaryStream) PeekUint16() (uint16, error) {
	err := s.ensureData(2)
	if err != nil {
		return 0, err
	}
	var val uint16
	err = binary.Read(bytes.NewReader(s.buf), s.byteOrder, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *BinaryStream) ReadUint8() (uint8, error) {
	x, err := s.PeekUint8()
	if err != nil {
		return 0, err
	}
	s.buf = s.buf[1:]
	return x, nil
}

func (s *BinaryStream) PeekUint8() (uint8, error) {
	err := s.ensureData(1)
	if err != nil {
		return 0, err
	}
	var val uint8
	err = binary.Read(bytes.NewReader(s.buf), s.byteOrder, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (s *BinaryStream) Skip(n int) error {
	err := s.ensureData(n)
	if err != nil {
		return err
	}
	s.buf = s.buf[n:]
	return nil
}

func (s *BinaryStream) ReadFixedString(n int) (string, error) {
	b, err := s.ReadBytes(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
