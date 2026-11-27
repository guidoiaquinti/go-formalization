package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executable = "godiff-image-color-alpha16-rgba"

func main() {
	const count = 1 << 16
	for input, model := range check.LeanRows(executable, count, nil) {
		r, g, b, a := (color.Alpha16{A: uint16(input)}).RGBA()
		channel := uint32(input)
		want := fmt.Sprintf("%d %d %d %d", channel, channel, channel, channel)
		if r != channel || g != channel || b != channel || a != channel || model != want {
			log.Fatalf("Alpha16.RGBA input %#x: Go=(%d %d %d %d), Lean=%q, expected=%q", input, r, g, b, a, model, want)
		}
	}
	fmt.Printf("%s: %s and Lean agree with the specification on all %d inputs.\n", executable, runtime.Version(), count)
}
