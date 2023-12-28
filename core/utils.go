package core

import (
	"assimp/common"
)

// -------------------------------------------------------------------------------
func ComputePositionEpsilon(pMeshes []*AiMesh) float64 {
	epsilon := 1e-4
	// calculate the position bounds so we have a reliable epsilon to check position differences against
	var mi, ma *common.AiVector3D
	minVec, maxVec := ChooseMinMaxAiVector3D()
	for a := 0; a < len(pMeshes); a++ {
		pMesh := pMeshes[a]
		tmi, tma := ArrayBounds(pMesh.Vertices, func(d1 *common.AiVector3D, d2 *common.AiVector3D) bool {
			return d1.Less(d2)
		})
		mi = tmi.(*common.AiVector3D)
		ma = tma.(*common.AiVector3D)
		if !minVec.Less(mi) {
			minVec = mi
		}
		if maxVec.Less(ma) {
			maxVec = ma
		}
	}
	return maxVec.Sub(minVec).Length() * epsilon
}

func ArrayBounds[T any](in []T, less func(T, T) bool) (min, max any) {
	var tmp T
	var t any = tmp
	switch t.(type) {
	case float32:
		min, max = ChooseMinMaxFloat()
	case float64:
		min, max = ChooseMinMaxDouble()
	case uint32:
		min, max = ChooseMinMaxUint32()
	case *common.AiVector3D:
		min, max = ChooseMinMaxAiVector3D()
	case *common.AiVector2D:
		min, max = ChooseMinMaxAiVector2D()
	case *common.AiColor4D:
		min, max = ChooseMinMaxAiColor4D()
	case *common.AiQuaternion:
		min, max = ChooseMinMaxAiQuaternion()
	case *AiVectorKey:
		min, max = ChooseMinMaxAiVectorKey()
	case *AiQuatKey:
		tmin, tmax := &AiQuatKey{}, &AiQuatKey{}
		ChooseMinMaxAiQuatKey(tmin, tmax)
		min, max = tmin, tmax
	case *AiVertexWeight:
		tmin, tmax := &AiVertexWeight{}, &AiVertexWeight{}
		ChooseMinMaxAiVertexWeight(tmin, tmax)
		min, max = tmin, tmax
	}
	for i := 0; i < len(in); i++ {
		if less(in[i], min.(T)) {
			min = in[i]
		}
		if !less(in[i], min.(T)) {
			max = in[i]
		}
	}
	return
}

// -------------------------------------------------------------------------------
// Start points for ArrayBounds<T> for all supported Ts
func ChooseMinMaxFloat() (min, max float32) {
	max = -1e10
	min = 1e10
	return
}

func ChooseMinMaxDouble() (min, max float64) {
	max = -1e10
	min = 1e10
	return
}

func ChooseMinMaxUint32() (min, max uint32) {
	max = 0
	min = (1<<4*8 - 1)
	return
}

func ChooseMinMaxAiVector3D() (min, max *common.AiVector3D) {
	max = common.NewAiVector3D3(-1e10, -1e10, -1e10)
	min = common.NewAiVector3D3(1e10, 1e10, 1e10)
	return
}

func ChooseMinMaxAiVector2D() (min, max *common.AiVector2D) {
	max = common.NewAiVector2D(-1e10, -1e10)
	min = common.NewAiVector2D(1e10, 1e10)
	return
}

func ChooseMinMaxAiColor4D() (min, max *common.AiColor4D) {
	max = common.NewAiColor4D(-1e10, -1e10, -1e10, -1e10)
	min = common.NewAiColor4D(1e10, 1e10, 1e10, 1e10)
	return
}

func ChooseMinMaxAiQuaternion() (min, max *common.AiQuaternion) {
	max = common.NewAiQuaternion(-1e10, -1e10, -1e10, -1e10)
	min = common.NewAiQuaternion(1e10, 1e10, 1e10, 1e10)
	return
}

func ChooseMinMaxAiVectorKey() (min, max *AiVectorKey) {
	min.Time, max.Time = ChooseMinMaxDouble()
	min.Value, max.Value = ChooseMinMaxAiVector3D()
	return
}

func ChooseMinMaxAiQuatKey(min, max *AiQuatKey) {
	min.Time, max.Time = ChooseMinMaxDouble()
	min.Value, max.Value = ChooseMinMaxAiQuaternion()
}

func ChooseMinMaxAiVertexWeight(min, max *AiVertexWeight) {
	min.VertexId, max.VertexId = ChooseMinMaxUint32()
	min.Weight, max.Weight = ChooseMinMaxFloat()
	return
}
