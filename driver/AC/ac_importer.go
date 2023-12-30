package AC

import (
	"assimp/common"
	"assimp/common/config"
	"assimp/common/logger"
	"assimp/common/reader"
	"assimp/core"
	"assimp/driver/base/iassimp"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

var (
	Desc = core.AiImporterDesc{
		"AC3D Importer",
		"",
		"",
		"",
		core.AiImporterFlags_SupportTextFlavour,
		0,
		0,
		0,
		0,
		[]string{".ac", ".acc", ".ac3d"},
		"AC3D",
	}
)

func NewAC3DImporter(data []byte) (iassimp.Loader, error) {
	r, err := reader.NewFileLineReader(data)
	if err != nil {
		return nil, err
	}
	im := &AC3DImporter{LineReader: r}
	return im, nil
}
func (ac *AC3DImporter) InitConfig(cfg *config.Config) {
	ac.configEvalSubdivision = cfg.SubDivision
}
func (ac *AC3DImporter) LoadObjectSection() (objects []*Object, err error) {
	if !ac.HasPrefix("OBJECT") {
		return objects, nil
	}
	ac.NumMeshes++
	name, err := ac.MustOneKeyString("OBJECT")
	if err != nil {
		return objects, err
	}
	obj := newObject()
	objects = append(objects, obj)
	var light *core.AiLight
	switch name {
	case "light":
		light = core.NewAiLight()
		ac.Lights = append(ac.Lights, light)
		light.Type = core.AiLightSource_POINT
		light.ColorSpecular = common.NewAiColor3D(1.0, 1.0, 1.0)
		light.ColorDiffuse = light.ColorSpecular
		light.AttenuationConstant = 1.0

		light.Name = fmt.Sprintf("ACLight_%s", len(ac.Lights)-1)
		obj.name = light.Name
		logger.Info("AC3D: Light source encountered")
		obj.Type = Light

	case "group":
		obj.Type = Group
	case "world":
		obj.Type = World
	default:
		obj.Type = Poly
	}
	for !ac.EOF() {
		if ac.HasPrefix("kids") {
			num, err := ac.MustOneKeyInt("kids")
			if err != nil {
				return objects, err
			}
			for i := 0; i < num; i++ {
				children, err := ac.LoadObjectSection()
				if err != nil {
					return nil, err
				}
				obj.children = append(obj.children, children...)
			}
			return objects, nil
		} else if ac.HasPrefix("name") {
			obj.name, err = ac.MustOneKeyString("name")
			if err != nil {
				return nil, err
			}
			obj.name = common.ClearQuotationMark(obj.name)
			if light != nil {
				light.Name = obj.name
			}
			continue
		} else if ac.HasPrefix("texture") {
			texture, err := ac.MustOneKeyString("texture")
			if err != nil {
				return nil, err
			}
			obj.textures = append(obj.textures, common.ClearQuotationMark(texture))
			continue
		} else if ac.HasPrefix("texrep") {
			obj.texRepeat, err = ac.NextKeyAiVector2d("texrep")
			if err != nil {
				return nil, err
			}
			if obj.texRepeat.X == 0 || obj.texRepeat.Y == 0 {
				obj.texRepeat = &common.AiVector2D{1, 1}
			}
		} else if ac.HasPrefix("texoff") {
			obj.texOffset, err = ac.NextKeyAiVector2d("texoff")
			if err != nil {
				return nil, err
			}
		} else if ac.HasPrefix("rot") {
			obj.rotation, err = ac.NextKeyAiMatrix3x3("loc")
			if err != nil {
				return nil, err
			}
		} else if ac.HasPrefix("loc") {
			obj.translation, err = ac.NextKeyAiVector3d("loc")
			if err != nil {
				return nil, err
			}
		} else if ac.HasPrefix("subdiv") {
			obj.subDiv, err = ac.MustOneKeyInt("subdiv")
			if err != nil {
				return nil, err
			}
			continue
		} else if ac.HasPrefix("crease") {
			obj.crease, err = ac.MustOneKeyFloat32("crease")
			if err != nil {
				return nil, err
			}
			continue
		} else if ac.HasPrefix("numvert") {
			vs, err := ac.NextLineVector3("numvert")
			if err != nil {
				return nil, err
			}
			obj.vertices = append(obj.vertices, vs...)
		} else if ac.HasPrefix("numsurf") {
			Q3DWorkAround := false

			num, err := ac.MustOneKeyInt("numsurf")
			if err != nil {
				return nil, err
			}
			for i := 0; i < num; i++ {
				// FIX: this can occur for some files - Quick 3D for
				// example writes no surf chunks
				if !ac.HasPrefix("SURF") { //TODO
					panic("not impl")
					Q3DWorkAround = true
				}

				var surf Surface
				obj.surfaces = append(obj.surfaces, &surf)
				surf.flags, err = ac.MustOneKeyInt("SURF", true)
				if err != nil {
					return nil, err
				}
				for {

					if ac.EOF() {
						return objects, errors.New("AC3D: Unexpected EOF: surface is incomplete")
					}
					if ac.HasPrefix("mat") {
						surf.mat, err = ac.MustOneKeyInt("mat")
						if err != nil {
							return nil, err
						}
						continue
					} else if ac.HasPrefix("refs") {
						if Q3DWorkAround {

						}
						vs, err := ac.NextLineVector3("refs")
						if err != nil {
							return nil, err
						}
						for _, v := range vs {
							var entry SurfaceEntry
							surf.entries = append(surf.entries, &entry)
							entry.First = int(v.X)
							entry.Second = common.AiVector2D{v.Y, v.Z}
						}
						obj.numRefs += len(vs)
					} else { // make sure the line is processed a second time
						break
					}
					ac.NextLine()
				}

			}
		}
		ac.NextLine()
	}
	return objects, err
}
func (ac *AC3DImporter) ConvertMaterial(object *Object,
	matSrc *Material,
	matDest *core.AiMaterial) {
	var s string
	if matSrc.name != "" {
		s = matSrc.name
		matDest.AddStringPropertyVar(core.AI_MATKEY_NAME, s)
	}
	if len(object.textures) > 0 {
		s = object.textures[0]
		matDest.AddStringPropertyVar(core.AI_MATKEY_TEXTURE_DIFFUSE(0), s)
		if object.texRepeat.X != 1 || object.texRepeat.Y != 1 || object.texOffset.X > 0 || object.texOffset.Y > 0 {
			var transform core.AiUVTransform
			transform.Scaling = object.texRepeat
			transform.Translation = object.texOffset
			matDest.AddAiUVTransformPropertyVar(core.AI_MATKEY_UVTRANSFORM_DIFFUSE(0), transform)
		}
	}

	matDest.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_DIFFUSE, matSrc.rgb)
	matDest.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_AMBIENT, matSrc.amb)
	matDest.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_EMISSIVE, matSrc.emis)
	matDest.AddAiColor3DPropertyVar(core.AI_MATKEY_COLOR_SPECULAR, matSrc.spec)
	n := int64(-1)
	if matSrc.shin != 0 {
		n = int64(core.AiShadingMode_Phong)
		matDest.AddFloat32PropertyVar(core.AI_MATKEY_SHININESS, matSrc.shin)
	} else {
		n = int64(core.AiShadingMode_Gouraud)
	}
	matDest.AddInt64PropertyVar(core.AI_MATKEY_SHADING_MODEL, n)
	f := 1.0 - matSrc.trans
	matDest.AddFloat32PropertyVar(core.AI_MATKEY_OPACITY, f)
}

func (ac *AC3DImporter) ConvertObjectSection(object *Object, meshes *[]*core.AiMesh,
	outMaterials *[]*core.AiMaterial,
	materials *[]*Material, parent *core.AiNode) (node *core.AiNode) {
	node = core.NewAiNode("")
	node.Parent = parent
	numMeshs := 0
	if len(object.vertices) > 0 {
		if len(object.surfaces) == 0 || object.numRefs == 0 {
			/* " An object with 7 vertices (no surfaces, no materials defined).
			   This is a good way of getting point data into AC3D.
			   The Vertex.create convex-surface/object can be used on these
			   vertices to 'wrap' a 3d shape around them "
			   (http://www.opencity.info/html/ac3dfileformat.html)

			   therefore: if no surfaces are defined return point data only
			*/
			logger.InfoF("AC3D: No surfaces defined in object definition a point list is returned")
			mesh := core.NewAiMesh()
			*meshes = append(*meshes, mesh)
			numVertices := len(object.vertices)
			numFaces := numVertices
			mesh.Faces = make([]*core.AiFace, numFaces)
			for i := range mesh.Faces {
				mesh.Faces[i] = &core.AiFace{}
			}
			mesh.Vertices = make([]*common.AiVector3D, numVertices)
			faces := 0
			verts := 0
			for i := 0; i < numVertices; i++ {
				mesh.Vertices[verts] = object.vertices[i]
				mesh.Faces[faces].Indices = make([]uint32, 1)
				mesh.Faces[faces].Indices[0] = uint32(i)
				i++
				faces++
				verts++
			}
			// use the primary material in this case. this should be the
			// default material if all objects of the file contain points
			// and no faces.
			mesh.MaterialIndex = 0
			var tmp core.AiMaterial
			*outMaterials = append(*outMaterials, &tmp)
			ac.ConvertMaterial(object, (*materials)[0], &tmp)
		} else {
			var needMat = make([]*common.Pair[int, int], len(*materials))
			for i := range needMat {
				needMat[i] = &common.Pair[int, int]{}
			}
			for _, v := range object.surfaces {
				idx := v.mat
				if idx >= len(needMat) {
					logger.Warn("AC3D: material index is out of range")
					idx = 0
				}
				if len(v.entries) == 0 {
					logger.Warn("AC3D: surface her zero vertex references")
				}
				// validate all vertex indices to make sure we won't crash here
				for it2 := 0; it2 < len(v.entries); it2++ {
					if v.entries[it2].First >= len(object.vertices) {
						logger.Warn("AC3D: Invalid vertex reference")
						v.entries[it2].First = 0
					}
				}
				if needMat[idx].First == 0 {
					numMeshs++
				}
				switch SurfaceType(v.GetType()) {
				// closed line
				case ClosedLine:
					needMat[idx].First += len(v.entries)
					needMat[idx].Second += len(v.entries) << 1
					break
					// unclosed line
				case OpenLine:
					needMat[idx].First += len(v.entries) - 1
					needMat[idx].Second += (len(v.entries) - 1) << 1
					break
					// triangle strip
				case TriangleStrip:
					needMat[idx].First += len(v.entries) - 2
					needMat[idx].Second += (len(v.entries) - 2) * 3
					break
				default:
					// Coerce unknowns to a polygon and warn
					logger.WarnF("AC3D: The type flag of a surface is unknown: %v", v.flags)
					v.flags &= ^(int(Mask))
					fallthrough
					// polygon
				case Polygon:
					// the number of faces increments by one, the number
					// of vertices by surface.numref.
					needMat[idx].First++
					needMat[idx].Second += len(v.entries)
				}
			}
			node.Meshes = make([]int32, numMeshs)
			mat := 0
			oldm := len(*meshes)
			cit := 0
			cend := len(needMat)
			for cit != cend {
				citv := needMat[cit]
				if citv.First == 0 {
					cit++
					mat++
					continue
				}
				mesh := core.NewAiMesh()
				*meshes = append(*meshes, mesh)
				mesh.MaterialIndex = int32(len(*outMaterials))
				var tmpMaterial core.AiMaterial
				*outMaterials = append(*outMaterials, &tmpMaterial)
				ac.ConvertMaterial(object, (*materials)[mat], &tmpMaterial)
				// allocate storage for vertices and normals
				numFaces := citv.First
				if numFaces == 0 {
					logger.FatalF("AC3D: No faces")
				} else if numFaces*int(unsafe.Sizeof(core.AiFace{})) > 256*1024*1024 {
					logger.FatalF("AC3D: Too many faces, would run out of memory")
				}

				mesh.Faces = make([]*core.AiFace, numFaces)
				for i := range mesh.Faces {
					mesh.Faces[i] = core.NewAiFace()
				}
				faces := 0
				vertices := 0
				numVertices := citv.Second
				if numVertices == 0 {
					logger.FatalF("AC3D: No vertices")
				} else if numVertices*4*3 > 256*1024*1024 {
					logger.FatalF("AC3D: Too many vertices, would run out of memory")
				}
				mesh.Vertices = make([]*common.AiVector3D, numVertices)
				cur := 0
				// allocate UV coordinates, but only if the texture name for the
				// surface is not empty
				uv := -1
				if len(object.textures) != 0 {
					uv = 0
					mesh.TextureCoords[0] = make([]*common.AiVector3D, len(mesh.Vertices))
					for i := range mesh.TextureCoords[0] {
						mesh.TextureCoords[0][i] = common.NewAiVector3D3(0, 0, 0)
					}
					mesh.NumUVComponents[0] = 2
				}
				uvs := mesh.TextureCoords[0]
				for it := 0; it < len(object.surfaces); it++ {
					itv := object.surfaces[it]
					if mat == itv.mat {
						src := object.surfaces[it]

						// closed polygon
						Type := SurfaceType(itv.GetType())
						if Type == Polygon {
							face := mesh.Faces[faces]
							faces++
							if 0 != len(src.entries) {
								face.Indices = make([]uint32, len(src.entries))
								for i := 0; i < len(face.Indices); i++ {
									entry := src.entries[i]
									face.Indices[i] = uint32(cur)
									cur++
									// copy vertex positions
									if (vertices - len(mesh.Vertices)) >= len(mesh.Vertices) {
										logger.FatalF("AC3D: Invalid number of vertices")
									}
									mesh.Vertices[vertices] = object.vertices[entry.First].Add(object.translation)

									// copy texture coordinates
									if uv != -1 {
										uvs[uv].X = entry.Second.X
										uvs[uv].Y = entry.Second.Y
										uv++
									}
									vertices++
								}
							}
						} else if Type == TriangleStrip {
							for i := 0; i < len(src.entries)-2; i++ {
								entry1 := src.entries[i]
								entry2 := src.entries[i+1]
								entry3 := src.entries[i+2]

								// skip degenerate triangles
								if object.vertices[entry1.First] == object.vertices[entry2.First] ||
									object.vertices[entry1.First] == object.vertices[entry3.First] ||
									object.vertices[entry2.First] == object.vertices[entry3.First] {
									continue
								}

								face := mesh.Faces[faces]
								faces++
								face.Indices = make([]uint32, 3)
								face.Indices[0] = uint32(cur)
								cur++
								face.Indices[1] = uint32(cur)
								cur++
								face.Indices[2] = uint32(cur)
								cur++
								if (i & 1) == 0 {
									mesh.Vertices[vertices] = object.vertices[entry1.First].Add(object.translation)
									vertices++
									if uv != -1 {
										uvs[uv].X = entry1.Second.X
										uvs[uv].Y = entry1.Second.Y
										uv++
									}

									mesh.Vertices[vertices] = object.vertices[entry2.First].Add(object.translation)
									vertices++
									if uv != -1 {
										uvs[uv].X = entry2.Second.X
										uvs[uv].Y = entry2.Second.Y
										uv++
									}

								} else {
									mesh.Vertices[vertices] = object.vertices[entry2.First].Add(object.translation)
									vertices++
									if uv != -1 {
										uvs[uv].X = entry2.Second.X
										uvs[uv].Y = entry2.Second.Y
										uv++
									}

									mesh.Vertices[vertices] = object.vertices[entry1.First].Add(object.translation)
									vertices++
									if uv != -1 {
										uvs[uv].X = entry1.Second.X
										uvs[uv].Y = entry1.Second.Y
										uv++
									}
								}
								if vertices-len(mesh.Vertices) >= len(mesh.Vertices) {
									logger.FatalF("AC3D: Invalid number of vertices")
								}
								mesh.Vertices[vertices] = object.vertices[entry3.First].Add(object.translation)
								vertices++
								if uv != -1 {
									uvs[uv].X = entry3.Second.X
									uvs[uv].Y = entry3.Second.Y
									uv++
								}

							}
						} else {

							it2 := 0
							// either a closed or an unclosed line
							tmp := len(itv.entries)
							if OpenLine == Type {
								tmp--
							}
							for m := 0; m < tmp; m++ {
								face := mesh.Faces[faces]
								faces++
								face.Indices = make([]uint32, 2)
								face.Indices[0] = uint32(cur)
								cur++
								face.Indices[1] = uint32(cur)
								cur++

								// copy vertex positions
								if it2 == len(itv.entries) {
									logger.FatalF("AC3D: Bad line")
								}
								common.AiAssert(itv.entries[it2].First < len(object.vertices))
								mesh.Vertices[vertices] = object.vertices[itv.entries[it2].First]
								vertices++

								// copy texture coordinates
								if uv != -1 {
									uvs[uv].X = itv.entries[it2].Second.X
									uvs[uv].Y = itv.entries[it2].Second.Y
									uv++
								}

								if ClosedLine == Type && tmp-1 == m {
									// if this is a closed line repeat its beginning now
									it2 = 0
								} else {
									it2++
								}

								// second point
								mesh.Vertices[vertices] = object.vertices[itv.entries[it2].First]
								vertices++
								if uv != -1 {
									uvs[uv].X = itv.entries[it2].Second.X
									uvs[uv].Y = itv.entries[it2].Second.Y
									uv++
								}
							}
						}
					}
				}
				cit++
				mat++
			}

			// Now apply catmull clark subdivision if necessary. We split meshes into
			// materials which is not done by AC3D during smoothing, so we need to
			// collect all meshes using the same material group.
			if object.subDiv != 0 {
				if ac.configEvalSubdivision {
					div := core.NewSubDivision(core.CATMULL_CLARKE)
					logger.InfoF("AC3D: Evaluating subdivision surface:%v ", object.name)
					cpy := make([]*core.AiMesh, len(*meshes)-oldm)
					tmp := (*meshes)[oldm:]
					div.Subdivide(tmp, cpy, object.subDiv, true)
					copy((*meshes)[oldm:], cpy)
					// previous meshes are deleted vy Subdivide().
				} else {
					logger.InfoF("AC3D: Letting the subdivision surface untouched due to my configuration: %v", object.name)
				}
			}
		}

	}
	if object.name != "" {
		node.Name = object.name
	} else {
		switch object.Type {
		case Group:
			node.Name = fmt.Sprintf("ACGroup_%v", ac.GroupsCounter)
			ac.GroupsCounter++
		case Poly:
			node.Name = fmt.Sprintf("ACPoly_%v", ac.GroupsCounter)
			ac.PolysCounter++
		case Light:
			node.Name = fmt.Sprintf("ACLight_%v", ac.GroupsCounter)
			ac.LightsCounter++
		case World:
			node.Name = fmt.Sprintf("ACWorld_%v", ac.GroupsCounter)
			ac.WorldsCounter++
		}
	}

	node.Transformation = common.NewAiMatrix4x4FromAiMatrix3x3(object.rotation)
	if object.Type == Group || object.numRefs == 0 {
		node.Transformation.A4 = object.translation.X
		node.Transformation.B4 = object.translation.Y
		node.Transformation.C4 = object.translation.Z
	}

	// add children to the object
	if len(object.children) > 0 {
		for i := 0; i < len(object.children); i++ {
			nodeTmp := ac.ConvertObjectSection(object.children[i], meshes, outMaterials, materials, node)
			node.Children = append(node.Children, nodeTmp)
		}
	}
	return
}

func (ac *AC3DImporter) Read(pScene *core.AiScene) (err error) {
	materials := make([]*Material, 0)
	rootObjects := make([]*Object, 0)
	for !ac.EOF() {
		if ac.HasPrefix("MATERIAL") {
			mat := newMaterial()
			materials = append(materials, mat)
			mat.name, err = ac.NextOneKeyString("MATERIAL")
			if err != nil {
				return err
			}
			mat.name = common.ClearQuotationMark(mat.name)
			mat.rgb, err = ac.NextKeyAiColor3d("rgb")
			if err != nil {
				return err
			}
			mat.amb, err = ac.NextKeyAiColor3d("amb")
			if err != nil {
				return err
			}
			mat.emis, err = ac.NextKeyAiColor3d("emis")
			if err != nil {
				return err
			}
			mat.spec, err = ac.NextKeyAiColor3d("spec")
			if err != nil {
				return err
			}
			mat.shin, err = ac.NextOneKeyFloat32("shi")
			if err != nil {
				return err
			}
			mat.trans, err = ac.NextOneKeyFloat32("trans")
			if err != nil {
				return err
			}
			ac.NextLine()
		}
		objs, err := ac.LoadObjectSection()
		if err != nil {
			return err
		}
		rootObjects = append(rootObjects, objs...)
	}
	if len(rootObjects) == 0 || ac.NumMeshes == 0 {
		return errors.New("AC3D: No meshes have been loaded")
	}
	if len(materials) == 0 {
		logger.Warn("AC3D: No material has been found")
		materials = append(materials, newMaterial())
	}
	ac.NumMeshes += (ac.NumMeshes >> 2) + 1
	var root *Object
	if len(rootObjects) == 1 {
		root = rootObjects[0]
	} else {
		root = newObject()
	}
	var meshes []*core.AiMesh
	var omaterials []*core.AiMaterial
	node := ac.ConvertObjectSection(root, &meshes, &omaterials, &materials, nil)
	pScene.RootNode = node
	if strings.HasPrefix(pScene.RootNode.Name, "Node") {
		pScene.RootNode.Name = "<AC3DWorld>"
	}
	if len(meshes) == 0 {
		return errors.New("An unknown error occurred during converting")
	}
	pScene.Meshes = meshes
	pScene.Materials = omaterials
	pScene.Lights = ac.Lights
	return nil
}
