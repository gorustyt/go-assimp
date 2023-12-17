package reader

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
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

	GetFloat32() (v float32, err error)
	GetFloat64() (v float64, err error)

	Peek(n int) ([]byte, error)
	Discard(n int) error
	SetCurPos(pos int)
	GetCurPos() int
	Remain() int
	ResetData()
}

type streamReader struct {
	*BaseReader
	binary.ByteOrder
	mCurrent int
	offset   bool
	mBuffer  int
}

func newStreamReader(b *BaseReader) StreamReader {
	return &streamReader{BaseReader: b, ByteOrder: binary.LittleEndian}
}

func (s *streamReader) SetCurPos(pos int) {
	s.offset = true
	s.mCurrent = s.mBuffer + pos
}
func (s *streamReader) GetCurPos() int {
	return s.mCurrent - s.mBuffer
}
func (s *streamReader) Remain() int {
	return len(s.data) - s.GetCurPos()
}

func (s *streamReader) ResetData() {
	s.data = s.data[s.mCurrent:]
	s.mCurrent = 0
	s.mBuffer = 0
}

func (s *streamReader) Discard(n int) error {
	if s.mCurrent+n > len(s.data) {
		return io.EOF
	}
	s.incPos(n)
	return nil
}

func (s *streamReader) incPos(n int) {
	if !s.offset {
		s.mCurrent += n
		s.mBuffer += n
	} else {
		s.mCurrent += n
	}

}
func (s *streamReader) Peek(n int) ([]byte, error) {
	return s.read(n, true)
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
	bytes, err := s.GetNBytes(2)
	if err != nil {
		return 0, err
	}
	return s.ByteOrder.Uint16(bytes), err
}

func (s *streamReader) GetUInt8() (v uint8, err error) {
	bytes, err := s.GetNBytes(1)
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

func (s *streamReader) GetFloat64() (v float64, err error) {
	res, err := s.GetNBytes(8)
	if err != nil {
		return v, err
	}
	bits := binary.LittleEndian.Uint64(res)
	float := math.Float64frombits(bits)
	return float, nil
}

func (s *streamReader) GetFloat32() (v float32, err error) {
	res, err := s.GetNBytes(4)
	if err != nil {
		return v, err
	}
	bits := binary.LittleEndian.Uint32(res)
	float := math.Float32frombits(bits)
	return float, nil
}
func (s *streamReader) GetString(n int) (string, error) {
	res, err := s.GetNBytes(n)
	if err != nil {
		return "", err
	}
	return string(res), err
}

func (s *streamReader) GetNBytes(n int) (res []byte, err error) {
	res, err = s.read(n, false)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *streamReader) read(n int, isPeek bool) (res []byte, err error) {
	if n > 1024 {
		return res, errors.New("streamReader limit bytes 1024")
	}
	if s.mCurrent+n > len(s.data) {
		return res, io.EOF
	}
	res = make([]byte, n)
	copy(res, s.data[s.mCurrent:s.mCurrent+n])
	if !isPeek {
		s.incPos(n)
	}
	return res, nil
}

func (s *streamReader) ResetGzipReader() error {
	//r, err := gzip.NewReader(s.Reader)
	//if err != nil {
	//	return err
	//}
	//s.Reader = bufio.NewReader(r)
	return nil
}
