package common

import "assimp/common/pb_msg"

type AiColor3D struct {
	R, G, B float32
}

func (ai *AiColor3D) Clone() *AiColor3D {
	if ai == nil {
		return nil
	}
	return &AiColor3D{R: ai.R, G: ai.G, B: ai.B}
}
func (ai *AiColor3D) ToPbMsg() *pb_msg.AiColor3D {
	if ai == nil {
		return nil
	}
	return &pb_msg.AiColor3D{R: ai.R, G: ai.G, B: ai.B}
}

func (ai *AiColor3D) FromPbMsg(data *pb_msg.AiColor3D) *AiColor3D {
	if data == nil {
		return nil
	}
	ai.R = data.R
	ai.G = data.G
	ai.B = data.B
	return ai
}

func NewAiColor3D(R, G, B float32) *AiColor3D {
	return &AiColor3D{R: R, G: G, B: B}
}

/** Component-wise multiplication */
func (ai *AiColor3D) MulColor3D(c *AiColor3D) *AiColor3D {
	return NewAiColor3D(ai.R*c.R, ai.G*c.G, ai.B*c.B)
}

/** Multiply with a scalar */
func (ai *AiColor3D) Mul(f float32) *AiColor3D {
	return NewAiColor3D(ai.R*f, ai.G*f, ai.B*f)
}
