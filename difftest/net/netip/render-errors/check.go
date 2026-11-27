package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Input, Message, Suffix []byte
	Want                   string
}

func main() {
	f, err := os.CreateTemp("", "golean-netip-render-*")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, len(c.Input), " ", len(c.Message))
			for _, field := range [][]byte{c.Input, c.Message, c.Suffix} {
				for _, b := range field {
					fmt.Fprint(&input, " ", b)
				}
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-netip-render-errors", len(cases), &input) {
			c := cases[i]
			c.Want = row
			must(encoder.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows recorded\n", label, len(cases))
	}
	fixed := []fixture{{}, {[]byte("IPv6\""), []byte{'c', 0xff}, []byte("x\n"), ""}}
	for b := range 256 {
		fixed = append(fixed, fixture{Input: []byte{byte(b)}, Message: []byte{byte(b)}, Suffix: []byte{byte(b)}})
	}
	for r := rune(1); r <= unicode.MaxRune; r++ {
		if strconv.IsPrint(r) != strconv.IsPrint(r-1) {
			fixed = append(fixed, fixture{Input: []byte(string(r)), Message: []byte("print boundary"), Suffix: []byte(string(r - 1))})
		}
	}
	compare("boundary", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		random := func(mode int) []byte {
			size := rng.Intn(12)
			if mode == 0 {
				raw := make([]byte, size)
				rng.Read(raw)
				return raw
			}
			var raw []byte
			for range size {
				raw = utf8.AppendRune(raw, rune(rng.Intn(unicode.MaxRune+1)))
			}
			return raw
		}
		cases := make([]fixture, 20000)
		for i := range cases {
			cases[i] = fixture{Input: random(i % 2), Message: random(0), Suffix: random((i / 2) % 2)}
		}
		compare(fmt.Sprint("seed ", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_NETIP_RENDER_ERRORS", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanRenderErrors", "difftest/net/netip/render-errors/testdata/render_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
