// Package binary64 provides a byte-order oracle for raw Go float64 storage.
package binary64

import (
	"encoding/binary"
	"unsafe"
)

// Store copies explicit native-endian bytes into a float64 object.
func Store(word uint64) float64 {
	var f float64
	binary.NativeEndian.PutUint64((*[8]byte)(unsafe.Pointer(&f))[:], word)
	return f
}

// Load reads explicit native-endian bytes from a float64 object.
func Load(f float64) uint64 {
	return binary.NativeEndian.Uint64((*[8]byte)(unsafe.Pointer(&f))[:])
}

// Boundaries crosses both signs with transitions between the binary64 classes.
func Boundaries() []uint64 {
	magnitudes := []uint64{0, 1, 0x000fffffffffffff, 0x0010000000000000,
		0x3ff0000000000000, 0x7fefffffffffffff, 0x7ff0000000000000,
		0x7ff0000000000001, 0x7ff7ffffffffffff, 0x7ff8000000000000,
		0x7ff8000000000001, 0x7fffffffffffffff}
	var words []uint64
	for _, sign := range []uint64{0, 1 << 63} {
		for _, magnitude := range magnitudes {
			words = append(words, sign|magnitude)
		}
	}
	return words
}
