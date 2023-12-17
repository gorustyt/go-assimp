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
func (ai *AiMatrix4x4) AddMatrix4x4(m *AiMatrix4x4) *AiMatrix4x4 {
	return NewAiMatrix4x4FromValues(
		m.A1+ai.A1,
		m.A2+ai.A2,
		m.A3+ai.A3,
		m.A4+ai.A4,
		m.B1+ai.B1,
		m.B2+ai.B2,
		m.B3+ai.B3,
		m.B4+ai.B4,
		m.C1+ai.C1,
		m.C2+ai.C2,
		m.C3+ai.C3,
		m.C4+ai.C4,
		m.D1+ai.D1,
		m.D2+ai.D2,
		m.D3+ai.D3,
		m.D4+ai.D4,
	)
}
func (ai *AiMatrix4x4) MulFloat32(aFloat float32) *AiMatrix4x4 {
	return NewAiMatrix4x4FromValues(
		ai.A1*aFloat,
		ai.A2*aFloat,
		ai.A3*aFloat,
		ai.A4*aFloat,
		ai.B1*aFloat,
		ai.B2*aFloat,
		ai.B3*aFloat,
		ai.B4*aFloat,
		ai.C1*aFloat,
		ai.C2*aFloat,
		ai.C3*aFloat,
		ai.C4*aFloat,
		ai.D1*aFloat,
		ai.D2*aFloat,
		ai.D3*aFloat,
		ai.D4*aFloat,
	)
}
func (ai *AiMatrix4x4) MulAiMatrix4x4(m *AiMatrix4x4) *AiMatrix4x4 {
	res := NewAiMatrix4x4FromValues(
		m.A1*ai.A1+m.B1*ai.A2+m.C1*ai.A3+m.D1*ai.A4,
		m.A2*ai.A1+m.B2*ai.A2+m.C2*ai.A3+m.D2*ai.A4,
		m.A3*ai.A1+m.B3*ai.A2+m.C3*ai.A3+m.D3*ai.A4,
		m.A4*ai.A1+m.B4*ai.A2+m.C4*ai.A3+m.D4*ai.A4,
		m.A1*ai.B1+m.B1*ai.B2+m.C1*ai.B3+m.D1*ai.B4,
		m.A2*ai.B1+m.B2*ai.B2+m.C2*ai.B3+m.D2*ai.B4,
		m.A3*ai.B1+m.B3*ai.B2+m.C3*ai.B3+m.D3*ai.B4,
		m.A4*ai.B1+m.B4*ai.B2+m.C4*ai.B3+m.D4*ai.B4,
		m.A1*ai.C1+m.B1*ai.C2+m.C1*ai.C3+m.D1*ai.C4,
		m.A2*ai.C1+m.B2*ai.C2+m.C2*ai.C3+m.D2*ai.C4,
		m.A3*ai.C1+m.B3*ai.C2+m.C3*ai.C3+m.D3*ai.C4,
		m.A4*ai.C1+m.B4*ai.C2+m.C4*ai.C3+m.D4*ai.C4,
		m.A1*ai.D1+m.B1*ai.D2+m.C1*ai.D3+m.D1*ai.D4,
		m.A2*ai.D1+m.B2*ai.D2+m.C2*ai.D3+m.D2*ai.D4,
		m.A3*ai.D1+m.B3*ai.D2+m.C3*ai.D3+m.D3*ai.D4,
		m.A4*ai.D1+m.B4*ai.D2+m.C4*ai.D3+m.D4*ai.D4)
	return res
}
func (ai *AiMatrix4x4) NotEqual(m *AiMatrix4x4, tepsilon float32) bool {
	return ai.Equal(m, tepsilon)
}
func (ai *AiMatrix4x4) Equal(m *AiMatrix4x4, tepsilon float32) bool {
	epsilon := float64(tepsilon)
	return math.Abs(float64(ai.A1-m.A1)) <= epsilon && math.Abs(float64(ai.A2-m.A2)) <= epsilon &&
		math.Abs(float64(ai.A3-m.A3)) <= epsilon &&
		math.Abs(float64(ai.A4-m.A4)) <= epsilon &&
		math.Abs(float64(ai.B1-m.B1)) <= epsilon &&
		math.Abs(float64(ai.B2-m.B2)) <= epsilon &&
		math.Abs(float64(ai.B3-m.B3)) <= epsilon &&
		math.Abs(float64(ai.B4-m.B4)) <= epsilon &&
		math.Abs(float64(ai.C1-m.C1)) <= epsilon &&
		math.Abs(float64(ai.C2-m.C2)) <= epsilon &&
		math.Abs(float64(ai.C3-m.C3)) <= epsilon &&
		math.Abs(float64(ai.C4-m.C4)) <= epsilon &&
		math.Abs(float64(ai.D1-m.D1)) <= epsilon &&
		math.Abs(float64(ai.D2-m.D2)) <= epsilon &&
		math.Abs(float64(ai.D3-m.D3)) <= epsilon &&
		math.Abs(float64(ai.D4-m.D4)) <= epsilon
}

func (ai *AiMatrix4x4) Index(p_iIndex int, j int) float32 {
	if p_iIndex > 3 {
		return 0
	}
	switch p_iIndex {
	case 0:
		switch j {
		case 0:
			return ai.A1
		case 1:
			return ai.A2
		case 2:
			return ai.A3
		case 3:
			return ai.A4
		}
	case 1:
		switch j {
		case 0:
			return ai.B1
		case 1:
			return ai.B2
		case 2:
			return ai.B3
		case 3:
			return ai.B4
		}
	case 2:
		switch j {
		case 0:
			return ai.C1
		case 1:
			return ai.C2
		case 2:
			return ai.C3
		case 3:
			return ai.C4
		}
	case 3:
		switch j {
		case 0:
			return ai.D1
		case 1:
			return ai.D2
		case 2:
			return ai.D3
		case 3:
			return ai.D4
		}
	default:
		break
	}
	return -1
}

func (ai *AiMatrix4x4) Set(p_iIndex int, j int, value float32) {
	if p_iIndex > 3 {
		return
	}
	switch p_iIndex {
	case 0:
		switch j {
		case 0:
			ai.A1 = value
		case 1:
			ai.A2 = value
		case 2:
			ai.A3 = value
		case 3:
			ai.A4 = value
		}
	case 1:
		switch j {
		case 0:
			ai.B1 = value
		case 1:
			ai.B2 = value
		case 2:
			ai.B3 = value
		case 3:
			ai.B4 = value
		}
	case 2:
		switch j {
		case 0:
			ai.C1 = value
		case 1:
			ai.C2 = value
		case 2:
			ai.C3 = value
		case 3:
			ai.C4 = value
		}
	case 3:
		switch j {
		case 0:
			ai.D1 = value
		case 1:
			ai.D2 = value
		case 2:
			ai.D3 = value
		case 3:
			ai.D4 = value
		}
	default:
		break
	}
	return
}
func (ai *AiMatrix4x4) FromEulerAnglesXYZ(x, y, z float32) *AiMatrix4x4 {
	_this := *ai
	cx := float32(math.Cos(float64(x)))
	sx := math.Sin(float64(x))
	cy := math.Cos(float64(y))
	sy := math.Sin(float64(y))
	cz := math.Cos(float64(z))
	sz := math.Sin(float64(z))

	// mz*my*mx
	_this.A1 = float32(cz) * float32(cy)
	_this.A2 = float32(cz)*float32(sy)*float32(sx) - float32(sz)*float32(cx)
	_this.A3 = float32(sz)*float32(sx) + float32(cz)*float32(sy)*float32(cx)

	_this.B1 = float32(sz) * float32(cy)
	_this.B2 = float32(cz)*float32(cx) + float32(sz)*float32(sy)*float32(sx)
	_this.B3 = float32(sz)*float32(sy)*float32(cx) - float32(cz)*float32(sx)

	_this.C1 = -float32(sy)
	_this.C2 = float32(cy) * float32(sx)
	_this.C3 = float32(cy) * float32(cx)

	return &_this
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
func (ai *AiMatrix4x4) Transpose() *AiMatrix4x4 {
	// (TReal&) don't remove, GCC complains cause of packed fields
	ai.B1, ai.A2 = ai.A2, ai.B1
	ai.C1, ai.A3 = ai.A3, ai.C1
	ai.C2, ai.B3 = ai.B3, ai.C2
	ai.D1, ai.A4 = ai.A4, ai.D1
	ai.D2, ai.B4 = ai.B4, ai.D2
	ai.D3, ai.C4 = ai.C4, ai.D3
	return ai
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
