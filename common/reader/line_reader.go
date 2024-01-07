package reader

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/logger"
	"io"
	"strconv"
	"strings"
)

type LineReader interface {
	NextLine()
	GetLine() string
	GetLineNum() int
	NextLineVector3(verticesKey string) (vertices []*common.AiVector3D, err error)
	ReadLineAiVector3d() (res *common.AiVector3D, err error)
	HasPrefix(prefix string) bool
	NextKeyString(key string, index int) (res []string, err error)
	NextOneKeyInt(key string) (res int, err error)
	NextOneKeyFloat32(key string) (res float32, err error)
	NextKeyAiColor3d(key string) (res *common.AiColor3D, err error)
	NextKeyAiMatrix3x3(key string) (res *common.AiMatrix3x3, err error)
	EOF() bool

	MustOneKeyString(key string) (string, error)
	MustOneKeyInt(key string, options ...bool) (int, error)
	NextKeyAiVector2d(key string) (res *common.AiVector2D, err error)
	MustOneKeyFloat32(key string) (float32, error)
	NextKeyAiVector3d(key string) (res *common.AiVector3D, err error)
	NextOneKeyString(key string) (res string, err error)
}

func newLineReader(r *BaseReader) LineReader {
	return &lineReader{Reader: bufio.NewReader(bytes.NewReader(r.data)), BaseReader: r}
}

type lineReader struct {
	*BaseReader
	*bufio.Reader
	lineNum  int
	curLine  string
	curLines []string
	curIndex int
	eof      bool
	err      error
}

func (r *lineReader) NextLine() {
	r.curLine = ""
	r.curLines = r.curLines[:0]
	r.curIndex = 0
ReadLine:
	line, isPrefix, err := r.ReadLine()
	r.curLine += string(line)
	if err == io.EOF {
		r.eof = true
	} else if err != nil {
		r.err = err
	} else if isPrefix {
		goto ReadLine
	}
	r.SplitLine()
	r.lineNum++

}
func (r *lineReader) EOF() bool {
	return r.eof
}

func (r *lineReader) SplitLine() {
	r.curLines = strings.Split(r.curLine, " ")
	var ss []string
	for _, v := range r.curLines {
		if v == "" {
			logger.Info("skip space .")
			continue
		}
		ss = append(ss, v)
	}
	r.curLines = ss
}
func (r *lineReader) GetLine() string {
	return r.curLine
}

func (r *lineReader) GetLineNum() int {
	return r.lineNum
}

func (r *lineReader) MustOneKeyFloat32(key string) (float32, error) {
	values, err := r.nextKey(key, 1)
	if err != nil {
		return 0, err
	}
	if !r.IsEndLine(1) {
		return 0, ErrBadParams
	}
	res, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return 0, err
	}
	r.NextLine()
	return float32(res), nil
}

func (r *lineReader) IsEndLine(index int) bool {
	if r.curIndex+index >= len(r.curLines) {
		return true
	}
	return false
}

func (r *lineReader) checkIndexValid(index int) bool {
	if r.curIndex+index > len(r.curLines) {
		return false
	}
	return true
}

func (r *lineReader) MustOneKeyInt(key string, options ...bool) (res int, err error) {
	values, err := r.nextKey(key, 1)
	if err != nil {
		return 0, err
	}
	if !r.IsEndLine(1) {
		return 0, ErrBadParams
	}
	if len(options) > 0 && options[0] {

		o, err := hex.DecodeString(strings.TrimPrefix(values[0], "0x"))
		if err != nil {
			return 0, err
		}
		res = int(o[0])
	} else {
		res, err = strconv.Atoi(values[0])
	}
	if err != nil {
		return 0, err
	}
	r.NextLine()
	return int(res), nil
}
func (r *lineReader) MustOneKeyString(key string) (string, error) {
	values, err := r.nextKey(key, 1)
	if err != nil {
		return "", err
	}
	if !r.IsEndLine(1) {
		return "", ErrBadParams
	}
	r.NextLine()
	return values[0], nil
}
func (r *lineReader) NextKeyAiColor3d(key string) (res *common.AiColor3D, err error) {
	data, err := r.NextKeyFloat32(key, 3)
	if err != nil {
		return res, nil
	}
	return &common.AiColor3D{data[0], data[1], data[2]}, err
}

func (r *lineReader) NextKeyAiMatrix3x3(key string) (res *common.AiMatrix3x3, err error) {
	res = &common.AiMatrix3x3{}
	data, err := r.NextKeyFloat32(key, 9)
	if err != nil {
		return res, nil
	}
	res.A1 = data[0]
	res.A2 = data[1]
	res.A3 = data[3]
	res.B1 = data[4]
	res.B2 = data[5]
	res.B3 = data[6]
	res.C1 = data[1]
	res.C2 = data[2]
	res.C3 = data[3]
	return res, err
}

func (r *lineReader) NextKeyAiVector2d(key string) (res *common.AiVector2D, err error) {
	data, err := r.NextKeyFloat32(key, 2)
	if err != nil {
		return res, nil
	}
	return &common.AiVector2D{data[0], data[1]}, err
}

func (r *lineReader) NextKeyAiVector3d(key string) (res *common.AiVector3D, err error) {
	data, err := r.NextKeyFloat32(key, 3)
	if err != nil {
		return res, nil
	}
	return &common.AiVector3D{data[0], data[1], data[2]}, err
}

func (r *lineReader) NextOneKeyFloat32(key string) (res float32, err error) {
	values, err := r.NextKeyFloat32(key, 1)
	if err != nil {
		return res, err
	}
	return values[0], nil
}

func (r *lineReader) NextOneKeyInt(key string) (res int, err error) {
	values, err := r.NextKeyInt(key, 1)
	if err != nil {
		return res, err
	}
	return values[0], nil
}

func (r *lineReader) NextOneKeyString(key string) (res string, err error) {
	values, err := r.NextKeyString(key, 1)
	if err != nil {
		return res, err
	}
	return values[0], nil
}

func (r *lineReader) NextKeyFloat32(key string, index int) (res []float32, err error) {
	values, err := r.nextKey(key, index)
	if err != nil {
		return res, err
	}
	for _, v := range values {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return res, err
		}
		res = append(res, float32(value))
	}
	return res, err
}

func (r *lineReader) NextKeyInt(key string, index int) (res []int, err error) {
	values, err := r.nextKey(key, index)
	if err != nil {
		return res, err
	}
	for _, v := range values {
		value, err := strconv.Atoi(v)
		if err != nil {
			return res, err
		}
		res = append(res, value)
	}
	return res, err
}

func (r *lineReader) NextKeyString(key string, index int) (res []string, err error) {
	return r.nextKey(key, index)
}

func (r *lineReader) HasPrefix(prefix string) bool {
	return strings.HasPrefix(r.curLine, prefix)
}

func (r *lineReader) nextKey(key string, index int) (values []string, err error) {
	if r.eof {
		return values, io.EOF
	}
	if r.err != nil {
		return values, r.err
	}
	if !r.checkIndexValid(index + 1) {
		return values, ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return values, fmt.Errorf("ErrBadKey:%v", key)
	}
	values = r.curLines[r.curIndex+1 : r.curIndex+index+1]
	r.curIndex += index + 1
	return values, err
}

func (r *lineReader) ReadLineFloat32() (res []float32, err error) {
	for _, v := range r.curLines {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return res, err
		}
		res = append(res, float32(value))
	}
	return res, err
}

func (r *lineReader) ReadLineAiVector3d() (res *common.AiVector3D, err error) {
	values, err := r.ReadLineFloat32()
	if err != nil {
		return res, err
	}
	if len(values) < 3 {
		return res, fmt.Errorf("expect find element 3 but found:%v %v", len(values), values)
	}
	if len(values) > 3 {
		logger.WarnF("expect find element 3 but found:%v %v", len(values), values)
	}

	return &common.AiVector3D{values[0], values[1], values[2]}, err
}

func (r *lineReader) NextLineVector3(verticesKey string) (vertices []*common.AiVector3D, err error) {
	if !r.HasPrefix(verticesKey) {
		return vertices, ErrBadParams
	}
	num, err := r.MustOneKeyInt(verticesKey)
	if err != nil {
		return nil, err
	}
	for i := 0; i < num; i++ {
		if i == 131 {
			_ = num
		}
		v, err := r.ReadLineAiVector3d()
		if err != nil {
			return nil, err
		}
		vertices = append(vertices, v)
		if i != num-1 {
			r.NextLine()
		}
	}
	if len(vertices) != num {
		return vertices, ErrBadParams
	}
	return vertices, nil
}
