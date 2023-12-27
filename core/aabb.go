package core

import (
	"assimp/common"
	"assimp/common/pb_msg"
)

type AiAABB struct {
	Min *common.AiVector3D
	Max *common.AiVector3D
}

func (ai *AiAABB) ToPbMsg() *pb_msg.AiAABB {
	return &pb_msg.AiAABB{
		Min: ai.Min.ToPbMsg(),
		Max: ai.Max.ToPbMsg(),
	}
}

func (ai *AiAABB) Clone() *AiAABB {
	tmp := *ai
	return &tmp
}
