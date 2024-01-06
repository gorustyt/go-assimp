package common

import "assimp/common/pb_msg"

type AiVertexWeight struct {
	//! Index of the vertex which is influenced by the bone.
	VertexId uint32

	//! The strength of the influence in the range (0...1).
	//! The influence from all bones at one vertex amounts to 1.
	Weight float32
}

func (ai *AiVertexWeight) FromPbMsg(p *pb_msg.AiVertexWeight) *AiVertexWeight {
	if p == nil {
		return nil
	}
	ai.VertexId = p.VertexId
	ai.Weight = p.Weight
	return ai
}

func (ai *AiVertexWeight) Clone() *AiVertexWeight {
	if ai == nil {
		return nil
	}
	tmp := *ai
	return &tmp
}

func (ai *AiVertexWeight) ToPbMsg() *pb_msg.AiVertexWeight {
	if ai == nil {
		return nil
	}
	return &pb_msg.AiVertexWeight{
		VertexId: ai.VertexId,
		Weight:   ai.Weight,
	}
}

func (ai *AiVertexWeight) BoundMin(b *AiVertexWeight) *AiVertexWeight {
	v := &AiVertexWeight{
		VertexId: Min(ai.VertexId, b.VertexId),
		Weight:   Min(ai.Weight, b.Weight),
	}
	return v
}

func (ai *AiVertexWeight) BoundMax(a, b *AiVertexWeight) *AiVertexWeight {
	v := &AiVertexWeight{
		VertexId: Max(ai.VertexId, b.VertexId),
		Weight:   Max(ai.Weight, b.Weight),
	}
	return v
}
