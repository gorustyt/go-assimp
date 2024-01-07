package common

import "github.com/gorustyt/go-assimp/common/pb_msg"

type AiColor4D struct {
	R, G, B, A float32
}

func (ai *AiColor4D) Clone() *AiColor4D {
	if ai == nil {
		return nil
	}
	tmp := *ai
	return &tmp
}
func (ai *AiColor4D) BoundMin(a, b *AiColor4D) *AiColor4D {
	return NewAiColor4D(Min(a.R, b.R), Min(a.G, b.G), Min(a.B, b.B), Min(a.A, b.A))
}

func (ai *AiColor4D) BoundMax(a, b *AiColor4D) *AiColor4D {
	return NewAiColor4D(Max(a.R, b.R), Max(a.G, b.G), Max(a.B, b.B), Max(a.A, b.A))
}

func (ai *AiColor4D) ToPbMsg() *pb_msg.AiColor4D {
	if ai == nil {
		return nil
	}
	return &pb_msg.AiColor4D{R: ai.R, G: ai.G, B: ai.B, A: ai.A}
}
func (ai *AiColor4D) FromPbMsg(data *pb_msg.AiColor4D) *AiColor4D {
	if data == nil {
		return nil
	}
	ai.R = data.R
	ai.G = data.G
	ai.B = data.B
	ai.A = data.A
	return ai
}
func NewAiColor4D(R, G, B, A float32) *AiColor4D {
	return &AiColor4D{R: R, G: G, B: B, A: A}
}

func NewAiColor4D0() *AiColor4D {
	return &AiColor4D{}
}

// ------------------------------------------------------------------------------------------------

func (ai *AiColor4D) Add(o *AiColor4D) *AiColor4D {
	ai.R += o.R
	ai.G += o.G
	ai.B += o.B
	ai.A += o.A

	tmp := *ai
	return &tmp
}

// ------------------------------------------------------------------------------------------------

func (ai *AiColor4D) Sub(o *AiColor4D) *AiColor4D {
	tmp := *ai
	tmp.R -= o.R
	tmp.G -= o.G
	tmp.B -= o.B
	tmp.A -= o.A
	return &tmp
}

func (ai *AiColor4D) Mul(v2 *AiColor4D) *AiColor4D {
	return NewAiColor4D(ai.R*v2.R, ai.G*v2.G, ai.B*v2.B, ai.A*v2.A)
}

// ------------------------------------------------------------------------------------------------

func (ai *AiColor4D) MulValue(f float32) *AiColor4D {
	tmp := *ai
	tmp.R *= f
	tmp.G *= f
	tmp.B *= f
	tmp.A *= f
	return &tmp
}

// ------------------------------------------------------------------------------------------------

func (ai *AiColor4D) DivValue(f float32) *AiColor4D {
	tmp := *ai
	tmp.R /= f
	tmp.G /= f
	tmp.B /= f
	tmp.A /= f

	return &tmp
}

// ------------------------------------------------------------------------------------------------
func (ai *AiColor4D) Index(i int) float32 {
	switch i {
	case 0:
		return ai.R
	case 1:
		return ai.G
	case 2:
		return ai.B
	case 3:
		return ai.A
	default:
		break
	}
	return ai.R
}

// ------------------------------------------------------------------------------------------------

func (ai *AiColor4D) Equal(other *AiColor4D) bool {
	return ai.R == other.R && ai.G == other.G && ai.B == other.B && ai.A == other.A
}

// ------------------------------------------------------------------------------------------------
func (ai *AiColor4D) NotEqual(other *AiColor4D) bool {
	return ai.R != other.R || ai.G != other.G || ai.B != other.B || ai.A != other.A
}
