package common

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	ErrBadParams = errors.New("error bad params")
	ErrBadKey    = errors.New("error bad key")
)

type AiReader struct {
	reader   bufio.Reader
	lineNum  int
	curLine  string
	curLines []string
	curIndex int
	Eof      bool
	err      error
}

func NewReader(reader bufio.Reader) *AiReader {
	r := &AiReader{reader: reader}
	r.NextLine()
	return r
}
func (r *AiReader) EOF() bool {
	return r.Eof
}
func (r *AiReader) NextLine() {
	r.curLine = ""
	r.curLines = r.curLines[:0]
	r.curIndex = 0
ReadLine:
	line, isPrefix, err := r.reader.ReadLine()
	r.curLine += string(line)
	if err == io.EOF {
		r.Eof = true
	} else if err != nil {
		r.err = err
	} else if isPrefix {
		goto ReadLine
	}
	r.curLines = strings.Split(r.curLine, " ")
	r.lineNum++

}

func (r *AiReader) GetLine() string {
	return r.curLine
}

func (r *AiReader) GetLineNum() int {
	return r.lineNum
}

func (r *AiReader) MustOneKeyfloat32(key string) (float32, error) {
	values, err := r.nextKey(key, 1)
	if err != nil {
		return 0, err
	}
	if !r.IsEndLine(1) {
		return 0, ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return 0, ErrBadKey
	}
	res, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return 0, err
	}
	r.NextLine()
	return res, nil
}

func (r *AiReader) IsEndLine(index int) bool {
	if r.curIndex+index == len(r.curLines) {
		return true
	}
	return false
}

func (r *AiReader) checkIndexValid(index int) bool {
	if r.curIndex+index > len(r.curLines) {
		return false
	}
	return true
}

func (r *AiReader) MustOneKeyInt(key string) (int, error) {
	values, err := r.nextKey(key, 1)
	if err != nil {
		return 0, err
	}
	if !r.IsEndLine(1) {
		return 0, ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return 0, ErrBadKey
	}
	res, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, err
	}
	r.NextLine()
	return res, nil
}
func (r *AiReader) MustOneKeyString(key string) (string, error) {
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
func (r *AiReader) NextKeyAiColor3d(key string) (res AiColor3D, err error) {
	data, err := r.NextKeyfloat32(key, 3)
	if err != nil {
		return res, nil
	}
	return AiColor3D{data[0], data[1], data[2]}, err
}

func (r *AiReader) NextKeyAiMatrix3x3(key string) (res AiMatrix3x3, err error) {
	data, err := r.NextKeyfloat32(key, 9)
	if err != nil {
		return res, nil
	}
	for i, v := range data {
		res[i] = v
	}
	return res, err
}

func (r *AiReader) NextKeyAiVector2d(key string) (res AiVector2D, err error) {
	data, err := r.NextKeyfloat32(key, 2)
	if err != nil {
		return res, nil
	}
	return AiVector2D{data[0], data[1]}, err
}

func (r *AiReader) NextKeyAiVector3d(key string) (res AiVector3D, err error) {
	data, err := r.NextKeyfloat32(key, 3)
	if err != nil {
		return res, nil
	}
	return AiVector3D{data[0], data[1], data[2]}, err
}

func (r *AiReader) NextOneKeyfloat32(key string) (res float32, err error) {
	values, err := r.NextKeyfloat32(key, 1)
	if err != nil {
		return res, err
	}
	return values[0], nil
}

func (r *AiReader) NextOneKeyInt(key string) (res int, err error) {
	values, err := r.NextKeyInt(key, 1)
	if err != nil {
		return res, err
	}
	return values[0], nil
}

func (r *AiReader) NextOneKeyString(key string) (res string, err error) {
	values, err := r.NextKeyString(key, 1)
	if err != nil {
		return res, err
	}
	return values[0], nil
}

func (r *AiReader) NextKeyfloat32(key string, index int) (res []float32, err error) {
	values, err := r.nextKey(key, index)
	if err != nil {
		return res, err
	}
	for _, v := range values {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return res, err
		}
		res = append(res, value)
	}
	return res, err
}

func (r *AiReader) NextKeyInt(key string, index int) (res []int, err error) {
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

func (r *AiReader) NextKeyString(key string, index int) (res []string, err error) {
	return r.nextKey(key, index)
}

func (r *AiReader) HasPrefix(prefix string) bool {
	return strings.HasPrefix(r.curLine, prefix)
}

func (r *AiReader) nextKey(key string, index int) (values []string, err error) {
	if r.Eof {
		return values, io.EOF
	}
	if r.err != nil {
		return values, r.err
	}
	if r.checkIndexValid(index + 1) {
		return values, ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return values, ErrBadKey
	}
	values = r.curLines[r.curIndex+1 : r.curIndex+index]
	r.curIndex += index + 1
	return values, err
}

func (r *AiReader) ReadLinefloat32() (res []float32, err error) {
	for _, v := range r.curLines {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return res, err
		}
		res = append(res, value)
	}
	return res, err
}

func (r *AiReader) ReadLineAiVector3d() (res AiVector3D, err error) {
	values, err := r.ReadLinefloat32()
	if err != nil {
		return res, err
	}
	if len(values) != 3 {
		return res, ErrBadParams
	}

	return AiVector3D{values[0], values[1], values[2]}, err
}

func (r *AiReader) NextLineVector3(verticesKey string) (vertices []AiVector3D, err error) {
	if !r.HasPrefix(verticesKey) {
		return vertices, ErrBadParams
	}
	num, err := r.MustOneKeyInt(verticesKey)
	if err != nil {
		return nil, err
	}
	for i := 0; i < num; i++ {
		v, err := r.ReadLineAiVector3d()
		if err != nil {
			return nil, err
		}
		vertices = append(vertices, v)
		r.NextLine()
	}
	if len(vertices) != num {
		return vertices, ErrBadParams
	}
	return vertices, nil
}
