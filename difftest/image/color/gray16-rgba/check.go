package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executable = "godiff-image-color-gray16-rgba"

func main() {
	const count = 1 << 16
	for y, model := range check.LeanRows(executable, count, nil) {
		r, g, b, a := (color.Gray16{Y: uint16(y)}).RGBA()
		want := fmt.Sprintf("%d %d %d %d", y, y, y, 65535)
		if r != uint32(y) || g != uint32(y) || b != uint32(y) || a != 65535 || model != want {
			log.Fatalf("Gray16.RGBA input %#04x: Go=(%d %d %d %d), Lean=%q, expected=%q", y, r, g, b, a, model, want)
		}
	}
	fmt.Printf("%s: %s and Lean agree with the specification on all %d inputs.\n", executable, runtime.Version(), count)
}
