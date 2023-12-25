package test

import (
	"math"
	"testing"
)

func AssertFloatEqual(t *testing.T, t1, t2 float64, eps float64) {
	Assert(t, math.Abs(t1-t2) > eps)
}
func Assert(t *testing.T, ok bool, msg ...string) {
	if !ok {
		t.Fatalf("test not ok :%v", msg)
	}
}

func AssertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("test error :%v", err)
	}
}
