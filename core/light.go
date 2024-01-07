package core

import (
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"math"
)

type AiLightSourceType int

const (
	AiLightSource_UNDEFINED AiLightSourceType = 0x0

	//! A directional light source has a well-defined direction
	//! but is infinitely far away. That's quite a good
	//! approximation for sun light.
	AiLightSource_DIRECTIONAL AiLightSourceType = 0x1

	//! A point light source has a well-defined position
	//! in space but no direction - it emits light in all
	//! directions. A normal bulb is a point light.
	AiLightSource_POINT AiLightSourceType = 0x2

	//! A spot light source emits light in a specific
	//! angle. It has a position and a direction it is pointing to.
	//! A good example for a spot light is a light spot in
	//! sport arenas.
	AiLightSource_SPOT AiLightSourceType = 0x3

	//! The generic light level of the world, including the bounces
	//! of all other light sources.
	//! Typically, there's at most one ambient light in a scene.
	//! This light type doesn't have a valid position, direction, or
	//! other properties, just a color.
	AiLightSource_AMBIENT AiLightSourceType = 0x4

	//! An area light is a rectangle with predefined size that uniformly
	//! emits light from one of its sides. The position is center of the
	//! rectangle and direction is its normal vector.
	AiLightSource_AREA AiLightSourceType = 0x5
)

type AiLight struct {
	/** The name of the light source.
	 *
	 *  There must be a node in the scene-graph with the same name.
	 *  This node specifies the position of the light in the scene
	 *  hierarchy and can be animated.
	 */
	Name string

	/** The type of the light source.
	 *
	 * aiLightSource_UNDEFINED is not a valid value for this member.
	 */
	Type AiLightSourceType

	/** Position of the light source in space. Relative to the
	 *  transformation of the node corresponding to the light.
	 *
	 *  The position is undefined for directional lights.
	 */
	Position *common.AiVector3D

	/** Direction of the light source in space. Relative to the
	 *  transformation of the node corresponding to the light.
	 *
	 *  The direction is undefined for point lights. The vector
	 *  may be normalized, but it needn't.
	 */
	Direction *common.AiVector3D

	/** Up direction of the light source in space. Relative to the
	 *  transformation of the node corresponding to the light.
	 *
	 *  The direction is undefined for point lights. The vector
	 *  may be normalized, but it needn't.
	 */
	Up *common.AiVector3D

	/** Constant light attenuation factor.
	 *
	 *  The intensity of the light source at a given distance 'd' from
	 *  the light's position is
	 *  @code
	 *  Atten = 1/( att0 + att1 * d + att2 * d*d)
	 *  @endcode
	 *  This member corresponds to the att0 variable in the equation.
	 *  Naturally undefined for directional lights.
	 */
	AttenuationConstant float32

	/** Linear light attenuation factor.
	 *
	 *  The intensity of the light source at a given distance 'd' from
	 *  the light's position is
	 *  @code
	 *  Atten = 1/( att0 + att1 * d + att2 * d*d)
	 *  @endcode
	 *  This member corresponds to the att1 variable in the equation.
	 *  Naturally undefined for directional lights.
	 */
	AttenuationLinear float32

	/** Quadratic light attenuation factor.
	 *
	 *  The intensity of the light source at a given distance 'd' from
	 *  the light's position is
	 *  @code
	 *  Atten = 1/( att0 + att1 * d + att2 * d*d)
	 *  @endcode
	 *  This member corresponds to the att2 variable in the equation.
	 *  Naturally undefined for directional lights.
	 */
	AttenuationQuadratic float32

	/** Diffuse color of the light source
	 *
	 *  The diffuse light color is multiplied with the diffuse
	 *  material color to obtain the final color that contributes
	 *  to the diffuse shading term.
	 */
	ColorDiffuse *common.AiColor3D

	/** Specular color of the light source
	 *
	 *  The specular light color is multiplied with the specular
	 *  material color to obtain the final color that contributes
	 *  to the specular shading term.
	 */
	ColorSpecular *common.AiColor3D

	/** Ambient color of the light source
	 *
	 *  The ambient light color is multiplied with the ambient
	 *  material color to obtain the final color that contributes
	 *  to the ambient shading term. Most renderers will ignore
	 *  this value it, is just a remaining of the fixed-function pipeline
	 *  that is still supported by quite many file formats.
	 */
	ColorAmbient *common.AiColor3D

	/** Inner angle of a spot light's light cone.
	 *
	 *  The spot light has maximum influence on objects inside this
	 *  angle. The angle is given in radians. It is 2PI for point
	 *  lights and undefined for directional lights.
	 */
	AngleInnerCone float32

	/** Outer angle of a spot light's light cone.
	 *
	 *  The spot light does not affect objects outside this angle.
	 *  The angle is given in radians. It is 2PI for point lights and
	 *  undefined for directional lights. The outer angle must be
	 *  greater than or equal to the inner angle.
	 *  It is assumed that the application uses a smooth
	 *  interpolation between the inner and the outer cone of the
	 *  spot light.
	 */
	AngleOuterCone float32

	/** Size of area light source. */
	Size *common.AiVector2D
}

func (ai *AiLight) FromPbMsg(p *pb_msg.AiLight) *AiLight {
	if ai == nil {
		return nil
	}
	ai.Name = p.Name
	ai.Type = AiLightSourceType(p.Type)
	ai.Position = (&common.AiVector3D{}).FromPbMsg(p.Position)
	ai.Direction = (&common.AiVector3D{}).FromPbMsg(p.Direction)
	ai.Up = (&common.AiVector3D{}).FromPbMsg(p.Up)
	ai.AttenuationConstant = p.AttenuationConstant
	ai.AttenuationLinear = p.AttenuationLinear
	ai.AttenuationQuadratic = p.AttenuationQuadratic
	ai.ColorDiffuse = (&common.AiColor3D{}).FromPbMsg(p.ColorDiffuse)
	ai.ColorSpecular = (&common.AiColor3D{}).FromPbMsg(p.ColorSpecular)
	ai.ColorAmbient = (&common.AiColor3D{}).FromPbMsg(p.ColorAmbient)
	ai.AngleInnerCone = p.AngleInnerCone
	ai.AngleOuterCone = p.AngleOuterCone
	ai.Size = (&common.AiVector2D{}).FromPbMsg(p.Size)
	return ai
}
func (ai *AiLight) Clone() *AiLight {
	if ai == nil {
		return nil
	}
	r := &AiLight{}
	r.Name = ai.Name
	r.Type = ai.Type
	r.Position = ai.Position.Clone()
	r.Direction = ai.Direction.Clone()
	r.Up = ai.Up.Clone()
	r.AttenuationConstant = ai.AttenuationConstant
	r.AttenuationLinear = ai.AttenuationLinear
	r.AttenuationQuadratic = ai.AttenuationQuadratic
	r.ColorDiffuse = ai.ColorDiffuse.Clone()
	r.ColorSpecular = ai.ColorSpecular.Clone()
	r.ColorAmbient = ai.ColorAmbient.Clone()
	r.AngleInnerCone = ai.AngleInnerCone
	r.AngleOuterCone = ai.AngleOuterCone
	r.Size = ai.Size.Clone()
	return r
}
func (ai *AiLight) ToPbMsg() *pb_msg.AiLight {
	if ai == nil {
		return nil
	}
	r := &pb_msg.AiLight{}
	r.Name = ai.Name
	r.Type = int32(ai.Type)
	r.Position = ai.Position.ToPbMsg()
	r.Direction = ai.Direction.ToPbMsg()
	r.Up = ai.Up.ToPbMsg()
	r.AttenuationConstant = ai.AttenuationConstant
	r.AttenuationLinear = ai.AttenuationLinear
	r.AttenuationQuadratic = ai.AttenuationQuadratic
	r.ColorDiffuse = ai.ColorDiffuse.ToPbMsg()
	r.ColorSpecular = ai.ColorSpecular.ToPbMsg()
	r.ColorAmbient = ai.ColorAmbient.ToPbMsg()
	r.AngleInnerCone = ai.AngleInnerCone
	r.AngleOuterCone = ai.AngleOuterCone
	r.Size = ai.Size.ToPbMsg()
	return r
}
func NewAiLight() *AiLight {
	return &AiLight{
		Up:                common.NewAiVector3D3(0, 0, 0),
		Direction:         common.NewAiVector3D3(0, 0, 0),
		Position:          common.NewAiVector3D3(0, 0, 0),
		ColorDiffuse:      common.NewAiColor3D(0, 0, 0),
		ColorSpecular:     common.NewAiColor3D(0, 0, 0),
		ColorAmbient:      common.NewAiColor3D(0, 0, 0),
		AttenuationLinear: 1,
		Type:              AiLightSource_UNDEFINED,
		AngleInnerCone:    math.Pi * 2,
		AngleOuterCone:    math.Pi * 2,
	}
}
