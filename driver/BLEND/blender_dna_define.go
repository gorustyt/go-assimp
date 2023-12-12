package BLEND

import (
	"assimp/common/reader"
	"errors"
)

type DNAConvert interface {
	Convert(db *FileDatabase) error
}

type DNAConverterFactory func() DNAConvert
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
	//reader  reader
	entries []*FileBlockHead

	_stats         *Statistics
	_cacheArrays   []*ObjectCache
	_cache         *ObjectCache
	next_cache_idx int
}

func NewFileDatabase() *FileDatabase {
	f := &FileDatabase{}

	return f
}

func (db *FileDatabase) stats() *Statistics {
	return db._stats
}

type ObjectCache struct {
}

func NewObjectCache() {

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
	reader.StreamReader
	name    string
	fields  []*Field
	indices map[string]int32
	size    int32
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
	if s.Remain() < s.current.size {
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
