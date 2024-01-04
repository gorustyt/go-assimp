package BLEND

import (
	"assimp/common"
	"assimp/common/logger"
	"assimp/core"
	"errors"
	"math"
)

// add all available modifiers here
const (
	TYPE_CatmullClarke = 0x0
	TYPE_Simple        = 0x1
)

var creators = []func() IBlenderModifier{
	func() IBlenderModifier { return &BlenderModifier_Mirror{} },
	func() IBlenderModifier { return &BlenderModifier_Subdivision{} },
	nil, // sentinel
}

// -------------------------------------------------------------------------------------------
/**
 *  Manage all known modifiers and instance and apply them if necessary
 */
// -------------------------------------------------------------------------------------------
type BlenderModifierShowcase struct {
	cached_modifiers []IBlenderModifier
}

func (b *BlenderModifierShowcase) ApplyModifiers(out core.AiNode,
	conv_data *ConversionData,
	in *Scene,
	orig_object *Object,
) error {
	cnt := 0
	ful := 0

	// NOTE: this cast is potentially unsafe by design, so we need to perform type checks before
	// we're allowed to dereference the pointers without risking to crash. We might still be
	// invoking UB btw - we're assuming that the ModifierData member of the respective modifier
	// structures is at offset sizeof(vftable) with no padding.
	var cur IElemBase
	if orig_object.modifiers != nil && orig_object.modifiers.first != nil {
		cur = orig_object.modifiers.first
	}
	var curv *SharedModifierData
	if cur != nil {
		curv = cur.(*SubsurfModifierData).SharedModifierData
	}
	for ; cur != nil && curv != nil; cur = curv.modifier.next {
		curv = cur.(*SubsurfModifierData).SharedModifierData
		if cur.GetDnaType() == "" {
			return errors.New("invalid GetDnaType")
		}

		s := conv_data.db.dna.Get(cur.GetDnaType())
		if s == nil {
			logger.WarnF("BlendModifier: could not resolve DNA name: ", cur.GetDnaType())
			ful++
			continue
		}

		// this is a common trait of all XXXMirrorData structures in BlenderDNA
		f := s.Get("modifier")
		if f == nil || f.offset != 0 {
			logger.Info("BlendModifier: expected a `modifier` member at offset 0")
			ful++
			continue
		}

		s = conv_data.db.dna.Get(f.Type)
		if s == nil || s.name != "ModifierData" {
			logger.WarnF("BlendModifier: expected a ModifierData structure as first member")
			ful++
			continue
		}

		// now, we can be sure that we should be fine to dereference *cur* as
		// ModifierData (with the above note).
		dat := curv.modifier
		curgod := 0
		curmod := 0
		endmod := len(b.cached_modifiers)

		for creators[curgod] != nil { // allocate modifiers on the fly
			if curmod == endmod {
				b.cached_modifiers = append(b.cached_modifiers, creators[curgod]())

				endmod = len(b.cached_modifiers)
				curmod = endmod - 1
			}

			modifier := b.cached_modifiers[curmod]
			if modifier.IsActive(dat) {
				modifier.DoIt(out, conv_data, cur, in, orig_object)
				cnt++

				curgod = -1
				break
			}
			curgod++
			curmod++
		}
		if curgod == -1 || creators[curgod] == nil {
			logger.WarnF("Couldn't find a handler for modifier: ", dat.name)
		}
		ful++
	}

	// Even though we managed to resolve some or all of the modifiers on this
	// object, we still can't say whether our modifier implementations were
	// able to fully do their job.
	if ful != 0 {
		logger.DebugF("BlendModifier: found handlers for%v of%v modifiers on `%v`, check log messages above for errors", cnt, ful, orig_object.id.name)
	}
	return nil
}

type IBlenderModifier interface {
	IsActive(modin *ModifierData) bool
	DoIt(out core.AiNode,
		conv_data *ConversionData,
		orig_modifier IElemBase,
		in *Scene,
		orig_object *Object,
	)
}

type BlenderModifier struct {
}

func (b *BlenderModifier) IsActive(modin *ModifierData) bool {
	return false
}

func (b *BlenderModifier) DoIt(out core.AiNode,
	conv_data *ConversionData,
	orig_modifier IElemBase,
	in *Scene,
	orig_object *Object,
) {
	logger.InfoF("This modifier is not supported, skipping: ", orig_modifier.GetDnaType())
	return
}

type BlenderModifier_Mirror struct {
	*BlenderModifier
}

func (b *BlenderModifier_Mirror) IsActive(modin *ModifierData) bool {
	return modin.Type == int32(eModifierType_Mirror)
}

func (b *BlenderModifier_Mirror) DoIt(out core.AiNode,
	conv_data *ConversionData,
	orig_modifier IElemBase,
	in *Scene,
	orig_object *Object,
) {
	// hijacking the ABI, see the big note in BlenderModifierShowcase::ApplyModifiers()
	mir := orig_modifier.(*MirrorModifierData)
	common.AiAssert(mir.modifier.Type == int32(eModifierType_Mirror))
	mirror_ob := mir.mirror_ob
	// XXX not entirely correct, mirroring on two axes results in 4 distinct objects in blender ...

	// take all input meshes and clone them
	for i := 0; i < len(out.Meshes); i++ {
		mesh := conv_data.meshes[out.Meshes[i]].Clone()
		xs := float32(1)
		if mir.flag&Flags_AXIS_X != 0 {
			xs = -1.
		}
		ys := float32(1)
		if mir.flag&Flags_AXIS_Y != 0 {
			ys = -1.
		}
		zs := float32(1)
		if mir.flag&Flags_AXIS_Z != 0 {
			zs = -1.
		}
		if mirror_ob != nil {
			center := common.NewAiVector3D3(mirror_ob.obmat[3][0], mirror_ob.obmat[3][1], mirror_ob.obmat[3][2])
			for j := 0; j < len(mesh.Vertices); j++ {
				v := mesh.Vertices[j]

				v.X = center.X + xs*(center.X-v.X)
				v.Y = center.Y + ys*(center.Y-v.Y)
				v.Z = center.Z + zs*(center.Z-v.Z)
			}
		} else {
			for j := 0; j < len(mesh.Vertices); j++ {
				v := mesh.Vertices[j]
				v.X *= xs
				v.Y *= ys
				v.Z *= zs
			}
		}

		if len(mesh.Normals) > 0 {
			for j := 0; j < len(mesh.Vertices); j++ {
				v := mesh.Normals[j]
				v.X *= xs
				v.Y *= ys
				v.Z *= zs
			}
		}

		if len(mesh.Tangents) > 0 {
			for j := 0; j < len(mesh.Vertices); j++ {
				v := mesh.Tangents[j]
				v.X *= xs
				v.Y *= ys
				v.Z *= zs
			}
		}

		if len(mesh.Bitangents) > 0 {
			for j := 0; j < len(mesh.Vertices); j++ {
				v := mesh.Bitangents[j]
				v.X *= xs
				v.Y *= ys
				v.Z *= zs
			}
		}

		us := float32(1.)
		if mir.flag&Flags_MIRROR_U != 0 {
			us = -1.
		}
		vs := float32(1.)
		if mir.flag&Flags_MIRROR_V != 0 {
			vs = -1.
		}

		for n := 0; mesh.HasTextureCoords(uint32(n)); n++ {
			for j := 0; j < len(mesh.Vertices); j++ {
				v := mesh.TextureCoords[n][j]
				v.X *= us
				v.Y *= vs
			}
		}

		// Only reverse the winding order if an odd number of axes were mirrored.
		if xs*ys*zs < 0 {
			for j := 0; j < len(mesh.Faces); j++ {
				face := mesh.Faces[j]
				for fi := 0; fi < len(face.Indices)/2; fi++ {
					face.Indices[fi], face.Indices[len(face.Indices)-1-fi] = face.Indices[len(face.Indices)-1-fi], face.Indices[fi]

				}

			}
		}

		conv_data.meshes = append(conv_data.meshes, mesh)
	}
	nind := make([]int32, len(out.Meshes)*2)
	copy(nind, out.Meshes[:len(out.Meshes)])
	for i := 0; i < len(out.Meshes); i++ {
		nind[len(out.Meshes)+i] = int32(len(out.Meshes)) + out.Meshes[i]
	}
	out.Meshes = nind
	logger.InfoF("BlendModifier: Applied the `Mirror` modifier to `%v %v ",
		orig_object.id.name, "`")
	return
}

type BlenderModifier_Subdivision struct {
	*BlenderModifier
}

func (b *BlenderModifier_Subdivision) IsActive(modin *ModifierData) bool {
	return modin.Type == int32(eModifierType_Subsurf)
}
func (b *BlenderModifier_Subdivision) DoIt(out core.AiNode,
	conv_data *ConversionData,
	orig_modifier IElemBase,
	in *Scene,
	orig_object *Object,
) {
	// hijacking the ABI, see the big note in BlenderModifierShowcase::ApplyModifiers()
	mir := orig_modifier.(*SubsurfModifierData)
	common.AiAssert(mir.modifier.Type == int32(eModifierType_Subsurf))

	var algo core.Algorithm
	switch mir.subdivType {
	case TYPE_CatmullClarke:
		algo = core.CATMULL_CLARKE
	case TYPE_Simple:
		logger.Warn("BlendModifier: The `SIMPLE` subdivision algorithm is not currently implemented, using Catmull-Clarke")
		algo = core.CATMULL_CLARKE
	default:
		logger.WarnF("BlendModifier: Unrecognized subdivision algorithm: %v", mir.subdivType)
		return
	}

	subd := core.NewSubDivision(algo)
	if len(conv_data.meshes) == 0 {
		return
	}
	meshIndex := len(conv_data.meshes) - len(out.Meshes)
	if meshIndex >= len(conv_data.meshes) {
		logger.Error("Invalid index detected.")
		return
	}
	meshes := conv_data.meshes[len(conv_data.meshes)-len(out.Meshes):]
	tempmeshes := make([]*core.AiMesh, len(out.Meshes))
	subd.Subdivide(meshes, tempmeshes, int(math.Max(float64(mir.renderLevels), float64(mir.levels))), true)
	copy(conv_data.meshes[:len(out.Meshes)], tempmeshes)
	logger.InfoF("BlendModifier: Applied the `Subdivision` modifier to `%v %v",
		orig_object.id.name, "`")
	return
}
