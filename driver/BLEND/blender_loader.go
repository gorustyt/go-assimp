package BLEND

import (
	"assimp/common"
	"assimp/common/logger"
	"assimp/common/reader"
	"assimp/core"
	"assimp/driver/base/iassimp"
	"errors"
	"fmt"
	"math"
	"sort"
	"unsafe"
)

var Desc = core.AiImporterDesc{
	"Blender 3D Importer (http://www.blender3d.org)",
	"",
	"",
	"No animation support yet",
	core.AiImporterFlags_SupportBinaryFlavour,
	0,
	0,
	2,
	50,
	[]string{".blend"},
	"BLENDER",
}

type BlenderImporter struct {
	reader.StreamReader
	modifier_cache *BlenderModifierShowcase
}

func (b *BlenderImporter) checkMagic() ([]byte, bool, error) {
	magic, err := b.Peek(7)
	if err != nil {
		return magic, false, err
	}
	if string(magic[:]) == Desc.Magic {
		return magic, true, nil
	}
	return magic, false, nil
}

func (b *BlenderImporter) ParseMagic() error {
	magic, ok, err := b.checkMagic()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	// Check for presence of the gzip header. If yes, assume it is a
	// compressed blend file and try uncompressing it, else fail. This is to
	// avoid uncompressing random files which our loader might end up with.
	if magic[0] != 0x1f || magic[1] != 0x8b {
		return errors.New("BLENDER magic bytes are missing, couldn't find GZIP header either")
	}

	logger.Info("Found no BLENDER magic word but a GZIP header, might be a compressed file")
	if magic[2] != 8 {
		return errors.New("Unsupported GZIP compression method")
	}
	err = b.ResetGzipReader()
	if err != nil {
		return err
	}
	magic, ok, err = b.checkMagic()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("magic not found ")
	}
	return nil
}

func (b *BlenderImporter) CanRead(checkSig bool) bool {
	return b.ParseMagic() == nil
}

func (b *BlenderImporter) Read(pScene *core.AiScene) (err error) {
	err = b.ParseMagic()
	if err != nil {
		return err
	}
	err = b.Discard(7)
	if err != nil {
		return err
	}
	file := NewFileDatabase(b.StreamReader)
	buffer, err := b.GetNBytes(1)
	if err != nil {
		return err
	}
	file.i64bit = buffer[0] == '-'
	buffer, err = b.GetNBytes(1)
	if err != nil {
		return err
	}
	file.little = buffer[0] == 'v'
	buffer, err = b.GetNBytes(3)
	if err != nil {
		return err
	}
	b.ChangeBytesOrder(file.little)
	logger.InfoF("Blender version is:%v.%v (64bit:%v , little endian:%v)",
		string(buffer[0]), uint8(buffer[0])+1,
		file.i64bit,
		file.little)
	err = b.ParseBlendFile(file)
	if err != nil {
		return err
	}
	var scene Scene
	err = b.ExtractScene(&scene, file)
	if err != nil {
		return err
	}
	return b.ConvertBlendFile(pScene, &scene, file)
}

func (b *BlenderImporter) ExtractScene(out *Scene, file *FileDatabase) error {
	var block *FileBlockHead
	it, ok := file.dna.indices["Scene"]
	if !ok {
		return errors.New("there is no `Scene` structure record")
	}

	ss := file.dna.structures[it]

	// we need a scene somewhere to start with.
	for _, bl := range file.entries {

		// Fix: using the DNA index is more reliable to locate scenes
		//if (bl.id == "SC") {

		if int(bl.dna_index) == it {
			block = bl
			break
		}
	}

	if block == nil {
		return errors.New("there is not a single `Scene` record to load")
	}
	file.SetCurPos(block.start)
	err := out.Convert(file, ss)
	if err != nil {
		return err
	}

	logger.InfoF(
		"(Stats) Fields read:%v, pointers resolved:%v, cache hits: %v, cached objects: ", file.stats().fields_read,
		file.stats().pointers_resolved,
		file.stats().cache_hits,
		file.stats().cached_objects)
	return nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) ParseSubCollection(in *Scene, root *core.AiNode, collection *Collection, conv_data *ConversionData) error {
	//*Object
	var root_objects = common.NewQueue[*Object]()
	// Count number of objects
	for cur := collection.gobject.first.(*CollectionObject); cur != nil; cur = cur.next {
		if cur.ob != nil {
			root_objects.PushBack(cur.ob)
		}
	}
	//
	var root_children = common.NewQueue[*Collection]()
	// Count number of child nodes
	for cur := collection.children.first.(*CollectionChild); cur != nil; cur = cur.next {
		if cur.collection != nil {
			root_children.PushBack(cur.collection)
		}
	}
	root.Children = make([]*core.AiNode, root_objects.Size()+root_children.Size())

	for i := 0; i < root_objects.Size(); i++ {
		var err error
		root.Children[i], err = b.ConvertNode(in, root_objects.Index(i), conv_data, common.NewAiMatrix4x4Identify())
		if err != nil {
			return err
		}
		root.Children[i].Parent = root
	}

	// For each subcollection create a new node to represent it
	iterator := root_objects.Size()
	for cur := collection.children.first.(*CollectionChild); cur != nil; cur = cur.next {
		if cur.collection != nil {
			root.Children[iterator] = core.NewAiNode(cur.collection.id.name[:2]) // skip over the name prefix 'OB'
			root.Children[iterator].Parent = root
			err := b.ParseSubCollection(in, root.Children[iterator], cur.collection, conv_data)
			if err != nil {
				return err
			}
		}
		iterator += 1
	}
	return nil
}

func (b *BlenderImporter) ConvertBlendFile(out *core.AiScene, in *Scene, file *FileDatabase) error {
	conv := newConversionData(file)
	out.RootNode = core.NewAiNode("<BlenderRoot>")
	root := out.RootNode
	// Iterate over all objects directly under master_collection,
	// If in.master_collection == null, then we're parsing something older.
	if in.master_collection != nil {
		err := b.ParseSubCollection(in, root, in.master_collection, conv)
		if err != nil {
			return err
		}
	} else {
		no_parents := common.NewQueue[*Object]()
		for cur := in.base.first.(*Base); cur != nil; cur = cur.next {
			if cur.object != nil {
				if cur.object.parent == nil {
					no_parents.PushBack(cur.object)
				} else {
					conv.objects = append(conv.objects, cur.object)
				}
			}
		}
		for cur := in.basact; cur != nil; cur = cur.next {
			if cur.object != nil {
				if cur.object.parent != nil {
					conv.objects = append(conv.objects, cur.object)
				}
			}
		}

		if no_parents.Size() == 0 {
			return errors.New("expected at least one object with no parent")
		}

		root.Children = make([]*core.AiNode, no_parents.Size())
		for i := 0; i < len(root.Children); i++ {
			var err error
			root.Children[i], err = b.ConvertNode(in, no_parents.Index(i), conv, common.NewAiMatrix4x4Identify())
			if err != nil {
				return err
			}
			root.Children[i].Parent = root
		}
	}

	err := b.BuildMaterials(conv)
	if err != nil {
		return err
	}
	if len(conv.meshes) > 0 {
		out.Meshes = make([]*core.AiMesh, len(conv.meshes))
		copy(out.Meshes, conv.meshes)
		conv.meshes = conv.meshes[:0]
	}

	if len(conv.lights) > 0 {
		out.Lights = make([]*core.AiLight, len(conv.lights))
		copy(out.Lights, conv.lights)
		conv.lights = conv.lights[:0]
	}

	if len(conv.cameras) > 0 {
		out.Cameras = make([]*core.AiCamera, len(conv.cameras))
		copy(out.Cameras, conv.cameras)
		conv.cameras = conv.cameras[:0]
	}

	if len(conv.materials) > 0 {
		out.Materials = make([]*core.AiMaterial, len(conv.materials))
		copy(out.Materials, conv.materials)
		conv.materials = conv.materials[:0]
	}

	if len(conv.textures) > 0 {
		out.Textures = make([]*core.AiTexture, len(conv.textures))
		copy(out.Textures, conv.textures)
		conv.textures = conv.textures[:0]
	}

	// acknowledge that the scene might come out incomplete
	// by Assimp's definition of `complete`: blender scenes
	// can consist of thousands of cameras or lights with
	// not a single mesh between them.
	if len(out.Meshes) == 0 {
		out.Flags |= core.AI_SCENE_FLAGS_INCOMPLETE
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) ConvertNode(in *Scene, obj *Object, conv_data *ConversionData, parentTransform *common.AiMatrix4x4) (*core.AiNode, error) {
	children := common.NewQueue[*Object]()
	for it := 0; it < len(conv_data.objects); it++ {
		object := conv_data.objects[it]
		if object.parent == obj {
			children.PushBack(object)
			var objs []*Object
			for i, v := range conv_data.objects {
				if i == it {
					continue
				}
				objs = append(objs, v)
			}
			conv_data.objects = objs
			it++
			continue
		}
	}

	node := core.NewAiNode(obj.id.name[2:]) // skip over the name prefix 'OB'
	if obj.data != nil {
		switch obj.Type {
		case Type_EMPTY:
			break // do nothing

			// supported object types
		case Type_MESH:
			old := len(conv_data.meshes)

			err := b.CheckActualType(obj.data, "Mesh")
			if err != nil {
				return nil, err
			}
			err = b.ConvertMesh(in, obj, obj.data.(*Mesh), conv_data, conv_data.meshes)
			if err != nil {
				return nil, err
			}
			if len(conv_data.meshes) > old {
				node.Meshes = make([]int, len(conv_data.meshes)-old)
				for i := 0; i < len(node.Meshes); i++ {
					node.Meshes[i] = i + old
				}
			}
			break
		case Type_LAMP:
			err := b.CheckActualType(obj.data, "Lamp")
			if err != nil {
				return nil, err
			}
			mesh := b.ConvertLight(in, obj, obj.data.(*Lamp), conv_data)
			if mesh != nil {
				conv_data.lights = append(conv_data.lights, mesh)
			}
			break
		case Type_CAMERA:
			err := b.CheckActualType(obj.data, "Camera")
			if err != nil {
				return nil, err
			}
			mesh := b.ConvertCamera(in, obj, obj.data.(*Camera), conv_data)
			if mesh != nil {
				conv_data.cameras = append(conv_data.cameras, mesh)
			}
			break

			// unsupported object types / log, but do not break
		case Type_CURVE:
			b.NotSupportedObjectType(obj, "Curve")
			break
		case Type_SURF:
			b.NotSupportedObjectType(obj, "Surface")
			break
		case Type_FONT:
			b.NotSupportedObjectType(obj, "Font")
			break
		case Type_MBALL:
			b.NotSupportedObjectType(obj, "MetaBall")
			break
		case Type_WAVE:
			b.NotSupportedObjectType(obj, "Wave")
			break
		case Type_LATTICE:
			b.NotSupportedObjectType(obj, "Lattice")
			break

			// invalid or unknown type
		default:
			break
		}
	}

	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			node.Transformation.Set(y, x, obj.obmat[x][y])
		}
	}

	m := parentTransform
	m = m.Inverse()

	node.Transformation = m.MulAiMatrix4x4(node.Transformation)

	if children.Size() != 0 {
		node.Children = make([]*core.AiNode, children.Size())
		nd := 0
		for i := 0; i < children.Size(); i++ {
			nobj := children.Index(i)
			var err error
			node.Children[nd], err = b.ConvertNode(in, nobj, conv_data, node.Transformation.MulAiMatrix4x4(parentTransform))
			if err != nil {
				return nil, err
			}
			node.Children[nd].Parent = node
			nd++
		}
	}

	// apply modifiers
	err := b.modifier_cache.ApplyModifiers(*node, conv_data, in, obj)
	return node, err
}

func (b *BlenderImporter) ParseBlendFile(file *FileDatabase) error {
	b.ResetData()
	dna_reader := NewDNAParser(file, b.StreamReader)
	var dna *DNA
	parser := NewSectionParser(b.StreamReader, file.i64bit)
	for {
		err := parser.Next()
		if err != nil {
			return err
		}
		head := parser.current.clone()
		if head.id == "ENDB" {
			break // only valid end of the file
		} else if head.id == "DNA1" {
			err = dna_reader.Parse()
			if err != nil {
				return err
			}
			dna = dna_reader.GetDNA()
			continue
		}
		if err != nil {
			return err
		}
		file.entries = append(file.entries, head)
	}
	if dna == nil {
		return errors.New("SDNA not found")
	}
	sort.Slice(file.entries, func(i, j int) bool {
		return file.entries[i].address.val < file.entries[j].address.val
	})
	return nil
}
func NewBlenderImporter(data []byte) (iassimp.Loader, error) {
	r, err := reader.NewFileStreamReader(data)
	if err != nil {
		return nil, err
	}
	return &BlenderImporter{StreamReader: r}, nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) ResolveImage(out *core.AiMaterial, mat *Material, tex *MTex, img *Image, conv_data *ConversionData) error {
	var name string

	// check if the file contents are bundled with the BLEND file
	if img.packedfile != nil {
		name += "*"
		var curTex = core.NewAiTexture()
		conv_data.textures = append(conv_data.textures, curTex)
		// usually 'img.name' will be the original file name of the embedded textures,
		// so we can extract the file extension from it.
		s := len(img.name)
		e := s
		for s >= len(img.name) && img.name[s] != '.' {
			s--
		}
		tmp := img.name[s:]
		curTex.AchFormatHint[0] = tmp[1]
		if s+1 > e {
			curTex.AchFormatHint[0] = '\x00'
		}
		curTex.AchFormatHint[1] = tmp[2]
		if s+2 > e {
			curTex.AchFormatHint[1] = '\x00'
		}
		curTex.AchFormatHint[2] = tmp[3]
		if s+3 > e {
			curTex.AchFormatHint[2] = '\x00'
		}
		curTex.AchFormatHint[3] = '\x00'

		// tex.mHeight = 0;
		curTex.Width = int(img.packedfile.size)
		ch := make([]byte, curTex.Width)

		conv_data.db.SetCurPos(img.packedfile.data.val)
		ch, err := conv_data.db.GetNBytes(curTex.Width)
		if err != nil {
			return err
		}

		curTex.PcData = *(*[]*core.AiTexel)(unsafe.Pointer(&ch))

		logger.InfoF("Reading embedded texture, original file was %v", img.name)
	} else {
		name = img.name
	}

	texture_type := core.AiTextureType_UNKNOWN
	map_type := tex.mapto

	if map_type&MapType_COL != 0 {
		texture_type = core.AiTextureType_DIFFUSE
	} else if map_type&MapType_NORM != 0 {
		if tex.tex.imaflag&ImageFlags_NORMALMAP != 0 {
			texture_type = core.AiTextureType_NORMALS
		} else {
			texture_type = core.AiTextureType_HEIGHT
		}
		out.AddFloat32PropertyVar(core.AI_MATKEY_BUMPSCALING, tex.norfac)
	} else if map_type&MapType_COLSPEC != 0 {
		texture_type = core.AiTextureType_SPECULAR
	} else if map_type&MapType_COLMIR != 0 {
		texture_type = core.AiTextureType_REFLECTION
	} else if map_type&MapType_REF != 0 {

	} else if map_type&MapType_SPEC != 0 {
		texture_type = core.AiTextureType_SHININESS
	} else if map_type&MapType_EMIT != 0 {
		texture_type = core.AiTextureType_EMISSIVE
	} else if map_type&MapType_ALPHA != 0 {
	} else if map_type&MapType_HAR != 0 {

	} else if map_type&MapType_RAYMIRR != 0 {

	} else if map_type&MapType_TRANSLU != 0 {

	} else if map_type&MapType_AMB != 0 {
		texture_type = core.AiTextureType_AMBIENT
	} else if map_type&MapType_DISPLACE != 0 {
		texture_type = core.AiTextureType_DISPLACEMENT
	} else if map_type&MapType_WARP != 0 {

	}

	out.AddStringPropertyVar(core.AI_MATKEY_TEXTURE(texture_type, int(conv_data.next_texture[texture_type])), name)
	conv_data.next_texture[texture_type]++
	return nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) AddSentinelTexture(out *core.AiMaterial, mat *Material, tex *MTex, conv_data *ConversionData) {
	name := fmt.Sprintf("Procedural,num=%v,type=%s", conv_data.sentinel_cnt, tex.tex.Type.GetTextureTypeDisplayString())
	conv_data.sentinel_cnt++
	out.AddStringPropertyVar(core.AI_MATKEY_TEXTURE_DIFFUSE(int(conv_data.next_texture[core.AiTextureType_DIFFUSE])), name)
	conv_data.next_texture[core.AiTextureType_DIFFUSE]++

}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) ResolveTexture(out *core.AiMaterial, mat *Material, tex *MTex, conv_data *ConversionData) error {
	rtex := tex.tex
	if rtex == nil || rtex.Type == 0 {
		return nil
	}

	// We can't support most of the texture types because they're mostly procedural.
	// These are substituted by a dummy texture.
	var dispnam = ""
	switch rtex.Type {
	// these are listed in blender's UI
	case Type_CLOUDS:
	case Type_WOOD:
	case Type_MARBLE:
	case Type_MAGIC:
	case Type_BLEND:
	case Type_STUCCI:
	case Type_NOISE:
	case Type_PLUGIN:
	case Type_MUSGRAVE:
	case Type_VORONOI:
	case Type_DISTNOISE:
	case Type_ENVMAP:

		// these do no appear in the UI, why?
	case Type_POINTDENSITY:
	case Type_VOXELDATA:

		logger.WarnF("Encountered a texture with an unsupported type: ", dispnam)
		b.AddSentinelTexture(out, mat, tex, conv_data)
		break

	case Type_IMAGE:
		if rtex.ima == nil {
			return errors.New("A texture claims to be an Image, but no image reference is given")
		}
		err := b.ResolveImage(out, mat, tex, rtex.ima, conv_data)
		if err != nil {
			return err
		}
		break

	default:
		return errors.New("invalid Type")
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) BuildDefaultMaterial(conv_data *ConversionData) error {
	// add a default material if necessary
	index := -1
	for _, mesh := range conv_data.meshes {
		if mesh.MaterialIndex == -1 {

			if index == -1 {
				// Setup a default material.
				p := &Material{ElemBase: &ElemBase{}}
				if len(core.AI_DEFAULT_MATERIAL_NAME) < len(p.id.name)-2 {
					return errors.New("invalid length")
				}
				p.id.name = p.id.name[:2] + core.AI_DEFAULT_MATERIAL_NAME
				// Note: MSVC11 does not zero-initialize Material here, although it should.
				// Thus all relevant fields should be explicitly initialized. We cannot add
				// a default constructor to Material since the DNA codegen does not support
				// parsing it.
				p.b = 0.6
				p.g = p.b
				p.r = p.g
				p.specb = 0.6
				p.specg = p.specb
				p.specr = p.specg
				p.ambb = 0.0
				p.ambg = p.ambb
				p.ambr = p.ambg
				p.mirb = 0.0
				p.mirg = p.mirb
				p.mirr = p.mirg
				p.emit = 0.
				p.alpha = 0.
				p.har = 0

				index = conv_data.materials_raw.Size()
				conv_data.materials_raw.PushBack(p)
				logger.Info("Adding default material")
			}
			mesh.MaterialIndex = index
		}
	}
	return nil
}

func (b *BlenderImporter) AddBlendParams(result *core.AiMaterial, source *Material) {
	diffuseColor := common.NewAiColor3D(source.r, source.g, source.b)
	result.AddAiColor3DPropertyVar(core.NewAiMaterialProperty("$mat.blend.diffuse.color", 0, 0), diffuseColor)

	diffuseIntensity := source.ref
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.diffuse.intensity", 0, 0), diffuseIntensity)

	diffuseShader := source.diff_shader
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.diffuse.shader", 0, 0), int64(diffuseShader))

	diffuseRamp := 0
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.diffuse.ramp", 0, 0), int64(diffuseRamp))

	specularColor := common.NewAiColor3D(source.specr, source.specg, source.specb)
	result.AddAiColor3DPropertyVar(core.NewAiMaterialProperty("$mat.blend.specular.color", 0, 0), specularColor)

	specularIntensity := source.spec
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.specular.intensity", 0, 0), specularIntensity)

	specularShader := source.spec_shader
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.specular.shader", 0, 0), int64(specularShader))

	specularRamp := 0
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.specular.ramp", 0, 0), int64(specularRamp))

	specularHardness := source.har
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.specular.hardness", 0, 0), int64(specularHardness))

	transparencyUse := 0
	if source.mode&MA_TRANSPARENCY != 0 {
		transparencyUse = 1
	}
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.use", 0, 0), int64(transparencyUse))

	transparencyMethod := 0
	if source.mode&MA_RAYTRANSP != 0 {
		transparencyMethod = 2
	} else {
		if source.mode&MA_ZTRANSP != 0 {
			transparencyMethod = 1
		}
	}
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.method", 0, 0), int64(transparencyMethod))

	transparencyAlpha := source.alpha
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.alpha", 0, 0), transparencyAlpha)

	transparencySpecular := source.spectra
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.specular", 0, 0), transparencySpecular)

	transparencyFresnel := source.fresnel_tra
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.fresnel", 0, 0), transparencyFresnel)

	transparencyBlend := source.fresnel_tra_i
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.blend", 0, 0), transparencyBlend)

	transparencyIor := source.ang
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.ior", 0, 0), transparencyIor)

	transparencyFilter := source.filter
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.filter", 0, 0), transparencyFilter)

	transparencyFalloff := source.tx_falloff
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.falloff", 0, 0), transparencyFalloff)

	transparencyLimit := source.tx_limit
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.limit", 0, 0), transparencyLimit)

	transparencyDepth := source.ray_depth_tra
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.depth", 0, 0), int64(transparencyDepth))

	transparencyGlossAmount := source.gloss_tra
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.glossAmount", 0, 0), transparencyGlossAmount)

	transparencyGlossThreshold := source.adapt_thresh_tra
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.glossThreshold", 0, 0), transparencyGlossThreshold)

	transparencyGlossSamples := source.samp_gloss_tra
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.transparency.glossSamples", 0, 0), int64(transparencyGlossSamples))

	mirrorUse := 0
	if source.mode&MA_RAYMIRROR != 0 {
		mirrorUse = 1
	}
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.use", 0, 0), int64(mirrorUse))

	mirrorReflectivity := source.ray_mirror
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.reflectivity", 0, 0), mirrorReflectivity)

	mirrorColor := common.NewAiColor3D(source.mirr, source.mirg, source.mirb)
	result.AddAiColor3DPropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.color", 0, 0), mirrorColor)

	mirrorFresnel := source.fresnel_mir
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.fresnel", 0, 0), mirrorFresnel)

	mirrorBlend := source.fresnel_mir_i
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.blend", 0, 0), mirrorBlend)

	mirrorDepth := source.ray_depth
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.depth", 0, 0), int64(mirrorDepth))

	mirrorMaxDist := source.dist_mir
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.maxDist", 0, 0), mirrorMaxDist)

	mirrorFadeTo := source.fadeto_mir
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.fadeTo", 0, 0), int64(mirrorFadeTo))

	mirrorGlossAmount := source.gloss_mir
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.glossAmount", 0, 0), mirrorGlossAmount)

	mirrorGlossThreshold := source.adapt_thresh_mir
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.glossThreshold", 0, 0), mirrorGlossThreshold)

	mirrorGlossSamples := source.samp_gloss_mir
	result.AddInt64PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.glossSamples", 0, 0), int64(mirrorGlossSamples))

	mirrorGlossAnisotropic := source.aniso_gloss_mir
	result.AddFloat32PropertyVar(core.NewAiMaterialProperty("$mat.blend.mirror.glossAnisotropic", 0, 0), mirrorGlossAnisotropic)
}

func (b *BlenderImporter) BuildMaterials(conv_data *ConversionData) error {

	err := b.BuildDefaultMaterial(conv_data)
	if err != nil {
		return err
	}
	for i := 0; i < conv_data.materials_raw.Size(); i++ {
		mat := conv_data.materials_raw.Index(i)
		// reset per material global counters
		for i := 0; i < len(conv_data.next_texture)/4; i++ {
			conv_data.next_texture[i] = 0
		}

		mout := core.AiMaterial{}
		conv_data.materials = append(conv_data.materials, &mout)
		// For any new material field handled here, the default material above must be updated with an appropriate default value.

		// set material name
		name := mat.id.name[2:] // skip over the name prefix 'MA'
		mout.AddStringPropertyVar(core.AI_MATKEY_NAME, name)

		// basic material colors
		col := common.NewAiColor3D(mat.r, mat.g, mat.b)
		if mat.r != 0 || mat.g != 0 || mat.b != 0 {

			// Usually, zero diffuse color means no diffuse color at all in the equation.
			// So we omit this member to express this intent.
			mout.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_DIFFUSE, col)

			if mat.emit != 0 {
				emit_col := common.NewAiColor3D(mat.emit*mat.r, mat.emit*mat.g, mat.emit*mat.b)
				mout.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_EMISSIVE, emit_col)
			}
		}

		col = common.NewAiColor3D(mat.specr, mat.specg, mat.specb)
		mout.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_SPECULAR, col)

		// is hardness/shininess set?
		if mat.har != 0 {
			har := mat.har
			mout.AddInt64PropertyVar(core.AI_MATKEY_SHININESS, int64(har))
		}

		col = common.NewAiColor3D(mat.ambr, mat.ambg, mat.ambb)
		mout.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_AMBIENT, col)

		// is mirror enabled?
		if mat.mode&MA_RAYMIRROR != 0 {
			ray_mirror := mat.ray_mirror
			mout.AddFloat32PropertyVar(core.AI_MATKEY_REFLECTIVITY, ray_mirror)
		}

		col = common.NewAiColor3D(mat.mirr, mat.mirg, mat.mirb)
		mout.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_REFLECTIVE, col)

		for i := 0; i < int(unsafe.Sizeof(mat.mtex)/unsafe.Sizeof(mat.mtex[0])); i++ {
			if mat.mtex[i] == nil {
				continue
			}

			err = b.ResolveTexture(&mout, mat, mat.mtex[i], conv_data)
			if err != nil {
				return err
			}
		}

		b.AddBlendParams(&mout, mat)
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) CheckActualType(dt IElemBase, check string) error {
	if dt.GetDnaType() != check {
		return fmt.Errorf("Expected object at%v  to be of type `%v`, but it claims to be a `%v`instead", dt, check,
			dt.GetDnaType())
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) NotSupportedObjectType(obj *Object, Type string) {
	logger.WarnF("Object `%v` - type is unsupported: `%V`, skipping", obj.id.name, Type)
}

var (
	TODO_FIX_BMESH_CONVERSION = false
)

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) ConvertMesh(in *Scene, obj *Object, mesh *Mesh,
	conv_data *ConversionData, temp []*core.AiMesh) error {
	// TODO: Resolve various problems with BMesh triangulation before re-enabling.
	//       See issues #400, #373, #318  #315 and #132.
	//if TODO_FIX_BMESH_CONVERSION{
	//	mesh := newBMeshConverter();
	//	if (BMeshConverter.ContainsBMesh()) {
	//		mesh = BMeshConverter.TriangulateBMesh();
	//	}
	//}
	if (mesh.totface == 0 && mesh.totloop == 0) || mesh.totvert == 0 {
		return nil
	}

	// some sanity checks
	if int(mesh.totface) > len(mesh.mface) {
		return errors.New("Number of faces is larger than the corresponding array")
	}

	if int(mesh.totvert) > len(mesh.mvert) {
		return errors.New("Number of vertices is larger than the corresponding array")
	}

	if int(mesh.totloop) > len(mesh.mloop) {
		return errors.New("Number of vertices is larger than the corresponding array")
	}
	getVoVn := func(o *core.AiMesh) (vo, vn *common.AiVector3D) {
		vo = common.NewAiVector3D()
		o.Vertices = append(o.Vertices, vo)
		vn = common.NewAiVector3D()
		o.Normals = append(o.Normals, vn)
		return
	}
	// collect per-submesh numbers
	per_mat := map[int]int{}
	per_mat_verts := map[int]int{}
	for i := 0; i < int(mesh.totface); i++ {

		mf := mesh.mface[i]
		per_mat[int(mf.mat_nr)]++
		tmp := 3
		if mf.v4 != 0 {
			tmp = 4
		}
		per_mat_verts[int(mf.mat_nr)] += tmp

	}

	for i := 0; i < int(mesh.totpoly); i++ {
		mp := mesh.mpoly[i]
		per_mat[int(mp.mat_nr)]++
		per_mat_verts[int(mp.mat_nr)] += int(mp.totloop)
	}

	// ... and allocate the corresponding meshes
	old := len(temp)
	mat_num_to_mesh_idx := map[int]int{}
	for k, v := range per_mat {

		mat_num_to_mesh_idx[k] = len(temp)
		var out = core.NewAiMesh()
		temp = append(temp, out)
		out.Vertices = make([]*common.AiVector3D, per_mat_verts[k])
		out.Normals = make([]*common.AiVector3D, per_mat_verts[k])

		//out.mNumFaces = 0
		//out.mNumVertices = 0
		out.Faces = make([]*core.AiFace, v)

		// all sub-meshes created from this mesh are named equally. this allows
		// curious users to recover the original adjacency.
		out.Name = mesh.id.name[2:]
		// skip over the name prefix 'ME'

		// resolve the material reference and add this material to the set of
		// output materials. The (temporary) material index is the index
		// of the material entry within the list of resolved materials.
		if len(mesh.mat) > 0 {

			if k >= len(mesh.mat) {
				return errors.New("material index is out of range")
			}

			mat := mesh.mat[k]
			has := conv_data.materials_raw.FindIndex(mat, func(v1, v2 *Material) bool {
				return v1 == v2
			})
			if has != -1 {
				out.MaterialIndex = has
			} else {
				out.MaterialIndex = conv_data.materials_raw.Size()
				conv_data.materials_raw.PushBack(mat)
			}
		} else {
			out.MaterialIndex = -1
		}

	}

	for i := 0; i < int(mesh.totface); i++ {

		mf := mesh.mface[i]

		out := temp[mat_num_to_mesh_idx[int(mf.mat_nr)]]
		var f = core.NewAiFace()
		out.Faces = append(out.Faces, f)
		tmp := 3
		if mf.v4 != 0 {
			tmp = 4
		}
		f.Indices = make([]int, tmp)

		vo, vn := getVoVn(out)
		// XXX we can't fold this easily, because we are restricted
		// to the member names from the BLEND file (v1,v2,v3,v4)
		// which are assigned by the genblenddna.py script and
		// cannot be changed without breaking the entire
		// import process.

		if mf.v1 >= mesh.totvert {
			return errors.New("Vertex index v1 out of range")
		}
		v := mesh.mvert[mf.v1]
		vo.X = v.co[0]
		vo.Y = v.co[1]
		vo.Z = v.co[2]
		vn.X = v.no[0]
		vn.Y = v.no[1]
		vn.Z = v.no[2]
		f.Indices[0] = len(out.Vertices)
		vo, vn = getVoVn(out)

		//  if (f.mNumIndices >= 2) {
		if mf.v2 >= mesh.totvert {
			return errors.New("Vertex index v2 out of range")
		}
		v = mesh.mvert[mf.v2]
		vo.X = v.co[0]
		vo.Y = v.co[1]
		vo.Z = v.co[2]
		vn.X = v.no[0]
		vn.Y = v.no[1]
		vn.Z = v.no[2]
		f.Indices[1] = len(out.Vertices)

		vo, vn = getVoVn(out)
		if mf.v3 >= mesh.totvert {
			return errors.New("Vertex index v3 out of range")
		}
		//  if (f.mNumIndices >= 3) {
		v = mesh.mvert[mf.v3]
		vo.X = v.co[0]
		vo.Y = v.co[1]
		vo.Z = v.co[2]
		vn.X = v.no[0]
		vn.Y = v.no[1]
		vn.Z = v.no[2]
		f.Indices[2] = len(out.Vertices)

		vo, vn = getVoVn(out)

		if mf.v4 >= mesh.totvert {
			return errors.New("Vertex index v4 out of range")
		}
		//  if (f.mNumIndices >= 4) {
		if mf.v4 != 0 {
			v = mesh.mvert[mf.v4]
			vo.X = v.co[0]
			vo.Y = v.co[1]
			vo.Z = v.co[2]
			vn.X = v.no[0]
			vn.Y = v.no[1]
			vn.Z = v.no[2]
			f.Indices[3] = len(out.Vertices)
			vo, vn = getVoVn(out)

			out.PrimitiveTypes |= int(core.AiPrimitiveType_POLYGON)
		} else {
			out.PrimitiveTypes |= int(core.AiPrimitiveType_TRIANGLE)
		}

		//  }
		//  }
		//  }
	}

	for i := 0; i < int(mesh.totpoly); i++ {

		mf := mesh.mpoly[i]

		out := temp[mat_num_to_mesh_idx[int(mf.mat_nr)]]
		var f = core.NewAiFace()
		out.Faces = append(out.Faces, f)
		f.Indices = make([]int, mf.totloop)
		vo, vn := getVoVn(out)
		// XXX we can't fold this easily, because we are restricted
		// to the member names from the BLEND file (v1,v2,v3,v4)
		// which are assigned by the genblenddna.py script and
		// cannot be changed without breaking the entire
		// import process.
		for j := 0; j < int(mf.totloop); j++ {
			loop := mesh.mloop[int(mf.loopstart)+j]

			if loop.v >= mesh.totvert {
				return errors.New("Vertex index out of range")
			}

			v := mesh.mvert[loop.v]

			vo.X = v.co[0]
			vo.Y = v.co[1]
			vo.Z = v.co[2]
			vn.X = v.no[0]
			vn.Y = v.no[1]
			vn.Z = v.no[2]
			f.Indices[j] = len(out.Vertices)
			vo, vn = getVoVn(out)
		}
		if mf.totloop == 3 {
			out.PrimitiveTypes |= int(core.AiPrimitiveType_TRIANGLE)
		} else {
			out.PrimitiveTypes |= int(core.AiPrimitiveType_POLYGON)
		}
	}

	// TODO should we create the TextureUVMapping map in Convert<Material> to prevent redundant processing?

	// create texture <. uvname mapping for all materials
	// key is texture number, value is data *

	// key is material number, value is the TextureUVMapping for the material
	matTexUvMappings := map[uint32]map[uint32]*MLoopUV{}
	for m := 0; m < len(mesh.mat); m++ {
		// get material by index
		pMat := mesh.mat[m]
		texuv := map[uint32]*MLoopUV{}
		maxTex := unsafe.Sizeof(pMat.mtex) / unsafe.Sizeof(pMat.mtex[0])
		for t := 0; t < int(maxTex); t++ {
			if pMat.mtex[t] != nil && pMat.mtex[t].uvname[0] != 0 {
				// get the CustomData layer for given uvname and correct type
				pLoop := getCustomDataLayerData(mesh.ldata, CD_MLOOPUV, pMat.mtex[t].uvname)
				if pLoop != nil {
					texuv[uint32(t)] = pLoop.(*MLoopUV)
				}
			}
		}
		if len(texuv) > 0 {
			matTexUvMappings[uint32(m)] = texuv
		}
	}

	// collect texture coordinates, they're stored in a separate per-face buffer
	if len(mesh.mtface) != 0 || len(mesh.mloopuv) != 0 {
		if int(mesh.totface) > len(mesh.mtface) {
			return errors.New("Number of UV faces is larger than the corresponding UV face array (#1)")
		}
		for it := old; it != len(temp); it++ {
			if 0 == len(temp[it].Vertices) {
				return errors.New("invalid type ")
			}
			if 0 == len(temp[it].Faces) {
				return errors.New("invalid type ")
			}
			itMatTexUvMapping, ok := matTexUvMappings[uint32(temp[it].MaterialIndex)]
			if !ok {
				// default behaviour like before
				temp[it].TextureCoords[0] = make([]*common.AiVector3D, len(temp[it].Vertices))
			} else {
				// create texture coords for every mapped tex
				for i := 0; i < len(itMatTexUvMapping); i++ {
					temp[it].TextureCoords[i] = make([]*common.AiVector3D, len(temp[it].Vertices))
				}
			}
		}

		for i := 0; i < int(mesh.totface); i++ {
			v := mesh.mtface[i]
			out := temp[mat_num_to_mesh_idx[int(mesh.mface[i].mat_nr)]]
			f := core.NewAiFace()
			out.Faces = append(out.Faces, f)
			vo1 := common.NewAiVector3D()
			out.TextureCoords[0] = append(out.TextureCoords[0], vo1)
			for j := 0; j < len(f.Indices); j++ {
				vo1.X = v.uv[j][0]
				vo1.Y = v.uv[j][1]
			}
		}

		for i := 0; i < int(mesh.totpoly); i++ {
			v := mesh.mpoly[i]
			out := temp[mat_num_to_mesh_idx[int(v.mat_nr)]]

			f := core.NewAiFace()
			out.Faces = append(out.Faces, f)
			itMatTexUvMapping, ok := matTexUvMappings[uint32(v.mat_nr)]
			if !ok {
				// old behavior
				vo1 := common.NewAiVector3D()
				out.TextureCoords[0] = append(out.TextureCoords[0], vo1)
				for j := 0; j < len(f.Indices); j++ {
					uv := mesh.mloopuv[int(v.loopstart)+j]
					vo1.X = uv.uv[0]
					vo1.Y = uv.uv[1]
				}
			} else {
				// create textureCoords for every mapped tex
				for m := 0; m < len(itMatTexUvMapping); m++ {
					tm := matTexUvMappings[uint32(m)]
					vo1 := common.NewAiVector3D()
					out.TextureCoords[0] = append(out.TextureCoords[0], vo1)
					j := 0
					for ; j < len(f.Indices); j++ {
						uv := tm[uint32(int(v.loopstart)+j)]
						vo1.X = uv.uv[0]
						vo1.Y = uv.uv[1]
					}
					// only update written mNumVertices in last loop
					// TODO why must the numVertices be incremented here?
					if m == len(itMatTexUvMapping)-1 {
						//TODO out.NumVertices += j
					}
				}
			}
		}
	}

	// collect texture coordinates, old-style (marked as deprecated in current blender sources)
	if len(mesh.tface) > 0 {
		if int(mesh.totface) > len(mesh.tface) {
			return errors.New("Number of faces is larger than the corresponding UV face array (#2)")
		}
		for it := old; it != len(temp); it++ {
			if 0 == len(temp[it].Vertices) {
				return errors.New("invalid ")
			}
			if 0 == len(temp[it].Faces) {
				return errors.New("invalid ")
			}

			temp[it].TextureCoords[0] = make([]*common.AiVector3D, len(temp[it].Vertices))
		}

		for i := 0; i < int(mesh.totface); i++ {
			v := mesh.tface[i]

			out := temp[mat_num_to_mesh_idx[int(mesh.mface[i].mat_nr)]]
			f := core.NewAiFace()
			out.Faces = append(out.Faces, f)
			vo1 := common.NewAiVector3D()
			out.TextureCoords[0] = append(out.TextureCoords[0], vo1)
			for j := 0; j < len(f.Indices); j++ {
				vo1.X = float32(v.uv[j][0])
				vo1.Y = float32(v.uv[j][1])
			}
		}
	}

	// collect vertex colors, stored separately as well
	if len(mesh.mcol) != 0 || len(mesh.mloopcol) != 0 {
		if int(mesh.totface) > (len(mesh.mcol) / 4) {
			return errors.New("Number of faces is larger than the corresponding color face array")
		}
		for it := old; it != len(temp); it++ {
			if 0 == len(temp[it].Vertices) {
				return errors.New("invalid ")
			}
			if 0 == len(temp[it].Faces) {
				return errors.New("invalid ")
			}

			temp[it].Colors[0] = make([]*common.AiColor4D, len(temp[it].Vertices))
		}

		for i := 0; i < int(mesh.totface); i++ {

			out := temp[mat_num_to_mesh_idx[int(mesh.mface[i].mat_nr)]]
			f := core.NewAiFace()
			out.Faces = append(out.Faces, f)
			vo := common.NewAiColor4D0()
			out.Colors[0] = append(out.Colors[0], vo)
			for n := 0; n < len(f.Indices); n++ {
				col := mesh.mcol[(i<<2)+n]
				vo.R = float32(col.r)
				vo.G = float32(col.g)
				vo.B = float32(col.b)
				vo.A = float32(col.a)
			}
			for n := len(f.Indices); n < 4; n++ {
			}
		}

		for i := 0; i < int(mesh.totpoly); i++ {
			v := mesh.mpoly[i]
			out := temp[mat_num_to_mesh_idx[int(v.mat_nr)]]
			f := core.NewAiFace()
			out.Faces = append(out.Faces, f)
			vo := common.NewAiColor4D0()
			out.Colors[0] = append(out.Colors[0], vo)
			scaleZeroToOne := float32(1. / 255.)
			for j := 0; j < len(f.Indices); j++ {
				col := mesh.mloopcol[int(v.loopstart)+j]
				vo.R = float32(col.r) * scaleZeroToOne
				vo.G = float32(col.g) * scaleZeroToOne
				vo.B = float32(col.b) * scaleZeroToOne
				vo.A = float32(col.a) * scaleZeroToOne
			}
		}
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) ConvertCamera(in *Scene, obj *Object, cam *Camera, conv_data *ConversionData) *core.AiCamera {
	out := core.NewAiCamera()
	out.Name = obj.id.name[2:]
	out.Position = common.NewAiVector3D3(0., 0., 0.)
	out.Up = common.NewAiVector3D3(0., 1., 0.)
	out.LookAt = common.NewAiVector3D3(0., 0., -1.)
	if cam.sensor_x != 0 && cam.lens != 0 {
		out.HorizontalFOV = 2. * float32(math.Atan2(float64(cam.sensor_x), float64(2.*cam.lens)))
	}
	out.ClipPlaneNear = cam.clipsta
	out.ClipPlaneFar = cam.clipend

	return out
}

// ------------------------------------------------------------------------------------------------
func (b *BlenderImporter) ConvertLight(in *Scene, obj *Object, lamp *Lamp, conv_data *ConversionData) (out *core.AiLight) {
	out = core.NewAiLight()
	out.Name = obj.id.name[2:]

	switch lamp.Type {
	case LampType_Local:
		out.Type = core.AiLightSource_POINT
		break
	case LampType_Spot:
		out.Type = core.AiLightSource_SPOT

		// blender orients directional lights as facing toward -z
		out.Direction = common.NewAiVector3D3(0., 0., -1.)
		out.Up = common.NewAiVector3D3(0., 1., 0.)
		out.AngleInnerCone = lamp.spotsize * (1.0 - lamp.spotblend)
		out.AngleOuterCone = lamp.spotsize
		break
	case LampType_Sun:
		out.Type = core.AiLightSource_DIRECTIONAL

		// blender orients directional lights as facing toward -z
		out.Direction = common.NewAiVector3D3(0., 0., -1.)
		out.Up = common.NewAiVector3D3(0., 1., 0.)
		break

	case LampType_Area:
		out.Type = core.AiLightSource_AREA
		if lamp.area_shape == 0 {
			out.Size = common.NewAiVector2D(lamp.area_size, lamp.area_size)
		} else {
			out.Size = common.NewAiVector2D(lamp.area_size, lamp.area_sizey)
		}

		// blender orients directional lights as facing toward -z
		out.Direction = common.NewAiVector3D3(0., 0., -1.)
		out.Up = common.NewAiVector3D3(0., 1., 0.)
		break

	default:
		break
	}

	out.ColorAmbient = common.NewAiColor3D(lamp.r, lamp.g, lamp.b).Mul(lamp.energy)
	out.ColorSpecular = common.NewAiColor3D(lamp.r, lamp.g, lamp.b).Mul(lamp.energy)
	out.ColorDiffuse = common.NewAiColor3D(lamp.r, lamp.g, lamp.b).Mul(lamp.energy)

	// If default values are supplied, compute the coefficients from light's max distance
	// Read this: https://imdoingitwrong.wordpress.com/2011/01/31/light-attenuation/
	//
	if lamp.constant_coefficient == 1.0 && lamp.linear_coefficient == 0.0 && lamp.quadratic_coefficient == 0.0 && lamp.dist > 0.0 {
		out.AttenuationConstant = 1.0
		out.AttenuationLinear = 2.0 / lamp.dist
		out.AttenuationQuadratic = 1.0 / (lamp.dist * lamp.dist)
	} else {
		out.AttenuationConstant = lamp.constant_coefficient
		out.AttenuationLinear = lamp.linear_coefficient
		out.AttenuationQuadratic = lamp.quadratic_coefficient
	}

	return out
}
