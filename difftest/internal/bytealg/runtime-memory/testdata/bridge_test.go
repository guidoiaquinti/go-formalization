package bytealg

import "unsafe"

func GoLeanMakeNoZero(length int) []byte {
	return MakeNoZero(length)
}

func GoLeanRuntimeMemequal(left, right []byte, size int) bool {
	return abigen_runtime_memequal(unsafe.Pointer(unsafe.SliceData(left)),
		unsafe.Pointer(unsafe.SliceData(right)), uintptr(size))
}
