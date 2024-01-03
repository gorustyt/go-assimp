package BLEND

import (
	"assimp/common"
	"assimp/core"
	"sort"
)

type TempArray struct {
}
type ConversionData struct {
	objects []*Object

	meshes    []*core.AiMesh
	cameras   []*core.AiCamera
	lights    []*core.AiLight
	materials []*core.AiMaterial
	textures  []*core.AiTexture

	// set of all materials referenced by at least one mesh in the scene
	materials_raw common.Queue[*Material]

	// counter to name sentinel textures inserted as substitutes for procedural textures.
	sentinel_cnt int32

	// next texture ID for each texture type, respectively
	next_texture []int32

	// original file data
	db *FileDatabase
}

func (s *ConversionData) InsertObject(obj *Object) {
	for _, v := range s.objects {
		if v == obj {
			return
		}
	}
	s.objects = append(s.objects, obj)
	sort.Slice(s.objects, func(i, j int) bool {
		return s.objects[i].id.name < s.objects[j].id.name
	})
}

func newConversionData(db *FileDatabase) *ConversionData {
	return &ConversionData{
		next_texture: make([]int32, core.AiTextureType_UNKNOWN+1),
		db:           db,
	}
}
