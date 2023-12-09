package common

import "assimp/common/pb_msg"

type AiVector3D struct {
	X, Y, Z float32
}

func (ai AiVector3D) Empty() bool {
	return ai.X == 0 && ai.Y == 0 && ai.Z == 0
}
func (ai AiVector3D) ToPbMsg() *pb_msg.AiVector3D {
	return &pb_msg.AiVector3D{X: ai.X, Y: ai.Y, Z: ai.Z}
}

func (ai AiVector3D) Add(o AiVector3D) AiVector3D {
	ai.X += o.X
	ai.Y += o.Y
	ai.Z += o.Z
	return ai
}

func (ai AiVector3D) Sub(o AiVector3D) AiVector3D {
	ai.X -= o.X
	ai.Y -= o.Y
	ai.Z -= o.Z
	return ai
}

func (ai AiVector3D) Mul(f float32) AiVector3D {
	ai.X *= f
	ai.Y *= f
	ai.Z *= f
	return ai
}

func (ai AiVector3D) Div(f float32) AiVector3D {
	if f == 0 {
		return ai
	}
	invF := float32(1.0) / f
	ai.X *= invF
	ai.Y *= invF
	ai.Z *= invF
	return ai
}
