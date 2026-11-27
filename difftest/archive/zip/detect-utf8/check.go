package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type observation struct {
	Raw            []byte
	Valid, Require bool
}

func oracle(raw []byte) observation {
	result := observation{Raw: raw, Valid: utf8.Valid(raw)}
	if !result.Valid {
		return result
	}
	for _, b := range raw {
		if b < 32 || b > 125 || b == '\\' {
			result.Require = true
			break
		}
	}
	return result
}

func main() {
	file, err := os.CreateTemp("", "golean-zip-detect-utf8")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	writer := bufio.NewWriter(file)
	enc := json.NewEncoder(writer)
	compare := func(label string, cases [][]byte) {
		var input bytes.Buffer
		for _, raw := range cases {
			fmt.Fprint(&input, "0")
			for _, b := range raw {
				fmt.Fprintf(&input, " %d", b)
			}
			fmt.Fprintln(&input)
		}
		for i, line := range check.LeanRows("godiff-zip-detect-utf8", len(cases), &input) {
			expected := oracle(cases[i])
			if want := fmt.Sprintf("%t %t", expected.Valid, expected.Require); line != want {
				log.Fatalf("%s %x: Lean %q, want %q", label, cases[i], line, want)
			}
			must(enc.Encode(expected))
		}
		fmt.Printf("%s: %d observations\n", label, len(cases))
	}
	fixed := [][]byte{nil, []byte("é\\x80"), []byte("\\x80é"), []byte("é�"), []byte("\\xed\\xa0\\x80"), []byte("\\xf4\\x90\\x80\\x80")}
	for _, tail := range []string{"", "ok", "é", "\xff"} {
		fixed = append(fixed, []byte("\\\x80"+tail))
	}
	for b := range 256 {
		fixed = append(fixed, []byte{byte(b)})
	}
	compare("fixed", fixed)
	pairs := make([][]byte, 0, 256*256)
	for b := range 256 {
		for c := range 256 {
			pairs = append(pairs, []byte{byte(b), byte(c)})
		}
	}
	compare("all two-byte strings", pairs)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([][]byte, 20_000)
		for i := range cases {
			if i%2 == 0 {
				cases[i] = make([]byte, r.Intn(34))
				_, err := r.Read(cases[i])
				must(err)
			} else {
				var payload []byte
				for n := r.Intn(12); n > 0; n-- {
					rn := rune(r.Intn(0x110000))
					if r.Intn(2) == 0 {
						rn = rune(r.Intn(128))
					}
					payload = utf8.AppendRune(payload, rn)
				}
				cases[i] = payload
			}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_ZIP_DETECT_UTF8_CASES", file.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanDetectUTF8", "difftest/archive/zip/detect-utf8/testdata/detect_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
