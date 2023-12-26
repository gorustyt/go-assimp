package ASSETBIN

import (
	"assimp/common"
	"assimp/common/pb_msg"
	"assimp/core"
	"encoding/binary"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math"
)

func GetAiMaterialPropertyData(p *core.AiMaterialProperty, tType uint32, data []byte) (err error) {
	pType := core.AiPropertyTypeInfo(tType)
	var v proto.Message
	switch pType {
	case core.AiPTI_Float:
		p.Type = pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat64
		res, err := AiGetMaterialFloatArray(tType, data, nil)
		if err != nil {
			return err
		}
		v = &pb_msg.AiMaterialPropertyFloat64{
			Data: res,
		}
	case core.AiPTI_Double:
		p.Type = pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat64
		res, err := AiGetMaterialFloatArray(tType, data, nil)
		if err != nil {
			return err
		}
		v = &pb_msg.AiMaterialPropertyFloat64{
			Data: res,
		}
	case core.AiPTI_String:
		p.Type = pb_msg.AiMaterialPropertyType_AiPropertyTypeString
		res, err := AiGetPropertyString(tType, data)
		if err != nil {
			return err
		}
		v = &pb_msg.AiMaterialPropertyString{
			Data: []string{res},
		}
	case core.AiPTI_Integer:
		p.Type = pb_msg.AiMaterialPropertyType_AiPropertyTypeInt

		v = &pb_msg.AiMaterialPropertyInt64{
			Data: []int64{},
		}
	case core.AiPTI_Buffer:
		p.Type = pb_msg.AiMaterialPropertyType_AiPropertyTypeBytes
	}
	p.Data, err = proto.Marshal(v)
	return err
}

// ------------------------------------------------------------------------------------------------
// Get an array of floating-point values from the material.
func AiGetMaterialFloatArray(tType uint32, data []byte, pMax *int) (res []float64, err error) {
	pType := core.AiPropertyTypeInfo(tType)
	// data is given in floats, convert to ai_real
	iWrite := 0
	if core.AiPTI_Float == pType || core.AiPTI_Buffer == pType {
		iWrite = len(data) / 4
		if pMax != nil {
			iWrite = int(math.Min(float64(*pMax), float64(iWrite)))
		}

		for a := 0; a < iWrite; a++ {
			tmp := binary.LittleEndian.Uint32(data[a*4 : a*4+4])
			res = append(res, float64(math.Float32frombits(tmp)))
		}

		if pMax != nil {
			*pMax = iWrite
		}
	} else if core.AiPTI_Double == pType { // data is given in doubles, convert to float
		iWrite = len(data) / 8
		if pMax != nil {
			iWrite = int(math.Min(float64(*pMax), float64(iWrite)))

		}
		for a := 0; a < iWrite; a++ {
			tmp := binary.LittleEndian.Uint64(data[a*8 : a*8+8])
			res = append(res, math.Float64frombits(tmp))
		}
		if pMax != nil {
			*pMax = iWrite
		}
	} else if core.AiPTI_Integer == pType { // data is given in ints, convert to float
		iWrite = len(data) / 4
		if pMax != nil {
			iWrite = int(math.Min(float64(*pMax), float64(iWrite)))

		}
		for a := 0; a < iWrite; a++ {
			tmp := binary.LittleEndian.Uint32(data[a*4 : a*4+4])
			res = append(res, float64(tmp))
		}
		if pMax != nil {
			*pMax = iWrite
		}
	} else { // a string ... read floats separated by spaces
		if pMax != nil {
			iWrite = *pMax
		}
		// strings are zero-terminated with a 32 bit length prefix, so this is safe
		cur := 4
		if (len(data) < 5) || data[len(data)-1] != 0 {
			return res, errors.New("AiGetMaterialFloatArray  invalid length")
		}
		for a := 0; ; a++ {
			var v float64
			v, cur, err = common.FastAtorealMove(data[cur:])
			if err != nil {
				return nil, err
			}
			res = append(res, v)
			if a == iWrite-1 {
				break
			}
			if data[cur] == ' ' || data[cur] == '\t' {
				return res, fmt.Errorf("Material property  is a string; failed to parse a float array out of it.")
			}
		}

		if pMax != nil {
			*pMax = iWrite
		}
	}
	return res, nil
}

// ------------------------------------------------------------------------------------------------
// Get an array if integers from the material
func AiGetMaterialIntegerArray(tType uint32, data []byte, pMax *int) (res []int64, err error) {
	pType := core.AiPropertyTypeInfo(tType)
	// data is given in ints, simply copy it
	iWrite := 0
	if core.AiPTI_Integer == pType || core.AiPTI_Buffer == pType {
		iWrite = int(math.Max(float64(len(data)/4), float64(1)))
		if pMax != nil {
			iWrite = int(math.Min(float64(*pMax), float64(iWrite)))
		}
		if 1 == len(data) {
			// bool type, 1 byte
			res = append(res, int64(data[0]))
		} else {
			for a := 0; a < iWrite; a++ {
				tmp := binary.LittleEndian.Uint32(data[a*4 : a*4+4])
				res = append(res, int64(tmp))
			}
		}
		if pMax != nil {
			*pMax = iWrite
		}
	} else if core.AiPTI_Float == pType { // data is given in floats convert to int
		iWrite = len(data) / 4
		if pMax != nil {
			iWrite = int(math.Min(float64(*pMax), float64(iWrite)))

		}
		for a := 0; a < iWrite; a++ {
			tmp := binary.LittleEndian.Uint32(data[a*4 : a*4+4])
			res = append(res, int64(math.Float32frombits(tmp)))
		}
		if pMax != nil {
			*pMax = iWrite
		}
	} else { // it is a string ... no way to read something out of this
		if pMax != nil {
			iWrite = *pMax
		}
		// strings are zero-terminated with a 32 bit length prefix, so this is safe
		data = data[4:]
		if (len(data) < 5) || data[len(data)-1] != 0 {
			return res, errors.New("AiGetMaterialIntegerArray  invalid length")
		}
		for a := 0; ; a++ {
			tmp, index := common.StrTol10(string(data))
			res = append(res, int64(tmp))
			if a == iWrite-1 {
				break
			}
			if data[index] == ' ' || data[index] == '\t' {
				return res, errors.New("material property is a string; failed to parse an integer array out of it")
			}
		}

		if pMax != nil {
			*pMax = iWrite
		}
	}
	return res, nil
}

// ------------------------------------------------------------------------------------------------
// Get a aiUVTransform (5 floats) from the material
func AiGetMaterialUVTransform(pType uint32, data []byte) (v *core.AiUVTransform, err error) {
	iMax := 5
	res, err := AiGetMaterialFloatArray(pType, data, &iMax)
	if err != nil {
		return nil, err
	}
	v = &core.AiUVTransform{
		Translation: common.NewAiVector2D(float32(res[0]), float32(res[1])),
		Scaling:     common.NewAiVector2D(float32(res[2]), float32(res[3])),
		Rotation:    float32(res[4]),
	}
	return
}

// ------------------------------------------------------------------------------------------------
// Get a color (3 or 4 floats) from the material
func AiGetMaterialColor(pType uint32, data []byte) (v *common.AiColor4D, err error) {
	v = &common.AiColor4D{}
	iMax := 4
	res, err := AiGetMaterialFloatArray(pType, data, &iMax)
	if err != nil {
		return nil, err
	}
	// if no alpha channel is defined: set it to 1.0
	if 3 == iMax {
		v.A = 1.0
	}
	v.R = float32(res[0])
	v.G = float32(res[1])
	v.B = float32(res[2])
	return v, err
}

// ------------------------------------------------------------------------------------------------
// Get a string from the material
func AiGetPropertyString(pType uint32, data []byte) (res string, err error) {
	if core.AiPTI_String == core.AiPropertyTypeInfo(pType) {
		if len(data) < 5 {
			return res, errors.New("GetPropertyString invalid Data Length")
		}
		// The string is stored as 32 but length prefix followed by zero-terminated UTF8 data
		res = string(data[4:])
	} else {
		// TODO - implement lexical cast as well
		return res, fmt.Errorf("Material property %v  was found, but is no string")
	}
	return res, nil
}
