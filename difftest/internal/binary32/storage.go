// Package binary32 provides a byte-order oracle for raw Go float32 storage.
package binary32

import (
	"encoding/binary"
	"unsafe"
)

// Store copies explicit native-endian bytes into a float32 object.
func Store(word uint32) float32 {
	var f float32
	binary.NativeEndian.PutUint32((*[4]byte)(unsafe.Pointer(&f))[:], word)
	return f
}

// Load reads explicit native-endian bytes from a float32 object.
func Load(f float32) uint32 {
	return binary.NativeEndian.Uint32((*[4]byte)(unsafe.Pointer(&f))[:])
}

// Boundaries crosses both signs with transitions between the binary32 classes.
func Boundaries() []uint32 {
	magnitudes := []uint32{0, 1, 0x007fffff, 0x00800000,
		0x3f800000, 0x7f7fffff, 0x7f800000,
		0x7f800001, 0x7fbfffff, 0x7fc00000,
		0x7fc00001, 0x7fffffff}
	var words []uint32
	for _, sign := range []uint32{0, 1 << 31} {
		for _, magnitude := range magnitudes {
			words = append(words, sign|magnitude)
		}
	}
	return words
}
