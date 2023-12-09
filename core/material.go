package core

import (
	"assimp/common"
	"assimp/common/pb_msg"
	"google.golang.org/protobuf/proto"
)

var (
	AI_MATKEY_NAME                    = newAiMaterialProperty("?mat.name", 0, 0)
	AI_MATKEY_TWOSIDED                = newAiMaterialProperty("$mat.twosided", 0, 0)
	AI_MATKEY_SHADING_MODEL           = newAiMaterialProperty("$mat.shadingm", 0, 0)
	AI_MATKEY_ENABLE_WIREFRAME        = newAiMaterialProperty("$mat.wireframe", 0, 0)
	AI_MATKEY_BLEND_FUNC              = newAiMaterialProperty("$mat.blend", 0, 0)
	AI_MATKEY_OPACITY                 = newAiMaterialProperty("$mat.opacity", 0, 0)
	AI_MATKEY_TRANSPARENCYFACTOR      = newAiMaterialProperty("$mat.transparencyfactor", 0, 0)
	AI_MATKEY_BUMPSCALING             = newAiMaterialProperty("$mat.bumpscaling", 0, 0)
	AI_MATKEY_SHININESS               = newAiMaterialProperty("$mat.shininess", 0, 0)
	AI_MATKEY_REFLECTIVITY            = newAiMaterialProperty("$mat.reflectivity", 0, 0)
	AI_MATKEY_SHININESS_STRENGTH      = newAiMaterialProperty("$mat.shinpercent", 0, 0)
	AI_MATKEY_REFRACTI                = newAiMaterialProperty("$mat.refracti", 0, 0)
	AI_MATKEY_COLOR_DIFFUSE           = newAiMaterialProperty("$clr.diffuse", 0, 0)
	AI_MATKEY_COLOR_AMBIENT           = newAiMaterialProperty("$clr.ambient", 0, 0)
	AI_MATKEY_COLOR_SPECULAR          = newAiMaterialProperty("$clr.specular", 0, 0)
	AI_MATKEY_COLOR_EMISSIVE          = newAiMaterialProperty("$clr.emissive", 0, 0)
	AI_MATKEY_COLOR_TRANSPARENT       = newAiMaterialProperty("$clr.transparent", 0, 0)
	AI_MATKEY_COLOR_REFLECTIVE        = newAiMaterialProperty("$clr.reflective", 0, 0)
	AI_MATKEY_GLOBAL_BACKGROUND_IMAGE = newAiMaterialProperty("?bg.global", 0, 0)
	AI_MATKEY_GLOBAL_SHADERLANG       = newAiMaterialProperty("?sh.lang", 0, 0)
	AI_MATKEY_SHADER_VERTEX           = newAiMaterialProperty("?sh.vs", 0, 0)
	AI_MATKEY_SHADER_FRAGMENT         = newAiMaterialProperty("?sh.fs", 0, 0)
	AI_MATKEY_SHADER_GEO              = newAiMaterialProperty("?sh.gs", 0, 0)
	AI_MATKEY_SHADER_TESSELATION      = newAiMaterialProperty("?sh.ts", 0, 0)
	AI_MATKEY_SHADER_PRIMITIVE        = newAiMaterialProperty("?sh.ps", 0, 0)
	AI_MATKEY_SHADER_COMPUTE          = newAiMaterialProperty("?sh.cs", 0, 0)

	// ---------------------------------------------------------------------------
	// PBR material support
	// --------------------
	// Properties defining PBR rendering techniques

	AI_MATKEY_USE_COLOR_MAP = newAiMaterialProperty("$mat.useColorMap", 0, 0)

	// Metallic/Roughness Workflow
	// ---------------------------
	// Base RGBA color factor. Will be multiplied by final base color texture values if extant
	// Note: Importers may choose to copy this into AI_MATKEY_COLOR_DIFFUSE for compatibility
	// with renderers and formats that do not support Metallic/Roughness PBR
	AI_MATKEY_BASE_COLOR = newAiMaterialProperty("$clr.base", 0, 0)
	//AI_MATKEY_BASE_COLOR_TEXTURE =newAiMaterialProperty(AiTextureType_BASE_COLOR, 0)
	AI_MATKEY_USE_METALLIC_MAP = newAiMaterialProperty("$mat.useMetallicMap", 0, 0)
	// Metallic factor. 0.0 = Full Dielectric, 1.0 = Full Metal
	AI_MATKEY_METALLIC_FACTOR = newAiMaterialProperty("$mat.metallicFactor", 0, 0)
	//AI_MATKEY_METALLIC_TEXTURE=newAiMaterialProperty( AiTextureType_METALNESS, 0)
	AI_MATKEY_USE_ROUGHNESS_MAP = newAiMaterialProperty("$mat.useRoughnessMap", 0, 0)
	// Roughness factor. 0.0 = Perfectly Smooth, 1.0 = Completely Rough
	AI_MATKEY_ROUGHNESS_FACTOR = newAiMaterialProperty("$mat.roughnessFactor", 0, 0)
	//AI_MATKEY_ROUGHNESS_TEXTURE =newAiMaterialProperty(AiTextureType_DIFFUSE_ROUGHNESS, 0)
	// Anisotropy factor. 0.0 = isotropic, 1.0 = anisotropy along tangent direction,
	// -1.0 = anisotropy along bitangent direction
	AI_MATKEY_ANISOTROPY_FACTOR = newAiMaterialProperty("$mat.anisotropyFactor", 0, 0)

	// Specular/Glossiness Workflow
	// ---------------------------
	// Diffuse/Albedo Color. Note: Pure Metals have a diffuse of {0,0,0}
	// AI_MATKEY_COLOR_DIFFUSE
	// Specular Color.
	// Note: Metallic/Roughness may also have a Specular Color
	// AI_MATKEY_COLOR_SPECULAR
	AI_MATKEY_SPECULAR_FACTOR = newAiMaterialProperty("$mat.specularFactor", 0, 0)
	// Glossiness factor. 0.0 = Completely Rough, 1.0 = Perfectly Smooth
	AI_MATKEY_GLOSSINESS_FACTOR = newAiMaterialProperty("$mat.glossinessFactor", 0, 0)

	// Sheen
	// -----
	// Sheen base RGB color. Default {0,0,0}
	AI_MATKEY_SHEEN_COLOR_FACTOR = newAiMaterialProperty("$clr.sheen.factor", 0, 0)
	// Sheen Roughness Factor.
	AI_MATKEY_SHEEN_ROUGHNESS_FACTOR = newAiMaterialProperty("$mat.sheen.roughnessFactor", 0, 0)
	//AI_MATKEY_SHEEN_COLOR_TEXTURE=newAiMaterialProperty( AiTextureType_SHEEN, 0)
	//AI_MATKEY_SHEEN_ROUGHNESS_TEXTURE =newAiMaterialProperty(AiTextureType_SHEEN, 1)

	// Clearcoat
	// ---------
	// Clearcoat layer intensity. 0.0 = none (disabled)
	AI_MATKEY_CLEARCOAT_FACTOR           = newAiMaterialProperty("$mat.clearcoat.factor", 0, 0)
	AI_MATKEY_CLEARCOAT_ROUGHNESS_FACTOR = newAiMaterialProperty("$mat.clearcoat.roughnessFactor", 0, 0)
	//AI_MATKEY_CLEARCOAT_TEXTURE =newAiMaterialProperty(AiTextureType_CLEARCOAT, 0)
	//AI_MATKEY_CLEARCOAT_ROUGHNESS_TEXTURE =newAiMaterialProperty(AiTextureType_CLEARCOAT, 1)
	//AI_MATKEY_CLEARCOAT_NORMAL_TEXTURE =newAiMaterialProperty(AiTextureType_CLEARCOAT, 2)

	// Transmission
	// ------------
	// https://github.com/KhronosGroup/glTF/tree/master/extensions/2.0/Khronos/KHR_materials_transmission
	// Base percentage of light transmitted through the surface. 0.0 = Opaque, 1.0 = Fully transparent
	AI_MATKEY_TRANSMISSION_FACTOR = newAiMaterialProperty("$mat.transmission.factor", 0, 0)
	// Texture defining percentage of light transmitted through the surface.
	// Multiplied by AI_MATKEY_TRANSMISSION_FACTOR
	//AI_MATKEY_TRANSMISSION_TEXTURE =newAiMaterialProperty(AiTextureType_TRANSMISSION, 0)

	// Volume
	// ------------
	// https://github.com/KhronosGroup/glTF/tree/main/extensions/2.0/Khronos/KHR_materials_volume
	// The thickness of the volume beneath the surface. If the value is 0 the material is thin-walled. Otherwise the material is a volume boundary.
	AI_MATKEY_VOLUME_THICKNESS_FACTOR = newAiMaterialProperty("$mat.volume.thicknessFactor", 0, 0)
	// Texture that defines the thickness.
	// Multiplied by AI_MATKEY_THICKNESS_FACTOR
	//AI_MATKEY_VOLUME_THICKNESS_TEXTURE =newAiMaterialProperty(AiTextureType_TRANSMISSION, 1)
	// Density of the medium given as the average distance that light travels in the medium before interacting with a particle.
	AI_MATKEY_VOLUME_ATTENUATION_DISTANCE = newAiMaterialProperty("$mat.volume.attenuationDistance", 0, 0)
	// The color that white light turns into due to absorption when reaching the attenuation distance.
	AI_MATKEY_VOLUME_ATTENUATION_COLOR = newAiMaterialProperty("$mat.volume.attenuationColor", 0, 0)

	// Emissive
	// --------
	AI_MATKEY_USE_EMISSIVE_MAP   = newAiMaterialProperty("$mat.useEmissiveMap", 0, 0)
	AI_MATKEY_EMISSIVE_INTENSITY = newAiMaterialProperty("$mat.emissiveIntensity", 0, 0)
	AI_MATKEY_USE_AO_MAP         = newAiMaterialProperty("$mat.useAOMap", 0, 0)

	// ---------------------------------------------------------------------------
	// Pure key names for all texture-related properties
	//! @cond MATS_DOC_FULL
	_AI_MATKEY_TEXTURE_BASE       = "$tex.file"
	_AI_MATKEY_UVWSRC_BASE        = "$tex.uvwsrc"
	_AI_MATKEY_TEXOP_BASE         = "$tex.op"
	_AI_MATKEY_MAPPING_BASE       = "$tex.mapping"
	_AI_MATKEY_TEXBLEND_BASE      = "$tex.blend"
	_AI_MATKEY_MAPPINGMODE_U_BASE = "$tex.mapmodeu"
	_AI_MATKEY_MAPPINGMODE_V_BASE = "$tex.mapmodev"
	_AI_MATKEY_TEXMAP_AXIS_BASE   = "$tex.mapaxis"
	_AI_MATKEY_UVTRANSFORM_BASE   = "$tex.uvtrafo"
	_AI_MATKEY_TEXFLAGS_BASE      = "$tex.flags"
	//! @endcond
)

// ---------------------------------------------------------------------------
func AI_MATKEY_TEXTURE(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_TEXTURE_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_TEXTURE_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_TEXTURE_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_TEXTURE_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_TEXTURE_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_TEXTURE_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_NORMALS, N)
}

func AI_MATKEY_TEXTURE_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_TEXTURE_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_SHININESS, N)
}

func AI_MATKEY_TEXTURE_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_OPACITY, N)
}

func AI_MATKEY_TEXTURE_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_TEXTURE_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_TEXTURE_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_TEXTURE(AiTextureType_REFLECTION, N)
}

//! @endcond

// ---------------------------------------------------------------------------
func AI_MATKEY_UVWSRC(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_UVWSRC_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_UVWSRC_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_UVWSRC_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_UVWSRC_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_UVWSRC_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_UVWSRC_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_NORMALS, N)
}

func AI_MATKEY_UVWSRC_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_UVWSRC_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_SHININESS, N)
}

func AI_MATKEY_UVWSRC_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_OPACITY, N)
}

func AI_MATKEY_UVWSRC_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_UVWSRC_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_UVWSRC_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_UVWSRC(AiTextureType_REFLECTION, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_TEXOP(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_TEXOP_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_TEXOP_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_TEXOP_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_TEXOP_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_TEXOP_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_TEXOP_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_NORMALS, N)
}

func AI_MATKEY_TEXOP_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_TEXOP_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_SHININESS, N)
}

func AI_MATKEY_TEXOP_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_OPACITY, N)
}

func AI_MATKEY_TEXOP_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_TEXOP_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_TEXOP_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_TEXOP(AiTextureType_REFLECTION, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_MAPPING(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_MAPPING_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_MAPPING_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_DIFFUSE, N)
}
func AI_MATKEY_MAPPING_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_SPECULAR, N)
}
func AI_MATKEY_MAPPING_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_AMBIENT, N)
}
func AI_MATKEY_MAPPING_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_EMISSIVE, N)
}
func AI_MATKEY_MAPPING_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_NORMALS, N)
}
func AI_MATKEY_MAPPING_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_HEIGHT, N)
}
func AI_MATKEY_MAPPING_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_SHININESS, N)
}
func AI_MATKEY_MAPPING_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_OPACITY, N)
}
func AI_MATKEY_MAPPING_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_DISPLACEMENT, N)
}
func AI_MATKEY_MAPPING_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_LIGHTMAP, N)
}
func AI_MATKEY_MAPPING_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPING(AiTextureType_REFLECTION, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_TEXBLEND(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_TEXBLEND_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_TEXBLEND_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_TEXBLEND_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_TEXBLEND_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_TEXBLEND_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_TEXBLEND_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_NORMALS, N)
}

func AI_MATKEY_TEXBLEND_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_TEXBLEND_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_SHININESS, N)
}

func AI_MATKEY_TEXBLEND_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_OPACITY, N)
}

func AI_MATKEY_TEXBLEND_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_TEXBLEND_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_TEXBLEND_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_TEXBLEND(AiTextureType_REFLECTION, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_MAPPINGMODE_U(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_MAPPINGMODE_U_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_MAPPINGMODE_U_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_MAPPINGMODE_U_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_MAPPINGMODE_U_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_MAPPINGMODE_U_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_MAPPINGMODE_U_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_NORMALS, N)
}

func AI_MATKEY_MAPPINGMODE_U_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_MAPPINGMODE_U_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_SHININESS, N)
}

func AI_MATKEY_MAPPINGMODE_U_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_OPACITY, N)
}

func AI_MATKEY_MAPPINGMODE_U_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_MAPPINGMODE_U_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_MAPPINGMODE_U_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_U(AiTextureType_REFLECTION, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_MAPPINGMODE_V(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_MAPPINGMODE_V_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_MAPPINGMODE_V_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_MAPPINGMODE_V_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_MAPPINGMODE_V_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_MAPPINGMODE_V_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_MAPPINGMODE_V_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_NORMALS, N)
}

func AI_MATKEY_MAPPINGMODE_V_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_MAPPINGMODE_V_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_SHININESS, N)
}

func AI_MATKEY_MAPPINGMODE_V_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_OPACITY, N)
}

func AI_MATKEY_MAPPINGMODE_V_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_MAPPINGMODE_V_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_MAPPINGMODE_V_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_MAPPINGMODE_V(AiTextureType_REFLECTION, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_TEXMAP_AXIS(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_TEXMAP_AXIS_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_TEXMAP_AXIS_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_TEXMAP_AXIS_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_TEXMAP_AXIS_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_TEXMAP_AXIS_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_TEXMAP_AXIS_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_NORMALS, N)
}

func AI_MATKEY_TEXMAP_AXIS_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_TEXMAP_AXIS_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_SHININESS, N)
}

func AI_MATKEY_TEXMAP_AXIS_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_OPACITY, N)
}

func AI_MATKEY_TEXMAP_AXIS_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_TEXMAP_AXIS_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_TEXMAP_AXIS_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_TEXMAP_AXIS(AiTextureType_REFLECTION, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_UVTRANSFORM(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_UVTRANSFORM_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_UVTRANSFORM_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_DIFFUSE, N)
}

func AI_MATKEY_UVTRANSFORM_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_SPECULAR, N)
}

func AI_MATKEY_UVTRANSFORM_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_AMBIENT, N)
}

func AI_MATKEY_UVTRANSFORM_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_EMISSIVE, N)
}

func AI_MATKEY_UVTRANSFORM_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_NORMALS, N)
}

func AI_MATKEY_UVTRANSFORM_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_HEIGHT, N)
}

func AI_MATKEY_UVTRANSFORM_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_SHININESS, N)
}

func AI_MATKEY_UVTRANSFORM_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_OPACITY, N)
}

func AI_MATKEY_UVTRANSFORM_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_DISPLACEMENT, N)
}

func AI_MATKEY_UVTRANSFORM_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_LIGHTMAP, N)
}

func AI_MATKEY_UVTRANSFORM_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_REFLECTION, N)
}

func AI_MATKEY_UVTRANSFORM_UNKNOWN(N int) AiMaterialProperty {
	return AI_MATKEY_UVTRANSFORM(AiTextureType_UNKNOWN, N)
}

// ! @endcond
// ---------------------------------------------------------------------------
func AI_MATKEY_TEXFLAGS(Type AiTextureType, N int) AiMaterialProperty {
	return newAiMaterialProperty(_AI_MATKEY_TEXFLAGS_BASE, int(Type), N)
}

// For backward compatibility and simplicity
// ! @cond MATS_DOC_FULL
func AI_MATKEY_TEXFLAGS_DIFFUSE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_DIFFUSE, N)
}
func AI_MATKEY_TEXFLAGS_SPECULAR(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_SPECULAR, N)
}
func AI_MATKEY_TEXFLAGS_AMBIENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_AMBIENT, N)
}
func AI_MATKEY_TEXFLAGS_EMISSIVE(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_EMISSIVE, N)
}
func AI_MATKEY_TEXFLAGS_NORMALS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_NORMALS, N)
}
func AI_MATKEY_TEXFLAGS_HEIGHT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_HEIGHT, N)
}
func AI_MATKEY_TEXFLAGS_SHININESS(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_SHININESS, N)
}
func AI_MATKEY_TEXFLAGS_OPACITY(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_OPACITY, N)
}
func AI_MATKEY_TEXFLAGS_DISPLACEMENT(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_DISPLACEMENT, N)
}
func AI_MATKEY_TEXFLAGS_LIGHTMAP(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_LIGHTMAP, N)
}
func AI_MATKEY_TEXFLAGS_REFLECTION(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_REFLECTION, N)
}
func AI_MATKEY_TEXFLAGS_UNKNOWN(N int) AiMaterialProperty {
	return AI_MATKEY_TEXFLAGS(AiTextureType_UNKNOWN, N)
}

type AiTextureType int

func (ai AiTextureType) ToString() string {
	switch ai {
	case AiTextureType_NONE:
		return "n/a"
	case AiTextureType_DIFFUSE:
		return "Diffuse"
	case AiTextureType_SPECULAR:
		return "Specular"
	case AiTextureType_AMBIENT:
		return "Ambient"
	case AiTextureType_EMISSIVE:
		return "Emissive"
	case AiTextureType_OPACITY:
		return "Opacity"
	case AiTextureType_NORMALS:
		return "Normals"
	case AiTextureType_HEIGHT:
		return "Height"
	case AiTextureType_SHININESS:
		return "Shininess"
	case AiTextureType_DISPLACEMENT:
		return "Displacement"
	case AiTextureType_LIGHTMAP:
		return "Lightmap"
	case AiTextureType_REFLECTION:
		return "Reflection"
	case AiTextureType_BASE_COLOR:
		return "BaseColor"
	case AiTextureType_NORMAL_CAMERA:
		return "NormalCamera"
	case AiTextureType_EMISSION_COLOR:
		return "EmissionColor"
	case AiTextureType_METALNESS:
		return "Metalness"
	case AiTextureType_DIFFUSE_ROUGHNESS:
		return "DiffuseRoughness"
	case AiTextureType_AMBIENT_OCCLUSION:
		return "AmbientOcclusion"
	case AiTextureType_SHEEN:
		return "Sheen"
	case AiTextureType_CLEARCOAT:
		return "Clearcoat"
	case AiTextureType_TRANSMISSION:
		return "Transmission"
	case AiTextureType_UNKNOWN:
		return "Unknown"
	default:
		break
	}
	return "BUG"
}

// ---------------------------------------------------------------------------
/** @brief Defines the purpose of a texture
 *
 *  This is a very difficult topic. Different 3D packages support different
 *  kinds of textures. For very common texture types, such as bumpmaps, the
 *  rendering results depend on implementation details in the rendering
 *  pipelines of these applications. Assimp loads all texture references from
 *  the model file and tries to determine which of the predefined texture
 *  types below is the best choice to match the original use of the texture
 *  as closely as possible.<br>
 *
 *  In content pipelines you'll usually define how textures have to be handled,
 *  and the artists working on models have to conform to this specification,
 *  regardless which 3D tool they're using.
 */

const (
	/** Dummy value.
	 *
	 *  No texture, but the value to be used as 'texture semantic'
	 *  (#aiMaterialProperty::mSemantic) for all material properties
	 *  *not* related to textures.
	 */
	AiTextureType_NONE AiTextureType = 0

	/** LEGACY API MATERIALS
	 * Legacy refers to materials which
	 * Were originally implemented in the specifications around 2000.
	 * These must never be removed, as most engines support them.
	 */

	/** The texture is combined with the result of the diffuse
	 *  lighting equation.
	 *  OR
	 *  PBR Specular/Glossiness
	 */
	AiTextureType_DIFFUSE AiTextureType = 1

	/** The texture is combined with the result of the specular
	 *  lighting equation.
	 *  OR
	 *  PBR Specular/Glossiness
	 */
	AiTextureType_SPECULAR AiTextureType = 2

	/** The texture is combined with the result of the ambient
	 *  lighting equation.
	 */
	AiTextureType_AMBIENT AiTextureType = 3

	/** The texture is added to the result of the lighting
	 *  calculation. It isn't influenced by incoming light.
	 */
	AiTextureType_EMISSIVE AiTextureType = 4

	/** The texture is a height map.
	 *
	 *  By convention, higher gray-scale values stand for
	 *  higher elevations from the base height.
	 */
	AiTextureType_HEIGHT AiTextureType = 5

	/** The texture is a (tangent space) normal-map.
	 *
	 *  Again, there are several conventions for tangent-space
	 *  normal maps. Assimp does (intentionally) not
	 *  distinguish here.
	 */
	AiTextureType_NORMALS AiTextureType = 6

	/** The texture defines the glossiness of the material.
	 *
	 *  The glossiness is in fact the exponent of the specular
	 *  (phong) lighting equation. Usually there is a conversion
	 *  function defined to map the linear color values in the
	 *  texture to a suitable exponent. Have fun.
	 */
	AiTextureType_SHININESS AiTextureType = 7

	/** The texture defines per-pixel opacity.
	 *
	 *  Usually 'white' means opaque and 'black' means
	 *  'transparency'. Or quite the opposite. Have fun.
	 */
	AiTextureType_OPACITY AiTextureType = 8

	/** Displacement texture
	 *
	 *  The exact purpose and format is application-dependent.
	 *  Higher color values stand for higher vertex displacements.
	 */
	AiTextureType_DISPLACEMENT AiTextureType = 9

	/** Lightmap texture (aka Ambient Occlusion)
	 *
	 *  Both 'Lightmaps' and dedicated 'ambient occlusion maps' are
	 *  covered by this material property. The texture contains a
	 *  scaling value for the final color value of a pixel. Its
	 *  intensity is not affected by incoming light.
	 */
	AiTextureType_LIGHTMAP AiTextureType = 10

	/** Reflection texture
	 *
	 * Contains the color of a perfect mirror reflection.
	 * Rarely used, almost never for real-time applications.
	 */
	AiTextureType_REFLECTION AiTextureType = 11

	/** PBR Materials
	 * PBR definitions from maya and other modelling packages now use this standard.
	 * This was originally introduced around 2012.
	 * Support for this is in game engines like Godot, Unreal or Unity3D.
	 * Modelling packages which use this are very common now.
	 */

	AiTextureType_BASE_COLOR        AiTextureType = 12
	AiTextureType_NORMAL_CAMERA     AiTextureType = 13
	AiTextureType_EMISSION_COLOR    AiTextureType = 14
	AiTextureType_METALNESS         AiTextureType = 15
	AiTextureType_DIFFUSE_ROUGHNESS AiTextureType = 16
	AiTextureType_AMBIENT_OCCLUSION AiTextureType = 17

	/** PBR Material Modifiers
	 * Some modern renderers have further PBR modifiers that may be overlaid
	 * on top of the 'base' PBR materials for additional realism.
	 * These use multiple texture maps, so only the base type is directly defined
	 */

	/** Sheen
	 * Generally used to simulate textiles that are covered in a layer of microfibers
	 * eg velvet
	 * https://github.com/KhronosGroup/glTF/tree/master/extensions/2.0/Khronos/KHR_materials_sheen
	 */
	AiTextureType_SHEEN AiTextureType = 19

	/** Clearcoat
	 * Simulates a layer of 'polish' or 'lacquer' layered on top of a PBR substrate
	 * https://autodesk.github.io/standard-surface/#closures/coating
	 * https://github.com/KhronosGroup/glTF/tree/master/extensions/2.0/Khronos/KHR_materials_clearcoat
	 */
	AiTextureType_CLEARCOAT AiTextureType = 20

	/** Transmission
	 * Simulates transmission through the surface
	 * May include further information such as wall thickness
	 */
	AiTextureType_TRANSMISSION AiTextureType = 21

	/** Unknown texture
	 *
	 *  A texture reference that does not match any of the definitions
	 *  above is considered to be 'unknown'. It is still imported,
	 *  but is excluded from any further post-processing.
	 */
	AiTextureType_UNKNOWN AiTextureType = 18
)

// ---------------------------------------------------------------------------
/** @brief Defines all shading models supported by the library
 *
 *  Property: #AI_MATKEY_SHADING_MODEL
 *
 *  The list of shading modes has been taken from Blender.
 *  See Blender documentation for more information. The API does
 *  not distinguish between "specular" and "diffuse" shaders (thus the
 *  specular term for diffuse shading models like Oren-Nayar remains
 *  undefined). <br>
 *  Again, this value is just a hint. Assimp tries to select the shader whose
 *  most common implementation matches the original rendering results of the
 *  3D modeler which wrote a particular model as closely as possible.
 *
 */

type AiShadingMode int

const (
	/** Flat shading. Shading is done on per-face base,
	 *  diffuse only. Also known as 'faceted shading'.
	 */
	AiShadingMode_Flat AiShadingMode = 0x1

	/** Simple Gouraud shading.
	 */
	AiShadingMode_Gouraud AiShadingMode = 0x2

	/** Phong-Shading -
	 */
	AiShadingMode_Phong AiShadingMode = 0x3

	/** Phong-Blinn-Shading
	 */
	AiShadingMode_Blinn AiShadingMode = 0x4

	/** Toon-Shading per pixel
	 *
	 *  Also known as 'comic' shader.
	 */
	AiShadingMode_Toon AiShadingMode = 0x5

	/** OrenNayar-Shading per pixel
	 *
	 *  Extension to standard Lambertian shading, taking the
	 *  roughness of the material into account
	 */
	AiShadingMode_OrenNayar AiShadingMode = 0x6

	/** Minnaert-Shading per pixel
	 *
	 *  Extension to standard Lambertian shading, taking the
	 *  "darkness" of the material into account
	 */
	AiShadingMode_Minnaert AiShadingMode = 0x7

	/** CookTorrance-Shading per pixel
	 *
	 *  Special shader for metallic surfaces.
	 */
	AiShadingMode_CookTorrance AiShadingMode = 0x8

	/** No shading at all. Constant light influence of 1.0.
	 * Also known as "Unlit"
	 */
	AiShadingMode_NoShading AiShadingMode = 0x9
	AiShadingMode_Unlit     AiShadingMode = AiShadingMode_NoShading // Alias

	/** Fresnel shading
	 */
	AiShadingMode_Fresnel AiShadingMode = 0xa

	/** Physically-Based Rendering (PBR) shading using
	 * Bidirectional scattering/reflectance distribution function (BSDF/BRDF)
	 * There are multiple methods under this banner, and model files may provide
	 * data for more than one PBR-BRDF method.
	 * Applications should use the set of provided properties to determine which
	 * of their preferred PBR rendering methods are likely to be available
	 * eg:
	 * - If AI_MATKEY_METALLIC_FACTOR is set, then a Metallic/Roughness is available
	 * - If AI_MATKEY_GLOSSINESS_FACTOR is set, then a Specular/Glossiness is available
	 * Note that some PBR methods allow layering of techniques
	 */
	AiShadingMode_PBR_BRDF AiShadingMode = 0xb
)

// ---------------------------------------------------------------------------
/** @brief Defines some mixed flags for a particular texture.
 *
 *  Usually you'll instruct your cg artists how textures have to look like ...
 *  and how they will be processed in your application. However, if you use
 *  Assimp for completely generic loading purposes you might also need to
 *  process these flags in order to display as many 'unknown' 3D models as
 *  possible correctly.
 *
 *  This corresponds to the #AI_MATKEY_TEXFLAGS property.
 */

type AiTextureFlags int

const (
	/** The texture's color values have to be inverted (component-wise 1-n)
	 */
	AiTextureFlags_Invert AiTextureFlags = 0x1

	/** Explicit request to the application to process the alpha channel
	 *  of the texture.
	 *
	 *  Mutually exclusive with #aiTextureFlags_IgnoreAlpha. These
	 *  flags are set if the library can say for sure that the alpha
	 *  channel is used/is not used. If the model format does not
	 *  define this, it is left to the application to decide whether
	 *  the texture alpha channel - if any - is evaluated or not.
	 */
	AiTextureFlags_UseAlpha AiTextureFlags = 0x2

	/** Explicit request to the application to ignore the alpha channel
	 *  of the texture.
	 *
	 *  Mutually exclusive with #aiTextureFlags_UseAlpha.
	 */
	AiTextureFlags_IgnoreAlpha AiTextureFlags = 0x4
)

// ---------------------------------------------------------------------------
/** @brief Defines alpha-blend flags.
 *
 *  If you're familiar with OpenGL or D3D, these flags aren't new to you.
 *  They define *how* the final color value of a pixel is computed, basing
 *  on the previous color at that pixel and the new color value from the
 *  material.
 *  The blend formula is:
 *  @code
 *    SourceColor * SourceBlend + DestColor * DestBlend
 *  @endcode
 *  where DestColor is the previous color in the frame-buffer at this
 *  position and SourceColor is the material color before the transparency
 *  calculation.<br>
 *  This corresponds to the #AI_MATKEY_BLEND_FUNC property.
 */
type AiBlendMode int

const (
	/**
	 *  Formula:
	 *  @code
	 *  SourceColor*SourceAlpha + DestColor*(1-SourceAlpha)
	 *  @endcode
	 */
	AiBlendMode_Default AiBlendMode = 0x0

	/** Additive blending
	 *
	 *  Formula:
	 *  @code
	 *  SourceColor*1 + DestColor*1
	 *  @endcode
	 */
	AiBlendMode_Additive AiBlendMode = 0x1

	// we don't need more for the moment, but we might need them
	// in future versions ...

)

// ---------------------------------------------------------------------------

type AiMaterial struct {
	/** List of all material properties loaded. */
	Properties []AiMaterialProperty
	/** Storage allocated */
	NumAllocated int
}

func (ai *AiMaterial) AddFloat32PropertyVar(pro AiMaterialProperty, data ...float32) {
	pro = pro.ResetData()
	pro.DataType = pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32
	bytesData, err := proto.Marshal(&pb_msg.AiMaterialPropertyFloat32{Data: data})
	if err != nil {
		panic(err)
	}
	pro.Data = bytesData
	ai.AddProperty(pro)
}
func (ai *AiMaterial) AddAiUVTransformPropertyVar(pro AiMaterialProperty, data AiUVTransform) {
	pro = pro.ResetData()
	pro.DataType = pb_msg.AiMaterialPropertyType_AiPropertyTypeAiUVTransform
	bytesData, err := proto.Marshal(data.ToPbMsg())
	if err != nil {
		panic(err)
	}
	pro.Data = bytesData
	ai.AddProperty(pro)
}
func (ai *AiMaterial) AddAiColorDPropertyVar(pro AiMaterialProperty, data common.AiColor3D) {
	pro = pro.ResetData()
	pro.DataType = pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D
	bytesData, err := proto.Marshal(data.ToPbMsg())
	if err != nil {
		panic(err)
	}
	pro.Data = bytesData
	ai.AddProperty(pro)
}
func (ai *AiMaterial) AddAiVector3DPropertyVar(pro AiMaterialProperty, data common.AiVector3D) {
	pro = pro.ResetData()
	pro.DataType = pb_msg.AiMaterialPropertyType_AiPropertyTypeVector3D
	bytesData, err := proto.Marshal(data.ToPbMsg())
	if err != nil {
		panic(err)
	}
	pro.Data = bytesData
	ai.AddProperty(pro)
}

func (ai *AiMaterial) AddInt64PropertyVar(pro AiMaterialProperty, data ...int64) {
	pro = pro.ResetData()
	pro.DataType = pb_msg.AiMaterialPropertyType_AiPropertyTypeInt
	bytesData, err := proto.Marshal(&pb_msg.AiMaterialPropertyInt64{Data: data})
	if err != nil {
		panic(err)
	}
	pro.Data = bytesData
	ai.AddProperty(pro)
}

func (ai *AiMaterial) AddStringPropertyVar(pro AiMaterialProperty, data ...string) {
	pro = pro.ResetData()
	pro.DataType = pb_msg.AiMaterialPropertyType_AiPropertyTypeString
	bytesData, err := proto.Marshal(&pb_msg.AiMaterialPropertyString{Data: data})
	if err != nil {
		panic(err)
	}
	pro.Data = bytesData
	ai.AddProperty(pro)
}

func (ai *AiMaterial) AddProperty(pro AiMaterialProperty) {
	for i, v := range ai.Properties {
		if v.Key == pro.Key && v.Semantic == pro.Semantic && v.Index == pro.Index {
			ai.Properties[i] = pro
			return
		}
	}
	ai.Properties = append(ai.Properties, pro)
}

type AiMaterialProperty struct {
	/** Specifies the name of the property (key)
	 *  Keys are generally case insensitive.
	 */
	Key string

	/** Textures: Specifies their exact usage semantic.
	 * For non-texture properties, this member is always 0
	 * (or, better-said, #AiTextureType_NONE).
	 */
	Semantic int

	/** Textures: Specifies the index of the texture.
	 *  For non-texture properties, this member is always 0.
	 */
	Index int
	/** Type information for the property.
	 *
	 * Defines the data layout inside the data buffer. This is used
	 * by the library internally to perform debug checks and to
	 * utilize proper type conversions.
	 * (It's probably a hacky solution, but it works.)
	 */
	DataType pb_msg.AiMaterialPropertyType
	Data     []byte
}

func (p AiMaterialProperty) ResetData() AiMaterialProperty {
	n := p
	n.Data = n.Data[:0]
	return n
}
func (p AiMaterialProperty) Clone() AiMaterialProperty {
	n := p
	n.Data = make([]byte, len(p.Data))
	copy(n.Data, p.Data)
	return n
}

func newAiMaterialProperty(key string, Type int, index int) AiMaterialProperty {
	return AiMaterialProperty{Key: key, Semantic: Type, Index: index}
}

// ---------------------------------------------------------------------------
/** @brief Defines how an UV channel is transformed.
 *
 *  This is just a helper structure for the #AI_MATKEY_UVTRANSFORM key.
 *  See its documentation for more details.
 *
 *  Typically you'll want to build a matrix of this information. However,
 *  we keep separate scaling/translation/rotation values to make it
 *  easier to process and optimize UV transformations internally.
 */
type AiUVTransform struct {
	/** Translation on the u and v axes.
	 *
	 *  The default value is (0|0).
	 */
	Translation common.AiVector2D

	/** Scaling on the u and v axes.
	 *
	 *  The default value is (1|1).
	 */
	Scaling common.AiVector2D

	/** Rotation - in counter-clockwise direction.
	 *
	 *  The rotation angle is specified in radians. The
	 *  rotation center is 0.5f|0.5f. The default value
	 *  0.f.
	 */
	Rotation float32
}

func (ai AiUVTransform) ToPbMsg() *pb_msg.AiUVTransform {
	return &pb_msg.AiUVTransform{
		Rotation: ai.Rotation,
	}
}
