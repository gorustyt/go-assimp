package common

import (
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"math"
)

type AiVector2D struct {
	X, Y float32
}

func (ai *AiVector2D) SquareLength() float32 {
	return ai.X*ai.X + ai.Y*ai.Y
}
func (ai *AiVector2D) Length() float32 {
	return float32(math.Sqrt(float64(ai.SquareLength())))
}
func (ai *AiVector2D) Div(f float32) *AiVector2D {
	return NewAiVector2D(ai.X/f, ai.Y/f)
}
func (ai *AiVector2D) Mul(v2 *AiVector2D) float32 {
	return ai.X*v2.X + ai.Y*v2.Y
}

func (ai *AiVector2D) Normalize() *AiVector2D {
	return ai.Div(ai.Length())
}

func (ai *AiVector2D) Sub(v2 *AiVector2D) *AiVector2D {
	return NewAiVector2D(ai.X-v2.X, ai.Y-v2.Y)
}
func (ai *AiVector2D) BoundMin(b *AiVector2D) *AiVector2D {
	return NewAiVector2D(Min(ai.X, b.X), Min(ai.Y, b.Y))
}

func (ai *AiVector2D) BoundMax(b *AiVector2D) *AiVector2D {
	return NewAiVector2D(Max(ai.X, b.X), Max(ai.Y, b.Y))
}
func (ai *AiVector2D) Clone() *AiVector2D {
	if ai == nil {
		return nil
	}
	return &AiVector2D{X: ai.X, Y: ai.Y}
}
func (ai *AiVector2D) ToPbMsg() *pb_msg.AiVector2D {
	if ai == nil {
		return nil
	}
	return &pb_msg.AiVector2D{X: ai.X, Y: ai.Y}
}
func (ai *AiVector2D) FromPbMsg(data *pb_msg.AiVector2D) *AiVector2D {
	if data == nil {
		return nil
	}
	ai.X = data.X
	ai.Y = data.Y
	return ai
}

func (ai *AiVector2D) Set(pX, pY float32) {
	ai.X = pX
	ai.Y = pY
}

func NewAiVector2D(X, Y float32) *AiVector2D {
	return &AiVector2D{X: X, Y: Y}
}
