package BLEND

import (
	"assimp/common"
	"assimp/common/logger"
	"errors"
	"reflect"
)

// --------------------------------------------------------------------------------
func IsPtrError(err error, fn func()) error {
	if err != nil && !errors.Is(err, ErrorPtrZero) {
		return err
	}
	if err == nil {
		fn()
	}
	return nil
}

func getName(bytes []byte) string {
	for i, v := range bytes {
		if v == 0 {
			return string(bytes[:i])
		}
	}
	return string(bytes)
}

func (dest *Object) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	temp := int32(0)
	err = s.ReadField(&temp, "type", db)
	if err != nil {
		return err
	}
	dest.Type = ObjectType(temp)
	if err != nil {
		return err
	}
	tmp := make([][]float32, 4)
	for i := range tmp {
		tmp[i] = make([]float32, 4)
	}
	err = s.ReadFieldArray2(SliceToAny2(tmp), "obmat", db)
	if err != nil {
		return err
	}
	for i, v := range tmp {
		copy(dest.obmat[i][:], v)
	}
	tmp = make([][]float32, 4)
	for i := range tmp {
		tmp[i] = make([]float32, 4)
	}
	err = s.ReadFieldArray2(SliceToAny2(tmp), "parentinv", db)
	if err != nil {
		return err
	}
	for i, v := range tmp {
		copy(dest.parentinv[i][:], v)
	}
	tmp2 := make([]uint8, 32)
	err = s.ReadFieldArray(SliceToAny(tmp2), "parsubstr", db)
	if err != nil {
		return err
	}
	dest.parsubstr = getName(tmp2)
	out, _, err := s.ReadFieldPtr("*parent", db)
	if err := IsPtrError(err, func() {
		dest.parent = out.(*Object)
	}); err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*track", db)
	if err := IsPtrError(err, func() {
		dest.track = value.(*Object)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*proxy", db)
	if err := IsPtrError(err, func() {
		if value != nil && !reflect.ValueOf(value).IsNil() {
			dest.proxy = value.(*Object)
		}
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*proxy_from", db)
	if err := IsPtrError(err, func() {
		if value != nil && !reflect.ValueOf(value).IsNil() {
			dest.proxy_from = value.(*Object)
		}
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*proxy_group", db)
	if err := IsPtrError(err, func() {
		if value != nil && !reflect.ValueOf(value).IsNil() {
			dest.proxy_group = value.(*Object)
		}
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*dup_group", db)
	if err := IsPtrError(err, func() {
		if value != nil && !reflect.ValueOf(value).IsNil() {
			dest.dup_group = value.(*Group)
		}
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*data", db)
	if err := IsPtrError(err, func() {
		dest.data = value.(IElemBase)
	}); err != nil {
		return err
	}
	dest.modifiers = &ListBase{}
	err = s.ReadField(dest.modifiers, "modifiers", db)
	if err != nil {
		return err
	}
	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Group) Convert(
	db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.layer, "layer", db)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*gobject", db)
	if err = IsPtrError(err, func() {
		dest.gobject = value.(*GroupObject)
	}); err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *CollectionObject) Convert(
	db *FileDatabase, s *Structure) (err error) {

	value, _, err := s.ReadFieldPtr("*next", db)
	if err = IsPtrError(err, func() {
		dest.next = value.(*CollectionObject)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*ob", db)
	if err = IsPtrError(err, func() {
		dest.ob = value.(*Object)
	}); err != nil {
		return err
	}
	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *CollectionChild) Convert(db *FileDatabase, s *Structure) (err error) {

	value, _, err := s.ReadFieldPtr("*prev", db)
	if err = IsPtrError(err, func() {
		dest.prev = value.(*CollectionChild)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*next", db)
	if err = IsPtrError(err, func() {
		dest.next = value.(*CollectionChild)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*collection", db)
	if err = IsPtrError(err, func() {
		dest.collection = value.(*Collection)
	}); err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Collection) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.gobject, "gobject", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.children, "children", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

// --------------------------------------------------------------------------------
func SliceToAny[T any](in []T) (out []any) {
	for i := range in {
		out = append(out, &in[i])
	}
	return out
}

func SliceToAny2[T any](in [][]T) (out [][]any) {
	out = make([][]any, len(in))
	for i, v := range in {
		out[i] = make([]any, len(v))
		for j := range v {
			out[i][j] = &in[i][j]
		}
	}
	return out
}

func (dest *MTex) Convert(db *FileDatabase, s *Structure) (err error) {

	temp_short := int32(0)
	err = s.ReadField(&temp_short, "mapto", db)
	if err != nil {
		return err
	}
	dest.mapto = MTexMapType(temp_short)
	if err != nil {
		return err
	}
	temp := int32(0)
	err = s.ReadField(&temp, "blendtype", db)
	if err != nil {
		return err
	}
	dest.blendtype = MTexBlendType(temp)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*object", db)
	if err = IsPtrError(err, func() {
		dest.object = value.(*Object)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*tex", db)
	if err = IsPtrError(err, func() {
		dest.tex = value.(*Tex)
	}); err != nil {
		return err
	}

	tmp := make([]uint8, 32)
	err = s.ReadFieldArray(SliceToAny(tmp), "uvname", db)
	if err != nil {
		return err
	}
	dest.uvname = getName(tmp)
	err = s.ReadField(&temp, "projx", db)
	if err != nil {
		return err
	}
	dest.projx = MTexProjection(temp)
	if err != nil {
		return err
	}
	err = s.ReadField(&temp, "projy", db)
	if err != nil {
		return err
	}
	dest.projy = MTexProjection(temp)
	if err != nil {
		return err
	}
	err = s.ReadField(&temp, "projz", db)
	if err != nil {
		return err
	}
	dest.projz = MTexProjection(temp)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mapping, "mapping", db)
	if err != nil {
		return err
	}
	tmp1 := dest.ofs[:]
	err = s.ReadFieldArray(SliceToAny[float32](tmp1), "ofs", db)
	if err != nil {
		return err
	}
	for i := range dest.ofs {
		dest.size[i] = tmp1[i]
	}
	tmp1 = dest.size[:]
	err = s.ReadFieldArray(SliceToAny[float32](tmp1), "size", db)
	if err != nil {
		return err
	}
	for i := range dest.size {
		dest.size[i] = tmp1[i]
	}
	err = s.ReadField(&dest.rot, "rot", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.texflag, "texflag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.colormodel, "colormodel", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pmapto, "pmapto", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pmaptoneg, "pmaptoneg", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.r, "r", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.g, "g", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.b, "b", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.k, "k", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.colspecfac, "colspecfac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mirrfac, "mirrfac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.alphafac, "alphafac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.difffac, "difffac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.specfac, "specfac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.emitfac, "emitfac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.hardfac, "hardfac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.norfac, "norfac", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *TFace) Convert(db *FileDatabase, s *Structure) (err error) {
	tmp := make([][]float64, 4)
	for i := range tmp {
		tmp[i] = make([]float64, 2)
	}
	err = s.ReadFieldArray2(SliceToAny2(tmp), "uv", db)
	if err != nil {
		return err
	}
	for i := range tmp {
		copy(dest.uv[i][:], tmp[i])
	}

	tmp1 := dest.col[:]
	err = s.ReadFieldArray(SliceToAny(tmp1), "col", db)
	if err != nil {
		return err
	}
	for i := range tmp {
		dest.col[i] = tmp1[i]
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mode, "mode", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.tile, "tile", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.unwrap, "unwrap", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *SubsurfModifierData) Convert(db *FileDatabase, s *Structure) (err error) {
	modTmp := &ModifierData{}
	err = s.ReadField(modTmp, "modifier", db)
	if err != nil {
		return err
	}
	dest.SharedModifierData = &SharedModifierData{}
	dest.modifier = modTmp
	err = s.ReadField(&dest.subdivType, "subdivType", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.levels, "levels", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.renderLevels, "renderLevels", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flags, "flags", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MFace) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.v1, "v1", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.v2, "v2", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.v3, "v3", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.v4, "v4", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mat_nr, "mat_nr", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Lamp) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	temp := int32(0)
	err = s.ReadField(&temp, "type", db)
	if err != nil {
		return err
	}
	dest.Type = LampType(temp)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flags, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.colormodel, "colormodel", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totex, "totex", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.r, "r", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.g, "g", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.b, "b", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.k, "k", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.energy, "energy", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.dist, "dist", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.spotsize, "spotsize", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.spotblend, "spotblend", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.constant_coefficient, "coeff_const", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.linear_coefficient, "coeff_lin", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.quadratic_coefficient, "coeff_quad", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.att1, "att1", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.att2, "att2", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&temp, "falloff_type", db)
	if err != nil {
		return err
	}
	dest.falloff_type = LampFalloffType(temp)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sun_brightness, "sun_brightness", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.area_size, "area_size", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.area_sizey, "area_sizey", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.area_sizez, "area_sizez", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.area_shape, "area_shape", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MDeformWeight) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.def_nr, "def_nr", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.weight, "weight", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *PackedFile) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.size, "size", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.seek, "seek", db)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldFileOffsetPtr("*data", db)
	if err != nil {
		return err
	}
	dest.data = value.(*FileOffset)
	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Base) Convert(
	db *FileDatabase, s *Structure) (err error) {
	// note: as per https://github.com/assimp/assimp/issues/128,
	// reading the Object linked list recursively is prone to stack overflow.
	// This structure converter is therefore a hand-written exception that
	// does it iteratively.

	initial_pos := db.GetCurPos()
	if err != nil {
		return err
	}

	todo := common.NewPair(dest, initial_pos)
	if err != nil {
		return err
	}
	for {

		cur_dest := todo.First
		if err != nil {
			return err
		}
		db.SetCurPos(todo.Second)
		if err != nil {
			return err
		}

		// we know that this is a double-linked, circular list which we never
		// traverse backwards, so don't bother resolving the back links.
		cur_dest.prev = nil
		if err != nil {
			return err
		}

		value, _, err := s.ReadFieldPtr("*object", db)
		if err = IsPtrError(err, func() {
			cur_dest.object = value.(*Object)
		}); err != nil {
			return err
		}

		// the return value of value,err =s.ReadFieldPtr indicates whether the object
		// was already cached. In this case, we don't need to resolve
		// it again.
		value, fromCache, err1 := s.ReadFieldPtr("*next", db, true)
		if err1 == nil {
			cur_dest.next = value.(*Base)
		} else if !errors.Is(err1, ErrorPtrZero) {
			return err1
		}
		if cur_dest.next != nil && !fromCache {
			todo = common.NewPair(cur_dest.next, db.GetCurPos())
			if err != nil {
				return err
			}
			continue
		}

		break
	}
	db.SetCurPos(initial_pos + int(s.size))
	return nil
}

//--------------------------------------------------------------------------------

func (dest *MTFace) Convert(db *FileDatabase, s *Structure) (err error) {
	tmp := make([][]float32, 4)
	for i := range tmp {
		tmp[i] = make([]float32, 2)
	}
	err = s.ReadFieldArray2(SliceToAny2(tmp), "uv", db)
	if err != nil {
		return err
	}
	for i := range tmp {
		copy(dest.uv[i][:], tmp[i])
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mode, "mode", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.tile, "tile", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.unwrap, "unwrap", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

// --------------------------------------------------------------------------------
func SliceToT[T any](out []any) (value []T) {
	for _, v := range out {
		va, ok := v.(T)
		if ok {
			value = append(value, va)
		} else {
			var tmp T
			value = append(value, tmp)
		}

	}
	return value
}

func (dest *Material) Convert(db *FileDatabase, s *Structure) (err error) {
	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.r, "r", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.g, "g", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.b, "b", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.specr, "specr", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.specg, "specg", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.specb, "specb", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.har, "har", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ambr, "ambr", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ambg, "ambg", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ambb, "ambb", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mirr, "mirr", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mirg, "mirg", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mirb, "mirb", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.emit, "emit", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ray_mirror, "ray_mirror", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.alpha, "alpha", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ref, "ref", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.translucency, "translucency", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mode, "mode", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.roughness, "roughness", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.darkness, "darkness", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.refrac, "refrac", db)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*group", db)
	if err = IsPtrError(err, func() {
		if value == nil {
			dest.group = nil
		} else {
			dest.group = value.(*Group)
		}

	}); err != nil {
		return err
	}

	err = s.ReadField(&dest.diff_shader, "diff_shader", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.spec_shader, "spec_shader", db)
	if err != nil {
		return err
	}

	value1, _, err := s.ReadFieldPtrArray(len(dest.mtex), "*mtex", db)
	if err = IsPtrError(err, func() {
		if value1 == nil {
			dest.group = nil
		} else {
			copy(dest.mtex[:], SliceToT[*MTex](value1))
		}
	}); err != nil {
		return err
	}

	err = s.ReadField(&dest.amb, "amb", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ang, "ang", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.spectra, "spectra", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.spec, "spec", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.zoffs, "zoffs", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.add, "add", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.fresnel_mir, "fresnel_mir", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.fresnel_mir_i, "fresnel_mir_i", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.fresnel_tra, "fresnel_tra", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.fresnel_tra_i, "fresnel_tra_i", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.filter, "filter", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.tx_limit, "tx_limit", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.tx_falloff, "tx_falloff", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.gloss_mir, "gloss_mir", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.gloss_tra, "gloss_tra", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.adapt_thresh_mir, "adapt_thresh_mir", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.adapt_thresh_tra, "adapt_thresh_tra", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.aniso_gloss_mir, "aniso_gloss_mir", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.dist_mir, "dist_mir", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.hasize, "hasize", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flaresize, "flaresize", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.subsize, "subsize", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flareboost, "flareboost", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.strand_sta, "strand_sta", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.strand_end, "strand_end", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.strand_ease, "strand_ease", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.strand_surfnor, "strand_surfnor", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.strand_min, "strand_min", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.strand_widthfade, "strand_widthfade", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sbias, "sbias", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.lbias, "lbias", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.shad_alpha, "shad_alpha", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.param, "param", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.rms, "rms", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.rampfac_col, "rampfac_col", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.rampfac_spec, "rampfac_spec", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.friction, "friction", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.fh, "fh", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.reflect, "reflect", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.fhdist, "fhdist", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.xyfrict, "xyfrict", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_radius, "sss_radius", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_col, "sss_col", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_error, "sss_error", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_scale, "sss_scale", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_ior, "sss_ior", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_colfac, "sss_colfac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_texfac, "sss_texfac", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_front, "sss_front", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_back, "sss_back", db)
	if err != nil {
		return err
	}

	err = s.ReadField(&dest.material_type, "material_type", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ray_depth, "ray_depth", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ray_depth_tra, "ray_depth_tra", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.samp_gloss_mir, "samp_gloss_mir", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.samp_gloss_tra, "samp_gloss_tra", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.fadeto_mir, "fadeto_mir", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.shade_flag, "shade_flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flarec, "flarec", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.starc, "starc", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.linec, "linec", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ringc, "ringc", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pr_lamp, "pr_lamp", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pr_texture, "pr_texture", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ml_flag, "ml_flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.diff_shader, "diff_shader", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.spec_shader, "spec_shader", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.texco, "texco", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mapto, "mapto", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.ramp_show, "ramp_show", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pad3, "pad3", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.dynamode, "dynamode", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pad2, "pad2", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_flag, "sss_flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sss_preset, "sss_preset", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.shadowonly_flag, "shadowonly_flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.index, "index", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.vcol_alpha, "vcol_alpha", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pad4, "pad4", db)
	if err != nil {
		return err
	}

	err = s.ReadField(&dest.seed1, "seed1", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.seed2, "seed2", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MTexPoly) Convert(db *FileDatabase, s *Structure) (err error) {
	value, _, err := s.ReadFieldPtr("*tpage", db)
	if err = IsPtrError(err, func() {
		dest.tpage = value.(*Image)
	}); err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.transp, "transp", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mode, "mode", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.tile, "tile", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pad, "pad", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Mesh) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totface, "totface", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totedge, "totedge", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totvert, "totvert", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totloop, "totloop", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totpoly, "totpoly", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.subdiv, "subdiv", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.subdivr, "subdivr", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.subsurftype, "subsurftype", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.smoothresh, "smoothresh", db)
	if err != nil {
		return err
	}
	value1, err := s.ReadFieldPtrSlice("*mface", db)
	if err != nil {
		return err
	}
	dest.mface = SliceToT[*MFace](value1)
	value1, err = s.ReadFieldPtrSlice("*mtface", db)
	if err != nil {
		return err
	}
	dest.mtface = SliceToT[*MTFace](value1)
	value1, err = s.ReadFieldPtrSlice("*tface", db)
	if err != nil {
		return err
	}
	dest.tface = SliceToT[*TFace](value1)
	value1, err = s.ReadFieldPtrSlice("*mvert", db)
	if err != nil {
		return err
	}
	dest.mvert = SliceToT[*MVert](value1)
	value1, err = s.ReadFieldPtrSlice("*medge", db)
	if err != nil {
		return err
	}
	dest.medge = SliceToT[*MEdge](value1)
	value1, err = s.ReadFieldPtrSlice("*mloop", db)
	if err != nil {
		return err
	}
	dest.mloop = SliceToT[*MLoop](value1)
	value1, err = s.ReadFieldPtrSlice("*mloopuv", db)
	if err != nil {
		return err
	}
	dest.mloopuv = SliceToT[*MLoopUV](value1)
	value1, err = s.ReadFieldPtrSlice("*mloopcol", db)
	if err != nil {
		return err
	}
	dest.mloopcol = SliceToT[*MLoopCol](value1)
	value1, err = s.ReadFieldPtrSlice("*mpoly", db)
	if err != nil {
		return err
	}
	dest.mpoly = SliceToT[*MPoly](value1)
	value1, err = s.ReadFieldPtrSlice("*mtpoly", db)
	if err != nil {
		return err
	}
	dest.mtpoly = SliceToT[*MTexPoly](value1)
	value1, err = s.ReadFieldPtrSlice("*dvert", db)
	if err != nil {
		return err
	}
	dest.dvert = SliceToT[*MDeformVert](value1)
	value1, err = s.ReadFieldPtrSlice("*mcol", db)
	if err != nil {
		return err
	}
	dest.mcol = SliceToT[*MCol](value1)
	value2, _, err := s.ReadFieldPtrPtr("**mat", db)
	if err := IsPtrError(err, func() {
		tmp := value2.([]any)
		for _, v := range tmp {
			dest.mat = append(dest.mat, v.(*Material))
		}
	}); err != nil {
		return err
	}
	dest.vdata = &CustomData{}
	err = s.ReadField(dest.vdata, "vdata", db)
	if err != nil {
		return err
	}
	dest.edata = &CustomData{}
	err = s.ReadField(dest.edata, "edata", db)
	if err != nil {
		return err
	}
	dest.fdata = &CustomData{}
	err = s.ReadField(dest.fdata, "fdata", db)
	if err != nil {
		return err
	}
	dest.pdata = &CustomData{}
	err = s.ReadField(dest.pdata, "pdata", db)
	if err != nil {
		return err
	}
	dest.ldata = &CustomData{}
	err = s.ReadField(dest.ldata, "ldata", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MDeformVert) Convert(db *FileDatabase, s *Structure) (err error) {

	value, err := s.ReadFieldPtrSlice("*dw", db)
	if err != nil {
		return err
	}
	dest.dw = SliceToT[*MDeformWeight](value)
	err = s.ReadField(&dest.totweight, "totweight", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *World) Convert(db *FileDatabase, s *Structure) (err error) {
	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MLoopCol) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.r, "r", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.g, "g", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.b, "b", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.a, "a", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MVert) Convert(
	db *FileDatabase, s *Structure) (err error) {
	tmp := dest.co[:]
	err = s.ReadFieldArray(SliceToAny(tmp), "co", db)
	if err != nil {
		return err
	}
	for i := range dest.co {
		dest.co[i] = tmp[i]
	}
	tmp = dest.no[:]
	err = s.ReadFieldArray(SliceToAny(tmp), "no", db)
	if err != nil {
		return err
	}
	for i := range dest.no {
		dest.no[i] = tmp[i]
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	//err =s.ReadField(&dest.mat_nr,"mat_nr",db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.bweight, "bweight", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MEdge) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.v1, "v1", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.v2, "v2", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.crease, "crease", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.bweight, "bweight", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MLoopUV) Convert(db *FileDatabase, s *Structure) (err error) {
	tmp := dest.uv[:]
	err = s.ReadFieldArray(SliceToAny(tmp), "uv", db)
	if err != nil {
		return err
	}
	for i, v := range tmp {
		dest.uv[i] = v
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *GroupObject) Convert(db *FileDatabase, s *Structure) (err error) {

	value, _, err := s.ReadFieldPtr("*prev", db)
	if err = IsPtrError(err, func() {
		dest.prev = value.(*GroupObject)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*next", db)
	if err = IsPtrError(err, func() {
		dest.next = value.(*GroupObject)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*ob", db)
	if err = IsPtrError(err, func() {
		dest.ob = value.(Object)
	}); err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *ListBase) Convert(db *FileDatabase, s *Structure) (err error) {

	value, _, err := s.ReadFieldPtr("*first", db)
	if err = IsPtrError(err, func() {
		dest.first = value.(IElemBase)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*last", db)
	if err = IsPtrError(err, func() {
		dest.last = value.(IElemBase)
	}); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MLoop) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.v, "v", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.e, "e", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *ModifierData) Convert(db *FileDatabase, s *Structure) (err error) {

	value, _, err := s.ReadFieldPtr("*next", db)
	if err = IsPtrError(err, func() {
		dest.next = value.(IElemBase)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*prev", db)
	if err = IsPtrError(err, func() {
		dest.prev = value.(IElemBase)
	}); err != nil {
		return err
	}
	err = s.ReadField(&dest.Type, "type", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mode, "mode", db)
	if err != nil {
		return err
	}
	tmp := make([]uint8, 32)
	err = s.ReadFieldArray(SliceToAny(tmp), "name", db)
	if err != nil {
		return err
	}
	dest.name = getName(tmp)
	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *ID) Convert(db *FileDatabase, s *Structure) (err error) {
	tmp := make([]uint8, 1024)
	err = s.ReadFieldArray(SliceToAny(tmp), "name", db)
	dest.name = getName(tmp)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MCol) Convert(
	db *FileDatabase, s *Structure) (err error) {
	err = s.ReadField(&dest.r, "r", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.g, "g", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.b, "b", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.a, "a", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MPoly) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.loopstart, "loopstart", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totloop, "totloop", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.mat_nr, "mat_nr", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Scene) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*camera", db)
	if err = IsPtrError(err, func() {
		dest.camera = value.(*Object)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*world", db)
	if err = IsPtrError(err, func() {
		dest.world = value.(*World)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*basact", db)
	if err = IsPtrError(err, func() {
		dest.basact = value.(*Base)
	}); err != nil {
		return err
	}

	value, _, err = s.ReadFieldPtr("*master_collection", db)
	if err = IsPtrError(err, func() {
		if value != nil {
			dest.master_collection = value.(*Collection)
		}
	}); err != nil {
		return err
	}

	err = s.ReadField(&dest.base, "base", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Library) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	tmp := make([]uint8, 240)
	err = s.ReadFieldArray(SliceToAny(tmp), "name", db)
	if err != nil {
		return err
	}
	dest.name = getName(tmp)
	tmp = make([]uint8, 240)
	err = s.ReadFieldArray(SliceToAny(tmp), "filename", db)
	if err != nil {
		return err
	}
	dest.filename = getName(tmp)
	value, _, err := s.ReadFieldPtr("*parent", db)
	if err = IsPtrError(err, func() {
		dest.parent = value.(*Library)
	}); err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Tex) Convert(db *FileDatabase, s *Structure) (err error) {
	temp_short := int16(0)
	err = s.ReadField(&temp_short, "imaflag", db)
	if err != nil {
		return err
	}
	dest.imaflag = TexImageFlags(temp_short)
	if err != nil {
		return err
	}
	temp := int32(0)
	err = s.ReadField(&temp, "type", db)
	if err != nil {
		return err
	}
	dest.Type = TexType(temp)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*ima", db)
	if err = IsPtrError(err, func() {
		dest.ima = value.(*Image)
	}); err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Camera) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	temp := int32(0)
	err = s.ReadField(&temp, "type", db)
	if err != nil {
		return err
	}
	dest.Type = CameraType(temp)
	if err != nil {
		return err
	}
	err = s.ReadField(&temp, "flag", db)
	if err != nil {
		return err
	}
	dest.flag = CameraType(temp)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.lens, "lens", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.sensor_x, "sensor_x", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.clipsta, "clipsta", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.clipend, "clipend", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *MirrorModifierData) Convert(
	db *FileDatabase, s *Structure) (err error) {
	modTmp := &ModifierData{}
	dest.SharedModifierData = &SharedModifierData{}
	err = s.ReadField(modTmp, "modifier", db)
	if err != nil {
		return err
	}
	dest.modifier = modTmp
	err = s.ReadField(&dest.axis, "axis", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.tolerance, "tolerance", db)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*mirror_ob", db)
	if err = IsPtrError(err, func() {
		dest.mirror_ob = value.(*Object)
	}); err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *Image) Convert(db *FileDatabase, s *Structure) (err error) {

	err = s.ReadField(&dest.id, "id", db)
	if err != nil {
		return err
	}
	tmp := make([]uint8, 240)
	err = s.ReadFieldArray(SliceToAny(tmp), "name", db)
	if err != nil {
		return err
	}
	dest.name = getName(tmp)
	err = s.ReadField(&dest.ok, "ok", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.source, "source", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.Type, "type", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pad, "pad", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.pad1, "pad1", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.lastframe, "lastframe", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.tpageflag, "tpageflag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totbind, "totbind", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.xrep, "xrep", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.yrep, "yrep", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.twsta, "twsta", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.twend, "twend", db)
	if err != nil {
		return err
	}
	value, _, err := s.ReadFieldPtr("*packedfile", db)
	if err := IsPtrError(err, func() {
		dest.packedfile = value.(*PackedFile)
	}); err != nil {
		return err
	}
	err = s.ReadField(&dest.lastupdate, "lastupdate", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.lastused, "lastused", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.animspeed, "animspeed", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.gen_x, "gen_x", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.gen_y, "gen_y", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.gen_type, "gen_type", db)
	if err != nil {
		return err
	}

	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *CustomData) Convert(db *FileDatabase, s *Structure) (err error) {
	tmp := dest.typemap[:]
	err = s.ReadFieldArray(SliceToAny(tmp), "typemap", db)
	if err != nil {
		return err
	}
	for i, v := range tmp {
		dest.typemap[i] = v
	}
	err = s.ReadField(&dest.totlayer, "totlayer", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.maxlayer, "maxlayer", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.totsize, "totsize", db)
	if err != nil {
		return err
	}
	out, err := s.ReadFieldPtrSlice("*layers", db)
	if err != nil {
		return err
	}
	dest.layers = SliceToT[*CustomDataLayer](out)
	return db.Discard(int(s.size))
}

//--------------------------------------------------------------------------------

func (dest *CustomDataLayer) Convert(
	db *FileDatabase, s *Structure) (err error) {
	var tmp1 int32
	err = s.ReadField(&tmp1, "type", db)
	if err != nil {
		return err
	}
	dest.Type = CustomDataType(tmp1)
	err = s.ReadField(&dest.offset, "offset", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.flag, "flag", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.active, "active", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.active_rnd, "active_rnd", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.active_clone, "active_clone", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.active_mask, "active_mask", db)
	if err != nil {
		return err
	}
	err = s.ReadField(&dest.uid, "uid", db)
	if err != nil {
		return err
	}
	tmp := make([]uint8, 64)
	err = s.ReadFieldArray(SliceToAny(tmp), "name", db)
	if err != nil {
		return err
	}
	dest.name = getName(tmp)
	datas, err := s.ReadCustomDataPtr(int(dest.Type), "*data", db)
	if err != nil {
		return err
	}
	if len(datas) > 0 {
		logger.Warn("CustomDataLayer find cnt:%v, default get index 0")
		dest.data = datas[0]
	}
	return db.Discard(int(s.size))
}

// --------------------------------------------------------
/** Fill the @c converters member with converters for all
 *  known data types. The implementation of this method is
 *  in BlenderScene.cpp and is machine-generated.
 *  Converters are used to quickly handle objects whose
 *  exact data type is a runtime-property and not yet
 *  known at compile time (consider Object::data).*/
func (d *DNA) RegisterConverters() {
	d.converters["Object"] = func() Converter {
		return &Object{ElemBase: &ElemBase{}}
	}
	d.converters["Group"] = func() Converter {
		return &Group{ElemBase: &ElemBase{}}
	}
	d.converters["MTex"] = func() Converter {
		return &MTex{ElemBase: &ElemBase{}}
	}
	d.converters["TFace"] = func() Converter {
		return &TFace{ElemBase: &ElemBase{}}
	}
	d.converters["SubsurfModifierData"] = func() Converter {
		return &SubsurfModifierData{ElemBase: &ElemBase{}}
	}
	d.converters["MFace"] = func() Converter {
		return &MFace{ElemBase: &ElemBase{}}
	}
	d.converters["Lamp"] = func() Converter {
		return &Lamp{ElemBase: &ElemBase{}}
	}
	d.converters["MDeformWeight"] = func() Converter {
		return &MDeformWeight{ElemBase: &ElemBase{}}
	}
	d.converters["PackedFile"] = func() Converter {
		return &PackedFile{ElemBase: &ElemBase{}}
	}
	d.converters["Base"] = func() Converter {
		return &Base{ElemBase: &ElemBase{}}
	}
	d.converters["MTFace"] = func() Converter {
		return &MTFace{ElemBase: &ElemBase{}}
	}
	d.converters["Material"] = func() Converter {
		return &Material{ElemBase: &ElemBase{}}
	}
	d.converters["MTexPoly"] = func() Converter {
		return &MTexPoly{ElemBase: &ElemBase{}}
	}
	d.converters["Mesh"] = func() Converter {
		return &Mesh{ElemBase: &ElemBase{}}
	}
	d.converters["MDeformVert"] = func() Converter {
		return &MDeformVert{ElemBase: &ElemBase{}}
	}
	d.converters["World"] = func() Converter {
		return &World{ElemBase: &ElemBase{}}
	}
	d.converters["MLoopCol"] = func() Converter {
		return &MLoopCol{ElemBase: &ElemBase{}}
	}
	d.converters["MVert"] = func() Converter {
		return &MVert{ElemBase: &ElemBase{}}
	}
	d.converters["MEdge"] = func() Converter {
		return &MEdge{ElemBase: &ElemBase{}}
	}
	d.converters["MLoopUV"] = func() Converter {
		return &MLoopUV{ElemBase: &ElemBase{}}
	}
	d.converters["GroupObject"] = func() Converter {
		return &GroupObject{ElemBase: &ElemBase{}}
	}
	d.converters["ListBase"] = func() Converter {
		return &ListBase{ElemBase: &ElemBase{}}
	}
	d.converters["MLoop"] = func() Converter {
		return &MLoop{ElemBase: &ElemBase{}}
	}
	d.converters["ModifierData"] = func() Converter {
		return &ModifierData{ElemBase: &ElemBase{}}
	}
	d.converters["ID"] = func() Converter {
		return &ID{ElemBase: &ElemBase{}}
	}
	d.converters["MCol"] = func() Converter {
		return &MCol{ElemBase: &ElemBase{}}
	}
	d.converters["MPoly"] = func() Converter {
		return &MPoly{ElemBase: &ElemBase{}}
	}
	d.converters["Scene"] = func() Converter {
		return &Scene{ElemBase: &ElemBase{}}
	}
	d.converters["Library"] = func() Converter {
		return &Library{ElemBase: &ElemBase{}}
	}
	d.converters["Tex"] = func() Converter {
		return &Tex{ElemBase: &ElemBase{}}
	}
	d.converters["Camera"] = func() Converter {
		return &Camera{ElemBase: &ElemBase{}}
	}
	d.converters["MirrorModifierData"] = func() Converter {
		return &MirrorModifierData{ElemBase: &ElemBase{}}
	}
	d.converters["Image"] = func() Converter {
		return &Image{ElemBase: &ElemBase{}}
	}
	d.converters["CustomData"] = func() Converter {
		return &CustomData{ElemBase: &ElemBase{}}
	}
	d.converters["CustomDataLayer"] = func() Converter {
		return &CustomDataLayer{ElemBase: &ElemBase{}}
	}
	d.converters["Collection"] = func() Converter {
		return &Collection{ElemBase: &ElemBase{}}
	}
	d.converters["CollectionChild"] = func() Converter {
		return &CollectionChild{ElemBase: &ElemBase{}}
	}
	d.converters["CollectionObject"] = func() Converter {
		return &CollectionObject{ElemBase: &ElemBase{}}
	}
}
