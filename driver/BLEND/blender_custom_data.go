package BLEND

import (
	"fmt"
	"reflect"
)

type CustomDataType int32

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

var customDataTypeDescriptions = [CD_NUMTYPES]func() Converter{
	func() Converter { return &MVert{ElemBase: &ElemBase{}} },
	nil,
	nil,
	func() Converter { return &MEdge{ElemBase: &ElemBase{}} },
	func() Converter { return &MFace{ElemBase: &ElemBase{}} },
	func() Converter { return &MTFace{ElemBase: &ElemBase{}} },
	nil,
	nil,
	nil,
	nil,

	nil,
	nil,
	nil,
	nil,
	nil,
	func() Converter { return &MTexPoly{ElemBase: &ElemBase{}} },
	func() Converter { return &MLoopUV{ElemBase: &ElemBase{}} },
	func() Converter { return &MLoopCol{ElemBase: &ElemBase{}} },
	nil,
	nil,

	nil,
	nil,
	nil,
	nil,
	nil,
	func() Converter { return &MPoly{ElemBase: &ElemBase{}} },
	func() Converter { return &MLoop{ElemBase: &ElemBase{}} },
	nil,
	nil,
	nil,

	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,

	nil,
	nil}

func isValidCustomDataType(cdtype int) bool {
	return cdtype >= 0 && cdtype < int(CD_NUMTYPES)
}

func readCustomData(cdtype int, cnt int, db *FileDatabase) (out []IElemBase, err error) {
	if !isValidCustomDataType(cdtype) {
		return out, fmt.Errorf("CustomData.type %v out of index", cdtype)
	}
	cdtd := customDataTypeDescriptions[cdtype]
	if cdtd != nil && cnt > 0 {
		// allocate cnt elements and parse them from file
		for i := 0; i < cnt; i++ {
			v := cdtd()
			out = append(out, v)
			ss := db.dna.IndexByString(reflect.TypeOf(v).Elem().Name())
			err = v.Convert(db, ss)
			if err != nil {
				return nil, err
			}
		}
		return out, nil
	}
	return out, nil
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
