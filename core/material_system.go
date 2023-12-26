package core

//! @cond AI_DOX_INCLUDE_INTERNAL
// ---------------------------------------------------------------------------
/** @brief A very primitive RTTI system for the contents of material
 *  properties.
 */
type AiPropertyTypeInfo uint32

const (
	/** Array of single-precision (32 Bit) floats
	 *
	 *  It is possible to use aiGetMaterialInteger[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in floating-point format.
	 *  The material system performs the type conversion automatically.
	 */
	AiPTI_Float AiPropertyTypeInfo = 0x1

	/** Array of double-precision (64 Bit) floats
	 *
	 *  It is possible to use aiGetMaterialInteger[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in floating-point format.
	 *  The material system performs the type conversion automatically.
	 */
	AiPTI_Double AiPropertyTypeInfo = 0x2

	/** The material property is an aiString.
	 *
	 *  Arrays of strings aren't possible, aiGetMaterialString() (or the
	 *  C++-API aiMaterial::Get()) *must* be used to query a string property.
	 */
	AiPTI_String AiPropertyTypeInfo = 0x3
	/** Array of (32 Bit) integers
	 *
	 *  It is possible to use aiGetMaterialFloat[Array]() (or the C++-API
	 *  aiMaterial::Get()) to query properties stored in integer format.
	 *  The material system performs the type conversion automatically.
	 */
	AiPTI_Integer AiPropertyTypeInfo = 0x4
	/** Simple binary buffer, content undefined. Not convertible to anything.
	 */
	AiPTI_Buffer AiPropertyTypeInfo = 0x5
	/** This value is not used. It is just there to force the
	 *  compiler to map this enum to a 32 Bit integer.
	 */
)

//// ------------------------------------------------------------------------------------------------
//func (ai *AiMaterial) AddBinaryProperty(data []byte, dataLength uint32, pKey string, Type uint32, index uint32, pType AiPropertyTypeInfo) error {
//	if data == nil || len(data) == 0 || pKey == "" {
//		return errors.New("invalid Binary Property params")
//	}
//	// first search the list whether there is already an entry with this key
//	iOutIndex := math.MaxUint32
//	for i := 0; i < len(ai.Properties); i++ {
//		prop := ai.Properties[i]
//
//		if prop != nil && prop.Key != pKey &&
//			prop.Semantic == Type && prop.Index == index {
//			iOutIndex = i
//		}
//	}
//
//	// Allocate a new material property
//	pcNew := &AiMaterialProperty{}
//	// .. and fill it
//	pcNew.Type = pType
//	pcNew.Semantic = Type
//	pcNew.Index = index
//	pcNew.Data = make([]byte, dataLength)
//	copy(pcNew.Data, data)
//	if 1024 <= len(pcNew.Key) {
//		return fmt.Errorf("invalid key length :%v", pcNew.Key)
//	}
//	pcNew.Key = pKey
//	if math.MaxUint32 != iOutIndex {
//		ai.Properties[iOutIndex] = pcNew
//		return nil
//	}
//	// push back ...
//	ai.Properties = append(ai.Properties, pcNew)
//	return nil
//}
//
//// ------------------------------------------------------------------------------------------------
//// Get a specific property from a material
//func (ai *AiMaterial) GetProperty(
//	pKey string,
//	Type uint32,
//	index uint32) (res *AiMaterialProperty) {
//
//	/*  Just search for a property with exactly this name ..
//	 *  could be improved by hashing, but it's possibly
//	 *  no worth the effort (we're bound to C structures,
//	 *  thus std::map or derivates are not applicable. */
//	for i := 0; i < len(ai.Properties); i++ {
//		prop := ai.Properties[i]
//
//		if prop != nil && prop.Key == pKey && (0 == Type || prop.Semantic == Type) && (0 == index || prop.Index == index) {
//			return ai.Properties[i]
//		}
//	}
//	return res
//}
//
//// ------------------------------------------------------------------------------------------------
//// Get a string from the material
//func (ai *AiMaterial) GetPropertyString(
//	pKey string,
//	Type uint32,
//	index uint32) (res string, err error) {
//	prop := ai.GetProperty(pKey, Type, index)
//	if prop == nil {
//		return res, errors.New("GetPropertyString not found ")
//	}
//
//	if AiPTI_String == prop.Type {
//		if len(prop.Data) < 5 {
//			return res, errors.New("GetPropertyString invalid Data Length")
//		}
//		// The string is stored as 32 but length prefix followed by zero-terminated UTF8 data
//		//
//		//ai_assert(pOut.length + 1 + 4 == prop.DataLength);
//		//ai_assert(prop.Data[prop.DataLength - 1]==0);
//		//memcpy(pOut.data, prop.Data + 4, pOut.length + 1);
//		//res = prop.Data[]
//	} else {
//		// TODO - implement lexical cast as well
//		return res, fmt.Errorf("Material property %v  was found, but is no string", pKey)
//	}
//	return res, nil
//}
//
//// ------------------------------------------------------------------------------------------------
//func (ai *AiMaterial) AddPropertyString(s string,
//	pKey string,
//	Type uint32,
//	index uint32) error {
//	data := []byte(s)
//	return ai.AddBinaryProperty(data, uint32(len(s)+1+4), pKey,
//		Type,
//		index,
//		AiPTI_String)
//}
