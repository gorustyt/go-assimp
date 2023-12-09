package test

import "testing"

func Assert(t *testing.T, ok bool, msg ...string) {
	if !ok {
		t.Fatalf("test not ok :%v", msg)
	}
}
