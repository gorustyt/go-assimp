package BLEND

type DNA struct {
	structures []*Structure
	indices    map[string]int
	//converters map[string]FactoryPair
}

type FileDatabase struct {
	i64bit bool
	little bool

	dna     DNA
	reader  reader
	entries []FileBlockHead

	_stats         *Statistics
	_cacheArrays   []*ObjectCache
	_cache         *ObjectCache
	next_cache_idx int
}

func NewFileDatabase() *FileDatabase {
	f := &FileDatabase{}

	return f
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

type Structure struct {
}

// -------------------------------------------------------------------------------
/** Describes a master file block header. Each master file sections holds n
 *  elements of a certain SDNA structure (or otherwise unspecified data). */
// -------------------------------------------------------------------------------

type FileBlockHead struct {
	// points right after the header of the file block
	id   string
	size int
	// original memory address of the data
	// index into DNA
	dna_index int

	// number of structure instances to follow
	num int
}

// -------------------------------------------------------------------------------
/** Utility to read all master file blocks in turn. */
// -------------------------------------------------------------------------------

type SectionParser struct {
	current FileBlockHead
	//stream  StreamReaderAny
	ptr64 bool
}

// -------------------------------------------------------------------------------
/** Factory to extract a #DNA from the DNA1 file block in a BLEND file. */
// -------------------------------------------------------------------------------

type DNAParser struct {
	db *FileDatabase
}

func NewDNAParser(db *FileDatabase) *DNAParser {
	return &DNAParser{db: db}
}

/** Obtain a reference to the extracted DNA information */
func (d *DNAParser) GetDNA() DNA {
	return d.db.dna
}
