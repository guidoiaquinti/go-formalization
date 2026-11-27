package main

import (
	"fmt"
	"log"
	"math/bits"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type span struct{ count, first, width int }

func table(parts []span, capacity int) (string, int) {
	var words []int
	changed := 0
	for _, part := range parts {
		for i := range part.count {
			forward := part.first + i
			reverse := int(bits.Reverse16(uint16(forward)) >> (16 - part.width))
			if reverse != forward {
				changed++
			}
			words = append(words, reverse*256+part.width)
		}
	}
	var row strings.Builder
	fmt.Fprintf(&row, "%d %d true", len(words), capacity)
	for _, w := range words {
		fmt.Fprintf(&row, " %d", w)
	}
	for i := len(words); i < capacity; i++ {
		row.WriteString(" 0")
	}
	return row.String(), changed
}

func main() {
	literal, badLiteral := table([]span{{144, 48, 8}, {112, 400, 9}, {24, 0, 7}, {6, 192, 8}}, 512)
	distance, badDistance := table([]span{{30, 0, 5}}, 32)
	if badLiteral != 266 || badDistance != 23 {
		log.Fatalf("unreversed differences = %d/%d", badLiteral, badDistance)
	}
	expected := []string{literal, distance}
	for i, got := range check.LeanRows("godiff-flate-fixed-huffman", len(expected), strings.NewReader("0\n1\n")) {
		if got != expected[i] {
			log.Fatalf("fixed table %d differs from independent enumeration", i)
		}
	}
	file, err := os.CreateTemp("", "golean-flate-fixed-review-")
	must(err)
	defer os.Remove(file.Name())
	_, err = fmt.Fprintln(file, strings.Join(expected, "\n"))
	must(err)
	must(file.Close())
	must(os.Setenv("GOLEAN_FLATE_FIXED_ROWS", file.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanFixedHuffman", "difftest/compress/flate/fixed-huffman/testdata/fixed_test.go"))
	fmt.Printf("Fixed Huffman review: 316 populated and 228 spare words matched; unreversed mutants differ at %d literal and %d distance words\n", badLiteral, badDistance)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
