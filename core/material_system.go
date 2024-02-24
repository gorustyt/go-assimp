package core

import (
	"fmt"
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"strings"
)

//! @cond AI_DOX_INCLUDE_INTERNAL
// ---------------------------------------------------------------------------
/** @brief A very primitive RTTI system for the contents of material
 *  properties.
 */
type AiPropertyTypeInfo uint32

const (
	/** Array of single-precision (32 Bit) floats
	 *
	 *  It is possible to use aiGetMaterialInteger[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in floating-point format.
	 *  The material system performs the type conversion automatically.
	 */
	AiPTI_Float AiPropertyTypeInfo = 0x1

	/** Array of double-precision (64 Bit) floats
	 *
	 *  It is possible to use aiGetMaterialInteger[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in floating-point format.
	 *  The material system performs the type conversion automatically.
	 */
	AiPTI_Double AiPropertyTypeInfo = 0x2

	/** The material property is an aiString.
	 *
	 *  Arrays of strings aren't possible, aiGetMaterialString() (or the
	 *  C++-API aiMaterial::Get()) *must* be used to query a string property.
	 */
	AiPTI_String AiPropertyTypeInfo = 0x3
	/** Array of (32 Bit) integers
	 *
	 *  It is possible to use aiGetMaterialFloat[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in integer format.
	 *  The material system performs the type conversion automatically.
	 */
	AiPTI_Integer AiPropertyTypeInfo = 0x4
	/** Simple binary buffer, content undefined. Not convertible to anything.
	 */
	AiPTI_Buffer AiPropertyTypeInfo = 0x5
	/** This value is not used. It is just there to force the
	 *  compiler to map this enum to a 32 Bit integer.
	 */
)

func init() {
	RegisterPropertyTypeInfo("default", AI_MATKEY_OBJ_ILLUM.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeInt)

	RegisterPropertyTypeInfo("default", AI_MATKEY_NAME.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeString)

	RegisterPropertyTypeInfo("default", AI_MATKEY_COLOR_AMBIENT.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D)
	RegisterPropertyTypeInfo("default", AI_MATKEY_COLOR_DIFFUSE.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D)
	RegisterPropertyTypeInfo("default", AI_MATKEY_COLOR_EMISSIVE.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D)
	RegisterPropertyTypeInfo("default", AI_MATKEY_COLOR_REFLECTIVE.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D)
	RegisterPropertyTypeInfo("default", AI_MATKEY_COLOR_SPECULAR.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D, pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D)

	RegisterPropertyTypeInfo("default", AI_MATKEY_SHININESS.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32)

	RegisterPropertyTypeInfo("default", AI_MATKEY_UVTRANSFORM_AMBIENT(0).Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeAiUVTransform)
	RegisterPropertyTypeInfo("default", AI_MATKEY_TEXTURE_DIFFUSE(0).Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeString)
	RegisterPropertyTypeInfo("default", AI_MATKEY_MAPPINGMODE_U_DIFFUSE(0).Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32)
	RegisterPropertyTypeInfo("default", AI_MATKEY_MAPPINGMODE_V_DIFFUSE(0).Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32)
	RegisterPropertyTypeInfo("default", AI_MATKEY_OPACITY.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32)
	RegisterPropertyTypeInfo("default", AI_MATKEY_REFLECTIVITY.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32)
	RegisterPropertyTypeInfo("default", AI_MATKEY_SHADING_MODEL.Key, pb_msg.AiMaterialPropertyType_AiPropertyTypeInt)
}

func GetAiPropertyTypeInfo(key string) (res []pb_msg.AiMaterialPropertyType) {
	for k, v := range aiPropertyTypeInfoMap {
		if strings.HasPrefix(k, key) {
			res = append(res, v...)
		}
	}
	return res
}

func UnRegisterPropertyTypeSource(source string) {
	for k := range aiPropertyTypeInfoMap {
		if strings.HasSuffix(k, fmt.Sprintf(".source.%v", source)) {
			delete(aiPropertyTypeInfoMap, k)
		}
	}
}

func RegisterPropertyTypeInfo(source, key string, values ...pb_msg.AiMaterialPropertyType) {
	key = fmt.Sprintf("%v.source.%v", key, source)
	v, ok := aiPropertyTypeInfoMap[key]
	if !ok {
		aiPropertyTypeInfoMap[key] = append(aiPropertyTypeInfoMap[key], values...)
		return
	}
	for _, t := range v {
		for _, value := range values {
			if t == value {
				continue
			}
			aiPropertyTypeInfoMap[key] = append(aiPropertyTypeInfoMap[key], value)
		}
	}
}

var aiPropertyTypeInfoMap = map[string][]pb_msg.AiMaterialPropertyType{}
