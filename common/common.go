package common

import "assimp/common/pb_msg"

type AiVector2D struct {
	X, Y float32
}

func (ai *AiVector2D) ToPbMsg() *pb_msg.AiVector2D {
	return &pb_msg.AiVector2D{X: ai.X, Y: ai.Y}
}

func (ai *AiVector2D) Set(pX, pY float32) {
	ai.X = pX
	ai.Y = pY
}

type AiColor4D struct {
	R, G, B, A float32
}

func (ai AiColor4D) Empty() bool {
	return ai.R == 0 && ai.G == 0 && ai.B == 0 && ai.A == 0
}
func (ai AiColor4D) ToPbMsg() *pb_msg.AiColor4D {
	return &pb_msg.AiColor4D{R: ai.R, G: ai.G, B: ai.B, A: ai.A}
}

type AiColor3D struct {
	R, G, B float32
}

func (ai AiColor3D) ToPbMsg() *pb_msg.AiColor3D {
	return &pb_msg.AiColor3D{R: ai.R, G: ai.G, B: ai.B}
}

func NewAiColor3D(R, G, B float32) AiColor3D {
	return AiColor3D{R: R, G: G, B: B}
}

type AiQuaternion struct {
	W, X, Y, Z float32
}

type AiPropertyStore struct {
	Sentinel uint8
}

// ------------------------------------------------------------------------------------
// Convert just one hex digit
// Return value is UINT_MAX if the input character is not a hex digit.
// ------------------------------------------------------------------------------------
func HexDigitToDecimal(in byte) (out uint) {
	if in >= '0' && in <= '9' {
		out = uint(in - '0')
	} else if in >= 'a' && in <= 'f' {
		out = 10 + uint(in-'a')
	} else if in >= 'A' && in <= 'F' {
		out = 10 + uint(in-'A')
		return out
	}

	// return value is UINT_MAX if the input is not a hex digit
	return out
}
