package common

type AiMatrix4x4 struct {
	A1, A2, A3, A4 float32
	B1, B2, B3, B4 float32
	C1, C2, C3, C4 float32
	D1, D2, D3, D4 float32
}

func NewAiMatrix4x4FromAiMatrix3x3(m AiMatrix3x3) AiMatrix4x4 {
	ai := AiMatrix4x4{}
	ai.A1 = m.A1
	ai.A2 = m.A2
	ai.A3 = m.A3
	ai.A4 = 0.0
	ai.B1 = m.B1
	ai.B2 = m.B2
	ai.B3 = m.B3
	ai.B4 = 0.0
	ai.C1 = m.C1
	ai.C2 = m.C2
	ai.C3 = m.C3
	ai.C4 = 0.0
	ai.D1 = 0.0
	ai.D2 = 0.0
	ai.D3 = 0.0
	ai.D4 = 1.0
	return ai
}
