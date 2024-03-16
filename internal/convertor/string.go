package convertor

import "unsafe"

// BytesToString 字节转字符串
// BytesToString converts bytes to string
func BytesToString(b []byte) string {
	// 使用unsafe.Pointer将字节切片转换为字符串
	// Use unsafe.Pointer to convert byte slice to string
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 字符串转字节
// StringToBytes converts string to bytes
func StringToBytes(s string) []byte {
	// 使用unsafe.Pointer将字符串转换为字节切片
	// Use unsafe.Pointer to convert string to byte slice
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
