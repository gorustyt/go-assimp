package test

import (
	"assimp/common"
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
