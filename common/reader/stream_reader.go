package reader

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
)

type StreamReader interface {
	ChangeBytesOrder(isLittle bool)
	GetInt64() (v int64, err error)
	GetInt32() (v int32, err error)
	GetInt16() (v int16, err error)
	GetInt8() (v int8, err error)
	GetNBytes(n int) (res []byte, err error)
	GetString(n int) (string, error)
	ResetGzipReader() error
	GetUInt64() (v uint64, err error)
	GetUInt32() (v uint32, err error)
	GetUInt16() (v uint16, err error)
	GetUInt8() (v uint8, err error)

	Remain() int32
	Peek(n int) ([]byte, error)
	Discard(n int) error
	GetReadNum() int32
}

type streamReader struct {
	*bufio.Reader
	buf []byte
	binary.ByteOrder
	r int32 //hasRead
}

func NewStreamReader(reader *bufio.Reader) StreamReader {
	return &streamReader{Reader: reader,
		buf:       make([]byte, 1024),
		ByteOrder: binary.LittleEndian}
}

func (s *streamReader) Remain() int32 {
	return int32(s.Reader.Size()) - s.r
}
func (s *streamReader) GetReadNum() int32 {
	return s.r
}
func (s *streamReader) Discard(n int) error {
	dis, err := s.Reader.Discard(n)
	if err != nil {
		return err
	}
	s.r += int32(dis)
	return nil
}

func (s *streamReader) ChangeBytesOrder(isLittle bool) {
	if isLittle {
		s.ByteOrder = binary.LittleEndian
	} else {
		s.ByteOrder = binary.BigEndian
	}

}
func (s *streamReader) GetUInt64() (v uint64, err error) {
	bytes, err := s.GetNBytes(8)
	if err != nil {
		return 0, err
	}
	return s.ByteOrder.Uint64(bytes), err
}

func (s *streamReader) GetUInt32() (v uint32, err error) {
	bytes, err := s.GetNBytes(4)
	if err != nil {
		return 0, err
	}
	return s.ByteOrder.Uint32(bytes), err
}

func (s *streamReader) GetUInt16() (v uint16, err error) {
	bytes, err := s.GetNBytes(4)
	if err != nil {
		return 0, err
	}
	return s.ByteOrder.Uint16(bytes), err
}

func (s *streamReader) GetUInt8() (v uint8, err error) {
	bytes, err := s.GetNBytes(4)
	if err != nil {
		return 0, err
	}
	return bytes[0], err
}
func (s *streamReader) GetInt64() (v int64, err error) {
	tmp, err := s.GetUInt64()
	if err != nil {
		return 0, err
	}
	return int64(tmp), err
}

func (s *streamReader) GetInt32() (v int32, err error) {
	tmp, err := s.GetUInt32()
	if err != nil {
		return 0, err
	}
	return int32(tmp), err
}

func (s *streamReader) GetInt16() (v int16, err error) {
	tmp, err := s.GetUInt16()
	if err != nil {
		return 0, err
	}
	return int16(tmp), err
}

func (s *streamReader) GetInt8() (v int8, err error) {
	tmp, err := s.GetUInt8()
	if err != nil {
		return 0, err
	}
	return int8(tmp), err
}

func (s *streamReader) GetString(n int) (string, error) {
	res, err := s.GetNBytes(n)
	if err != nil {
		return "", err
	}
	return string(res), err
}

func (s *streamReader) GetNBytes(n int) (res []byte, err error) {
	if n > 1024 {
		return nil, errors.New("streamReader limit bytes 1024")
	}
	_, err = io.ReadFull(s.Reader, s.buf[:n])
	if err != nil {
		return nil, err
	}
	res = make([]byte, n)
	copy(res, s.buf[:n])
	s.r += int32(n)
	return res, nil
}

func (s *streamReader) ResetGzipReader() error {
	r, err := gzip.NewReader(s.Reader)
	if err != nil {
		return err
	}
	s.Reader = bufio.NewReader(r)
	return nil
}
