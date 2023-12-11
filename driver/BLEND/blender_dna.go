package BLEND

// --------------------------------------------------------
/** Locate the DNA in the file and parse it. The input
 *  stream is expected to poto the beginning of the DN1
 *  chunk at the time this method is called and is
 *  undefined afterwards.
 *  @throw DeadlyImportError if the DNA cannot be read.
 *  @note The position of the stream pointer is undefined
 *    afterwards.*/

func (d *DNAParser) Parse() {

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
