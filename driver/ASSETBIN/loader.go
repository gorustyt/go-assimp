package ASSETBIN

import (
	"assimp/common"
	"assimp/common/config"
	"assimp/common/reader"
	"assimp/core"
	"assimp/driver/base/iassimp"
	"errors"
	"fmt"
)

var (
	Desc = core.AiImporterDesc{
		"Assimp Binary Importer",
		"Gargaj / Conspiracy",
		"",
		"",
		core.AiImporterFlags_SupportBinaryFlavour | core.AiImporterFlags_SupportCompressedFlavour,
		0,
		0,
		0,
		0,
		[]string{"assbin"},
		"ASSIMP.binary-dump.",
	}
)

const (
	// these are the magic chunk identifiers for the binary ASS file format
	ASSBIN_CHUNK_AICAMERA           = 0x1234
	ASSBIN_CHUNK_AILIGHT            = 0x1235
	ASSBIN_CHUNK_AITEXTURE          = 0x1236
	ASSBIN_CHUNK_AIMESH             = 0x1237
	ASSBIN_CHUNK_AINODEANIM         = 0x1238
	ASSBIN_CHUNK_AISCENE            = 0x1239
	ASSBIN_CHUNK_AIBONE             = 0x123a
	ASSBIN_CHUNK_AIANIMATION        = 0x123b
	ASSBIN_CHUNK_AINODE             = 0x123c
	ASSBIN_CHUNK_AIMATERIAL         = 0x123d
	ASSBIN_CHUNK_AIMATERIALPROPERTY = 0x123e

	ASSBIN_MESH_HAS_POSITIONS               = 0x1
	ASSBIN_MESH_HAS_NORMALS                 = 0x2
	ASSBIN_MESH_HAS_TANGENTS_AND_BITANGENTS = 0x4
	ASSBIN_MESH_HAS_TEXCOORD_BASE           = 0x100
	ASSBIN_MESH_HAS_COLOR_BASE              = 0x10000

	ASSBIN_VERSION_MAJOR = 1
	ASSBIN_VERSION_MINOR = 0
)

func ASSBIN_MESH_HAS_TEXCOORD(n int32) int32 {
	return ASSBIN_MESH_HAS_TEXCOORD_BASE << n
}
func ASSBIN_MESH_HAS_COLOR(n int32) int32 {
	return ASSBIN_MESH_HAS_COLOR_BASE << n
}

type AssBinImporter struct {
	reader.StreamReader
}

func (ai *AssBinImporter) CanRead(checkSig bool) bool {
	data, err := ai.Peek(32)
	if err != nil {
		return false
	}
	if string(data)[:19] != Desc.Magic {
		return false
	}
	return true
}

func NewAssBinImporter(data []byte) (iassimp.Loader, error) {
	r, err := reader.NewFileStreamReader(data)
	if err != nil {
		return nil, err
	}
	return &AssBinImporter{StreamReader: r}, nil
}
func (ai *AssBinImporter) InitConfig(cfg *config.Config) {

}
func (ai *AssBinImporter) Read(pScene *core.AiScene) (err error) {
	// signature
	err = ai.Discard(44)
	if err != nil {
		return err
	}
	versionMajor, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	versionMinor, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if versionMinor != ASSBIN_VERSION_MINOR || versionMajor != ASSBIN_VERSION_MAJOR {
		return errors.New("Invalid version, data format not compatible!")
	}
	versionRevision, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	_ = versionRevision
	compileFlags, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	_ = compileFlags
	shortened, err := ai.GetUInt16()
	if err != nil {
		return err
	}
	compressed, err := ai.GetUInt16()
	if err != nil {
		return err
	}

	if shortened != 0 {
		return errors.New("Shortened binaries are not supported!")
	}

	err = ai.Discard(256) // original filename
	if err != nil {
		return err
	}
	err = ai.Discard(128)
	if err != nil {
		return err
	} // options
	err = ai.Discard(64)
	if err != nil {
		return err
	} // padding

	if compressed != 0 {
		//TODO
	}
	return ai.ReadBinaryScene(pScene)
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiVector3D() (v *common.AiVector3D, err error) {
	v = &common.AiVector3D{}
	v.X, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	v.Y, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	v.Z, err = ai.ReadAiReal()
	return v, err
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiColor4D() (c *common.AiColor4D, err error) {
	c = &common.AiColor4D{}
	c.R, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	c.G, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	c.B, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	c.A, err = ai.ReadAiReal()
	return c, err
}
func (ai *AssBinImporter) ReadAiReal() (v float32, err error) {
	return ai.GetFloat32()
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiQuaternion() (v *common.AiQuaternion, err error) {
	v = &common.AiQuaternion{}
	v.W, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	v.X, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	v.Y, err = ai.ReadAiReal()
	if err != nil {
		return nil, err
	}
	v.Z, err = ai.ReadAiReal()
	return v, err
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiString() (string, error) {
	length, err := ai.GetUInt32()
	if err != nil {
		return "", err
	}
	if length > 0 {
		data, err := ai.GetNBytes(int(length))
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", nil
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiVertexWeight() (*common.AiVertexWeight, error) {
	var (
		w   common.AiVertexWeight
		err error
	)
	w.VertexId, err = ai.GetUInt32()
	if err != nil {
		return &w, err
	}
	w.Weight, err = ai.ReadAiReal()
	return &w, err
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiMatrix4x4() (m *common.AiMatrix4x4, err error) {
	m = &common.AiMatrix4x4{}
	for i := 0; i < 4; i++ {
		for i2 := 0; i2 < 4; i2++ {
			v, err := ai.ReadAiReal()
			if err != nil {
				return m, err
			}
			m.Set(i, i2, v)
		}
	}
	return m, nil
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiVectorKey() (v *common.AiVectorKey, err error) {
	v = &common.AiVectorKey{}
	v.Time, err = ai.GetFloat64()
	if err != nil {
		return nil, err
	}
	v.Value, err = ai.ReadAiVector3D()
	if err != nil {
		return nil, err
	}
	return v, err
}

// -----------------------------------------------------------------------------------

func (ai *AssBinImporter) ReadAiQuatKey() (v *common.AiQuatKey, err error) {
	v = &common.AiQuatKey{}
	v.Time, err = ai.GetFloat64()
	if err != nil {
		return nil, err
	}
	v.Value, err = ai.ReadAiQuaternion()
	return v, err
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryNode(parent *core.AiNode) (node *core.AiNode, err error) {
	magic, err := ai.GetUInt32()
	if err != nil {
		return node, err
	}
	if magic != ASSBIN_CHUNK_AINODE {
		return node, errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return nil, err
	}
	node = core.NewAiNode("")
	node.Name, err = ai.ReadAiString()
	if err != nil {
		return node, err
	}
	node.Transformation, err = ai.ReadAiMatrix4x4()
	if err != nil {
		return node, err
	}
	numChildren, err := ai.GetUInt32()
	if err != nil {
		return node, err
	}
	numMeshes, err := ai.GetUInt32()
	if err != nil {
		return node, err
	}
	nb_metadata, err := ai.GetUInt32()
	if err != nil {
		return node, err
	}
	if parent != nil {
		node.Parent = parent
	}

	if numMeshes > 0 {
		node.Meshes = make([]int32, numMeshes)
		for i := uint32(0); i < numMeshes; i++ {
			t, err := ai.GetUInt32()
			if err != nil {
				return node, err
			}
			node.Meshes[i] = int32(t)
		}
	}

	if numChildren > 0 {
		node.Children = make([]*core.AiNode, numChildren)
		for i := uint32(0); i < numChildren; i++ {
			node.Children[i], err = ai.ReadBinaryNode(node)
			if err != nil {
				return node, err
			}
		}
	}

	if nb_metadata > 0 {
		if node.MetaData == nil {
			node.MetaData = &core.AiMetadata{}
		}
		node.MetaData.Keys = make([]string, nb_metadata)
		node.MetaData.Values = make([]*core.AiMetadataEntry, nb_metadata)
		for i := range node.MetaData.Keys {
			node.MetaData.Keys[i], err = ai.ReadAiString()
			t, err := ai.GetUInt16()
			if err != nil {
				return node, err
			}
			node.MetaData.Values[i].Type = core.AiMetadataType(t)
			var data any
			switch node.MetaData.Values[i].Type {
			case core.AI_BOOL:
				data, err = ai.GetNBytes(1)
				if err != nil {
					return node, err
				}
			case core.AI_INT32:
				data, err = ai.GetUInt32()
				if err != nil {
					return node, err
				}
			case core.AI_UINT64:
				data, err = ai.GetUInt64()
				if err != nil {
					return node, err
				}
			case core.AI_FLOAT:
				data, err = ai.GetFloat32()
				if err != nil {
					return node, err
				}
			case core.AI_DOUBLE:
				data, err = ai.GetFloat64()
				if err != nil {
					return node, err
				}
			case core.AI_AISTRING:
				data, err = ai.ReadAiString()
				if err != nil {
					return node, err
				}
			case core.AI_AIVECTOR3D:
				data, err = ai.ReadAiVector3D()
				if err != nil {
					return node, err
				}
			default:
				return node, fmt.Errorf("invalid Ai Type:%v", node.MetaData.Values[i].Type)
			}

			node.MetaData.Values[i].Data = data
		}
	}
	return node, nil
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryBone(b *core.AiBone) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AIBONE {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	b.Name, err = ai.ReadAiString()
	if err != nil {
		return err
	}
	numWeights, err := ai.GetInt32()
	if err != nil {
		return err
	}
	b.OffsetMatrix, err = ai.ReadAiMatrix4x4()
	if err != nil {
		return err
	}
	// for the moment we write dumb min/max values for the bones, too.
	// maybe I'll add a better, hash-like solution later
	// else write as usual
	b.Weights = make([]*common.AiVertexWeight, numWeights)
	for i := range b.Weights {
		b.Weights[i], err = ai.ReadAiVertexWeight()
		if err != nil {
			return err
		}
	}
	return nil
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) fitsIntoUI16(NumVertices int) bool {
	return (NumVertices < (1 << 16))
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryMesh(mesh *core.AiMesh) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AIMESH {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	tmp, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	mesh.PrimitiveTypes = core.AiPrimitiveType(tmp)
	numVertices, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numFaces, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numBones, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	t, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	mesh.MaterialIndex = int32(t)
	// first of all, write bits for all existent vertex components
	c, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if c&ASSBIN_MESH_HAS_POSITIONS != 0 {
		// else write as usual
		mesh.Vertices = make([]*common.AiVector3D, numVertices)
		for i := range mesh.Vertices {
			mesh.Vertices[i], err = ai.ReadAiVector3D()
			if err != nil {
				return err
			}
		}
	}
	if c&ASSBIN_MESH_HAS_NORMALS != 0 {
		// else write as usual
		mesh.Normals = make([]*common.AiVector3D, numVertices)
		for i := range mesh.Normals {
			mesh.Normals[i], err = ai.ReadAiVector3D()
			if err != nil {
				return err
			}
		}
	}
	if c&ASSBIN_MESH_HAS_TANGENTS_AND_BITANGENTS != 0 {
		// else write as usual
		mesh.Tangents = make([]*common.AiVector3D, numVertices)
		for i := range mesh.Tangents {
			mesh.Tangents[i], err = ai.ReadAiVector3D()
			if err != nil {
				return err
			}
		}
		mesh.Bitangents = make([]*common.AiVector3D, numVertices)
		for i := range mesh.Bitangents {
			mesh.Bitangents[i], err = ai.ReadAiVector3D()
			if err != nil {
				return err
			}
		}
	}
	for n := 0; n < core.AI_MAX_NUMBER_OF_COLOR_SETS; n++ {
		if c&uint32(ASSBIN_MESH_HAS_COLOR(int32(n))) == 0 {
			break
		}

		// else write as usual
		mesh.Colors[n] = make([]*common.AiColor4D, numVertices)
		for i := range mesh.Colors[n] {
			mesh.Colors[n][i], err = ai.ReadAiColor4D()
			if err != nil {
				return err
			}
		}
	}
	for n := uint32(0); n < core.AI_MAX_NUMBER_OF_TEXTURECOORDS; n++ {
		if c&uint32(ASSBIN_MESH_HAS_TEXCOORD(int32(n))) == 0 {
			break
		}

		// write number of UV components
		mesh.NumUVComponents[n], err = ai.GetUInt32()

		// else write as usual
		mesh.TextureCoords[n] = make([]*common.AiVector3D, numVertices)
		for i := range mesh.TextureCoords[n] {
			mesh.TextureCoords[n][i], err = ai.ReadAiVector3D()
			if err != nil {
				return err
			}
		}
	}

	// write faces. There are no floating-point calculations involved
	// in these, so we can write a simple hash over the face data
	// to the dump file. We generate a single 32 Bit hash for 512 faces
	// using Assimp's standard hashing function.
	// else write as usual
	// if there are less than 2^16 vertices, we can simply use 16 bit integers ...
	mesh.Faces = make([]*core.AiFace, numFaces)
	for i := uint32(0); i < numFaces; i++ {
		mesh.Faces[i] = core.NewAiFace()
		f := mesh.Faces[i]
		numIndices, err := ai.GetUInt16()
		if err != nil {
			return err
		}
		f.Indices = make([]uint32, numIndices)

		for a := uint16(0); a < numIndices; a++ {
			// Check if unsigned  short ( 16 bit  ) are big enough for the indices
			if ai.fitsIntoUI16(int(numVertices)) {
				t, err := ai.GetUInt16()
				if err != nil {
					return err
				}
				f.Indices[a] = uint32(t)
			} else {
				f.Indices[a], err = ai.GetUInt32()
				if err != nil {
					return err
				}
			}
		}
	}

	// write bones
	if numBones > 0 {
		mesh.Bones = make([]*core.AiBone, numBones)
		for a := uint32(0); a < numBones; a++ {
			mesh.Bones[a] = &core.AiBone{}
			err = ai.ReadBinaryBone(mesh.Bones[a])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryMaterialProperty(prop *core.AiMaterialProperty) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AIMATERIALPROPERTY {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	prop.Key, err = ai.ReadAiString()
	if err != nil {
		return err
	}
	prop.Semantic, err = ai.GetUInt32()
	if err != nil {
		return err
	}
	prop.Index, err = ai.GetUInt32()
	if err != nil {
		return err
	}
	dataLength, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	_, err = ai.GetUInt32()
	if err != nil {
		return err
	}

	data, err := ai.GetNBytes(int(dataLength))
	if err != nil {
		return err
	}
	return parseAiMaterialByKey(prop, data)
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryMaterial(mat *core.AiMaterial) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AIMATERIAL {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	numProperties, err := ai.GetUInt32()
	if err != nil {
		return err
	}

	if numProperties > 0 {
		if len(mat.Properties) > 0 {
			mat.Properties = mat.Properties[:0]
		}
		mat.Properties = make([]*core.AiMaterialProperty, numProperties)
		for i := uint32(0); i < numProperties; i++ {
			mat.Properties[i] = &core.AiMaterialProperty{}
			err = ai.ReadBinaryMaterialProperty(mat.Properties[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryNodeAnim(nd *core.AiNodeAnim) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AINODEANIM {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	nd.NodeName, err = ai.ReadAiString()
	if err != nil {
		return err
	}
	numPositionKeys, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numRotationKeys, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numScalingKeys, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	t, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	t1, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	nd.PreState = core.AiAnimBehaviour(t)
	nd.PostState = core.AiAnimBehaviour(t1)

	if numPositionKeys > 0 {
		nd.PositionKeys = make([]*common.AiVectorKey, numPositionKeys)
		for i := range nd.PositionKeys {
			nd.PositionKeys[i], err = ai.ReadAiVectorKey()
			if err != nil {
				return err
			}
		}
	}
	if numRotationKeys > 0 {
		// else write as usual
		nd.RotationKeys = make([]*common.AiQuatKey, numRotationKeys)
		for i := range nd.RotationKeys {
			nd.RotationKeys[i], err = ai.ReadAiQuatKey()
			if err != nil {
				return err
			}
		}
	}
	if numScalingKeys > 0 {
		// else write as usual
		nd.ScalingKeys = make([]*common.AiVectorKey, numScalingKeys)
		for i := range nd.ScalingKeys {
			nd.ScalingKeys[i], err = ai.ReadAiVectorKey()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryAnim(anim *core.AiAnimation) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AIANIMATION {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	anim.Name, err = ai.ReadAiString()
	if err != nil {
		return err
	}
	anim.Duration, err = ai.GetFloat64()
	if err != nil {
		return err
	}
	anim.TicksPerSecond, err = ai.GetFloat64()
	if err != nil {
		return err
	}
	numChannels, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if numChannels > 0 {
		anim.Channels = make([]*core.AiNodeAnim, numChannels)
		for a := uint32(0); a < numChannels; a++ {
			anim.Channels[a] = &core.AiNodeAnim{}
			err = ai.ReadBinaryNodeAnim(anim.Channels[a])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryTexture(tex *core.AiTexture) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AITEXTURE {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	tex.Width, err = ai.GetUInt32()
	if err != nil {
		return err
	}
	tex.Height, err = ai.GetUInt32()
	if err != nil {
		return err
	}
	tex.AchFormatHint, err = ai.GetNBytes(core.HINTMAXTEXTURELEN - 1)
	if err != nil {
		return err
	}
	if tex.Height == 0 {
		tex.PcData = make([]*core.AiTexel, tex.Width)
	} else {
		tex.PcData = make([]*core.AiTexel, tex.Width*tex.Height)
	}

	for i := range tex.PcData {
		tex.PcData[i] = &core.AiTexel{}
		err = ai.ReadBinaryTex(tex.PcData[i])
		if err != nil {
			return err
		}
	}
	return nil
}
func (ai *AssBinImporter) ReadBinaryTex(t *core.AiTexel) (err error) {
	t.R, err = ai.GetUInt8()
	if err != nil {
		return err
	}
	t.G, err = ai.GetUInt8()
	if err != nil {
		return err
	}
	t.B, err = ai.GetUInt8()
	if err != nil {
		return err
	}
	t.A, err = ai.GetUInt8()
	return err
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryLight(l *core.AiLight) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AILIGHT {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	l.Name, err = ai.ReadAiString()
	if err != nil {
		return err
	}
	t, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	l.Type = core.AiLightSourceType(t)

	l.Position, err = ai.ReadAiVector3D()
	if err != nil {
		return err
	}
	l.Direction, err = ai.ReadAiVector3D()
	if err != nil {
		return err
	}
	l.Up, err = ai.ReadAiVector3D()
	if err != nil {
		return err
	}

	if l.Type != core.AiLightSource_DIRECTIONAL {
		l.AttenuationConstant, err = ai.GetFloat32()
		if err != nil {
			return err
		}
		l.AttenuationLinear, err = ai.GetFloat32()
		if err != nil {
			return err
		}
		l.AttenuationQuadratic, err = ai.GetFloat32()
		if err != nil {
			return err
		}
	}

	l.ColorDiffuse, err = ai.ReadAiColor3D()
	if err != nil {
		return err
	}
	l.ColorSpecular, err = ai.ReadAiColor3D()
	if err != nil {
		return err
	}
	l.ColorAmbient, err = ai.ReadAiColor3D()
	if err != nil {
		return err
	}
	if l.Type == core.AiLightSource_SPOT {
		l.AngleInnerCone, err = ai.GetFloat32()
		if err != nil {
			return err
		}
		l.AngleOuterCone, err = ai.GetFloat32()
		if err != nil {
			return err
		}
	}
	return nil
}
func (ai *AssBinImporter) ReadAiColor3D() (v *common.AiColor3D, err error) {
	v = &common.AiColor3D{}
	v.R, err = ai.GetFloat32()
	if err != nil {
		return v, err
	}
	v.G, err = ai.GetFloat32()
	if err != nil {
		return v, err
	}
	v.B, err = ai.GetFloat32()
	if err != nil {
		return v, err
	}
	return v, err
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryCamera(cam *core.AiCamera) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AICAMERA {
		return errors.New("Magic chunk identifiers are wrong!")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	cam.Name, err = ai.ReadAiString()
	if err != nil {
		return err
	}
	cam.Position, err = ai.ReadAiVector3D()
	if err != nil {
		return err
	}
	cam.LookAt, err = ai.ReadAiVector3D()
	if err != nil {
		return err
	}
	cam.Up, err = ai.ReadAiVector3D()
	if err != nil {
		return err
	}
	cam.HorizontalFOV, err = ai.GetFloat32()
	if err != nil {
		return err
	}
	cam.ClipPlaneNear, err = ai.GetFloat32()
	if err != nil {
		return err
	}
	cam.ClipPlaneFar, err = ai.GetFloat32()
	if err != nil {
		return err
	}
	cam.Aspect, err = ai.GetFloat32()
	return err
}

// -----------------------------------------------------------------------------------
func (ai *AssBinImporter) ReadBinaryScene(scene *core.AiScene) error {
	magic, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	if magic != ASSBIN_CHUNK_AISCENE {
		return errors.New("magic chunk identifiers are wrong")
	}
	err = ai.Discard(4)
	if err != nil {
		return err
	}
	scene.Flags, err = ai.GetUInt32()
	if err != nil {
		return err
	}
	numMeshes, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numMaterials, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numAnimations, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numTextures, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numLights, err := ai.GetUInt32()
	if err != nil {
		return err
	}
	numCameras, err := ai.GetUInt32()
	if err != nil {
		return err
	}

	// Read node graph
	scene.RootNode, err = ai.ReadBinaryNode(nil)
	if err != nil {
		return err
	}
	// Read all meshes
	if numMeshes > 0 {
		scene.Meshes = make([]*core.AiMesh, numMeshes)
		for i := uint32(0); i < numMeshes; i++ {
			scene.Meshes[i] = core.NewAiMesh()
			err = ai.ReadBinaryMesh(scene.Meshes[i])
			if err != nil {
				return err
			}
		}
	}

	// Read materials
	if numMaterials > 0 {
		scene.Materials = make([]*core.AiMaterial, numMaterials)
		for i := uint32(0); i < numMaterials; i++ {
			scene.Materials[i] = &core.AiMaterial{}
			err = ai.ReadBinaryMaterial(scene.Materials[i])
			if err != nil {
				return err
			}
		}
	}

	// Read all animations
	if numAnimations > 0 {
		scene.Animations = make([]*core.AiAnimation, numAnimations)
		for i := uint32(0); i < numAnimations; i++ {
			scene.Animations[i] = &core.AiAnimation{}
			err = ai.ReadBinaryAnim(scene.Animations[i])
			if err != nil {
				return err
			}
		}
	}

	// Read all textures
	if numTextures > 0 {
		scene.Textures = make([]*core.AiTexture, numTextures)
		for i := uint32(0); i < numTextures; i++ {
			scene.Textures[i] = core.NewAiTexture()
			err = ai.ReadBinaryTexture(scene.Textures[i])
			if err != nil {
				return err
			}
		}
	}

	// Read lights
	if numLights > 0 {
		scene.Lights = make([]*core.AiLight, numLights)
		for i := uint32(0); i < numLights; i++ {
			scene.Lights[i] = core.NewAiLight()
			err = ai.ReadBinaryLight(scene.Lights[i])
			if err != nil {
				return err
			}
		}
	}

	// Read cameras
	if numCameras > 0 {
		scene.Cameras = make([]*core.AiCamera, numCameras)

		for i := uint32(0); i < numCameras; i++ {
			scene.Cameras[i] = core.NewAiCamera()
			err = ai.ReadBinaryCamera(scene.Cameras[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
