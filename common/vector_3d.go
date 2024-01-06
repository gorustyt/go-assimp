package common

import (
	"assimp/common/pb_msg"
	"math"
)

type AiVector3D struct {
	X, Y, Z float32
}

func (ai *AiVector3D) BoundMin(b *AiVector3D) *AiVector3D {
	return NewAiVector3D3(Min(ai.X, b.X), Min(ai.Y, b.Y), Min(ai.Z, b.Z))
}

func (ai *AiVector3D) BoundMax(b *AiVector3D) *AiVector3D {
	return NewAiVector3D3(Max(ai.X, b.X), Max(ai.Y, b.Y), Max(ai.Z, b.Z))
}

func (ai *AiVector3D) FromPbMsg(data *pb_msg.AiVector3D) *AiVector3D {
	if data == nil {
		return nil
	}
	ai.X = data.X
	ai.Y = data.Y
	ai.Z = data.Z
	return ai
}
func (ai *AiVector3D) ToPbMsg() *pb_msg.AiVector3D {
	if ai == nil {
		return nil
	}
	return &pb_msg.AiVector3D{X: ai.X, Y: ai.Y, Z: ai.Z}
}

func (ai *AiVector3D) Clone() *AiVector3D {
	if ai == nil {
		return nil
	}
	tmp := *ai
	return &tmp
}

func (ai *AiVector3D) Add(o *AiVector3D) *AiVector3D {
	tmp := *ai
	tmp.X += o.X
	tmp.Y += o.Y
	tmp.Z += o.Z
	return &tmp
}

func (ai *AiVector3D) Sub(o *AiVector3D) *AiVector3D {
	tmp := *ai
	tmp.X -= o.X
	tmp.Y -= o.Y
	tmp.Z -= o.Z
	return &tmp
}

func (ai *AiVector3D) Mul(f float32) *AiVector3D {
	tmp := *ai
	tmp.X *= f
	tmp.Y *= f
	tmp.Z *= f
	return &tmp
}

func (ai *AiVector3D) Div(f float64) *AiVector3D {
	if f == 0 {
		tmp := *ai
		return &tmp
	}
	tmp := *ai
	invF := float32(1.0) / float32(f)
	tmp.X *= invF
	tmp.Y *= invF
	tmp.Z *= invF
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

func (ai *AiVector3D) SquareLength() float64 {
	return float64(ai.X*ai.X + ai.Y*ai.Y + ai.Z*ai.Z)
}

func (ai *AiVector3D) Length() float64 {
	return math.Sqrt(ai.SquareLength())
}

func (ai *AiVector3D) Normalize() *AiVector3D {
	return ai.Div(ai.Length())
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

func (ai *AiVector3D) NegationOperationSymbol(v2 *AiVector3D) *AiVector3D {
	return NewAiVector3D3(ai.Y*v2.Z-ai.Z*v2.Y, ai.Z*v2.X-ai.X*v2.Z, ai.X*v2.Y-ai.Y*v2.X)
}

func (ai *AiVector3D) MulAiVector3D(v2 *AiVector3D) float64 {
	return float64(ai.X*v2.X + ai.Y*v2.Y + ai.Z*v2.Z)
}

/** A time-value pair specifying a certain 3D vector for the given time. */
type AiVectorKey struct {
	/** The time of this key */
	Time float64

	/** The value of this key */
	Value *AiVector3D
}

func (ai *AiVectorKey) FromPbMsg(p *pb_msg.AiVectorKey) *AiVectorKey {
	if p == nil {
		return nil
	}
	ai.Time = p.Time
	ai.Value = (&AiVector3D{}).FromPbMsg(p.Value)
	return ai
}

func (ai *AiVectorKey) Clone() *AiVectorKey {
	if ai == nil {
		return nil
	}
	r := &AiVectorKey{
		Time:  ai.Time,
		Value: ai.Value.Clone(),
	}

	return r
}

func (ai *AiVectorKey) ToPbMsg() *pb_msg.AiVectorKey {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiVectorKey{
		Time:  ai.Time,
		Value: ai.Value.ToPbMsg(),
	}

	return r
}

func (ai *AiVectorKey) BoundMin(b *AiVectorKey) *AiVectorKey {
	return &AiVectorKey{
		Time:  Min(ai.Time, b.Time),
		Value: ai.Value.BoundMin(b.Value),
	}
}

func (ai *AiVectorKey) BoundMax(b *AiVectorKey) *AiVectorKey {
	return &AiVectorKey{
		Time:  Max(ai.Time, b.Time),
		Value: ai.Value.BoundMin(b.Value),
	}
}
