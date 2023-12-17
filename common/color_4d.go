package common

import "assimp/common/pb_msg"

type AiColor4D struct {
	R, G, B, A float32
}

func (ai AiColor4D) Empty() bool {
	return ai.R == 0 && ai.G == 0 && ai.B == 0 && ai.A == 0
}
func (ai AiColor4D) ToPbMsg() *pb_msg.AiColor4D {
	return &pb_msg.AiColor4D{R: ai.R, G: ai.G, B: ai.B, A: ai.A}
}

func NewAiColor4D(R, G, B, A float32) *AiColor4D {
	return &AiColor4D{R: R, G: G, B: B, A: A}
}

func NewAiColor4D0() *AiColor4D {
	return &AiColor4D{}
}
