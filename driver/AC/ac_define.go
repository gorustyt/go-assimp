package AC

import (
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/logger"
	"github.com/gorustyt/go-assimp/common/reader"
	"github.com/gorustyt/go-assimp/core"
	"strings"
)

type ObjectType int

const (
	World ObjectType = 0x0
	Poly  ObjectType = 0x1
	Group ObjectType = 0x2
	Light ObjectType = 0x4
)

// Represents an AC3D object
type Object struct {
	Type ObjectType
	// name of the object
	name string

	// object children
	children []*Object

	// texture to be assigned to all surfaces of the object
	// the .acc format supports up to 4 textures
	textures []string

	// texture repat factors (scaling for all coordinates)
	texRepeat, texOffset *common.AiVector2D

	// rotation matrix
	rotation *common.AiMatrix3x3

	// translation vector
	translation *common.AiVector3D

	// vertices
	vertices []*common.AiVector3D

	// surfaces
	surfaces []*Surface

	// number of indices (= num verts in verbose format)
	numRefs int

	// number of subdivisions to be performed on the
	// imported data
	subDiv int

	// max angle limit for smoothing
	crease float32
}

func newObject() *Object {
	return &Object{
		Type:        World,
		texRepeat:   common.NewAiVector2D(1, 1),
		texOffset:   common.NewAiVector2D(0, 0),
		rotation:    common.NewAiMatrix3x3(),
		translation: common.NewAiVector3D3(0, 0, 0),
	}
}

type SurfaceType int

// Type is low nibble of flags
const (
	Polygon       SurfaceType = 0x0
	ClosedLine    SurfaceType = 0x1
	OpenLine      SurfaceType = 0x2
	TriangleStrip SurfaceType = 0x4 // ACC extension (TORCS and Speed Dreams)
	Mask          SurfaceType = 0xf
)

type SurfaceEntry common.Pair[int, common.AiVector2D]
type Surface struct {
	mat, flags int
	entries    []*SurfaceEntry
}

func (s *Surface) GetType() int { return s.flags & int(Mask) }

type Material struct {
	// base color of the material
	rgb *common.AiColor3D

	// ambient color of the material
	amb *common.AiColor3D

	// emissive color of the material
	emis *common.AiColor3D

	// specular color of the material
	spec *common.AiColor3D

	// shininess exponent
	shin float32

	// transparency. 0 == opaque
	trans float32

	// name of the material. optional.
	name string
}

func newMaterial() *Material {
	return &Material{
		rgb:  common.NewAiColor3D(0.6, 0.6, 0.6),
		spec: common.NewAiColor3D(1, 1, 1),
		emis: common.NewAiColor3D(0, 0, 0),
		amb:  common.NewAiColor3D(0, 0, 0),
	}
}

type AC3DImporter struct {
	// points to the next data line
	buffer []byte

	// Configuration option: if enabled, up to two meshes
	// are generated per material: those faces who have
	// their bf cull flags set are separated.
	configSplitBFCull bool

	// Configuration switch: subdivision surfaces are only
	// evaluated if the value is true.
	configEvalSubdivision bool

	// counts how many objects we have in the tree.
	// basing on this information we can find a
	// good estimate how many meshes we'll have in the final scene.
	NumMeshes int

	// current list of light sources
	Lights []*core.AiLight

	// name counters
	LightsCounter, GroupsCounter, PolysCounter, WorldsCounter int
	reader.LineReader
}

func (im *AC3DImporter) CanRead(checkSig bool) bool {
	im.NextLine()
	if im.GetLineNum() == 1 && !strings.HasPrefix(im.GetLine(), Desc.Magic) {
		if !checkSig {
			logger.WarnF("not found magic expect:%v found:%v", Desc.Magic, im.GetLine())
		}
		return false
	}
	version := strings.TrimPrefix(im.GetLine(), Desc.Magic)
	logger.InfoF("importer:%v version:%v", Desc.Name, common.HexDigitToDecimal([]byte(version)[0]))
	im.NextLine()
	return true
}
