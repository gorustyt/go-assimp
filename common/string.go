package common

import (
	"assimp/common/logger"
	"fmt"
	"math"
	"strings"
)

const NumItems = 16

var fast_atof_table = [NumItems]float64{ // we write [16] here instead of [] to work around a swig bug
	0.0,
	0.1,
	0.01,
	0.001,
	0.0001,
	0.00001,
	0.000001,
	0.0000001,
	0.00000001,
	0.000000001,
	0.0000000001,
	0.00000000001,
	0.000000000001,
	0.0000000000001,
	0.00000000000001,
	0.000000000000001,
}

// ------------------------------------------------------------------------------------
// ! Provides a fast function for converting a string into a float,
// ! about 6 times faster than atof in win32.
// If you find any bugs, please send them to me, niko (at) irrlicht3d.org.
// ------------------------------------------------------------------------------------
// Number of relevant decimals for floating-point parsing.
const AI_FAST_ATOF_RELAVANT_DECIMALS = 15

func FastAtorealMove(cData []byte, check_commas ...bool) (res float64, index int, err error) {
	f := float64(0.0)
	check_comma := true
	if len(check_commas) > 0 {
		check_comma = check_commas[0]
	}
	c := 0
	inv := (cData[c] == '-')
	if inv || cData[c] == '+' {
		c++
	}

	if (cData[c:][0] == 'N' || cData[c:][0] == 'n') && string(cData[c:c+3]) == "nan" {
		res = math.NaN()
		c += 3
		return res, c, nil
	}

	if (cData[c:][0] == 'I' || cData[c:][0] == 'i') && string(cData[c:c+3]) == "inf" {
		res = math.Inf(1)
		if inv {
			res = math.Inf(-1)
		}
		c += 3
		if (cData[c:][0] == 'I' || cData[c:][0] == 'i') && string(cData[c:c+5]) == "inity" {
			c += 5
		}
		return res, c, nil
	}

	if !(cData[c:][0] >= '0' && cData[c:][0] <= '9') &&
		!((cData[c:][0] == '.' || (check_comma && cData[c:][0] == ',')) && cData[c:][1] >= '0' && cData[c:][1] <= '9') {
		// The string is known to be bad, so don't risk printing the whole thing.
		return res, c, fmt.Errorf("cannot parse string \"%v\" as a real number: does not start with digit or decimal point followed by digit", AiSAtrToPrintable(string(cData[:c])))
	}

	if cData[c] != '.' && (!check_comma || cData[c:][0] != ',') {
		var tmp uint64
		tmp, c, err = StrToul10_64(cData[c:], nil)
		f = math.Float64frombits(tmp)
	}

	if (cData[c] == '.' || (check_comma && cData[c:][0] == ',')) && cData[c:][1] >= '0' && cData[c:][1] <= '9' {
		c++

		// NOTE: The original implementation is highly inaccurate here. The precision of a single
		// IEEE 754 float is not high enough, everything behind the 6th digit tends to be more
		// inaccurate than it would need to be. Casting to double seems to solve the problem.
		// strtol_64 is used to prevent integer overflow.

		// Another fix: this tends to become 0 for long numbers if we don't limit the maximum
		// number of digits to be read. AI_FAST_ATOF_RELAVANT_DECIMALS can be a value between
		// 1 and 15.
		diff := uint32(AI_FAST_ATOF_RELAVANT_DECIMALS)
		var tmp uint64
		tmp, c, err = StrToul10_64(cData[c:], &diff)
		pl := math.Float64frombits(tmp)
		if err != nil {
			return 0, 0, err
		}

		pl *= fast_atof_table[diff]
		f += pl
	} else if cData[c] == '.' { // For backwards compatibility: eat trailing dots, but not trailing commas.
		c++
	}

	// A major 'E' must be allowed. Necessary for proper reading of some DXF files.
	// Thanks to Zhao Lei to point out that this if() must be outside the if (*c == '.' ..)
	if cData[c] == 'e' || cData[c] == 'E' {
		c++
		einv := (cData[c] == '-')
		if einv || cData[c] == '+' {
			c++
		}

		// The reason float constants are used here is that we've seen cases where compilers
		// would perform such casts on compile-time constants at runtime, which would be
		// bad considering how frequently fast_atoreal_move<float> is called in Assimp.
		var tmp uint64
		tmp, c, err = StrToul10_64(cData[:c], nil)
		exp := math.Float64frombits(tmp)
		if err != nil {
			return 0, 0, err
		}
		if einv {
			exp = -exp
		}
		f *= math.Pow(10.0, exp)
	}

	if inv {
		f = -f
	}
	res = f
	return res, c, nil
}

// ------------------------------------------------------------------------------------
// Special version of the function, providing higher accuracy and safety
// It is mainly used by fast_atof to prevent ugly and unwanted integer overflows.
// ------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------
// signed variant of strtoul10_64
// ------------------------------------------------------------------------------------

func StrTol10_64(cData []byte, max_inout *uint32) (value uint64, index int, err error) {
	in := 0
	inv := (cData[in] == '-')
	if inv || cData[in] == '+' {
		in++
	}

	value, index, err = StrToul10_64(cData[in:], max_inout)
	if inv {
		value = -value
	}
	return value, index, err
}

func StrToul10_64(cData []byte, max_inout *uint32) (value uint64, index int, err error) {
	cur := uint32(0)
	in := 0
	if cData[in] < '0' || cData[in] > '9' {
		// The string is known to be bad, so don't risk printing the whole thing.
		return value, in, fmt.Errorf("the string \"%v \" cannot be converted into a value", AiSAtrToPrintable(string(cData)))
	}

	for {
		if cData[in] < '0' || cData[in] > '9' {
			break
		}

		new_value := (value * 10) + (uint64(cData[in] - '0'))

		// numeric overflow, we rely on you
		if new_value < value {
			logger.WarnF("Converting the string \"%v \" into a value resulted in overflow.", in)
			return value, in, nil
		}

		value = new_value

		in++
		cur++

		if max_inout != nil && *max_inout == cur {
			for cData[in] >= '0' && cData[in] <= '9' { /* skip to end */
				in++
			}
			value = uint64(in)

			return value, in, nil
		}
	}
	if max_inout != nil {
		*max_inout = cur
	}

	return value, in, nil
}

// ---------------------------------------------------------------------------------
// / @brief  Make a string printable by replacing all non-printable characters with
// /         the specified placeholder character.
// / @param  in  The incoming string.
// / @param  placeholder  Placeholder character, default is a question mark.
// / @return The string, with all non-printable characters replaced.
func AiSAtrToPrintable(in string, placeholders ...uint8) string {
	placeholder := uint8('?')
	if len(placeholders) > 0 {
		placeholder = placeholders[0]
	}
	var res []uint8
	for i, v := range in {
		if isprint(v) {
			res = append(res, in[i])
		} else {
			res = append(res, placeholder)
		}
	}
	return string(res)
}

func isprint(c int32) bool {
	return int(c-' ') < 127-' ' //判断字符c是否为可打印字符（含空格）。当c为可打印字符（0x20-0x7e）时，返回非零值，不然返回零。
}

func ClearQuotationMark(s string) string { //去除引号
	s = strings.ReplaceAll(s, "\"", "")
	return s
}
