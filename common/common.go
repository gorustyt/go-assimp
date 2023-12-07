package common

type AiVector2D [2]float64
type AiVector3D [3]float64
type AiVector4D [4]float64
type AiColor4D struct {
	R, G, B, A float64
}
type AiColor3D struct {
	R, G, B float64
}

func NewAiColor3D(R, G, B float64) AiColor3D {
	return AiColor3D{R: R, G: G, B: B}
}

type AiMatrix4x4 [16]float64

type AiQuaternion struct {
	W, X, Y, Z float64
}

type AiPropertyStore struct {
	Sentinel uint8
}

type AiMatrix3x3 [9]float64
