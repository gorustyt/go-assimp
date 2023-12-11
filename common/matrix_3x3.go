package common

import "math"

type AiMatrix3x3 struct {
	A1, A2, A3 float32
	B1, B2, B3 float32
	C1, C2, C3 float32
}

func NewAiMatrix3x3() *AiMatrix3x3 {
	m := &AiMatrix3x3{
		A1: 1,
		B2: 1,
		C3: 1,
	}
	return m
}

func NewAiMatrix3x3WithValues(
	_a1, _a2, _a3,
	_b1, _b2, _b3,
	_c1, _c2, _c3 float32) *AiMatrix3x3 {
	m := &AiMatrix3x3{}
	m.A1 = _a1
	m.A2 = _a2
	m.A3 = _a3
	m.B1 = _b1
	m.B2 = _b2
	m.B3 = _b3
	m.C1 = _c1
	m.C2 = _c2
	m.C3 = _c3
	return m
}

func (ai *AiMatrix3x3) Determinant() float32 {
	return ai.A1*ai.B2*ai.C3 - ai.A1*ai.B3*ai.C2 + ai.A2*ai.B3*ai.C1 - ai.A2*ai.B1*ai.C3 + ai.A3*ai.B1*ai.C2 - ai.A3*ai.B2*ai.C1
}

func (ai *AiMatrix3x3) Inverse() *AiMatrix3x3 {
	// Compute the reciprocal determinant
	det := ai.Determinant()
	if det == 0.0 {
		// Matrix not invertible. Setting all elements to nan is not really
		// correct in a mathematical sense; but at least qnans are easy to
		// spot. XXX we might throw an exception instead, which would
		// be even much better to spot :/.
		nan := float32(math.NaN())
		return NewAiMatrix3x3WithValues(nan, nan, nan, nan, nan, nan, nan, nan, nan)
	}

	invdet := 1.0 / det

	ai.A1 = invdet * (ai.B2*ai.C3 - ai.B3*ai.C2)
	ai.A2 = -invdet * (ai.A2*ai.C3 - ai.A3*ai.C2)
	ai.A3 = invdet * (ai.A2*ai.B3 - ai.A3*ai.B2)
	ai.B1 = -invdet * (ai.B1*ai.C3 - ai.B3*ai.C1)
	ai.B2 = invdet * (ai.A1*ai.C3 - ai.A3*ai.C1)
	ai.B3 = -invdet * (ai.A1*ai.B3 - ai.A3*ai.B1)
	ai.C1 = invdet * (ai.B1*ai.C2 - ai.B2*ai.C1)
	ai.C2 = -invdet * (ai.A1*ai.C2 - ai.A2*ai.C1)
	ai.C3 = invdet * (ai.A1*ai.B2 - ai.A2*ai.B1)
	return ai
}

func (ai *AiMatrix3x3) Index(p_iIndex int) *float32 {
	switch p_iIndex {
	case 0:
		return &ai.A1
	case 1:
		return &ai.B1
	case 2:
		return &ai.C1
	default:
		break
	}
	return &ai.A1
}
