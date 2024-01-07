package pre_processing

import (
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/logger"
	"github.com/gorustyt/go-assimp/core"
	"log"
	"math"
)

type ScenePreprocessor struct {
	scene *core.AiScene
}

func NewScenePreprocessor(scene *core.AiScene) *ScenePreprocessor {
	return &ScenePreprocessor{scene: scene}
}

func (s *ScenePreprocessor) SetScene(scene *core.AiScene) {
	s.scene = scene
}

func (s *ScenePreprocessor) ProcessScene() {
	scene := s.scene
	if scene == nil {
		log.Fatal("empty scene")
	}
	// Process all meshes
	for i := 0; i < len(scene.Meshes); i++ {
		if nil == scene.Meshes[i] {
			continue
		}
		s.ProcessMesh(scene.Meshes[i])
	}

	// - nothing to do for nodes for the moment
	// - nothing to do for textures for the moment
	// - nothing to do for lights for the moment
	// - nothing to do for cameras for the moment

	// Process all animations
	for i := 0; i < len(scene.Animations); i++ {
		if nil == scene.Animations[i] {
			continue
		}
		s.ProcessAnimation(scene.Animations[i])
	}

	// Generate a default material if none was specified
	if len(scene.Materials) == 0 && len(scene.Meshes) != 0 {
		scene.Materials = make([]*core.AiMaterial, 2)
		helper := &core.AiMaterial{}
		scene.Materials[len(scene.Materials)] = helper
		clr := common.NewAiColor3D(0.6, 0.6, 0.6)
		helper.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_DIFFUSE, clr)

		// setup the default name to make this material identifiable
		name := core.AI_DEFAULT_MATERIAL_NAME
		helper.AddStringPropertyVar(core.AI_MATKEY_NAME, name)

		logger.DebugF("ScenePreprocessor: Adding default material \" %v \"", core.AI_DEFAULT_MATERIAL_NAME)

		for i := 0; i < len(scene.Meshes); i++ {
			if nil == scene.Meshes[i] {
				continue
			}
			scene.Meshes[i].MaterialIndex = int32(len(scene.Materials))
		}
	}
}

func (s *ScenePreprocessor) ProcessAnimation(anim *core.AiAnimation) {
	first := 10e10
	last := -10e10
	for i := 0; i < len(anim.Channels); i++ {
		channel := anim.Channels[i]

		//  If the exact duration of the animation is not given
		//  compute it now.
		if anim.Duration == -1. {
			// Position keys
			for j := 0; j < len(channel.PositionKeys); j++ {
				key := channel.PositionKeys[j]
				first = math.Min(first, key.Time)
				last = math.Max(last, key.Time)
			}

			// Scaling keys
			for j := 0; j < len(channel.ScalingKeys); j++ {
				key := channel.ScalingKeys[j]
				first = math.Min(first, key.Time)
				last = math.Max(last, key.Time)
			}

			// Rotation keys
			for j := 0; j < len(channel.RotationKeys); j++ {
				key := channel.RotationKeys[j]
				first = math.Min(first, key.Time)
				last = math.Max(last, key.Time)
			}
		}
	}

	if anim.Duration == -1. {
		logger.Debug("ScenePreprocessor: Setting animation duration")
		anim.Duration = last - math.Min(float64(first), 0.)
	}
}

func (s *ScenePreprocessor) ProcessMesh(mesh *core.AiMesh) {
	// If aiMesh::mNumUVComponents is *not* set assign the default value of 2
	for i := 0; i < int(core.AI_MAX_NUMBER_OF_TEXTURECOORDS); i++ {
		if len(mesh.TextureCoords[i]) == 0 {
			mesh.NumUVComponents[i] = 0
			continue
		}

		if mesh.NumUVComponents[i] == 0 {
			mesh.NumUVComponents[i] = 2
		}

		p := 0

		end := len(mesh.Vertices)

		// Ensure unused components are zeroed. This will make 1D texture channels work
		// as if they were 2D channels .. just in case an application doesn't handle
		// this case
		if 2 == mesh.NumUVComponents[i] {
			for ; p != end; p++ {
				mesh.TextureCoords[i][p].Z = 0.
			}
		} else if 1 == mesh.NumUVComponents[i] {
			for ; p != end; p++ {
				mesh.TextureCoords[i][p].Y = 0.
				mesh.TextureCoords[i][p].Z = mesh.TextureCoords[i][p].Y
			}
		} else if 3 == mesh.NumUVComponents[i] {
			// Really 3D coordinates? Check whether the third coordinate is != 0 for at least one element
			for ; p != end; p++ {
				if mesh.TextureCoords[i][p].Z != 0 {
					break
				}
			}
			if p == end {
				logger.Warn("ScenePreprocessor: UVs are declared to be 3D but they're obviously not. Reverting to 2D.")
				mesh.NumUVComponents[i] = 2
			}
		}
	}

	// If the information which primitive types are there in the
	// mesh is currently not available, compute it.
	if mesh.PrimitiveTypes == 0 {
		if mesh.Faces == nil {
			logger.Fatal("mesh.Faces == nil")
		}
		for a := 0; a < len(mesh.Faces); a++ {
			face := mesh.Faces[a]
			switch len(face.Indices) {
			case 3:
				mesh.PrimitiveTypes |= core.AiPrimitiveType_TRIANGLE
			case 2:
				mesh.PrimitiveTypes |= core.AiPrimitiveType_LINE
			case 1:
				mesh.PrimitiveTypes |= core.AiPrimitiveType_POINT

			default:
				mesh.PrimitiveTypes |= core.AiPrimitiveType_POLYGON
			}
		}
	}

	// If tangents and normals are given but no bitangents compute them
	if mesh.Tangents != nil && mesh.Normals != nil && mesh.Bitangents == nil {
		mesh.Bitangents = make([]*common.AiVector3D, len(mesh.Vertices))
		for i := 0; i < len(mesh.Vertices); i++ {
			mesh.Bitangents[i] = mesh.Normals[i].NegationOperationSymbol(mesh.Tangents[i])
		}
	}
}
