package common

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
	min = 1<<4*8 - 1
	return
}

func ChooseMinMaxAiVector2D() (min, max *AiVector2D) {
	max = NewAiVector2D(-1e10, -1e10)
	min = NewAiVector2D(1e10, 1e10)
	return
}

func ChooseMinMaxAiColor4D() (min, max *AiColor4D) {
	max = NewAiColor4D(-1e10, -1e10, -1e10, -1e10)
	min = NewAiColor4D(1e10, 1e10, 1e10, 1e10)
	return
}

func ChooseMinMaxAiQuaternion() (min, max *AiQuaternion) {
	max = NewAiQuaternion(-1e10, -1e10, -1e10, -1e10)
	min = NewAiQuaternion(1e10, 1e10, 1e10, 1e10)
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

func ChooseMinMaxAiVector3D() (min, max *AiVector3D) {
	max = NewAiVector3D3(-1e10, -1e10, -1e10)
	min = NewAiVector3D3(1e10, 1e10, 1e10)
	return
}

func ArrayBounds[T any](in []T, boundMin, boundMax func(a T, b T) T) (min, max any) {
	var tmp T
	var t any = tmp
	switch t.(type) {
	case float32:
		min, max = ChooseMinMaxFloat()
	case float64:
		min, max = ChooseMinMaxDouble()
	case uint32:
		min, max = ChooseMinMaxUint32()
	case *AiVector3D:
		min, max = ChooseMinMaxAiVector3D()
	case *AiVector2D:
		min, max = ChooseMinMaxAiVector2D()
	case *AiColor4D:
		min, max = ChooseMinMaxAiColor4D()
	case *AiQuaternion:
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
		min = boundMin(in[i], min.(T))
		max = boundMax(in[i], max.(T))
	}
	return
}
