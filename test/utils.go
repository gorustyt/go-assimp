package test

import (
	"assimp/common/logger"
	"assimp/core"
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

func DeepEqual(p1, p2 *core.AiScene) bool {
	for _, mesh := range p1.Meshes {
		mesh.Name = "" //assetBin 沒有這個字段
	}
	for _, mesh := range p1.Materials {
		for _, mesh1 := range p2.Materials {
			for i, v := range mesh.Properties {
				v1 := mesh1.Properties[i]
				if !deepEqual(v, v1) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!")
				}
			}
		}
	}
	return deepEqual(p1, p2)
}
