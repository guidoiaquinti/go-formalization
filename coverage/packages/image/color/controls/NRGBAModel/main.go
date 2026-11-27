package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"runtime"
	"strings"
)

type probe struct {
	calls int
}

func (p *probe) RGBA() (uint32, uint32, uint32, uint32) {
	p.calls++
	return 65538, 1, 0x8000ff00, 3
}

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}

	p8, p16 := new(probe), new(probe)
	want8 := color.NRGBA{R: 85, G: 85}
	want16 := color.NRGBA64{R: 21844, G: 21845, B: 85, A: 3}
	if got := color.NRGBAModel.Convert(p8); got != want8 || p8.calls != 1 {
		log.Fatalf("NRGBAModel converted to %#v with %d calls, want %#v with one call", got, p8.calls, want8)
	}
	if got := color.NRGBA64Model.Convert(p16); got != want16 || p16.calls != 1 {
		log.Fatalf("NRGBA64Model converted to %#v with %d calls, want %#v with one call", got, p16.calls, want16)
	}

	r, g, b, a := new(probe).RGBA()
	wide := func(channel uint32) uint64 { return uint64(channel) * 65535 / uint64(a) }
	mutant8 := color.NRGBA{R: uint8(wide(r) >> 8), G: uint8(wide(g) >> 8), B: uint8(wide(b) >> 8), A: uint8(a >> 8)}
	mutant16 := color.NRGBA64{R: uint16(wide(r)), G: uint16(wide(g)), B: uint16(wide(b)), A: uint16(a)}
	if mutant8 != (color.NRGBA{R: 170, G: 85, B: 171}) ||
		mutant16 != (color.NRGBA64{R: 43690, G: 21845, B: 43776, A: 3}) {
		log.Fatalf("unexpected widened-product outputs: %#v and %#v", mutant8, mutant16)
	}
	fmt.Printf("%s %s/%s: both adapters agree and reject the widened-product mutant\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
