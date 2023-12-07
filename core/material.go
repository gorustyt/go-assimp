package core

type AiPropertyTypeInfo int

const (
	/** Array of single-precision (32 Bit) floats
	 *
	 *  It is possible to use aiGetMaterialInteger[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in floating-point format.
	 *  The material system performs the type conversion automatically.
	 */
	aiPTI_Float AiPropertyTypeInfo = 0x1

	/** Array of double-precision (64 Bit) floats
	 *
	 *  It is possible to use aiGetMaterialInteger[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in floating-point format.
	 *  The material system performs the type conversion automatically.
	 */
	aiPTI_Double AiPropertyTypeInfo = 0x2

	/** The material property is an aiString.
	 *
	 *  Arrays of strings aren't possible, aiGetMaterialString() (or the
	 *  C++-API aiMaterial::Get()) *must* be used to query a string property.
	 */
	aiPTI_String AiPropertyTypeInfo = 0x3

	/** Array of (32 Bit) integers
	 *
	 *  It is possible to use aiGetMaterialFloat[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in integer format.
	 *  The material system performs the type conversion automatically.
	 */
	aiPTI_Integer AiPropertyTypeInfo = 0x4

	/** Simple binary buffer, content undefined. Not convertible to anything.
	 */
	aiPTI_Buffer AiPropertyTypeInfo = 0x5
)

type AiMaterial struct {
	/** List of all material properties loaded. */
	mProperties []*AiMaterialProperty

	/** Number of properties in the data base */
	NumProperties int

	/** Storage allocated */
	NumAllocated int
}

type AiMaterialProperty struct {
	/** Specifies the name of the property (key)
	 *  Keys are generally case insensitive.
	 */
	Key string

	/** Textures: Specifies their exact usage semantic.
	 * For non-texture properties, this member is always 0
	 * (or, better-said, #aiTextureType_NONE).
	 */
	Semantic int

	/** Textures: Specifies the index of the texture.
	 *  For non-texture properties, this member is always 0.
	 */
	Index int

	/** Size of the buffer mData is pointing to, in bytes.
	 *  This value may not be 0.
	 */
	DataLength int

	/** Type information for the property.
	 *
	 * Defines the data layout inside the data buffer. This is used
	 * by the library internally to perform debug checks and to
	 * utilize proper type conversions.
	 * (It's probably a hacky solution, but it works.)
	 */
	Type AiPropertyTypeInfo

	/** Binary buffer to hold the property's value.
	 * The size of the buffer is always mDataLength.
	 */
	Data []byte
}
