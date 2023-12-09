package common

func SwapBytes2(data []byte) {
	data[0], data[1] = data[1], data[0]
}
func SwapBytes4(data []byte) {
	data[0], data[3] = data[3], data[1]
	data[1], data[2] = data[2], data[1]
}

func SwapBytes8(data []byte) {
	data[0], data[7] = data[7], data[0]
	data[1], data[6] = data[6], data[1]
	data[2], data[5] = data[5], data[2]
	data[3], data[4] = data[4], data[3]
}

func SwapFloat32() {

}

func SwapFloat64() {

}

func SwapInt32() {

}

func SwapUint32() {

}

func SwapInt64() {

}

func SwapUint64() {

}

func SwapIn16() {

}

func SwapUint16(data uint16) {

}
