package common

import "assimp/common/pb_msg"

type AiQuaternion struct {
	W, X, Y, Z float32
}

func (ai *AiQuaternion) BoundMin(b *AiQuaternion) *AiQuaternion {
	return NewAiQuaternion(Min(ai.W, b.W), Min(ai.X, b.X), Min(ai.Y, b.Y), Min(ai.Z, b.Z))
}

func (ai *AiQuaternion) BoundMax(b *AiQuaternion) *AiQuaternion {
	return NewAiQuaternion(Max(ai.W, b.W), Max(ai.X, b.X), Max(ai.Y, b.Y), Max(ai.Z, b.Z))
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

type AiQuatKey struct {
	/** The time of this key */
	Time float64

	/** The value of this key */
	Value *AiQuaternion
}

func (ai *AiQuatKey) ToPbMsg() *pb_msg.AiQuatKey {
	r := &pb_msg.AiQuatKey{
		Time:  ai.Time,
		Value: ai.Value.ToPbMsg(),
	}
	return r
}

func (ai *AiQuatKey) BoundMin(b *AiQuatKey) *AiQuatKey {
	v := &AiQuatKey{
		Time:  Min(ai.Time, b.Time),
		Value: ai.Value.BoundMin(b.Value),
	}
	return v
}

func (ai *AiQuatKey) BoundMax(b *AiQuatKey) *AiQuatKey {
	v := &AiQuatKey{
		Time:  Max(ai.Time, b.Time),
		Value: ai.Value.BoundMax(b.Value),
	}
	return v
}
