package BLEND

import (
	"assimp/common"
	"assimp/common/logger"
	"assimp/common/reader"
	"errors"
	"fmt"
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
		for d.GetReadNum()&0x3 != 0 {
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
		st := &Structure{}
		n, err := d.GetInt16()
		if err != nil {
			return err
		}
		if int(n) >= len(types) {
			return fmt.Errorf("BlenderDNA: Invalid type index in structure name %v (there are only:%v  entries", n, len(types))
		}
		dna.indices[types[n].name] = len(dna.structures)
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
				return fmt.Errorf("BlenderDNA: Invalid name index in structure field ", j, " (there are only ", len(names), " entries)")
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
			if strings.HasPrefix(f.name, "]") {
				index := strings.Index(f.name, "[")
				if index == -1 {
					return fmt.Errorf("BlenderDNA: Encountered invalid array declaration ", f.name)
				}

				f.flags |= FieldFlag_Array
				f.array_sizes = ExtractArraySize(f.name)
				f.name = f.name[:index]

				f.size *= f.array_sizes[0] * f.array_sizes[1]
			}

			// maintain separate indexes
			st.indices[f.name] = int32(len(st.fields)) - 1
			offset += f.size
		}
		dna.structures = append(dna.structures)
		st.size = offset
	}
	logger.DebugF("BlenderDNA: Got %v structures with totally %v fields", len(dna.structures), fields)
	dna.AddPrimitiveStructures()
	dna.RegisterConverters()
	return nil
}

/**
*   @brief  read CustomData's data to ptr to mem
*   @param[out] out memory ptr to set
*   @param[in]  cdtype  to read
*   @param[in]  cnt cnt of elements to read
*   @param[in]  db to read elements from
*   @return true when ok
 */
func (d *DNAParser) ReadCustomData(out *ElemBase, cdtype int, cnt int, db *FileDatabase) bool {
	return true
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
	if index+1 == len(s) {
		return
	}
	array_sizes[0] = common.Strtoul10(s[index+1:])
	index = strings.Index(s, "[")
	if index+1 == len(s) {
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

// --------------------------------------------------------
/** Fill the @c converters member with converters for all
 *  known data types. The implementation of this method is
 *  in BlenderScene.cpp and is machine-generated.
 *  Converters are used to quickly handle objects whose
 *  exact data type is a runtime-property and not yet
 *  known at compile time (consider Object::data).*/
func (d *DNA) RegisterConverters() {
	d.converters["Object"] = func() any {
		return &Object{}
	}
	d.converters["Group"] = func() any {
		return &Group{}
	}
	d.converters["MTex"] = func() any {
		return &MTex{}
	}
	d.converters["TFace"] = func() any {
		return &TFace{}
	}
	d.converters["SubsurfModifierData"] = func() any {
		return &SubsurfModifierData{}
	}
	d.converters["MFace"] = func() any {
		return &MFace{}
	}
	d.converters["Lamp"] = func() any {
		return &Lamp{}
	}
	d.converters["MDeformWeight"] = func() any {
		return &MDeformWeight{}
	}
	d.converters["PackedFile"] = func() any {
		return &PackedFile{}
	}
	d.converters["Base"] = func() any {
		return &Base{}
	}
	d.converters["MTFace"] = func() any {
		return &MTFace{}
	}
	d.converters["Material"] = func() any {
		return &Material{}
	}
	d.converters["MTexPoly"] = func() any {
		return &MTexPoly{}
	}
	d.converters["Mesh"] = func() any {
		return &Mesh{}
	}
	d.converters["MDeformVert"] = func() any {
		return &MDeformVert{}
	}
	d.converters["World"] = func() any {
		return &World{}
	}
	d.converters["MLoopCol"] = func() any {
		return &MLoopCol{}
	}
	d.converters["MVert"] = func() any {
		return &MVert{}
	}
	d.converters["MEdge"] = func() any {
		return &MEdge{}
	}
	d.converters["MLoopUV"] = func() any {
		return &MLoopUV{}
	}
	d.converters["GroupObject"] = func() any {
		return &GroupObject{}
	}
	d.converters["ListBase"] = func() any {
		return &ListBase{}
	}
	d.converters["MLoop"] = func() any {
		return &MLoop{}
	}
	d.converters["ModifierData"] = func() any {
		return &ModifierData{}
	}
	d.converters["ID"] = func() any {
		return &ID{}
	}
	d.converters["MCol"] = func() any {
		return &MCol{}
	}
	d.converters["MPoly"] = func() any {
		return &MPoly{}
	}
	d.converters["Scene"] = func() any {
		return &Scene{}
	}
	d.converters["Library"] = func() any {
		return &Library{}
	}
	d.converters["Tex"] = func() any {
		return &Tex{}
	}
	d.converters["Camera"] = func() any {
		return &Camera{}
	}
	d.converters["MirrorModifierData"] = func() any {
		return &MirrorModifierData{}
	}
	d.converters["Image"] = func() any {
		return &Image{}
	}
	d.converters["CustomData"] = func() any {
		return &CustomData{}
	}
	d.converters["CustomDataLayer"] = func() any {
		return &CustomDataLayer{}
	}
	d.converters["Collection"] = func() any {
		return &Collection{}
	}
	d.converters["CollectionChild"] = func() any {
		return &CollectionChild{}
	}
	d.converters["CollectionObject"] = func() any {
		return &CollectionObject{}
	}
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
		logger.ErrorF("BlendDNA: Did not find a field named `%v ` in structure `%v", ss, s.name)
		return nil
	}

	return s.fields[it]
}

func (s *Structure) ReadField(out any, name string, db *FileDatabase) {
	old := reader.GetReadNum()
	f := s.IndexByString(name)
	// find the structure definition pertaining to this field
	ss := db.dna.IndexByString(f.Type)
	s.Convert(out, db)
	// and recover the previous stream position
	db.stats().fields_read++
}

// --------------------------------------------------------
/** Try to read an instance of the structure from the stream
 *  and attempt to convert to `T`. This is done by
 *  an appropriate specialization. If none is available,
 *  a compiler complain is the result.
 *  @param dest Destination value to be written
 *  @param db File database, including input stream. */

func (s *Structure) Convert(out any, db *FileDatabase) error {
	switch out.(type) {
	case *float64:

	}
	return nil
}

func (s *Structure) ConvertFloat64(db *FileDatabase) (dest float64, err error) {
	if s.name == "char" {
		v, err := s.GetInt8()
		if err != nil {
			return 0, err
		}
		dest = float64(v) / 255.
		return
	} else if s.name == "short" {
		v, err := s.GetInt16()
		if err != nil {
			return 0, err
		}
		dest = float64(v) / 32767.
		return
	}
	return 0, s.ConvertDispatcher(dest, db)
}

// ------------------------------------------------------------------------------------------------
func (s *Structure) ConvertDispatcher(out any, db *FileDatabase) error {
	if s.name == "int" {
		out = s.GetU4()
	} else if s.name == "short" {
		out = s.GetU2()
	} else if s.name == "char" {
		out = s.GetU1()
	} else if s.name == "float" {
		out = s.GetF4()
	} else if s.name == "double" {
		out = s.GetF8()
	}
	return fmt.Errorf("Unknown source for conversion to primitive data type: %v", s.name)
}
