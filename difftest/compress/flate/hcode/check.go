package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	word   uint32
	code   uint16
	length uint8
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-hcode")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	native := bufio.NewWriter(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.word, c.code, c.length)
		}
		for i, lean := range check.LeanRows("godiff-flate-hcode", len(cases), &input) {
			c := cases[i]
			packed := uint32(c.code)*256 + uint32(c.length)
			want := fmt.Sprintf("%d %d %d %t %d", packed, c.word%256, c.word/256, c.word == 0, packed)
			if lean != want {
				log.Fatalf("%s row %d %+v: Lean %q want %q", label, i, c, lean, want)
			}
			_, err := fmt.Fprintln(native, c.word, c.code, c.length, want)
			must(err)
		}
		fmt.Printf("%s: %d Lean inputs agree\n", label, len(cases))
	}
	var boundaries []fixture
	for _, word := range []uint32{0, 1, 255, 256, 65535, 65536, 0xff000000, math.MaxUint32} {
		for _, code := range []uint16{0, 1, 255, 256, 0x8000, math.MaxUint16} {
			for _, length := range []uint8{0, 1, 8, 15, 16, 17, 127, 255} {
				boundaries = append(boundaries, fixture{word, code, length})
			}
		}
	}
	compare("boundaries", boundaries)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{r.Uint32(), uint16(r.Uint32()), uint8(r.Uint32())}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(native.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_HCODE_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanHCode", "difftest/compress/flate/hcode/testdata/hcode_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
