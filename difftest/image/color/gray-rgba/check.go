package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executable = "godiff-image-color-gray-rgba"

func main() {
	const count = 1 << 8
	for y, model := range check.LeanRows(executable, count, nil) {
		r, g, b, a := (color.Gray{Y: uint8(y)}).RGBA()
		channel := uint32(y) * 257
		want := fmt.Sprintf("%d %d %d %d", channel, channel, channel, 65535)
		if r != channel || g != channel || b != channel || a != 65535 || model != want {
			log.Fatalf("Gray.RGBA input %#02x: Go=(%d %d %d %d), Lean=%q, expected=%q", y, r, g, b, a, model, want)
		}
	}
	fmt.Printf("%s: %s and Lean agree with the specification on all %d inputs.\n", executable, runtime.Version(), count)
}
