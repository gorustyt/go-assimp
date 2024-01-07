package ASSETBIN

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gorustyt/go-assimp/common"
	"github.com/gorustyt/go-assimp/common/logger"
	"github.com/gorustyt/go-assimp/common/pb_msg"
	"github.com/gorustyt/go-assimp/core"
	"google.golang.org/protobuf/proto"
	"math"
)

func getAiMaterialPropertyData(p *core.AiMaterialProperty, tType uint32, data []byte) (err error) {
	pType := core.AiPropertyTypeInfo(tType)
	var v proto.Message
	switch pType {
	case core.AiPTI_Float:
		p.Type = pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat64
		res, err := AiGetMaterialFloat(tType, data)
		if err != nil {
			return err
		}
		v = &pb_msg.AiMaterialPropertyFloat64{
			Data: res,
		}
	case core.AiPTI_Double:
		p.Type = pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat64
		res, err := AiGetMaterialFloat(tType, data)
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
		res, err := AiGetMaterialInteger(tType, data)
		if err != nil {
			return err
		}
		v = &pb_msg.AiMaterialPropertyInt64{
			Data: res,
		}
	default:
		return fmt.Errorf("invalid  prop name:%v", p.Key)
	}
	p.Data, err = proto.Marshal(v)
	return err
}

func parseAiMaterialByKey(p *core.AiMaterialProperty, data []byte) error {
	var err error
	if len(core.GetAiPropertyTypeInfo(p.Key)) == 0 {
		return fmt.Errorf("parseAiMaterialByKey length ==0 prop:%v", p.Key)
	}
	for _, v := range core.GetAiPropertyTypeInfo(p.Key) {
		switch v {
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeString:
			err = getAiMaterialPropertyData(p, uint32(core.AiPTI_String), data)
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat32:
			err = getAiMaterialPropertyData(p, uint32(core.AiPTI_Float), data)
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeFloat64:
			err = getAiMaterialPropertyData(p, uint32(core.AiPTI_Double), data)
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeInt:
			err = getAiMaterialPropertyData(p, uint32(core.AiPTI_Integer), data)
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeVector3D:
			var res *common.AiVector4D
			res, err = AiGetMaterialVector(uint32(core.AiPTI_Buffer), data)
			if err == nil {
				p.Data, err = proto.Marshal(res.ToPbMsg())
				if err != nil {
					return err
				}
				p.Type = v
			}

		case pb_msg.AiMaterialPropertyType_AiPropertyTypeVector4D:
			var res *common.AiVector4D
			res, err = AiGetMaterialVector(uint32(core.AiPTI_Buffer), data)
			if err == nil {
				tmp := common.NewAiColor3D(res[0], res[1], res[2])
				p.Data, err = proto.Marshal(tmp.ToPbMsg())
				if err != nil {
					return err
				}
				p.Type = v
			}
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeAiUVTransform:
			var res *core.AiUVTransform
			res, err = AiGetMaterialUVTransform(uint32(core.AiPTI_Buffer), data)
			if err == nil {
				p.Data, err = proto.Marshal(res.ToPbMsg())
				if err != nil {
					return err
				}
				p.Type = v
			}
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeColor3D:
			var res *common.AiColor4D
			res, err = AiGetMaterialColor(uint32(core.AiPTI_Buffer), data)
			if err == nil {
				tmp := common.NewAiColor3D(res.R, res.G, res.B)
				p.Data, err = proto.Marshal(tmp.ToPbMsg())
				if err != nil {
					return err
				}
				p.Type = v
			}
		case pb_msg.AiMaterialPropertyType_AiPropertyTypeColor4D:
			var res *common.AiColor4D
			res, err = AiGetMaterialColor(uint32(core.AiPTI_Buffer), data)
			if err == nil {
				p.Data, err = proto.Marshal(res.ToPbMsg())
				if err != nil {
					return err
				}
				p.Type = v
			}
		}
		if err != nil {
			logger.WarnF("not found this prop type:%v continue found", v.String())
			continue
		}

		return nil
	}
	return fmt.Errorf("not found prop name:%v", p.Key)
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
// Get a color (3 or 4 floats) from the material
func AiGetMaterialVector(pType uint32, data []byte) (v *common.AiVector4D, err error) {
	v = &common.AiVector4D{}
	iMax := 4
	res, err := AiGetMaterialFloatArray(pType, data, &iMax)
	if err != nil {
		return nil, err
	}
	// if no alpha channel is defined: set it to 1.0
	if 3 == iMax {
		v[3] = 1.0
	}
	v[0] = float32(res[0])
	v[1] = float32(res[1])
	v[2] = float32(res[2])
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
		res = string(data[4 : len(data)-1])
	} else {
		// TODO - implement lexical cast as well
		return res, fmt.Errorf("Material property %v  was found, but is no string")
	}
	return res, nil
}

// ---------------------------------------------------------------------------
/** @brief Retrieve an integer property with a specific key from a material
 *
 * See the sample for aiGetMaterialFloat for more information.*/
// ---------------------------------------------------------------------------
func AiGetMaterialInteger(pType uint32, data []byte) (res []int64, err error) {
	return AiGetMaterialIntegerArray(pType, data, nil)
}

// ---------------------------------------------------------------------------
/** @brief Retrieve a single float property with a specific key from the material.
*
* Pass one of the AI_MATKEY_XXX constants for the last three parameters (the
* example reads the #AI_MATKEY_SHININESS_STRENGTH property of the first diffuse texture)
* @code
* float specStrength = 1.f; // default value, remains unmodified if we fail.
* aiGetMaterialFloat(mat, AI_MATKEY_SHININESS_STRENGTH,
*    (float*)&specStrength);
* @endcode
*
* @param pMat Pointer to the input material. May not be NULL
* @param pKey Key to search for. One of the AI_MATKEY_XXX constants.
* @param pOut Receives the output float.
* @param type (see the code sample above)
* @param index (see the code sample above)
* @return Specifies whether the key has been found. If not, the output
*   float remains unmodified.*/
// ---------------------------------------------------------------------------
func AiGetMaterialFloat(pType uint32, data []byte) (res []float64, err error) {
	return AiGetMaterialFloatArray(pType, data, nil)
}
