package core

import (
	"assimp/common"
	"assimp/common/pb_msg"
)

type AiAABB struct {
	Min *common.AiVector3D
	Max *common.AiVector3D
}

func (ai *AiAABB) FromPbMsg(p *pb_msg.AiAABB) *AiAABB {
	if p == nil {
		return nil
	}
	ai.Min = (&common.AiVector3D{}).FromPbMsg(p.Min)
	ai.Max = (&common.AiVector3D{}).FromPbMsg(p.Max)
	return ai
}

func (ai *AiAABB) ToPbMsg() *pb_msg.AiAABB {
	if ai == nil {
		return nil
	}
	return &pb_msg.AiAABB{
		Min: ai.Min.ToPbMsg(),
		Max: ai.Max.ToPbMsg(),
	}
}

func (ai *AiAABB) Clone() *AiAABB {
	if ai == nil {
		return nil
	}
	tmp := *ai
	return &tmp
}
