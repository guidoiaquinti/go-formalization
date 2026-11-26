package main

import (
	"encoding/binary"
	"math/bits"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	check.Exhaustive("godiff", 1<<16, func(x int) (int, int) {
		var b [2]byte
		binary.LittleEndian.PutUint16(b[:], uint16(x))
		return int(bits.ReverseBytes16(uint16(x))), int(binary.BigEndian.Uint16(b[:]))
	})
}
