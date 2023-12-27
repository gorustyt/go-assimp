package common

import "assimp/common/pb_msg"

type AiColor4D struct {
	R, G, B, A float32
}

func (ai *AiColor4D) Clone() *AiColor4D {
	tmp := *ai
	return &tmp
}
func (ai *AiColor4D) Empty() bool {
	return ai.R == 0 && ai.G == 0 && ai.B == 0 && ai.A == 0
}
func (ai *AiColor4D) ToPbMsg() *pb_msg.AiColor4D {
	return &pb_msg.AiColor4D{R: ai.R, G: ai.G, B: ai.B, A: ai.A}
}
func (ai *AiColor4D) FromPbMsg(data *pb_msg.AiColor4D) *AiColor4D {
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
	ai.R -= o.R
	ai.G -= o.G
	ai.B -= o.B
	ai.A -= o.A
	tmp := *ai
	return &tmp
}

func (ai *AiColor4D) Mul(v2 *AiColor4D) *AiColor4D {
	return NewAiColor4D(ai.R*v2.R, ai.G*v2.G, ai.B*v2.B, ai.A*v2.A)
}

// ------------------------------------------------------------------------------------------------

func (ai *AiColor4D) MulValue(f float32) *AiColor4D {
	ai.R *= f
	ai.G *= f
	ai.B *= f
	ai.A *= f

	tmp := *ai
	return &tmp
}

// ------------------------------------------------------------------------------------------------

func (ai *AiColor4D) DivValue(f float32) *AiColor4D {
	ai.R /= f
	ai.G /= f
	ai.B /= f
	ai.A /= f
	tmp := *ai
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
