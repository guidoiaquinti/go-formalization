package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"runtime"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}

	values := []uint16{0, 1, 0x8000, 0xffff}
	var actual, swapped int
	for _, red := range values {
		for _, green := range values {
			for _, blue := range values {
				for _, alpha := range values {
					c := color.RGBA64{R: red, G: green, B: blue, A: alpha}
					r, g, b, a := c.RGBA()
					want := [4]uint32{uint32(red), uint32(green), uint32(blue), uint32(alpha)}
					if [4]uint32{r, g, b, a} != want {
						actual++
					}
					if [4]uint32{g, r, b, a} != want {
						swapped++
					}
				}
			}
		}
	}
	if actual != 0 || swapped != 192 {
		log.Fatalf("RGBA64.RGBA failures: Go=%d, red-green swap=%d, want 0 and 192", actual, swapped)
	}
	fmt.Printf("%s %s/%s: RGBA64.RGBA Go=%d/256, red-green swap=%d/256\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, actual, swapped)
}
