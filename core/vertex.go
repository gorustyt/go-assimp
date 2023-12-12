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
	return &Vertex{
		texcoords: make([]*common.AiVector3D, AI_MAX_NUMBER_OF_TEXTURECOORDS),
		colors:    make([]*common.AiColor4D, AI_MAX_NUMBER_OF_COLOR_SETS),
	}
}

// ----------------------------------------------------------------------------
/** Extract a particular vertex from a mesh and interleave all components */

func (v *Vertex) VertexFromAiMesh(msh *AiMesh, idx int) {
	common.AiAssert(idx < msh.NumVertices)
	v.position = msh.Vertices[idx]

	if msh.HasNormals() {
		v.normal = msh.Normals[idx]
	}

	if msh.HasTangentsAndBitangents() {
		v.tangent = msh.Tangents[idx]
		v.bitangent = msh.Bitangents[idx]
	}

	for i := 0; msh.HasTextureCoords(i); i++ {
		v.texcoords[i] = msh.TextureCoords[i][idx]
	}

	for i := 0; msh.HasVertexColors(i); i++ {
		v.colors[i] = msh.Colors[i][idx]
	}
}

// ----------------------------------------------------------------------------
/** Extract a particular vertex from a anim mesh and interleave all components */
func (v *Vertex) VertexFromAiAnimMesh(msh *AiAnimMesh, idx int) {
	common.AiAssert(idx < msh.NumVertices)
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

	for i := 0; msh.HasTextureCoords(i); i++ {
		v.texcoords[i] = msh.TextureCoords[i][idx]
	}

	for i := 0; msh.HasVertexColors(i); i++ {
		v.colors[i] = msh.Colors[i][idx]
	}
}
