package test

import (
	"assimp/common"
	"assimp/common/logger"
	"assimp/common/pb_msg"
	"assimp/core"
	"google.golang.org/protobuf/proto"
	"math"
	"testing"
)

func AssertFloatEqual(t *testing.T, t1, t2 float64, eps float64) {
	Assert(t, math.Abs(t1-t2) > eps)
}
func Assert(t *testing.T, ok bool, msg ...string) {
	if !ok {
		t.Fatalf("test not ok :%v", msg)
	}
}

func AssertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("test error :%v", err)
	}
}

func deepEqualMaterials(p1, p2 []*core.AiMaterial) {
	for i, m1 := range p1 {
		m2 := p2[i]
		for j, v := range m1.Properties {
			v1 := m2.Properties[j]
			if v.Key != v1.Key {
				logger.ErrorF("deepEqualMaterials key  not equal v1:%v v2%v  ", v.Key, v1.Key)
			}
			if v.Type != v1.Type {
				logger.ErrorF("deepEqualMaterials key  not equal v1:%v v2%v", v.Type, v1.Type)
			}
			switch v.Type {
			case pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32, pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat64:
				t1 := &pb_msg.AiMaterialPropertyFloat64{}
				t2 := &pb_msg.AiMaterialPropertyFloat64{}
				err := proto.Unmarshal(v.Data, t1)
				if err != nil {
					panic(err)
				}
				err = proto.Unmarshal(v1.Data, t2)
				if err != nil {
					panic(err)
				}
				if !deepEqual(v, v1) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!", v.Key, v1.Key)
				}
			case pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D:
				t1 := &pb_msg.AiColor3D{}
				t2 := &pb_msg.AiColor3D{}
				err := proto.Unmarshal(v.Data, t1)
				if err != nil {
					panic(err)
				}
				err = proto.Unmarshal(v1.Data, t2)
				if err != nil {
					panic(err)
				}
				if !deepEqual((&common.AiColor3D{}).FromPbMsg(t1), (&common.AiColor3D{}).FromPbMsg(t2)) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!", v.Key, v1.Key)
				}
			case pb_msg.AiMaterialPropertyType_AiPropertyTypeColor4D:
				t1 := &pb_msg.AiColor4D{}
				t2 := &pb_msg.AiColor4D{}
				err := proto.Unmarshal(v.Data, t1)
				if err != nil {
					panic(err)
				}
				err = proto.Unmarshal(v1.Data, t2)
				if err != nil {
					panic(err)
				}
				if !deepEqual((&common.AiColor4D{}).FromPbMsg(t1), (&common.AiColor4D{}).FromPbMsg(t2)) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!", v.Key, v1.Key)
				}
			case pb_msg.AiMaterialPropertyType_AiPropertyTypeVector4D:
				t1 := &pb_msg.AiVector4D{}
				t2 := &pb_msg.AiVector4D{}
				err := proto.Unmarshal(v.Data, t1)
				if err != nil {
					panic(err)
				}
				err = proto.Unmarshal(v1.Data, t2)
				if err != nil {
					panic(err)
				}
				if !deepEqual((&common.AiVector4D{}).FromPbMsg(t1), (&common.AiVector4D{}).FromPbMsg(t2)) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!", v.Key, v1.Key)
				}
			case pb_msg.AiMaterialPropertyType_AiPropertyTypeVector3D:
				t1 := &pb_msg.AiVector3D{}
				t2 := &pb_msg.AiVector3D{}
				err := proto.Unmarshal(v.Data, t1)
				if err != nil {
					panic(err)
				}
				err = proto.Unmarshal(v1.Data, t2)
				if err != nil {
					panic(err)
				}
				if !deepEqual((&common.AiVector3D{}).FromPbMsg(t1), (&common.AiVector3D{}).FromPbMsg(t2)) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!", v.Key, v1.Key)
				}
			case pb_msg.AiMaterialPropertyType_AiPropertyTypeAiUVTransform:
				t1 := &pb_msg.AiUVTransform{}
				t2 := &pb_msg.AiUVTransform{}
				err := proto.Unmarshal(v.Data, t1)
				if err != nil {
					panic(err)
				}
				err = proto.Unmarshal(v1.Data, t2)
				if err != nil {
					panic(err)
				}
				if !deepEqual((&core.AiUVTransform{}).FromPbMsg(t1), (&core.AiUVTransform{}).FromPbMsg(t2)) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!", v.Key, v1.Key)
				}
			default:
				if !deepEqual(v, v1) {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties not equal!", v.Key, v1.Key)
				}
			}
			v.Data = v.Data[:0] //校验数据完直接强制相等
			v1.Data = v.Data
		}
	}
}
func deepEqualMesh(p1, p2 []*core.AiMesh) {

	for _, vi := range p1 {
		for _, vj := range p2 {
			if !deepEqual(vi.Method, vj.Method) {
				logger.ErrorF("deepEqualMesh Method not equal vi:%v,vj:%v", vi.Method, vj.Method)
			}

			if !deepEqual(vi.PrimitiveTypes, vj.PrimitiveTypes) {
				logger.ErrorF("deepEqualMesh PrimitiveTypes not equal vi:%v,vj:%v", vi.PrimitiveTypes, vj.PrimitiveTypes)
			}

			if !deepEqual(vi.MaterialIndex, vj.MaterialIndex) {
				logger.ErrorF("deepEqualMesh MaterialIndex not equal vi:%v,vj:%v", vi.MaterialIndex, vj.MaterialIndex)
			}

			if !deepEqual(vi.Name, vj.Name) {
				logger.ErrorF("deepEqualMesh Name not equal vi:%v,vj:%v", vi.Name, vj.Name)
			}

			if !deepEqual(vi.TextureCoordsNames, vj.TextureCoordsNames) {
				logger.ErrorF("deepEqualMesh TextureCoordsNames not equal vi:%v,vj:%v", vi.TextureCoordsNames, vj.TextureCoordsNames)
			}

			if !deepEqual(vi.AABB, vj.AABB) {
				logger.ErrorF("deepEqualMesh AABB not equal vi:%v,vj:%v", *vi.AABB, *vj.AABB)
			}

			for iv := range vi.AnimMeshes {
				v1, v2 := vi.AnimMeshes[iv], vj.AnimMeshes[iv]
				if !deepEqual(v1, v2) {
					logger.ErrorF("deepEqualMesh Face not equal vi:%v,vj:%v", *v1, *v2)
				}
			}
			for iv := range vi.Faces {
				v1, v2 := vi.Faces[iv], vj.Faces[iv]
				if !deepEqual(v1, v2) {
					logger.ErrorF("deepEqualMesh Face not equal vi:%v,vj:%v", *v1, *v2)
				}
			}

			for iv := range vi.Vertices {
				v1, v2 := vi.Vertices[iv], vj.Vertices[iv]
				if !deepEqual(v1, v2) {
					logger.ErrorF("deepEqualMesh Vertic not equal vi:%v,vj:%v", *v1, *v2)
				}
			}
			for iv := range vi.Normals {
				if !deepEqual(vi.Normals[iv], vj.Normals[iv]) {
					logger.ErrorF("deepEqualMesh Normal not equal vi:%v,vj:%v", *vi.Normals[iv], *vj.Normals[iv])
				}
			}
			for iv := range vi.Tangents {
				v1, v2 := vi.Tangents[iv], vj.Tangents[iv]
				if !deepEqual(v1, v2) {
					logger.ErrorF("deepEqualMesh Tangent not equal vi:%v,vj:%v", *v1, *v2)
				}
			}

			for i := range vi.Colors {
				for j := range vi.Colors[i] {
					v1, v2 := vi.Colors[i][j], vj.Colors[i][j]
					if !deepEqual(v1, v2) {
						logger.ErrorF("deepEqualMesh Color not equal vi:%v,vj:%v", *v1, *v2)
					}
				}
			}

			for i := range vi.TextureCoords {
				for j := range vi.TextureCoords[i] {
					v1, v2 := vi.TextureCoords[i][j], vj.TextureCoords[i][j]
					if !deepEqual(v1, v2) {
						logger.ErrorF("deepEqualMesh TextureCoord not equal vi:%v,vj:%v", *v1, *v2)
					}
				}
			}
			for iv := range vi.NumUVComponents {
				v1, v2 := vi.NumUVComponents[iv], vj.NumUVComponents[iv]
				if !deepEqual(v1, v2) {
					logger.ErrorF("deepEqualMesh NumUVComponent not equal vi:%v,vj:%v", v1, v2)
				}
			}

			for iv := range vi.Bones {
				v1, v2 := vi.Bones[iv], vj.Bones[iv]
				if !deepEqual(v1, v2) {
					logger.ErrorF("deepEqualMesh Bone not equal vi:%v,vj:%v", *v1, *v2)
				}
			}

		}
	}
}
func deepEqualScene(p1, p2 *core.AiScene) {
	for _, mesh := range p1.Meshes {
		mesh.Name = "" //assetBin 沒有這個字段
	}
	if p1.Flags != p2.Flags {
		logger.ErrorF("deepScene Flags not equal v1:%v v2%v  ", p1.Flags, p2.Flags)
	}

	if p1.Name != p2.Name {
		logger.ErrorF("deepScene Name not equal v1:%v v2%v", p1.Name, p2.Name)
	}
}

func deepEqualNode(p1, p2 *core.AiNode, level int) {
	if p1.Name != p2.Name {
		logger.ErrorF("deepEqualNode Name not equal ")
	}
	for i := range p1.Children {
		deepEqualNode(p1.Children[i], p2.Children[i], level+1)
	}
	if !deepEqual(p1.Transformation, p2.Transformation) {
		logger.ErrorF("deepEqualNode Transformation not equal p1:%v,p2:%v", *p1.Transformation, *p2.Transformation)
	}
	if !deepEqual(p1.Meshes, p2.Meshes) {
		logger.ErrorF("deepEqualNode Meshes not equal p1:%v,p2:%v", p1.Meshes, p2.Meshes)
	}
	if !deepEqual(p1.MetaData, p2.MetaData) {
		logger.Error("deepEqualNode MetaData not equal ")
	}
}

func deepEqualTexture(p1, p2 []*core.AiTexture) {
	if !deepEqual(p1, p2) {
		logger.Error("deepScene Texture not equal ")
	}
}

func deepEqualAnimation(p1, p2 []*core.AiAnimation) {
	if !deepEqual(p1, p2) {
		logger.Error("deepScene Animation not equal ")
	}
}

func deepEqualAiLight(p1, p2 []*core.AiLight) {
	if !deepEqual(p1, p2) {
		logger.Error("deepScene Light not equal ")
	}
}

func deepEqualCamera(p1, p2 []*core.AiCamera) {
	if !deepEqual(p1, p2) {
		logger.Error("deepScene Camera not equal ")
	}
}

func deepEqualMetadata(p1, p2 []*core.AiMetadata) {
	if !deepEqual(p1, p2) {
		logger.Error("deepScene Metadata not equal ")
	}
}

func deepEqualAiSkeleton(p1, p2 []*core.AiSkeleton) {
	if !deepEqual(p1, p2) {
		logger.Error("deepScene Skeleton not equal ")
	}
}

func DeepEqual(p1, p2 *core.AiScene) bool {

	//只比较一部分，并且修改协议
	deepEqualScene(p1, p2)
	deepEqualNode(p1.RootNode, p2.RootNode, 0)
	deepEqualMaterials(p1.Materials, p2.Materials)
	deepEqualMesh(p1.Meshes, p2.Meshes)
	deepEqualTexture(p1.Textures, p2.Textures)
	deepEqualAnimation(p1.Animations, p2.Animations)
	deepEqualAiLight(p1.Lights, p2.Lights)
	deepEqualCamera(p1.Cameras, p2.Cameras)
	deepEqualMetadata(p1.MetaData, p2.MetaData)
	deepEqualAiSkeleton(p1.Skeletons, p2.Skeletons)
	//全比较
	return deepEqual(p1, p2)
}
