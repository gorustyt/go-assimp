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

func DeepEqual(p1, p2 *core.AiScene) bool {
	for _, mesh := range p1.Meshes {
		mesh.Name = "" //assetBin 沒有這個字段
	}
	for _, mesh := range p1.Materials {
		for _, mesh1 := range p2.Materials {
			for i, v := range mesh.Properties {
				v1 := mesh1.Properties[i]
				if v.Type != v1.Type {
					logger.ErrorF("v1Name:%v v2Name:%v  Properties Type not equal!", v.Key, v1.Key)
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

	return deepEqual(p1, p2)
}
