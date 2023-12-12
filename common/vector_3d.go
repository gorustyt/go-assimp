package common

import (
	"assimp/common/pb_msg"
	"math"
)

type AiVector3D struct {
	X, Y, Z float32
}

func (ai AiVector3D) Empty() bool {
	return ai.X == 0 && ai.Y == 0 && ai.Z == 0
}
func (ai AiVector3D) ToPbMsg() *pb_msg.AiVector3D {
	return &pb_msg.AiVector3D{X: ai.X, Y: ai.Y, Z: ai.Z}
}

func (ai AiVector3D) Add(o *AiVector3D) *AiVector3D {
	ai.X += o.X
	ai.Y += o.Y
	ai.Z += o.Z
	tmp := ai
	return &tmp
}

func (ai AiVector3D) Sub(o *AiVector3D) *AiVector3D {
	ai.X -= o.X
	ai.Y -= o.Y
	ai.Z -= o.Z
	tmp := ai
	return &tmp
}

func (ai AiVector3D) Mul(f float32) *AiVector3D {
	ai.X *= f
	ai.Y *= f
	ai.Z *= f
	tmp := ai
	return &tmp
}

func (ai AiVector3D) Div(f float32) *AiVector3D {
	if f == 0 {
		tmp := ai
		return &tmp
	}
	invF := float32(1.0) / f
	ai.X *= invF
	ai.Y *= invF
	ai.Z *= invF
	tmp := ai
	return &tmp
}

// / @brief  The default class constructor.
func NewAiVector3D() *AiVector3D {
	return &AiVector3D{}
}

// / @brief  The class constructor with the components.
// / @param  _x  The x-component for the vector.
// / @param  _y  The y-component for the vector.
// / @param  _z  The z-component for the vector.
func NewAiVector3D3(_x, _y, _z float32) *AiVector3D {
	t := NewAiVector3D()
	t.X = _x
	t.Y = _y
	t.Z = _z
	return t
}

// / @brief  The class constructor with a default value.
// / @param  _xyz  The value for x, y and z.
func NewAiVector3D1(_xyz float32) *AiVector3D {
	t := NewAiVector3D()
	t.X = _xyz
	t.Y = _xyz
	t.Z = _xyz
	return t
}

func (ai *AiVector3D) Set(pX, pY, pZ float32) {
	ai.X = pX
	ai.Y = pY
	ai.Z = pZ
}

func (ai *AiVector3D) SquareLength() float32 {
	return ai.X*ai.X + ai.Y*ai.Y + ai.Z*ai.Z
}

func (ai *AiVector3D) Length() float32 {
	return float32(math.Sqrt(float64(ai.SquareLength())))
}

func (ai *AiVector3D) Normalize() *AiVector3D {
	l := ai.Length()
	if l == 0 {
		return ai
	}
	ai.Div(ai.Length())
	return ai
}

func (ai *AiVector3D) MulMatrix3x3(mat *AiMatrix3x3) {

}
func (ai *AiVector3D) MulMatrix4x4(mat *AiMatrix4x4) {

}

func (ai *AiVector3D) Index(i int) float32 {
	switch i {
	case 0:
		return ai.X
	case 1:
		return ai.Y
	case 2:
		return ai.Z
	default:
		break
	}
	return ai.X
}

func (ai *AiVector3D) Equal(other *AiVector3D) bool {
	return ai.X == other.X && ai.Y == other.Y && ai.Z == other.Z
}

func (ai *AiVector3D) NotEqual(other *AiVector3D) bool {
	return ai.X != other.X || ai.Y != other.Y || ai.Z != other.Z
}

func (ai *AiVector3D) Equal1(other *AiVector3D, epsilon float32) {

}

func (ai *AiVector3D) Less(other *AiVector3D) bool {
	if ai.X != other.X {
		return ai.X < other.X
	} else {
		if ai.Y != other.Y {
			return ai.Y < other.Y
		} else {
			return ai.Z < other.Z
		}
	}
}

func NegationOperationSymbol(v1, v2 *AiVector3D) *AiVector3D {
	return NewAiVector3D3(v1.Y*v2.Z-v1.Z*v2.Y, v1.Z*v2.X-v1.X*v2.Z, v1.X*v2.Y-v1.Y*v2.X)
}

func MulAiVector3D(v1, v2 *AiVector3D) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// ------------------------------------------------------------------------------------------------
/** Transformation of a vector by a 4x4 matrix */

func Matrix4x4tMulAiVector3D(pMatrix *AiMatrix4x4, pVector *AiVector3D) *AiVector3D {
	var res AiVector3D
	res.X = pMatrix.A1*pVector.X + pMatrix.A2*pVector.Y + pMatrix.A3*pVector.Z + pMatrix.A4
	res.Y = pMatrix.B1*pVector.X + pMatrix.B2*pVector.Y + pMatrix.B3*pVector.Z + pMatrix.B4
	res.Z = pMatrix.C1*pVector.X + pMatrix.C2*pVector.Y + pMatrix.C3*pVector.Z + pMatrix.C4
	return &res
}
