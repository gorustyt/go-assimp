package common

import "assimp/common/pb_msg"

type AiVertexWeight struct {
	//! Index of the vertex which is influenced by the bone.
	VertexId uint32

	//! The strength of the influence in the range (0...1).
	//! The influence from all bones at one vertex amounts to 1.
	Weight float32
}

func (ai *AiVertexWeight) Clone() *AiVertexWeight {
	tmp := *ai
	return &tmp
}
func (ai *AiVertexWeight) ToPbMsg() *pb_msg.AiVertexWeight {
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
