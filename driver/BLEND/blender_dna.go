package BLEND

import (
	"assimp/common"
	"assimp/common/logger"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
)

// --------------------------------------------------------
/** Locate the DNA in the file and parse it. The input
 *  stream is expected to poto the beginning of the DN1
 *  chunk at the time this method is called and is
 *  undefined afterwards.
 *  @throw DeadlyImportError if the DNA cannot be read.
 *  @note The position of the stream pointer is undefined
 *    afterwards.*/

func (d *DNAParser) Parse() error {
	dna := NewDNA()
	d.db.dna = dna
	s, err := d.GetString(4)
	if err != nil {
		return err
	}
	if s != "SDNA" {
		return errors.New("BlenderDNA: Expected SDNA chunk")
	}
	s, err = d.GetString(4)
	if err != nil {
		return err
	}
	if s != "NAME" {
		return errors.New("BlenderDNA: Expected NAME field")
	}
	namesLen, err := d.GetInt32()
	if err != nil {
		return err
	}
	names := make([]string, namesLen)
	for i := range names {
		var tmp []byte
		for {
			c, err := d.GetUInt8()
			if err != nil {
				return err
			}
			if c == 0 {
				break
			}
			tmp = append(tmp, c)
		}
		names[i] += string(tmp)
	}
	skip := func() error {
		for d.GetCurPos()&0x3 != 0 {
			_, err = d.GetInt8()
			if err != nil {
				return err
			}
		}
		return nil
	}
	if err = skip(); err != nil {
		return err
	}
	s, err = d.GetString(4)
	if err != nil {
		return err
	}
	if s != "TYPE" {
		return errors.New("BlenderDNA: Expected TYPE field")
	}
	typesLen, err := d.GetInt32()
	if err != nil {
		return err
	}
	var types []*Type
	for i := int32(0); i < typesLen; i++ {
		types = append(types, &Type{})
	}
	for _, v := range types {
		var tmp []byte
		for {
			c, err := d.GetUInt8()
			if err != nil {
				return err
			}
			if c == 0 {
				break
			}
			tmp = append(tmp, c)
		}
		v.name = string(tmp)
	}
	if err = skip(); err != nil {
		return err
	}

	s, err = d.GetString(4)
	if err != nil {
		return err
	}
	if s != "TLEN" {
		return errors.New("BlenderDNA: Expected TYPE field")
	}
	for _, v := range types {
		tmp, err := d.GetInt16()
		if err != nil {
			return err
		}
		v.size = int32(tmp)
	}
	if err = skip(); err != nil {
		return err
	}
	s, err = d.GetString(4)
	if err != nil {
		return err
	}
	if s != "STRC" {
		return errors.New("BlenderDNA: Expected TYPE field")
	}
	fields := 0
	end, err := d.GetInt32()
	if err != nil {
		return err
	}
	for i := int32(0); i < end; i++ {
		st := NewStructure()

		n, err := d.GetInt16()
		if err != nil {
			return err
		}
		if int(n) >= len(types) {
			return fmt.Errorf("BlenderDNA: Invalid type index in structure name %v (there are only:%v  entries", n, len(types))
		}
		dna.indices[types[n].name] = len(dna.structures)
		dna.structures = append(dna.structures, st)
		st.name = types[n].name

		n, err = d.GetInt16()
		if err != nil {
			return err
		}
		offset := int32(0)
		for m := int16(0); m < n; m++ {

			j, err := d.GetInt16()
			if err != nil {
				return err
			}
			if int(j) >= len(types) {
				return fmt.Errorf("BlenderDNA: Invalid type index in structure field %v (there are only %v entries)", j, len(types))
			}
			var f Field
			st.fields = append(st.fields, &f)
			f.offset = offset
			f.Type = types[j].name
			f.size = types[j].size

			j, err = d.GetInt16()
			if err != nil {
				return err
			}
			if int(j) >= len(names) {
				return fmt.Errorf("BlenderDNA: Invalid name index in structure field %v  (there are only %v entries", j, len(names))
			}

			f.name = names[j]
			f.flags = 0

			// pointers always specify the size of the pointee instead of their own.
			// The pointer asterisk remains a property of the lookup name.
			if f.name[0] == '*' {
				f.size = 4
				if d.db.i64bit {
					f.size = 8
				}
				f.flags |= FieldFlag_Pointer
			}

			// arrays, however, specify the size of a single element so we
			// need to parse the (possibly multi-dimensional) array declaration
			// in order to obtain the actual size of the array in the file.
			// Also we need to alter the lookup name to include no array
			// brackets anymore or size fixup won't work (if our size does
			// not match the size read from the DNA).
			if strings.HasSuffix(f.name, "]") {
				index := strings.Index(f.name, "[")
				if index == -1 {
					return fmt.Errorf("BlenderDNA: Encountered invalid array declaration%v ", f.name)
				}

				f.flags |= FieldFlag_Array
				f.array_sizes = ExtractArraySize(f.name)
				f.name = f.name[:index]

				f.size *= f.array_sizes[0] * f.array_sizes[1]
			}
			// maintain separate indexes
			st.indices[f.name] = int32(len(st.fields)) - 1
			offset += f.size
			fields++
		}
		st.size = offset

	}
	logger.DebugF("BlenderDNA: Got %v structures with totally %v fields", len(dna.structures), fields)
	dna.AddPrimitiveStructures()
	dna.RegisterConverters()
	return nil
}

// -------------------------------------------------------------------------------
/** Represents a generic offset within a BLEND file */
// -------------------------------------------------------------------------------
type FileOffset struct {
	val int
}

// ------------------------------------------------------------------------------------------------
func ExtractArraySize(s string) (array_sizes [2]int32) {
	array_sizes[1] = 1
	array_sizes[0] = array_sizes[1]
	index := strings.Index(s, "[")
	if index == -1 {
		return
	}
	array_sizes[0] = common.Strtoul10(s[index+1:])
	s = s[index+1:]
	index = strings.Index(s, "[")
	if index == -1 {
		return
	}
	array_sizes[1] = common.Strtoul10(s[index+1:])
	return array_sizes
}

// --------------------------------------------------------
/** Add structure definitions for all the primitive types,
 *  i.e. integer, short, char, float */
func (d *DNA) AddPrimitiveStructures() {
	// NOTE: these are just dummies. Their presence enforces
	// Structure::Convert<target_type> to be called on these
	// empty structures. These converters are special
	// overloads which scan the name of the structure and
	// perform the required data type conversion if one
	// of these special names is found in the structure
	// in question.

	d.indices["int"] = len(d.structures)
	s := NewStructure()
	d.structures = append(d.structures, s)
	s.name = "int"
	s.size = 4

	d.indices["short"] = len(d.structures)
	s = NewStructure()
	d.structures = append(d.structures, s)
	s.name = "short"
	s.size = 2

	d.indices["char"] = len(d.structures)
	s = NewStructure()
	d.structures = append(d.structures, s)
	s.name = "char"
	s.size = 1

	d.indices["float"] = len(d.structures)
	s = NewStructure()
	d.structures = append(d.structures, s)
	s.name = "float"
	s.size = 4

	d.indices["double"] = len(d.structures)
	s = NewStructure()
	d.structures = append(d.structures, s)
	s.name = "double"
	s.size = 8

	// no long, seemingly.
}

// --------------------------------------------------------------------------------
func (d *DNA) Get(ss string) *Structure {
	return d.structures[d.indices[ss]]
}

// --------------------------------------------------------------------------------
func (d *DNA) IndexByString(ss string) *Structure {
	it, ok := d.indices[ss]
	if !ok {
		logger.ErrorF("BlendDNA: Did not find a structure named `%v`", ss)
		return nil
	}
	return d.structures[it]
}

// --------------------------------------------------------------------------------
func (d *DNA) Index(i int) *Structure {
	if i >= len(d.structures) {
		logger.ErrorF("BlendDNA: There is no structure with index `%v`", i)
	}

	return d.structures[i]
}

// ------------------------------------------------------------------------------------------------
func (d *DNA) GetBlobToStructureConverter(structure *Structure, db *FileDatabase) DNAConverterFactory {
	return d.converters[structure.name]
}

// --------------------------------------------------------------------------------
func (s *Structure) Get(ss string) *Field {
	v, ok := s.indices[ss]
	if !ok {
		return nil
	}
	return s.fields[v]
}

// --------------------------------------------------------------------------------
func (s *Structure) Index(i int) *Field {
	if i >= len(s.fields) {
		logger.ErrorF("BlendDNA: There is no field with index %v` ` in structure `%v`", i, s.name)
		return nil
	}

	return s.fields[i]
}

func (s *Structure) IndexByString(ss string) *Field {
	it, ok := s.indices[ss]
	if !ok {
		logger.WarnF("BlendDNA: Did not find a field named `%v ` in structure `%v", ss, s.name)
		return nil
	}

	return s.fields[it]
}

func (s *Structure) ReadField(out any, name string, db *FileDatabase) error {
	f := s.IndexByString(name)
	if f == nil {
		logger.WarnF("Structure :%v not found filed name:%v", s.name, name)
		return nil
	}
	oldPos := db.GetCurPos()
	defer db.SetCurPos(oldPos)
	err := db.Discard(int(f.offset))
	if err != nil {
		return err
	}
	// find the structure definition pertaining to this field
	ss := db.dna.IndexByString(f.Type)
	err = ss.Convert(out, db)
	if err != nil {
		return err
	}

	// and recover the previous stream position
	db.stats().fields_read++
	return nil
}

var emptyElemSlice = []IElemBase{emptyElem}
var emptyElem = (IElemBase)(&ElemBase{})
var emptyFileSet = &FileOffset{}

// --------------------------------------------------------------------------------
func (s *Structure) ReadFieldPtrPtr(name string, db *FileDatabase, non_recursives ...bool) (out any, err error) {
	return s.readFieldPtr(emptyElemSlice, name, db, non_recursives...)
}
func (s *Structure) ReadFieldPtr(name string, db *FileDatabase, non_recursives ...bool) (out any, err error) {
	return s.readFieldPtr(emptyElem, name, db, non_recursives...)
}
func (s *Structure) ReadFieldFileOffsetPtr(name string, db *FileDatabase, non_recursives ...bool) (out any, err error) {
	return s.readFieldPtr(emptyFileSet, name, db, non_recursives...)
}
func (s *Structure) readFieldPtr(in any, name string, db *FileDatabase, non_recursives ...bool) (out any, err error) {
	non_recursive := false
	if len(non_recursives) > 0 {
		non_recursive = non_recursives[0]
	}
	old := db.GetCurPos()
	if !non_recursive {
		// and recover the previous stream position
		defer db.SetCurPos(old)
	}

	var ptrval Pointer

	f := s.IndexByString(name)
	if f == nil {
		logger.WarnF("readFieldPtr Structure:%v not found name :%v", s.name, name)
		return
	}
	// sanity check, should never happen if the genblenddna script is right
	if (f.flags & FieldFlag_Pointer) == 0 {
		return out, fmt.Errorf("field `%v` of structure `%v` ought to be a pointer", name, s.name)
	}

	err = db.Discard(int(f.offset))
	if err != nil {
		return out, err
	}
	err = s.Convert(&ptrval, db)
	if err != nil {
		return out, err
	}
	// actually it is meaningless on which Structure the Convert is called
	// because the `Pointer` argument triggers a special implementation.

	// resolve the pointer and load the corresponding structure
	out, err = s.ResolvePointer(in, &ptrval, db, f, non_recursives...)
	if err != nil {
		return nil, err
	}

	db.stats().fields_read++
	return out, nil
}

func (s *Structure) ReadFieldPtrArray(count int, name string, db *FileDatabase) (out []any, err error) {
	return s.readFieldPtrArray(emptyElem, count, name, db)
}

func (s *Structure) ReadFieldFileOffsetPtrArray(count int, name string, db *FileDatabase) (out []any, err error) {
	return s.readFieldPtrArray((*FileOffset)(nil), count, name, db)
}

//--------------------------------------------------------------------------------

func (s *Structure) ReadFieldPtrSlice(name string, db *FileDatabase) (out []any, err error) {
	old := db.GetCurPos()
	defer db.SetCurPos(old)

	var ptrval Pointer
	f := s.IndexByString(name)
	if f == nil {
		return out, nil
	}
	err = db.Discard(int(f.offset))
	if err != nil {
		return out, err
	}
	// sanity check, should never happen if the genblenddna script is right
	if (f.flags & FieldFlag_Pointer) == 0 {
		return out, fmt.Errorf("Field `%v` of structure `%v` ought to be a pointer", name, s.name)
	}

	err = s.Convert(&ptrval, db)
	// actually it is meaningless on which Structure the Convert is called
	// because the `Pointer` argument triggers a special implementation.

	if ptrval.val != 0 {
		// find the file block the pointer is pointing to
		block, err := s.LocateFileBlockForAddress(&ptrval, db)
		if err != nil {
			return out, err
		}
		db.SetCurPos(block.start + int(ptrval.val-block.address.val))
		// FIXME: basically, this could cause problems with 64 bit pointers on 32 bit systems.
		// I really ought to improve StreamReader to work with 64 bit indices exclusively.

		ss := db.dna.IndexByString(f.Type)
		for i := int32(0); i < block.num; i++ {
			// allocate the object hull
			// continue conversion after allocating the required storage
			fa := db.dna.GetBlobToStructureConverter(ss, db)
			if fa == nil {
				// this might happen if DNA::RegisterConverters hasn't been called so far
				// or if the target type is not contained in `our` DNA.
				logger.WarnF("Failed to find a converter for the `%v` structure", s.name)
				return out, nil
			}
			oc := fa()
			err = oc.Convert(db, ss)
			if err != nil {
				return out, err
			}
			// store a pointer to the name string of the actual type
			// in the object itself. This allows the conversion code
			// to perform additional type checking.
			oc.SetDnaType(s.name)
			out = append(out, oc)
		}
	}

	db.SetCurPos(old)
	db.stats().fields_read++
	return out, nil
}

func (s *Structure) readFieldPtrArray(in any, count int, name string, db *FileDatabase) (out1 []any, err error) {
	ptrval := make([]Pointer, count)
	f := s.IndexByString(name)
	// sanity check, should never happen if the genblenddna script is right
	if (FieldFlag_Pointer | FieldFlag_Pointer) != (f.flags & (FieldFlag_Pointer | FieldFlag_Pointer)) {
		return out1, fmt.Errorf("field ` %v` of structure ` %v ` ought to be a pointer AND an array", name, s.name)
	}
	oldPos := db.GetCurPos()
	defer db.SetCurPos(oldPos)
	err = db.Discard(int(f.offset))
	if err != nil {
		return out1, err
	}
	// find the structure definition pertaining to this field
	i := 0
	for ; i < int(math.Min(float64(f.array_sizes[0]), float64(len(ptrval)))); i++ {
		err := s.Convert(&ptrval[i], db)
		if err != nil {
			return out1, err
		}
	}
	for ; i < count; i++ {
		ptrval[i] = Pointer{}
	}
	for i = 0; i < count; i++ {
		// resolve the pointer and load the corresponding structure
		v, err := s.ResolvePointer(in, &ptrval[i], db, f)
		if err != nil {
			if errors.Is(err, ErrorPtrZero) {
				continue
			}
			return out1, err
		}
		out1 = append(out1, v)
	}

	// and recover the previous stream position
	db.stats().fields_read++
	return out1, nil
}

func (st *Structure) ResolvePointer(in any, ptrval *Pointer, db *FileDatabase, f *Field, non_recursiveds ...bool) (out any, err error) {
	switch v := in.(type) {
	case []IElemBase:
		return st.ResolvePointerSlice(v, ptrval, db, f)
	case IElemBase:
		return st.ResolvePointerObject(v, ptrval, db, f, non_recursiveds...)
	case *FileOffset:
		return v, st.ResolvePointerFileOffset(v, ptrval, db, f)
	default:
		panic("not impl")
	}
}

func (st *Structure) ResolvePointerSlice(in []IElemBase, ptrval *Pointer, db *FileDatabase, f *Field) (out any, err error) {
	// This is a function overload, not a template specialization. According to
	// the partial ordering rules, it should be selected by the compiler
	// for array-of-pointer inputs, i.e. Object::mats.

	if ptrval.val == 0 {
		return out, ErrorPtrZero
	}

	// find the file block the pointer is pointing to
	block, err := st.LocateFileBlockForAddress(ptrval, db)
	if err != nil {
		return nil, err
	}
	tmp := 4
	if db.i64bit {
		tmp = 8
	}
	num := int(block.size) / (tmp)

	// keep the old stream position
	pold := db.GetCurPos()
	db.SetCurPos(block.start + int(ptrval.val-block.address.val))

	var res []any
	// allocate raw storage for the array
	for i := 0; i < num; i++ {
		var val Pointer
		err = st.Convert(&val, db)
		if err != nil {
			return res, err
		}
		// and resolve the pointees
		out, err = st.ResolvePointer(in[i], &val, db, f)
		if err != nil {
			return res, err
		}
		res = append(res, out)
	}

	db.SetCurPos(pold)
	return res, nil
}
func (st *Structure) ResolvePointerFileOffset(out *FileOffset, ptrval *Pointer, db *FileDatabase,
	f *Field) error {
	// Currently used exclusively by PackedFile::data to represent
	// a simple offset into the mapped BLEND file.

	if ptrval.val == 0 {
		return nil
	}

	// find the file block the pointer is pointing to
	block, err := st.LocateFileBlockForAddress(ptrval, db)
	if err != nil {
		return err
	}
	out.val = block.start + int(ptrval.val-block.address.val)
	return nil
}

func (st *Structure) resolvePointer(out IElemBase, ptrval *Pointer, db *FileDatabase, f *Field, non_recursiveds ...bool) (ok bool, err error) {
	non_recursived := false
	if len(non_recursiveds) > 0 {
		non_recursived = non_recursiveds[0]
	}
	if ptrval.val == 0 {
		return false, nil
	}
	s := db.dna.IndexByString(f.Type)
	// find the file block the pointer is pointing to
	block, err := st.LocateFileBlockForAddress(ptrval, db)
	if err != nil {
		return false, nil
	}
	// also determine the target type from the block header
	// and check if it matches the type which we expect.
	ss := db.dna.Index(int(block.dna_index))
	if ss != s {
		return false, fmt.Errorf("expected target to be of type `%v ` but seemingly it is a `%v ` instead", s.name, ss.name)
	}

	// try to retrieve the object from the cache
	out = db.cache().get(s, ptrval)
	if out != nil {
		return true, nil
	}
	pold := db.GetCurPos()
	db.SetCurPos(int(ptrval.val - block.address.val))
	// seek to this location, but save the previous stream pointer.
	// FIXME: basically, this could cause problems with 64 bit pointers on 32 bit systems.
	// I really ought to improve StreamReader to work with 64 bit indices exclusively.

	// continue conversion after allocating the required storage
	//num := block.size / ss.size
	// cache the object before we convert it to avoid cyclic recursion.
	db.cache().set(s, ptrval, out)

	// if the non_recursive flag is set, we don't do anything but leave
	// the cursor at the correct position to resolve the object.
	if !non_recursived {
		err = s.Convert(out, db)
		if err != nil {
			return false, err
		}
		db.SetCurPos(pold)
	}
	if out != nil {
		db.stats().pointers_resolved++
	}

	return false, nil
}

var (
	ErrorPtrZero = errors.New("error pt zero")
)

// --------------------------------------------------------------------------------
func (st *Structure) ResolvePointerObject(in IElemBase,
	ptrval *Pointer,
	db *FileDatabase,
	f *Field, non_recursiveds ...bool) (out IElemBase, err error) {
	// Special case when the data type needs to be determined at runtime.
	// Less secure than in the `strongly-typed` case.
	if ptrval.val == 0 {
		return in, ErrorPtrZero
	}

	// find the file block the pointer is pointing to
	block, err := st.LocateFileBlockForAddress(ptrval, db)
	if err != nil {
		return in, err
	}
	// determine the target type from the block header
	s := db.dna.Index(int(block.dna_index))

	// try to retrieve the object from the cache
	out = db.cache().get(s, ptrval)
	if out != nil {
		return out, nil
	}
	pold := db.GetCurPos()
	// seek to this location, but save the previous stream pointer.
	db.SetCurPos(block.start + int(ptrval.val-block.address.val))
	// FIXME: basically, this could cause problems with 64 bit pointers on 32 bit systems.
	// I really ought to improve StreamReader to work with 64 bit indices exclusively.

	// continue conversion after allocating the required storage
	fa := db.dna.GetBlobToStructureConverter(s, db)
	if fa == nil {
		// this might happen if DNA::RegisterConverters hasn't been called so far
		// or if the target type is not contained in `our` DNA.
		logger.WarnF("Failed to find a converter for the `%v` structure", s.name)
		return in, nil
	}

	// allocate the object hull
	oc := fa()

	// cache the object immediately to prevent infinite recursion in a
	// circular list with a single element (i.e. a self-referencing element).
	db.cache().set(s, ptrval, oc)
	if len(non_recursiveds) == 0 || len(non_recursiveds) > 0 && !non_recursiveds[0] {
		// and do the actual conversion
		err = oc.Convert(db, s)
		if err != nil {
			return oc, err
		}
		// store a pointer to the name string of the actual type
		// in the object itself. This allows the conversion code
		// to perform additional type checking.

		oc.SetDnaType(s.name)
		db.SetCurPos(pold)
	} else {

	}
	db.stats().pointers_resolved++
	return oc, err
}

func (st *Structure) ReadFieldArray(out []any, name string, db *FileDatabase) error {
	f := st.IndexByString(name)
	if f == nil {
		logger.WarnF("Structure:%v ReadFieldArray not found name:%v", st.name, name)
		return nil
	}
	s := db.dna.IndexByString(f.Type)
	oldPos := db.GetCurPos()
	defer db.SetCurPos(oldPos)
	err := db.Discard(int(f.offset))
	if err != nil {
		return err
	}
	// is the input actually an array?
	if f.flags&FieldFlag_Array == 0 {
		return fmt.Errorf("field `%v ` of structure `%v ` ought to be an array of size %v", name, st.name, len(out))
	}
	// find the structure definition pertaining to this field
	i := 0
	for ; i < int(math.Min(float64(f.array_sizes[0]), float64(len(out)))); i++ {
		err := s.Convert(out[i], db)
		if err != nil {
			return err
		}
	}
	// and recover the previous stream position
	db.stats().fields_read++
	return nil
}

func (s *Structure) ReadFieldArray2(out [][]any, name string, db *FileDatabase) error {
	M, N := len(out), len(out[0])
	f := s.IndexByString(name)
	ss := db.dna.IndexByString(f.Type)
	oldPos := db.GetCurPos()
	defer db.SetCurPos(oldPos)
	err := db.Discard(int(f.offset))
	if err != nil {
		return err
	}
	// is the input actually an array?
	if f.flags&FieldFlag_Array == 0 {
		return fmt.Errorf("field `%v ` of structure `%v ` ought to be an array of size %v*%v", name, s.name, M, N)
	}
	// size conversions are always allowed, regardless of error_policy
	i := 0.0
	for ; i < math.Min(float64(f.array_sizes[0]), float64(M)); i++ {
		j := 0.0
		for ; j < math.Min(float64(f.array_sizes[1]), float64(N)); j++ {
			err := ss.Convert(out[int(i)][int(j)], db)
			if err != nil {
				return err
			}
		}
		for ; j < float64(N); j++ {
			out[int(i)][int(j)] = nil
		}
	}
	for ; i < float64(M); i++ {
		out[int(i)] = nil
	}

	// and recover the previous stream position
	db.stats().fields_read++

	return nil
}

func (s *Structure) ReadCustomDataPtr(cdtype int, name string, db *FileDatabase) (out []IElemBase, err error) {
	ptrval := &Pointer{}
	f := s.IndexByString(name)
	oldPos := db.GetCurPos()
	defer db.SetCurPos(oldPos)
	err = db.Discard(int(f.offset))
	if err != nil {
		return out, err
	}
	// sanity check, should never happen if the genblenddna script is right
	if (f.flags & FieldFlag_Pointer) == 0 {
		return out, fmt.Errorf("field `%v ` of structure `%v ` ought to be a pointer", name, s.name)
	}
	err = s.ConvertPointer(ptrval, db)
	if err != nil {
		return out, err
	}
	if ptrval.val != 0 {
		// get block for ptr
		block, err := s.LocateFileBlockForAddress(ptrval, db)
		if err != nil {
			return out, err
		}
		db.SetCurPos(block.start + int(ptrval.val-block.address.val))
		// read block.num instances of given type to out
		out, err = readCustomData(cdtype, int(block.num), db)
		if err != nil {
			return out, err
		}
	}

	// and recover the previous stream position
	db.stats().fields_read++
	return out, nil
}

// --------------------------------------------------------
/** Try to read an instance of the structure from the stream
 *  and attempt to convert to `T`. This is done by
 *  an appropriate specialization. If none is available,
 *  a compiler complain is the result.
 *  @param dest Destination value to be written
 *  @param db File database, including input stream. */

func (s *Structure) Convert(out any, db *FileDatabase) error {
	switch v := out.(type) {
	case **CustomDataLayer:
		var c CustomDataLayer
		err := c.Convert(db, s)
		if err != nil {
			return err
		}
		*v = &c
		return nil
	case *Pointer:
		return s.ConvertPointer(v, db)
	case *int64, *uint64:
	case *int32, *uint, *uint32:
		return s.ConvertInt(v, db)
	case *uint16, *int16:
		return s.ConvertInt16(v, db)
	case *int8:
		return s.ConvertInt8(v, db)
	case *uint8:
		return s.ConvertUInt8(v, db)
	case *float32:
		return s.ConvertFloat32(v, db)
	case *float64:
		return s.ConvertFloat64(v, db)
	default:
		fa := db.dna.GetBlobToStructureConverter(s, db)
		if fa != nil {
			oc := fa()
			err := oc.Convert(db, s)
			if err != nil {
				return err
			}
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(oc).Elem())
			return nil
		}
	}
	return errors.New("not impl")
}

// --------------------------------------------------------------------------------
func (s *Structure) LocateFileBlockForAddress(ptrval *Pointer, db *FileDatabase) (*FileBlockHead, error) {
	// the file blocks appear in list sorted by
	// with ascending base addresses so we can run a
	// binary search to locate the pointer quickly.

	// NOTE: Blender seems to distinguish between side-by-side
	// data (stored in the same data block) and far pointers,
	// which are only used for structures starting with an ID.
	// We don't need to make this distinction, our algorithm
	// works regardless where the data is stored.
	var v *FileBlockHead
	index := common.LowerBound(0, len(db.entries), func(index int) bool {
		return db.entries[index].address.val < ptrval.val
	})

	if index >= len(db.entries) {
		// this is crucial, pointers may not be invalid.
		// this is either a corrupted file or an attempted attack.
		return v, fmt.Errorf("failure resolving pointer 0x %v , no file block falls into this address range", ptrval.val)
	}
	v = db.entries[index]
	if ptrval.val >= v.address.val+uint64(v.size) {
		return v, fmt.Errorf("failure resolving pointer 0x %v ,nearest file block starting at 0x %v ends at 0x:%v", ptrval.val,
			v.address.val,
			v.address.val+uint64(v.size))
	}
	return v, nil
}

func (s *Structure) ConvertFloat64(dest any, db *FileDatabase) (err error) {
	if s.name == "char" {
		v, err := db.GetInt8()
		if err != nil {
			return err
		}
		return ConvertValue(dest, float64(v)/255.)
	} else if s.name == "short" {
		v, err := db.GetInt16()
		if err != nil {
			return err
		}
		return ConvertValue(dest, float64(v)/32767.)
	}
	return s.ConvertDispatcher(dest, db)
}

func (s *Structure) ConvertInt(dest any, db *FileDatabase) error {
	return s.ConvertDispatcher(dest, db)
}

func (s *Structure) ConvertInt16(dest any, db *FileDatabase) error {
	// automatic rescaling from short to float and vice versa (seems to be used by normals)
	if s.name == "float" {
		f, err := db.GetFloat32()
		if err != nil {
			return err
		}
		if f > 1.0 {
			f = 1.0
		}
		//db.reader.IncPtr(-4);
		return ConvertValue(dest, f*32767.)
	} else if s.name == "double" {
		f, err := db.GetFloat64()
		if err != nil {
			return err
		}
		//db.reader.IncPtr(-8);
		return ConvertValue(dest, f*32767.)
	}
	return s.ConvertDispatcher(dest, db)
}

func (s *Structure) ConvertInt8(dest any, db *FileDatabase) error {
	// automatic rescaling from char to float and vice versa (seems useful for RGB colors)
	if s.name == "float" {
		f, err := db.GetFloat32()
		if err != nil {
			return err
		}
		return ConvertValue(dest, f*255.)
	} else if s.name == "double" {
		f, err := db.GetFloat64()
		if err != nil {
			return err
		}
		return ConvertValue(dest, f*255.)
	}
	return s.ConvertDispatcher(dest, db)
}

func (s *Structure) ConvertUInt8(dest any, db *FileDatabase) error {
	// automatic rescaling from char to float and vice versa (seems useful for RGB colors)
	if s.name == "float" {
		f, err := db.GetFloat32()
		if err != nil {
			return err
		}
		return ConvertValue(dest, f*255.)
	} else if s.name == "double" {
		f, err := db.GetFloat64()
		if err != nil {
			return err
		}
		return ConvertValue(dest, f*255.)
	}
	return s.ConvertDispatcher(dest, db)
}

func (s *Structure) ConvertFloat32(dest any, db *FileDatabase) error {
	// automatic rescaling from char to float and vice versa (seems useful for RGB colors)
	if s.name == "char" {
		i, err := db.GetInt8()
		if err != nil {
			return err
		}
		return ConvertValue(dest, float32(i)/255.)
		// automatic rescaling from short to float and vice versa (used by normals)
	} else if s.name == "short" {
		i, err := db.GetInt16()
		if err != nil {
			return err
		}
		return ConvertValue(dest, float32(i)/32767.)
	}
	return s.ConvertDispatcher(dest, db)
}

// ------------------------------------------------------------------------------------------------
func (s *Structure) ConvertPointer(dest *Pointer, db *FileDatabase) (err error) {
	if db.i64bit {
		dest.val, err = db.GetUInt64()
		//db.reader.IncPtr(-8);
		return err
	}
	v, err := db.GetUInt32()
	dest.val = uint64(v)
	//db.reader.IncPtr(-4);
	return err
}

func ConvertValue[T common.Number](dest any, out T) error {
	err := fmt.Errorf("invalid type %v", reflect.TypeOf(dest).Name())
	switch v := dest.(type) {
	case *float64:
		*v = float64(out)
	case *float32:
		*v = float32(out)
	case *int8:
		*v = int8(out)
	case *int16:
		*v = int16(out)
	case *int32:
		*v = int32(out)
	case *int64:
		*v = int64(out)
	case *uint8:
		*v = uint8(out)
	case *uint16:
		*v = uint16(out)
	case *uint32:
		*v = uint32(out)
	case *uint64:
		*v = uint64(out)
	default:
		return err
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
func (s *Structure) ConvertDispatcher(out any, db *FileDatabase) error {
	if s.name == "int" {
		v, err := db.GetUInt32()
		if err != nil {
			return err
		}
		return ConvertValue(out, v)
	} else if s.name == "short" {
		v, err := db.GetUInt16()
		if err != nil {
			return err
		}
		return ConvertValue(out, v)
	} else if s.name == "char" {
		v, err := db.GetUInt8()
		if err != nil {
			return err
		}
		return ConvertValue(out, v)
	} else if s.name == "float" {
		v, err := db.GetFloat32()
		if err != nil {
			return err
		}
		return ConvertValue(out, v)
	} else if s.name == "double" {
		v, err := db.GetFloat64()
		if err != nil {
			return err
		}
		return ConvertValue(out, v)
	}
	return fmt.Errorf("unknown source for conversion to primitive data type: %v", s.name)
}
