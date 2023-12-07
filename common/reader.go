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

func (r *AiReader) MustOneKeyFloat64(key string) (float64, error) {
	if !r.checkIndexEof(1) {
		return 0, ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return 0, ErrBadKey
	}
	res, err := strconv.ParseFloat(r.curLines[r.curIndex+1], 64)
	if err != nil {
		return 0, err
	}
	r.NextLine()
	return res, nil
}

func (r *AiReader) checkIndexEof(index int) bool {
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
	if !r.checkIndexEof(1) {
		return 0, ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return 0, ErrBadKey
	}
	res, err := strconv.Atoi(r.curLines[r.curIndex+1])
	if err != nil {
		return 0, err
	}
	r.NextLine()
	return res, nil
}
func (r *AiReader) MustOneKeyString(key string) (string, error) {
	if !r.checkIndexEof(1) {
		return "", ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return "", ErrBadKey
	}
	res := r.curLines[r.curIndex+1]
	r.NextLine()
	return res, nil
}

func (r *AiReader) NextKeyFloat64(key string, index int) (res []float64, err error) {
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

func (r *AiReader) nextKey(key string, index int) (values []string, err error) {
	if r.checkIndexValid(index + 1) {
		return values, ErrBadParams
	}
	if r.curLines[r.curIndex] != key {
		return values, ErrBadKey
	}
	values = r.curLines[r.curIndex : r.curIndex+index]
	r.curIndex += index + 1
	return values, err
}

func (r *AiReader) EOF() bool {
	return r.Eof || len(r.curLines) == r.curIndex
}
