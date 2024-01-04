package test

import (
	"assimp/common/logger"
	"bytes"
	"math"
	"reflect"
	"unsafe"
)

const Eps = 0.01

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

func deepValueEqual(v1, v2 reflect.Value, visited map[visit]bool) bool {
	if !v1.IsValid() || !v2.IsValid() {
		var t1, t2 string
		if v1.IsValid() {
			t1 = v1.Type().String()
		}
		if v2.IsValid() {
			t2 = v2.Type().String()
		}
		logger.ErrorF("v1.Slice.IsValid:%v != v2.Slice.IsValid:%v v1:%v v2:%v",
			t1,
			t2,
			v1.IsValid(),
			v2.IsValid())
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		logger.ErrorF("v1.Slice.Type:%v != v2.Slice.Type:%v v1:%v v2:%v",
			v1.Type().String(),
			v2.Type().String(),
			v1.Type().Kind(),
			v2.Type().Kind())
		return false
	}

	hard := func(v1, v2 reflect.Value) bool {
		switch v1.Kind() {
		case reflect.Pointer:
			v1.Pointer()
			fallthrough
		case reflect.Map, reflect.Slice, reflect.Interface:
			return !v1.IsNil() && !v2.IsNil()
		}
		return false
	}

	if hard(v1, v2) {
		ptrval := func(v reflect.Value) unsafe.Pointer {
			switch v.Kind() {
			case reflect.Pointer, reflect.Map:
				return v.UnsafePointer()
			default:
				return v.UnsafePointer()
			}
		}
		addr1 := ptrval(v1)
		addr2 := ptrval(v2)
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			logger.ErrorF("v1.Slice.IsNil:%v != v2.Slice.IsNil:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.IsNil(),
				v2.IsNil())
			return false
		}
		if v1.Len() != v2.Len() {
			logger.ErrorF("v1.Slice.Len:%v != v2.Slice.Len:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.Len(),
				v2.Len())
			return false
		}
		if v1.UnsafePointer() == v2.UnsafePointer() {
			return true
		}
		// Special case for []byte, which is common.
		if v1.Type().Elem().Kind() == reflect.Uint8 {
			return bytes.Equal(v1.Bytes(), v2.Bytes())
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited) {
				logger.ErrorF("not equal index:%v type1:%v type2:%v", i, v1.Type().String(), v2.Type().String())
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			logger.ErrorF("v1.Interface.IsNil:%v != v2.Interface.IsNil:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.IsNil(),
				v2.IsNil())
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited)
	case reflect.Pointer:
		if v1.UnsafePointer() == v2.UnsafePointer() {
			return true
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepValueEqual(v1.Field(i), v2.Field(i), visited) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.UnsafePointer() == v2.UnsafePointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || !deepValueEqual(val1, val2, visited) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v1.Int() != v2.Int() {
			logger.ErrorF("v1.Int:%v != v2.Int:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.Int(),
				v2.Int())
		}
		return v1.Int() == v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if v1.Uint() != v2.Uint() {
			logger.ErrorF("v1.uInt:%v != v2.uInt:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.Uint(),
				v2.Uint())
		}
		return v1.Uint() == v2.Uint()
	case reflect.String:
		if v1.String() != v2.String() {
			logger.ErrorF("v1.String:%v != v2.String:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.String(),
				v2.String())
		}
		return v1.String() == v2.String()
	case reflect.Bool:
		if v1.Bool() != v2.Bool() {
			logger.ErrorF("v1.String:%v != v2.String:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.Bool(),
				v2.Bool())
		}
		return v1.Bool() == v2.Bool()
	case reflect.Float32, reflect.Float64: //浮点数比较精度
		if math.Abs(v1.Float()-v2.Float()) > Eps {
			logger.ErrorF("v1.Float:%v != v2.Float:%v v1:%v v2:%v eps:%v delta:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.Float(),
				v2.Float(), Eps, math.Abs(v1.Float()-v2.Float()))
		}
		return math.Abs(v1.Float()-v2.Float()) <= Eps
	case reflect.Complex64, reflect.Complex128:
		if v1.Complex() != v2.Complex() {
			logger.ErrorF("v1.Complex:%v != v2.Complex:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.Complex(),
				v2.Complex())
		}
		return v1.Complex() == v2.Complex()
	default:
		if v1.Interface() != v2.Interface() {
			logger.ErrorF("v1.Interface:%v != v2.Interface:%v v1:%v v2:%v",
				v1.Type().String(),
				v2.Type().String(),
				v1.Interface(),
				v2.Interface())
		}
		// Normal equality suffices
		return v1.Interface() == v2.Interface()
	}
}
func deepEqual(x, y any) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		logger.ErrorF("v1.Type:%v != v2.Type:%v", v1.Type().String(), v2.Type().String())
		return false
	}
	return deepValueEqual(v1, v2, make(map[visit]bool))
}
