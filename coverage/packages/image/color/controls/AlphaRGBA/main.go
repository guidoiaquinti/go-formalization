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
	fmt.Println(runtime.Version(), runtime.GOOS, runtime.GOARCH)
	checkAll("Alpha.RGBA", 256, 257, func(x uint32) (uint32, uint32, uint32, uint32) {
		return (color.Alpha{A: uint8(x)}).RGBA()
	})
	checkAll("Alpha16.RGBA", 65536, 1, func(x uint32) (uint32, uint32, uint32, uint32) {
		return (color.Alpha16{A: uint16(x)}).RGBA()
	})
}

func checkAll(name string, size, scale uint32, rgba func(uint32) (uint32, uint32, uint32, uint32)) {
	var actualFailures, alphaOnlyFailures uint32
	for x := range size {
		want := [4]uint32{x * scale, x * scale, x * scale, x * scale}
		r, g, b, a := rgba(x)
		if [4]uint32{r, g, b, a} != want {
			actualFailures++
		}
		if [4]uint32{0, 0, 0, x * scale} != want {
			alphaOnlyFailures++
		}
	}
	if actualFailures != 0 || alphaOnlyFailures != size-1 {
		log.Fatalf("%s: Go failures=%d, alpha-only mutant failures=%d, want 0 and %d",
			name, actualFailures, alphaOnlyFailures, size-1)
	}
	fmt.Printf("%s: Go failures=%d, alpha-only mutant failures=%d/%d\n",
		name, actualFailures, alphaOnlyFailures, size)
}
