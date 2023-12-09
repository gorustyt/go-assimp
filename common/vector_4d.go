package common

import "assimp/common/pb_msg"

type AiVector4D [4]float32

func (ai AiVector4D) ToPbMsg() *pb_msg.AiVector4D {
	return &pb_msg.AiVector4D{X: ai[0], Y: ai[1], Z: ai[2], W: ai[3]}
}
