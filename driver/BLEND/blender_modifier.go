package BLEND

import (
	"assimp/common/logger"
	"assimp/core"
	"errors"
)

// add all available modifiers here

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
	cur := orig_object.modifiers.first.(*SharedModifierData)
	for ; cur != nil; cur = cur.modifier.next.(*SharedModifierData) {
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
		if s != nil || s.name != "ModifierData" {
			logger.WarnF("BlendModifier: expected a ModifierData structure as first member")
			ful++
			continue
		}

		// now, we can be sure that we should be fine to dereference *cur* as
		// ModifierData (with the above note).
		dat := cur.modifier
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
	logger.InfoF("This modifier is not supported, skipping: ", orig_modifier.GetDnaType())
	return
}

type BlenderModifier_Subdivision struct {
	*BlenderModifier
}

func (b *BlenderModifier_Subdivision) IsActive(modin *ModifierData) bool {
	return modin.Type == int32(eModifierType_Subsurf)
}
