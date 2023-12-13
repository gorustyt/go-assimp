package BLEND

import "fmt"

type CustomDataType int

const (
	CD_AUTO_FROM_NAME CustomDataType = -1
	CD_MVERT          CustomDataType = 0
	CD_MSTICKY        CustomDataType = 1 /* DEPRECATED */
	CD_MDEFORMVERT    CustomDataType = 2
	CD_MEDGE          CustomDataType = 3
	CD_MFACE          CustomDataType = 4
	CD_MTFACE         CustomDataType = 5
	CD_MCOL           CustomDataType = 6
	CD_ORIGINDEX      CustomDataType = 7
	CD_NORMAL         CustomDataType = 8
	/*	CD_POLYINDEX        = 9, */
	CD_PROP_FLT     CustomDataType = 10
	CD_PROP_INT     CustomDataType = 11
	CD_PROP_STR     CustomDataType = 12
	CD_ORIGSPACE    CustomDataType = 13 /* for modifier stack face location mapping */
	CD_ORCO         CustomDataType = 14
	CD_MTEXPOLY     CustomDataType = 15
	CD_MLOOPUV      CustomDataType = 16
	CD_MLOOPCOL     CustomDataType = 17
	CD_TANGENT      CustomDataType = 18
	CD_MDISPS       CustomDataType = 19
	CD_PREVIEW_MCOL CustomDataType = 20 /* for displaying weightpaint colors */
	/*	CD_ID_MCOL          = 21, */
	CD_TEXTURE_MLOOPCOL CustomDataType = 22
	CD_CLOTH_ORCO       CustomDataType = 23
	CD_RECAST           CustomDataType = 24

	/* BMESH ONLY START */
	CD_MPOLY            CustomDataType = 25
	CD_MLOOP            CustomDataType = 26
	CD_SHAPE_KEYINDEX   CustomDataType = 27
	CD_SHAPEKEY         CustomDataType = 28
	CD_BWEIGHT          CustomDataType = 29
	CD_CREASE           CustomDataType = 30
	CD_ORIGSPACE_MLOOP  CustomDataType = 31
	CD_PREVIEW_MLOOPCOL CustomDataType = 32
	CD_BM_ELEM_PYPTR    CustomDataType = 33
	/* BMESH ONLY END */

	CD_PAINT_MASK       CustomDataType = 34
	CD_GRID_PAINT_MASK  CustomDataType = 35
	CD_MVERT_SKIN       CustomDataType = 36
	CD_FREESTYLE_EDGE   CustomDataType = 37
	CD_FREESTYLE_FACE   CustomDataType = 38
	CD_MLOOPTANGENT     CustomDataType = 39
	CD_TESSLOOPNORMAL   CustomDataType = 40
	CD_CUSTOMLOOPNORMAL CustomDataType = 41

	CD_NUMTYPES CustomDataType = 42
)

/**
 *   @brief  read/convert of Structure array to memory
 */

func read[T any](s *Structure, p []T, db *FileDatabase) error {
	for i := 0; i < len(p); i++ {
		var r T
		err := s.Convert(r, db)
		if err != nil {
			return err
		}
		p[i] = r
	}
	return nil
}

type PRead func(pOut IElemBase, cnt int, db *FileDatabase) bool
type PCreate func(cnt int) []IElemBase
type PDestroy func(IElemBase)

func defaultPRead(pOut []*IElemBase, cnt int, db *FileDatabase) bool {
	return read(db.dna[], prt, db)
}

func defaultPCreate(cnt int) []IElemBase {
	return make([]IElemBase, cnt)
}
func defaultPDestroy(IElemBase) {

}

/**
 *   @brief  helper macro to define Structure type specific CustomDataTypeDescription
 *   @note   IMPL_STRUCT_READ for same ty must be used earlier to implement the typespecific read function
 */

func DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(ty IElemBase) *CustomDataTypeDescription {
	return newCustomDataTypeDescription(defaultPRead, defaultPCreate, defaultPDestroy)
}

func DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION() *CustomDataTypeDescription {
	return newCustomDataTypeDescription(nil, nil, nil)
}

/**
 *   @brief  descriptors for data pointed to from CustomDataLayer.data
 *   @note   some of the CustomData uses already well defined Structures
 *           other (like CD_ORCO, ...) uses arrays of rawtypes or even arrays of Structures
 *           use a special readfunction for that cases
 */
var customDataTypeDescriptions = []*CustomDataTypeDescription{
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MVert{ElemBase: &ElemBase{}}),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MEdge{ElemBase: &ElemBase{}}),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MFace{ElemBase: &ElemBase{}}),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MTFace{ElemBase: &ElemBase{}}),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),

	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MTexPoly{ElemBase: &ElemBase{}}),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MLoopUV{ElemBase: &ElemBase{}}),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MLoopCol{ElemBase: &ElemBase{}}),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),

	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MPoly{ElemBase: &ElemBase{}}),
	DECL_STRUCT_CUSTOMDATATYPEDESCRIPTION(&MLoop{ElemBase: &ElemBase{}}),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),

	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),

	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION(),
	DECL_UNSUPPORTED_CUSTOMDATATYPEDESCRIPTION()}

/**
 *   @brief  describes the size of data and the read function to be used for single CustomerData.type
 */
type CustomDataTypeDescription struct {
	Read    PRead   ///< function to read one CustomData type element
	Create  PCreate ///< function to allocate n type elements
	Destroy PDestroy
}

func newCustomDataTypeDescription(read PRead, create PCreate, destroy PDestroy) *CustomDataTypeDescription {
	return &CustomDataTypeDescription{
		Read:    read,
		Create:  create,
		Destroy: destroy,
	}
}
func isValidCustomDataType(cdtype int) bool {
	return cdtype >= 0 && cdtype < int(CD_NUMTYPES)
}

func readCustomData(out IElemBase, cdtype int, cnt int, db *FileDatabase) (ok bool, err error) {
	if !isValidCustomDataType(cdtype) {
		return ok, fmt.Errorf("CustomData.type %v out of index", cdtype)
	}
	cdtd := customDataTypeDescriptions[cdtype]
	if cdtd.Read != nil && cdtd.Create != nil && cdtd.Destroy != nil && cnt > 0 {
		// allocate cnt elements and parse them from file
		out.reset(cdtd.Create(cnt), cdtd.Destroy)
		return cdtd.Read(out, cnt, db), nil
	}
	return false, nil
}

func getCustomDataLayer(customdata *CustomData, cdtype CustomDataType, name string) *CustomDataLayer {
	for _, v := range customdata.layers {
		if v.Type == cdtype && name == v.name {
			return v
		}
	}
	return nil
}

func getCustomDataLayerData(customdata *CustomData, cdtype CustomDataType, name string) IElemBase {
	pLayer := getCustomDataLayer(customdata, cdtype, name)
	if pLayer != nil && pLayer.data != nil {
		return pLayer.data
	}
	return nil
}
