package core

import "assimp/common"

type Vertex struct {
	position           *common.AiVector3D
	normal             *common.AiVector3D
	tangent, bitangent *common.AiVector3D

	texcoords []*common.AiVector3D
	colors    []*common.AiColor4D
}

func NewVertex() *Vertex {
	v := &Vertex{
		position:  common.NewAiVector3D(),
		normal:    common.NewAiVector3D(),
		tangent:   common.NewAiVector3D(),
		bitangent: common.NewAiVector3D(),
		texcoords: make([]*common.AiVector3D, AI_MAX_NUMBER_OF_TEXTURECOORDS),
		colors:    make([]*common.AiColor4D, AI_MAX_NUMBER_OF_COLOR_SETS),
	}
	for i := range v.texcoords {
		v.texcoords[i] = common.NewAiVector3D()
	}
	for i := range v.colors {
		v.colors[i] = common.NewAiColor4D0()
	}
	return v
}

// ----------------------------------------------------------------------------
/** Extract a particular vertex from a mesh and interleave all components */

func NewVertexFromAiMesh(msh *AiMesh, idx int) *Vertex {
	v := NewVertex()
	common.AiAssert(idx < len(msh.Vertices))
	v.position = msh.Vertices[idx]

	if msh.HasNormals() {
		v.normal = msh.Normals[idx]
	}

	if msh.HasTangentsAndBitangents() {
		v.tangent = msh.Tangents[idx]
		v.bitangent = msh.Bitangents[idx]
	}

	for i := 0; msh.HasTextureCoords(uint32(i)); i++ {
		v.texcoords[i] = msh.TextureCoords[i][idx]
	}

	for i := 0; msh.HasVertexColors(i); i++ {
		v.colors[i] = msh.Colors[i][idx]
	}
	return v
}

// ----------------------------------------------------------------------------
/** Extract a particular vertex from a anim mesh and interleave all components */
func NewVertexFromAiAnimMesh(msh *AiAnimMesh, idx int) *Vertex {
	v := NewVertex()
	common.AiAssert(idx < len(msh.Vertices))
	if msh.HasPositions() {
		v.position = msh.Vertices[idx]
	}

	if msh.HasNormals() {
		v.normal = msh.Normals[idx]
	}

	if msh.HasTangentsAndBitangents() {
		v.tangent = msh.Tangents[idx]
		v.bitangent = msh.Bitangents[idx]
	}

	for i := 0; msh.HasTextureCoords(uint32(i)); i++ {
		v.texcoords[i] = msh.TextureCoords[i][idx]
	}

	for i := 0; msh.HasVertexColors(i); i++ {
		v.colors[i] = msh.Colors[i][idx]
	}
	return v
}

// ------------------------------------------------------------------------------------------------
func (v *Vertex) Add(v0 *Vertex, v1 *Vertex) *Vertex {
	return v.BinaryOpVertex(v0, v1, "+")
}

func (v *Vertex) Sub(v0 *Vertex, v1 *Vertex) *Vertex {
	return v.BinaryOpVertex(v0, v1, "-")
}

func (v *Vertex) Mul(v0 *Vertex, f float32) *Vertex {
	return v.BinaryOp(v0, f, "*")
}

func (v *Vertex) Div(v0 *Vertex, f float32) *Vertex {
	return v.BinaryOp(v0, 1./f, "*")
}

func (v *Vertex) Mul1(f float32, v0 *Vertex) *Vertex {
	return v.BinaryOp1(f, v0, "*")
}

// ----------------------------------------------------------------------------
// / Convert back to non-interleaved storage
func (v *Vertex) SortBack(out *AiMesh, idx int) {
	common.AiAssert(idx < len(out.Vertices))
	out.Vertices[idx] = v.position

	if out.HasNormals() {
		out.Normals[idx] = v.normal
	}

	if out.HasTangentsAndBitangents() {
		out.Tangents[idx] = v.tangent
		out.Bitangents[idx] = v.bitangent
	}

	for i := 0; out.HasTextureCoords(uint32(i)); i++ {
		out.TextureCoords[i][idx] = v.texcoords[i]
	}

	for i := 0; out.HasVertexColors(i); i++ {
		out.Colors[i][idx] = v.colors[i]
	}
}

// ----------------------------------------------------------------------------
// / Construct from two operands and a binary operation to combine them
func (v *Vertex) BinaryOpVertex(v0 *Vertex, v1 *Vertex, opStr string) *Vertex {
	// this is a heavy task for the compiler to optimize ... *pray*

	res := NewVertex()
	res.position = op(v0.position, v1.position, opStr).(*common.AiVector3D)
	res.normal = op(v0.normal, v1.normal, opStr).(*common.AiVector3D)
	res.tangent = op(v0.tangent, v1.tangent, opStr).(*common.AiVector3D)
	res.bitangent = op(v0.bitangent, v1.bitangent, opStr).(*common.AiVector3D)
	for i := 0; i < int(AI_MAX_NUMBER_OF_TEXTURECOORDS); i++ {
		res.texcoords[i] = op(v0.texcoords[i], v1.texcoords[i], opStr).(*common.AiVector3D)
	}
	for i := 0; i < AI_MAX_NUMBER_OF_COLOR_SETS; i++ {
		res.colors[i] = op(v0.colors[i], v1.colors[i], opStr).(*common.AiColor4D)
	}
	return res
}

// ----------------------------------------------------------------------------
// / This time binary arithmetic of v0 with a floating-point number
func (v *Vertex) BinaryOp(v0 *Vertex, f float32, opStr string) *Vertex {
	// this is a heavy task for the compiler to optimize ... *pray*

	res := NewVertex()
	res.position = op(v0.position, f, opStr).(*common.AiVector3D)
	res.normal = op(v0.normal, f, opStr).(*common.AiVector3D)
	res.tangent = op(v0.tangent, f, opStr).(*common.AiVector3D)
	res.bitangent = op(v0.bitangent, f, opStr).(*common.AiVector3D)

	for i := 0; i < int(AI_MAX_NUMBER_OF_TEXTURECOORDS); i++ {
		res.texcoords[i] = op(v0.texcoords[i], f, opStr).(*common.AiVector3D)
	}
	for i := 0; i < AI_MAX_NUMBER_OF_COLOR_SETS; i++ {
		res.colors[i] = op(v0.colors[i], f, opStr).(*common.AiColor4D)
	}
	return res
}

// ----------------------------------------------------------------------------
/** This time binary arithmetic of v0 with a floating-point number */
func (v *Vertex) BinaryOp1(f float32, v0 *Vertex, opStr string) *Vertex {
	// this is a heavy task for the compiler to optimize ... *pray*

	res := NewVertex()
	res.position = op(f, v0.position, opStr).(*common.AiVector3D)
	res.normal = op(f, v0.normal, opStr).(*common.AiVector3D)
	res.tangent = op(f, v0.tangent, opStr).(*common.AiVector3D)
	res.bitangent = op(f, v0.bitangent, opStr).(*common.AiVector3D)

	for i := 0; i < int(AI_MAX_NUMBER_OF_TEXTURECOORDS); i++ {
		res.texcoords[i] = op(f, v0.texcoords[i], opStr).(*common.AiVector3D)
	}
	for i := 0; i < AI_MAX_NUMBER_OF_COLOR_SETS; i++ {
		res.colors[i] = op(f, v0.colors[i], opStr).(*common.AiColor4D)
	}
	return res
}

func op(v1 interface{}, v2 interface{}, op string) any {
	_, ok := v2.(float32)
	if ok {
		v1, v2 = v2, v1
	}
	switch value := v1.(type) {
	case *common.AiVector3D:
		switch op {
		case "+":
			return value.Add(v2.(*common.AiVector3D))
		case "*":
			return value.MulAiVector3D(v2.(*common.AiVector3D))
		case "-":
			return value.Sub(v2.(*common.AiVector3D))
		}
	case *common.AiColor4D:
		switch op {
		case "+":
			return value.Add(v2.(*common.AiColor4D))
		case "*":
			return value.Mul(v2.(*common.AiColor4D))
		case "-":
			return value.Sub(v2.(*common.AiColor4D))
		}
	case float32:
		switch v2Value := v2.(type) {
		case *common.AiVector3D:
			return v2Value.Mul(value)
		case *common.AiColor4D:
			return v2Value.MulValue(value)
		}
	}
	return nil
}
