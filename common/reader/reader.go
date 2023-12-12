package reader

import (
	"assimp/common/logger"
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	ErrBadParams = errors.New("error bad params")
	ErrBadKey    = errors.New("error bad key")
)

type AiReader struct {
	data   []byte
	reader *bufio.Reader
}

func NewReader(data []byte) (r *AiReader, err error) {
	r = &AiReader{data: data}
	err = r.ConvertToUTF8()
	if err != nil {
		return r, err
	}
	r.reader = bufio.NewReader(bytes.NewReader(r.data))
	return r, err
}

func (r *AiReader) CheckMagicToken(magic []byte, offset int, size int) error {
	hasSize := r.reader.Size()
	logger.InfoF("init aiReader size:%v", hasSize)
	if size > 16 {
		logger.FatalF("not expect magic token")
	}
	data := make([]byte, 16)
	n, err := r.reader.Read(data[:size])
	if err != nil {
		return err
	}
	if n != size {
		logger.ErrorF("read size:%v not enough:%v", n, size)
		return ErrBadParams
	}
	for _ = range magic {
		// also check against big endian versions of tokens with size 2,4
		// that's just for convenience, the chance that we cause conflicts
		// is quite low and it can save some lines and prevent nasty bugs
		if 2 == size {
			var magic_u16_little uint16
			var magic_u16_big uint16
			var data_u16 uint16
			binary.BigEndian.PutUint16(data[:2], data_u16)
			binary.BigEndian.PutUint16(magic[:2], magic_u16_big)
			binary.LittleEndian.PutUint16(magic[:2], magic_u16_little)
			if data_u16 == magic_u16_little || data_u16 == magic_u16_big {
				return nil
			}
		} else if 4 == size {
			var magic_u32_little uint32
			var magic_u32_big uint32
			var data_u32 uint32
			binary.BigEndian.PutUint32(data[:4], data_u32)
			binary.BigEndian.PutUint32(magic[:4], magic_u32_big)
			binary.LittleEndian.PutUint32(magic[:4], magic_u32_little)
			if data_u32 == magic_u32_little || data_u32 == magic_u32_big {
				return nil
			}
		} else {
			// any length ... just compare
			if bytes.Equal(magic[:size], data[:size]) {
				return nil
			}
		}
		magic = magic[size:]
	}
	return ErrBadParams
}

func (r *AiReader) GetLineReader() LineReader {
	return NewLineReader(r.reader)
}

func (r *AiReader) GetReader() *bufio.Reader {
	return r.reader
}

func (r *AiReader) GetStreamReader() StreamReader {
	return NewStreamReader(r.reader)
}

// Convert to UTF8 data
func (r *AiReader) ConvertToUTF8() (err error) {
	//ConversionResult result;
	size := len(r.data)
	if size < 8 {
		logger.FatalF("File is too small")
	}
	data4Bytes := r.data[:4]
	// UTF 8 with BOM
	if r.data[0] == 0xEF && r.data[1] == 0xBB && r.data[2] == 0xBF {
		logger.Info("Found UTF-8 BOM ...")
		r.data = r.data[3:]
		return nil
	}

	// UTF 32 BE with BOM
	var buf [8]byte
	binary.PutVarint(buf[:], 0xFFFE0000)
	if bytes.Equal(data4Bytes, buf[:4]) {
		logger.FatalF("UTF-32 BOM not support .")
	}

	// UTF 32 LE with BOM

	buf = [8]byte{}
	binary.PutVarint(buf[:], 0x0000FFFE)
	if bytes.Equal(data4Bytes, buf[:4]) {
		logger.FatalF("UTF-32 BOM not support .")
	}

	// UTF 16 BE with BOM

	if r.data[0] == 0xfe && r.data[1] == 0xff {
		r.data, _, err = transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder(), r.data)
		if err != nil {
			return err
		}
	}

	//UTF 16 LE with BOM
	if r.data[0] == 0xff && r.data[1] == 0xfe {
		r.data, _, err = transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder(), r.data)
		if err != nil {
			return err
		}
	}

	logger.Info("this is utf8 data ,not to transform.")
	return nil
}
