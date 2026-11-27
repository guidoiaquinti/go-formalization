package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	f, err := os.CreateTemp("", "golean-flate-token-fields")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	native := bufio.NewWriter(f)
	compare := func(label string, words []uint32) {
		var input bytes.Buffer
		for _, word := range words {
			fmt.Fprintln(&input, word)
		}
		for i, row := range check.LeanRows("godiff-flate-token-fields", len(words), &input) {
			word := words[i]
			typ := word / 1073741824 * 1073741824
			length := (word - typ) / 4194304
			offset := word - typ - length*4194304
			literal := offset % 256
			want := fmt.Sprintf("%d %d %d %d", typ, literal, offset, length)
			if row != want {
				log.Fatalf("%s word %#08x: Lean=%q want=%q", label, word, row, want)
			}
			_, err := fmt.Fprintln(native, word, want)
			must(err)
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(words))
	}
	var edges []uint32
	for typ := range uint32(4) {
		for length := range uint32(256) {
			for _, offset := range []uint32{0, 1, 255, 256, 65535, 65536, 4194303} {
				edges = append(edges, typ<<30|length<<22|offset)
			}
		}
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		words := make([]uint32, 20_000)
		for i := range words {
			words[i] = r.Uint32()
		}
		compare(fmt.Sprintf("seed %d", seed), words)
	}
	must(native.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_TOKEN_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanTokenFields", "difftest/compress/flate/token-fields/testdata/token_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
