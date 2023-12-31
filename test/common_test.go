package test

import (
	"assimp/common"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"os"
	"testing"
)

func TestLowerBound(t *testing.T) {
	var arr = []int{1, 2, 3, 4, 5, 6, 7}
	for i, v := range arr {
		index := common.LowerBound(0, len(arr), func(index int) bool {
			return arr[index] < v
		})
		Assert(t, index == i)
	}

	for i, v := range arr {
		index := common.LowerBound(0, len(arr), func(index int) bool {
			return arr[index] < v+1
		})
		Assert(t, index == i+1)
	}
}

func TestDeepEqual(t *testing.T) {
	Assert(t, deepEqual(0.001, 0.0011111))
	Assert(t, !deepEqual(0.01, 0.0011111))
	Assert(t, deepEqual([]int{1, 2, 3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6}))
}

func TestUtf8(t *testing.T) {
	data, err := os.ReadFile("../example_data/AC/SphereWithLight_UTF16LE.ac")
	AssertError(t, err)

	data, _, err = transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder(), data)
	AssertError(t, err)
	err = os.WriteFile("./SphereWithLight_UTF16LE.ac", data, 0664)
	AssertError(t, err)
}