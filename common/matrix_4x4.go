package common

import "math"

type AiMatrix4x4 struct {
	A1, A2, A3, A4 float32
	B1, B2, B3, B4 float32
	C1, C2, C3, C4 float32
	D1, D2, D3, D4 float32
}

func NewAiMatrix4x4FromValues(_a1, _a2, _a3, _a4, _b1, _b2, _b3, _b4,
	_c1, _c2, _c3, _c4,
	_d1, _d2, _d3, _d4 float32) *AiMatrix4x4 {
	ai := &AiMatrix4x4{}
	ai.A1 = _a1
	ai.A2 = _a2
	ai.A3 = _a3
	ai.A4 = _a4
	ai.B1 = _b1
	ai.B2 = _b2
	ai.B3 = _b3
	ai.B4 = _b4

	ai.C1 = _c1
	ai.C2 = _c2
	ai.C3 = _c3
	ai.C4 = _c4
	ai.D1 = _d1
	ai.D2 = _d2
	ai.D3 = _d3
	ai.D4 = _d4
	return ai
}

func NewAiMatrix4x4FromAiMatrix3x3(m *AiMatrix3x3) *AiMatrix4x4 {
	ai := AiMatrix4x4{}
	ai.A1 = m.A1
	ai.A2 = m.A2
	ai.A3 = m.A3
	ai.A4 = 0.0
	ai.B1 = m.B1
	ai.B2 = m.B2
	ai.B3 = m.B3
	ai.B4 = 0.0
	ai.C1 = m.C1
	ai.C2 = m.C2
	ai.C3 = m.C3
	ai.C4 = 0.0
	ai.D1 = 0.0
	ai.D2 = 0.0
	ai.D3 = 0.0
	ai.D4 = 1.0
	return &ai
}

func (ai *AiMatrix4x4) Determinant() float32 {
	return ai.A1*ai.B2*ai.C3*ai.D4 -
		ai.A1*ai.B2*ai.C4*ai.D3 +
		ai.A1*ai.B3*ai.C4*ai.D2 -
		ai.A1*ai.B3*ai.C2*ai.D4 + ai.A1*ai.B4*ai.C2*ai.D3 -
		ai.A1*ai.B4*ai.C3*ai.D2 -
		ai.A2*ai.B3*ai.C4*ai.D1 +
		ai.A2*ai.B3*ai.C1*ai.D4 - ai.A2*ai.B4*ai.C1*ai.D3 +
		ai.A2*ai.B4*ai.C3*ai.D1 -
		ai.A2*ai.B1*ai.C3*ai.D4 +
		ai.A2*ai.B1*ai.C4*ai.D3 + ai.A3*ai.B4*ai.C1*ai.D2 -
		ai.A3*ai.B4*ai.C2*ai.D1 + ai.A3*ai.B1*ai.C2*ai.D4 -
		ai.A3*ai.B1*ai.C4*ai.D2 + ai.A3*ai.B2*ai.C4*ai.D1 - ai.A3*ai.B2*ai.C1*ai.D4 -
		ai.A4*ai.B1*ai.C2*ai.D3 + ai.A4*ai.B1*ai.C3*ai.D2 - ai.A4*ai.B2*ai.C3*ai.D1 +
		ai.A4*ai.B2*ai.C1*ai.D3 - ai.A4*ai.B3*ai.C1*ai.D2 + ai.A4*ai.B3*ai.C2*ai.D1
}
func (ai *AiMatrix4x4) Inverse() *AiMatrix4x4 {
	// Compute the reciprocal determinant
	det := ai.Determinant()
	if det == 0.0 {
		// Matrix not invertible. Setting all elements to nan is not really
		// correct in a mathematical sense but it is easy to debug for the
		// programmer.
		nan := float32(math.NaN())
		return NewAiMatrix4x4FromValues(
			nan, nan, nan, nan,
			nan, nan, nan, nan,
			nan, nan, nan, nan,
			nan, nan, nan, nan)
	}

	invdet := 1.0 / det

	ai.A1 = invdet * (ai.B2*(ai.C3*ai.D4-ai.C4*ai.D3) + ai.B3*(ai.C4*ai.D2-ai.C2*ai.D4) + ai.B4*(ai.C2*ai.D3-ai.C3*ai.D2))
	ai.A2 = -invdet * (ai.A2*(ai.C3*ai.D4-ai.C4*ai.D3) + ai.A3*(ai.C4*ai.D2-ai.C2*ai.D4) + ai.A4*(ai.C2*ai.D3-ai.C3*ai.D2))
	ai.A3 = invdet * (ai.A2*(ai.B3*ai.D4-ai.B4*ai.D3) + ai.A3*(ai.B4*ai.D2-ai.B2*ai.D4) + ai.A4*(ai.B2*ai.D3-ai.B3*ai.D2))
	ai.A4 = -invdet * (ai.A2*(ai.B3*ai.C4-ai.B4*ai.C3) + ai.A3*(ai.B4*ai.C2-ai.B2*ai.C4) + ai.A4*(ai.B2*ai.C3-ai.B3*ai.C2))
	ai.B1 = -invdet * (ai.B1*(ai.C3*ai.D4-ai.C4*ai.D3) + ai.B3*(ai.C4*ai.D1-ai.C1*ai.D4) + ai.B4*(ai.C1*ai.D3-ai.C3*ai.D1))
	ai.B2 = invdet * (ai.A1*(ai.C3*ai.D4-ai.C4*ai.D3) + ai.A3*(ai.C4*ai.D1-ai.C1*ai.D4) + ai.A4*(ai.C1*ai.D3-ai.C3*ai.D1))
	ai.B3 = -invdet * (ai.A1*(ai.B3*ai.D4-ai.B4*ai.D3) + ai.A3*(ai.B4*ai.D1-ai.B1*ai.D4) + ai.A4*(ai.B1*ai.D3-ai.B3*ai.D1))
	ai.B4 = invdet * (ai.A1*(ai.B3*ai.C4-ai.B4*ai.C3) + ai.A3*(ai.B4*ai.C1-ai.B1*ai.C4) + ai.A4*(ai.B1*ai.C3-ai.B3*ai.C1))
	ai.C1 = invdet * (ai.B1*(ai.C2*ai.D4-ai.C4*ai.D2) + ai.B2*(ai.C4*ai.D1-ai.C1*ai.D4) + ai.B4*(ai.C1*ai.D2-ai.C2*ai.D1))
	ai.C2 = -invdet * (ai.A1*(ai.C2*ai.D4-ai.C4*ai.D2) + ai.A2*(ai.C4*ai.D1-ai.C1*ai.D4) + ai.A4*(ai.C1*ai.D2-ai.C2*ai.D1))
	ai.C3 = invdet * (ai.A1*(ai.B2*ai.D4-ai.B4*ai.D2) + ai.A2*(ai.B4*ai.D1-ai.B1*ai.D4) + ai.A4*(ai.B1*ai.D2-ai.B2*ai.D1))
	ai.C4 = -invdet * (ai.A1*(ai.B2*ai.C4-ai.B4*ai.C2) + ai.A2*(ai.B4*ai.C1-ai.B1*ai.C4) + ai.A4*(ai.B1*ai.C2-ai.B2*ai.C1))
	ai.D1 = -invdet * (ai.B1*(ai.C2*ai.D3-ai.C3*ai.D2) + ai.B2*(ai.C3*ai.D1-ai.C1*ai.D3) + ai.B3*(ai.C1*ai.D2-ai.C2*ai.D1))
	ai.D2 = invdet * (ai.A1*(ai.C2*ai.D3-ai.C3*ai.D2) + ai.A2*(ai.C3*ai.D1-ai.C1*ai.D3) + ai.A3*(ai.C1*ai.D2-ai.C2*ai.D1))
	ai.D3 = -invdet * (ai.A1*(ai.B2*ai.D3-ai.B3*ai.D2) + ai.A2*(ai.B3*ai.D1-ai.B1*ai.D3) + ai.A3*(ai.B1*ai.D2-ai.B2*ai.D1))
	ai.D4 = invdet * (ai.A1*(ai.B2*ai.C3-ai.B3*ai.C2) + ai.A2*(ai.B3*ai.C1-ai.B1*ai.C3) + ai.A3*(ai.B1*ai.C2-ai.B2*ai.C1))
	return ai
}
