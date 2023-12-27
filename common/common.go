package common

import (
	"assimp/common/logger"
	"assimp/common/pb_msg"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"unsafe"
)

type AiQuaternion struct {
	W, X, Y, Z float32
}

func (ai *AiQuaternion) ToPbMsg() *pb_msg.AiQuaternion {
	return &pb_msg.AiQuaternion{
		W: ai.W,
		X: ai.X,
		Y: ai.Y,
		Z: ai.Z,
	}
}

type AiPropertyStore struct {
	Sentinel uint8
}

// ------------------------------------------------------------------------------------
// Convert just one hex digit
// Return value is UINT_MAX if the input character is not a hex digit.
// ------------------------------------------------------------------------------------
func HexDigitToDecimal(in byte) (out uint) {
	if in >= '0' && in <= '9' {
		out = uint(in - '0')
	} else if in >= 'a' && in <= 'f' {
		out = 10 + uint(in-'a')
	} else if in >= 'A' && in <= 'F' {
		out = 10 + uint(in-'A')
		return out
	}

	// return value is UINT_MAX if the input is not a hex digit
	return out
}

func GetBytesOrder() binary.ByteOrder {
	if IsLittleEndian() {
		return binary.LittleEndian
	}
	return binary.BigEndian
}

func IsLittleEndian() bool {
	n := 0x1234
	return *(*byte)(unsafe.Pointer(&n)) == 0x34
}

// ------------------------------------------------------------------------------------
// Convert a string in decimal format to a number
// ------------------------------------------------------------------------------------
func Strtoul10(in string) (value int32, index int) {
	for i, v := range in {
		index = i
		if v < '0' || v > '9' {
			break
		}
		value = (value * 10) + (v - '0')
	}
	return value, index
}

// ------------------------------------------------------------------------------------
// signed variant of strtoul10
// ------------------------------------------------------------------------------------
func StrTol10(inData string) (value int32, index int) {
	in := 0
	inv := (inData[in] == '-')
	if inv || inData[in] == '+' {
		in++
	}

	value, index = Strtoul10(inData)
	if inv {
		if value < math.MaxInt32 {
			value = -value
		} else {
			logger.WarnF("Converting the string \"%v \" into an inverted value resulted in overflow.", in)
		}
	}
	return value, in
}

// find >=
func LowerBound(begin, end int, less func(index int) bool) int {
	for begin < end {
		half := begin + (end-begin)>>1
		if less(half) {
			begin = half + 1
		} else {
			end = half
		}
	}
	return end
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}



// -------------------------------------------------------------------------------
/** @brief Find the min/max values of an array of Ts
 *  @param in Input array
 *  @param size Number of elements to process
 *  @param[out] min minimum value
 *  @param[out] max maximum value
 */

func  ArrayBounds[T any ]( in []T) {
MinMaxChooser<T>()(min, max);
for  i := 0; i < len(t); i++ {
min = std::min(in[i], min);
max = std::max(in[i], max);
}
}
