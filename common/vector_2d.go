package common

import "assimp/common/pb_msg"

type AiVector2D struct {
	X, Y float32
}

func (ai *AiVector2D) ToPbMsg() *pb_msg.AiVector2D {
	return &pb_msg.AiVector2D{X: ai.X, Y: ai.Y}
}
func (ai *AiVector2D) FromPbMsg(data *pb_msg.AiVector2D) *AiVector2D {
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
