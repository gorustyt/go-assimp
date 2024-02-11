package post_processing

import (
	"context"
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/logger"
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"github.com/gorustyt/go-assimp/core"
	"github.com/gorustyt/go-assimp/driver/base/iassimp"
	"google.golang.org/protobuf/proto"
)

var (
	_ iassimp.PostProcessing = (*MakeLeftHandedProcess)(nil)
	_ iassimp.PostProcessing = (*FlipWindingOrderProcess)(nil)
	_ iassimp.PostProcessing = (*FlipUVsProcess)(nil)
)

type FlipWindingOrderProcess struct {
}

func (m *FlipWindingOrderProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_FlipWindingOrder) != 0
}

func (m *FlipWindingOrderProcess) Execute(pScene *core.AiScene) {
	logger.DebugF("FlipWindingOrderProcess begin")
	for i := 0; i < len(pScene.Meshes); i++ {
		m.ProcessMesh(pScene.Meshes[i])
	}
	logger.DebugF("FlipWindingOrderProcess finished")
}

func (m *FlipWindingOrderProcess) SetupProperties(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

// ------------------------------------------------------------------------------------------------
// Converts a single mesh
func (m *FlipWindingOrderProcess) ProcessMesh(pMesh *core.AiMesh) {
	// invert the order of all faces in this mesh
	for a := 0; a < len(pMesh.Faces); a++ {
		face := pMesh.Faces[a]
		for b := 0; b < len(face.Indices)/2; b++ {
			face.Indices[b], face.Indices[len(face.Indices)-1-b] = face.Indices[len(face.Indices)-1-b], face.Indices[b]
		}
	}

	// invert the order of all components in this mesh anim meshes
	for m := 0; m < len(pMesh.AnimMeshes); m++ {
		animMesh := pMesh.AnimMeshes[m]
		numVertices := len(animMesh.Vertices)
		if animMesh.HasPositions() {
			for a := 0; a < numVertices; a++ {
				animMesh.Vertices[a], animMesh.Vertices[numVertices-1-a] = animMesh.Vertices[numVertices-1-a], animMesh.Vertices[a]
			}
		}
		if animMesh.HasNormals() {
			for a := 0; a < numVertices; a++ {
				animMesh.Normals[a], animMesh.Normals[numVertices-1-a] = animMesh.Normals[numVertices-1-a], animMesh.Normals[a]
			}
		}
		for i := uint32(0); i < core.AI_MAX_NUMBER_OF_TEXTURECOORDS; i++ {
			if animMesh.HasTextureCoords(i) {
				for a := 0; a < numVertices; a++ {
					animMesh.TextureCoords[i][a], animMesh.TextureCoords[i][numVertices-1-a] = animMesh.TextureCoords[i][numVertices-1-a], animMesh.TextureCoords[i][a]
				}
			}
		}
		if animMesh.HasTangentsAndBitangents() {
			for a := 0; a < numVertices; a++ {
				animMesh.Tangents[a], animMesh.Tangents[numVertices-1-a] = animMesh.Tangents[numVertices-1-a], animMesh.Tangents[a]
				animMesh.Bitangents[a], animMesh.Bitangents[numVertices-1-a] = animMesh.Bitangents[numVertices-1-a], animMesh.Bitangents[a]
			}
		}
		for v := 0; v < core.AI_MAX_NUMBER_OF_COLOR_SETS; v++ {
			if animMesh.HasVertexColors(v) {
				for a := 0; a < numVertices; a++ {
					animMesh.Colors[v][a], animMesh.Colors[v][numVertices-1-a] = animMesh.Colors[v][numVertices-1-a], animMesh.Colors[v][a]
				}
			}
		}
	}
}

type FlipUVsProcess struct {
}

func (m *FlipUVsProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_FlipUVs) != 0
}

func (m *FlipUVsProcess) Execute(pScene *core.AiScene) {
	logger.DebugF("FlipUVsProcess begin")
	for i := 0; i < len(pScene.Meshes); i++ {
		m.ProcessMesh(pScene.Meshes[i])
	}

	for i := 0; i < len(pScene.Materials); i++ {
		m.ProcessMaterial(pScene.Materials[i])
	}

	logger.DebugF("FlipUVsProcess finished")
}

func (m *FlipUVsProcess) SetupProperties(ctx context.Context) {

}

// ------------------------------------------------------------------------------------------------
// Converts a single material
func (m *FlipUVsProcess) ProcessMaterial(mat *core.AiMaterial) {

	for a := 0; a < len(mat.Properties); a++ {
		prop := mat.Properties[a]
		if prop == nil {
			logger.DebugF("Property is null")
			continue
		}

		// UV transformation key?
		if prop.Key == "$tex.uvtrafo" {
			prop.UpdateData(func(v proto.Message) {
				uv := v.(*pb_msg.AiUVTransform)
				// just flip it, that's everything
				uv.Translation.Y *= -1.
				uv.Rotation *= -1.
			})

		}
	}
}

// ------------------------------------------------------------------------------------------------
// Converts a single mesh
func (m *FlipUVsProcess) ProcessMesh(pMesh *core.AiMesh) {
	flipUVsAiMesh(pMesh)
	for idx := 0; idx < len(pMesh.AnimMeshes); idx++ {
		flipUVsAiAnimMesh(pMesh.AnimMeshes[idx])
	}
}

type MakeLeftHandedProcess struct {
}

func (m *MakeLeftHandedProcess) SetupProperties(ctx context.Context) {

}

// ------------------------------------------------------------------------------------------------
// Returns whether the processing step is present in the given flag field.
func (m *MakeLeftHandedProcess) IsActive(pFlags int) bool {
	return (iassimp.AiPostProcessSteps(pFlags) & iassimp.AiProcess_MakeLeftHanded) != 0
}

// ------------------------------------------------------------------------------------------------
// Executes the post processing step on the given imported data.
func (m *MakeLeftHandedProcess) Execute(pScene *core.AiScene) {
	// Check for an existent root node to proceed
	logger.DebugF("MakeLeftHandedProcess begin")

	// recursively convert all the nodes
	m.ProcessNode(pScene.RootNode, common.NewAiMatrix4x4Identify())

	// process the meshes accordingly
	for a := 0; a < len(pScene.Meshes); a++ {
		m.ProcessMesh(pScene.Meshes[a])
	}

	// process the materials accordingly
	for a := 0; a < len(pScene.Materials); a++ {
		m.ProcessMaterial(pScene.Materials[a])
	}

	// transform all animation channels as well
	for a := 0; a < len(pScene.Animations); a++ {
		anim := pScene.Animations[a]
		for b := 0; b < len(anim.Channels); b++ {
			nodeAnim := anim.Channels[b]
			m.ProcessAnimation(nodeAnim)
		}
	}

	// process the cameras accordingly
	for a := 0; a < len(pScene.Cameras); a++ {
		m.ProcessCamera(pScene.Cameras[a])
	}
	logger.DebugF("MakeLeftHandedProcess finished")
}

func flipUVsAiMesh(pMesh *core.AiMesh) {
	if pMesh == nil {
		return
	}
	// mirror texture y coordinate
	for tcIdx := uint32(0); tcIdx < core.AI_MAX_NUMBER_OF_TEXTURECOORDS; tcIdx++ {
		if !pMesh.HasTextureCoords(tcIdx) {
			break
		}

		for vIdx := 0; vIdx < len(pMesh.Vertices); vIdx++ {
			pMesh.TextureCoords[tcIdx][vIdx].Y = 1.0 - pMesh.TextureCoords[tcIdx][vIdx].Y
		}
	}
}

func flipUVsAiAnimMesh(pMesh *core.AiAnimMesh) {
	if pMesh == nil {
		return
	}
	// mirror texture y coordinate
	for tcIdx := uint32(0); tcIdx < core.AI_MAX_NUMBER_OF_TEXTURECOORDS; tcIdx++ {
		if !pMesh.HasTextureCoords(tcIdx) {
			break
		}

		for vIdx := 0; vIdx < len(pMesh.Vertices); vIdx++ {
			pMesh.TextureCoords[tcIdx][vIdx].Y = 1.0 - pMesh.TextureCoords[tcIdx][vIdx].Y
		}
	}
}

// ------------------------------------------------------------------------------------------------
// Recursively converts a node, all of its children and all of its meshes
func (m *MakeLeftHandedProcess) ProcessNode(pNode *core.AiNode, pParentGlobalRotation *common.AiMatrix4x4) {
	// mirror all base vectors at the local Z axis
	pNode.Transformation.C1 = -pNode.Transformation.C1
	pNode.Transformation.C2 = -pNode.Transformation.C2
	pNode.Transformation.C3 = -pNode.Transformation.C3
	pNode.Transformation.C4 = -pNode.Transformation.C4

	// now invert the Z axis again to keep the matrix determinant positive.
	// The local meshes will be inverted accordingly so that the result should look just fine again.
	pNode.Transformation.A3 = -pNode.Transformation.A3
	pNode.Transformation.B3 = -pNode.Transformation.B3
	pNode.Transformation.C3 = -pNode.Transformation.C3
	pNode.Transformation.D3 = -pNode.Transformation.D3 // useless, but anyways...

	// continue for all children
	for a := 0; a < len(pNode.Children); a++ {
		m.ProcessNode(pNode.Children[a], pParentGlobalRotation.MulAiMatrix4x4(pNode.Transformation))
	}
}

// ------------------------------------------------------------------------------------------------
// Converts a single mesh to left handed coordinates.
func (m *MakeLeftHandedProcess) ProcessMesh(pMesh *core.AiMesh) {
	if nil == pMesh {
		logger.ErrorF("Nullptr to mesh found.")
		return
	}
	// mirror positions, normals and stuff along the Z axis
	for a := 0; a < len(pMesh.Vertices); a++ {
		pMesh.Vertices[a].Z *= -1.0
		if pMesh.HasNormals() {
			pMesh.Normals[a].Z *= -1.0
		}
		if pMesh.HasTangentsAndBitangents() {
			pMesh.Tangents[a].Z *= -1.0
			pMesh.Bitangents[a].Z *= -1.0
		}
	}

	// mirror anim meshes positions, normals and stuff along the Z axis
	for m := 0; m < len(pMesh.AnimMeshes); m++ {
		for a := 0; a < len(pMesh.AnimMeshes[m].Vertices); a++ {
			pMesh.AnimMeshes[m].Vertices[a].Z *= -1.0
			if pMesh.AnimMeshes[m].HasNormals() {
				pMesh.AnimMeshes[m].Normals[a].Z *= -1.0
			}
			if pMesh.AnimMeshes[m].HasTangentsAndBitangents() {
				pMesh.AnimMeshes[m].Tangents[a].Z *= -1.0
				pMesh.AnimMeshes[m].Bitangents[a].Z *= -1.0
			}
		}
	}

	// mirror offset matrices of all bones
	for a := 0; a < len(pMesh.Bones); a++ {
		bone := pMesh.Bones[a]
		bone.OffsetMatrix.A3 = -bone.OffsetMatrix.A3
		bone.OffsetMatrix.B3 = -bone.OffsetMatrix.B3
		bone.OffsetMatrix.D3 = -bone.OffsetMatrix.D3
		bone.OffsetMatrix.C1 = -bone.OffsetMatrix.C1
		bone.OffsetMatrix.C2 = -bone.OffsetMatrix.C2
		bone.OffsetMatrix.C4 = -bone.OffsetMatrix.C4
	}

	// mirror bitangents as well as they're derived from the texture coords
	if pMesh.HasTangentsAndBitangents() {
		for a := 0; a < len(pMesh.Vertices); a++ {
			pMesh.Bitangents[a] = pMesh.Bitangents[a].Mul(-1.0)
		}

	}
}

// ------------------------------------------------------------------------------------------------
// Converts a single material to left handed coordinates.
func (m *MakeLeftHandedProcess) ProcessMaterial(mat *core.AiMaterial) {
	if nil == mat {
		logger.ErrorF("Nullptr to aiMaterial found.")
		return
	}

	for a := 0; a < len(mat.Properties); a++ {
		prop := mat.Properties[a]

		// Mapping axis for UV mappings?
		if prop.Key == "$tex.mapaxis" {
			prop.UpdateData(func(v proto.Message) {
				pff := v.(*pb_msg.AiVector3D)
				pff.Z *= -1.
			})
		}
	}
}

// ------------------------------------------------------------------------------------------------
// Converts the given animation to LH coordinates.
func (m *MakeLeftHandedProcess) ProcessAnimation(pAnim *core.AiNodeAnim) {
	// position keys
	for a := 0; a < len(pAnim.PositionKeys); a++ {
		pAnim.PositionKeys[a].Value.Z *= -1.0
	}

	// rotation keys
	for a := 0; a < len(pAnim.RotationKeys); a++ {
		pAnim.RotationKeys[a].Value.X *= -1.0
		pAnim.RotationKeys[a].Value.Y *= -1.0
	}
}

// ------------------------------------------------------------------------------------------------
// Converts a single camera to left handed coordinates.
func (m *MakeLeftHandedProcess) ProcessCamera(pCam *core.AiCamera) {
	pCam.LookAt = pCam.Position.Mul(2.0).Sub(pCam.LookAt)
}
