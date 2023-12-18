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

func (ai *AiMatrix3x3) Equal(m *AiMatrix3x3) bool {
	return ai.A1 == m.A1 && ai.A2 == m.A2 && ai.A3 == m.A3 &&
		ai.B1 == m.B1 && ai.B2 == m.B2 && ai.B3 == m.B3 &&
		ai.C1 == m.C1 && ai.C2 == m.C2 && ai.C3 == m.C3
}

func (ai *AiMatrix3x3) RotationZ(a float32, out *AiMatrix3x3) *AiMatrix3x3 {
	out.B2 = float32(math.Cos(float64(a)))
	out.A1 = out.B2
	out.B1 = float32(math.Sin(float64(a)))
	out.A2 = -out.B1

	out.C2 = 0.
	out.C1 = out.C2
	out.B3 = out.C1
	out.A3 = out.B3
	out.C3 = 1.

	return out
}

// ------------------------------------------------------------------------------------------------
// Returns a rotation matrix for a rotation around an arbitrary axis.

func (ai *AiMatrix3x3) Rotation(a float32, axis *AiVector3D, out AiMatrix3x3) AiMatrix3x3 {
	c := float32(math.Cos(float64(a)))
	s := float32(math.Sin(float64(a)))
	t := 1 - c
	x := axis.X
	y := axis.Y
	z := axis.Z

	// Many thanks to MathWorld and Wikipedia
	out.A1 = t*x*x + c
	out.A2 = t*x*y - s*z
	out.A3 = t*x*z + s*y
	out.B1 = t*x*y + s*z
	out.B2 = t*y*y + c
	out.B3 = t*y*z - s*x
	out.C1 = t*x*z - s*y
	out.C2 = t*y*z + s*x
	out.C3 = t*z*z + c

	return out
}

func (ai *AiMatrix3x3) Translation(v *AiVector2D) *AiMatrix3x3 {
	out := NewAiMatrix3x3()
	out.A3 = v.X
	out.B3 = v.Y
	return out
}

// ------------------------------------------------------------------------------------------------
/** Transformation of a vector by a 3x3 matrix */

func (ai *AiMatrix3x3) MulVector3d(pVector *AiVector3D) *AiVector3D {
	var res = NewAiVector3D()
	res.X = ai.A1*pVector.X + ai.A2*pVector.Y + ai.A3*pVector.Z
	res.Y = ai.B1*pVector.X + ai.B2*pVector.Y + ai.B3*pVector.Z
	res.Z = ai.C1*pVector.X + ai.C2*pVector.Y + ai.C3*pVector.Z
	return res
}

func (ai *AiMatrix3x3) EpsEqual(m *AiMatrix3x3, epsilon float64) bool {
	return math.Abs(float64(ai.A1-m.A1)) <= epsilon &&
		math.Abs(float64(ai.A2-m.A2)) <= epsilon &&
		math.Abs(float64(ai.A3-m.A3)) <= epsilon &&
		math.Abs(float64(ai.B1-m.B1)) <= epsilon &&
		math.Abs(float64(ai.B2-m.B2)) <= epsilon &&
		math.Abs(float64(ai.B3-m.B3)) <= epsilon &&
		math.Abs(float64(ai.C1-m.C1)) <= epsilon &&
		math.Abs(float64(ai.C2-m.C2)) <= epsilon &&
		math.Abs(float64(ai.C3-m.C3)) <= epsilon
}

func (ai *AiMatrix3x3) NotEqual(m *AiMatrix3x3) bool {
	return !ai.Equal(m)
}

func (ai *AiMatrix3x3) MulMatrix3x3(m *AiMatrix3x3) *AiMatrix3x3 {
	return NewAiMatrix3x3WithValues(
		m.A1*ai.A1+m.B1*ai.A2+m.C1*ai.A3,
		m.A2*ai.A1+m.B2*ai.A2+m.C2*ai.A3,
		m.A3*ai.A1+m.B3*ai.A2+m.C3*ai.A3,
		m.A1*ai.B1+m.B1*ai.B2+m.C1*ai.B3,
		m.A2*ai.B1+m.B2*ai.B2+m.C2*ai.B3,
		m.A3*ai.B1+m.B3*ai.B2+m.C3*ai.B3,
		m.A1*ai.C1+m.B1*ai.C2+m.C1*ai.C3,
		m.A2*ai.C1+m.B2*ai.C2+m.C2*ai.C3,
		m.A3*ai.C1+m.B3*ai.C2+m.C3*ai.C3)
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

func (ai *AiMatrix3x3) SetByIndex(i, j int, value float32) {
	switch i {
	case 0:
		switch j {
		case 0:
			ai.A1 = value
		case 1:
			ai.A2 = value
		case 2:
			ai.A3 = value
		}

	case 1:
		switch j {
		case 0:
			ai.B1 = value
		case 1:
			ai.B2 = value
		case 2:
			ai.B3 = value
		}
	case 2:
		switch j {
		case 0:
			ai.C1 = value
		case 1:
			ai.C2 = value
		case 2:
			ai.C3 = value
		}
	default:
		break
	}

}

func (ai *AiMatrix3x3) Index(p_iIndex int, j int) float32 {
	switch p_iIndex {
	case 0:
		switch j {
		case 0:
			return ai.A1
		case 1:
			return ai.A2
		case 2:
			return ai.A3
		}

	case 1:
		switch j {
		case 0:
			return ai.B1
		case 1:
			return ai.B2
		case 2:
			return ai.B3
		}
	case 2:
		switch j {
		case 0:
			return ai.C1
		case 1:
			return ai.C2
		case 2:
			return ai.C3
		}
	default:
		break
	}
	return -1
}
