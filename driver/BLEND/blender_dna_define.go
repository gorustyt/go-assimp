package BLEND

import (
	"assimp/common/reader"
	"errors"
)

type Converter interface {
	IElemBase
	Convert(db *FileDatabase, s *Structure) error
}
type DNAConverterFactory func() Converter
type DNA struct {
	structures []*Structure
	indices    map[string]int
	converters map[string]DNAConverterFactory
}

func NewDNA() *DNA {
	return &DNA{indices: map[string]int{}, converters: map[string]DNAConverterFactory{}}
}

type FileDatabase struct {
	i64bit bool
	little bool

	dna *DNA
	reader.StreamReader
	entries []*FileBlockHead

	_stats         *Statistics
	_cacheArrays   []*ObjectCache
	_cache         *ObjectCache
	next_cache_idx int32
}

func NewFileDatabase(reader reader.StreamReader) *FileDatabase {
	f := &FileDatabase{_stats: &Statistics{}, StreamReader: reader}

	return f
}

func (db *FileDatabase) stats() *Statistics {
	return db._stats
}

func (db *FileDatabase) cache() *ObjectCache {
	return db._cache
}

func (db *FileDatabase) cacheArray() []*ObjectCache {
	return db._cacheArrays
}

type ObjectCache struct {
	caches []map[*Pointer]IElemBase
	db     *FileDatabase
}

// --------------------------------------------------------------------------------
func (oc *ObjectCache) get(s *Structure, ptr *Pointer) IElemBase {
	if s.cache_idx == -1 {
		s.cache_idx = oc.db.next_cache_idx
		oc.db.next_cache_idx++
		oc.caches = make([]map[*Pointer]IElemBase, oc.db.next_cache_idx)
		return nil
	}
	it, ok := oc.caches[s.cache_idx][ptr]
	if !ok {
		oc.db.stats().cache_hits++
	}
	return it
	// otherwise, out remains untouched
}

// --------------------------------------------------------------------------------
func (oc *ObjectCache) set(s *Structure, ptr *Pointer, value IElemBase) {
	if s.cache_idx == -1 {
		s.cache_idx = oc.db.next_cache_idx
		oc.db.next_cache_idx++
		oc.caches = make([]map[*Pointer]IElemBase, oc.db.next_cache_idx)
	}
	if oc.caches[s.cache_idx] == nil {
		oc.caches[s.cache_idx] = map[*Pointer]IElemBase{}
	}
	oc.caches[s.cache_idx][ptr] = value
	oc.db.stats().cached_objects++

}
func NewObjectCache(db *FileDatabase) *ObjectCache {
	return &ObjectCache{db: db}
}

type Statistics struct {
	/** total number of fields we read */
	fields_read int
	/** total number of resolved pointers */
	pointers_resolved int
	/** number of pointers resolved from the cache */
	cache_hits int
	/** number of blocks (from  FileDatabase::entries)
	  we did actually read from. */
	// unsigned blocks_read;
	/** objects in FileData::cache */
	cached_objects int
}

type ElemBase struct {
	dna_type string
}

func (e *ElemBase) GetDnaType() string         { return e.dna_type }
func (e *ElemBase) SetDnaType(dna_type string) { e.dna_type = dna_type }

type IElemBase interface {
	GetDnaType() string
	SetDnaType(dna_type string)
}

// -------------------------------------------------------------------------------
/** Mixed flags for use in #Field */
// -------------------------------------------------------------------------------
const (
	FieldFlag_Pointer = 0x1
	FieldFlag_Array   = 0x2
)

type Field struct {
	name   string
	Type   string
	size   int32
	offset int32
	/** Size of each array dimension. For flat arrays,
	 *  the second dimension is set to 1. */
	array_sizes [2]int32
	/** Any of the #FieldFlags enumerated values */
	flags int32
}

type Structure struct {
	name    string
	fields  []*Field
	indices map[string]int32
	size    int32

	cache_idx int32
}

func NewStructure() *Structure {
	return &Structure{indices: map[string]int32{}}
}

// -------------------------------------------------------------------------------
/** Represents a generic pointer to a memory location, which can be either 32
 *  or 64 bits. These pointers are loaded from the BLEND file and finally
 *  fixed to point to the real, converted representation of the objects
 *  they used to point to.*/
// -------------------------------------------------------------------------------
type Pointer struct {
	val uint64
}

// -------------------------------------------------------------------------------
/** Describes a master file block header. Each master file sections holds n
 *  elements of a certain SDNA structure (or otherwise unspecified data). */
// -------------------------------------------------------------------------------

type FileBlockHead struct {
	// points right after the header of the file block
	id   string
	size int32
	// original memory address of the data
	start int
	// index into DNA
	dna_index int32
	// original memory address of the data
	address Pointer
	// number of structure instances to follow
	num int32
}

// -------------------------------------------------------------------------------
/** Utility to read all master file blocks in turn. */
// -------------------------------------------------------------------------------

type SectionParser struct {
	current *FileBlockHead
	reader.StreamReader
	ptr64 bool
}

func NewSectionParser(reader reader.StreamReader, ptr64 bool) *SectionParser {
	return &SectionParser{StreamReader: reader, ptr64: ptr64, current: &FileBlockHead{}}
}
func (s *SectionParser) GetCurrent() *FileBlockHead {
	return s.current
}

// Advance to the next section.
func (s *SectionParser) Next() error {
	tmp, err := s.GetNBytes(4)
	if err != nil {
		return err
	}
	if tmp[3] != 0 {
		s.current.id = string(tmp)
	} else if tmp[2] != 0 {
		s.current.id = string(tmp[:3])
	} else if tmp[1] != 0 {
		s.current.id = string(tmp[:2])
	} else {
		s.current.id = string(tmp[:1])
	}

	s.current.size, err = s.GetInt32()
	if err != nil {
		return err
	}
	if s.ptr64 {
		s.current.address.val, err = s.GetUInt64()
	} else {
		var tmp uint32
		tmp, err = s.GetUInt32()
		s.current.address.val = uint64(tmp)
	}
	if err != nil {
		return err
	}
	s.current.dna_index, err = s.GetInt32()
	if err != nil {
		return err
	}
	s.current.num, err = s.GetInt32()
	if err != nil {
		return err
	}
	s.current.start = s.GetCurPos()
	if s.Remain() < int(s.current.size) {
		return errors.New("BLEND: invalid size of file block")
	}
	return nil
}

// -------------------------------------------------------------------------------
/** Factory to extract a #DNA from the DNA1 file block in a BLEND file. */
// -------------------------------------------------------------------------------

type DNAParser struct {
	db *FileDatabase
	reader.StreamReader
}

func NewDNAParser(db *FileDatabase, reader reader.StreamReader) *DNAParser {
	return &DNAParser{db: db, StreamReader: reader}
}

/** Obtain a reference to the extracted DNA information */
func (d *DNAParser) GetDNA() *DNA {
	return d.db.dna
}

type Type struct {
	size int32
	name string
}
