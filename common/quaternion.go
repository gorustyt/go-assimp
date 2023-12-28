package common

import "assimp/common/pb_msg"

type AiQuaternion struct {
	W, X, Y, Z float32
}

func (ai *AiQuaternion) ToPbMsg() *pb_msg.AiQuaternion {
	return &pb_msg.AiQuaternion{
		W: ai.W,
		X: ai.X,
		Y: ai.Y,
		Z: ai.Z,
	}
}

func NewAiQuaternion(pw, px, py, pz float32) *AiQuaternion {
	return &AiQuaternion{
		W: pw,
		X: px,
		Y: py,
		Z: pz,
	}
}

func NewAiQuaternion0() *AiQuaternion {
	return &AiQuaternion{
		W: 1,
	}
}
