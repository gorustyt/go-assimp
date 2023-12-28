package common

import "assimp/common/pb_msg"

type AiVector4D [4]float32

func (ai AiVector4D) ToPbMsg() *pb_msg.AiVector4D {
	return &pb_msg.AiVector4D{X: ai[0], Y: ai[1], Z: ai[2], W: ai[3]}
}

func (ai *AiVector4D) FromPbMsg(data *pb_msg.AiVector4D) *AiVector4D {
	ai[0] = data.X
	ai[1] = data.Y
	ai[2] = data.Z
	ai[3] = data.W
	return ai
}

func NewAiVector4D() *AiVector4D {
	return &AiVector4D{}
}
