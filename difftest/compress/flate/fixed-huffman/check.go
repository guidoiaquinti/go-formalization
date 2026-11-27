package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

// Compute canonical codes from length counts for the complete alphabet.
// Only afterward discard the reserved symbols absent from the Go encoder.
func canonical(widths []int, symbols, capacity int) string {
	var count, next [16]int
	for _, width := range widths {
		count[width]++
	}
	for width := 1; width < len(next); width++ {
		next[width] = (next[width-1] + count[width-1]) * 2
	}
	var row strings.Builder
	fmt.Fprintf(&row, "%d %d true", symbols, capacity)
	for _, width := range widths[:symbols] {
		code := next[width]
		next[width]++
		reversed := 0
		for range width {
			reversed = 2*reversed + code%2
			code /= 2
		}
		fmt.Fprint(&row, " ", reversed*256+width)
	}
	for i := symbols; i < capacity; i++ {
		row.WriteString(" 0")
	}
	return row.String()
}

func main() {
	literals := make([]int, 288)
	for i := range literals {
		switch {
		case i < 144 || i >= 280:
			literals[i] = 8
		case i < 256:
			literals[i] = 9
		default:
			literals[i] = 7
		}
	}
	offsets := make([]int, 32)
	for i := range offsets {
		offsets[i] = 5
	}
	want := []string{canonical(literals, 286, 512), canonical(offsets, 30, 32)}
	rows := check.LeanRows("godiff-flate-fixed-huffman", len(want), strings.NewReader("0\n1\n"))
	for i, row := range rows {
		if row != want[i] {
			log.Fatalf("Lean fixed table %d differs from canonical assignment", i)
		}
	}
	f, err := os.CreateTemp("", "golean-flate-fixed")
	must(err)
	defer os.Remove(f.Name())
	_, err = fmt.Fprintln(f, strings.Join(want, "\n"))
	must(err)
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_FIXED_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanFixedHuffman", "difftest/compress/flate/fixed-huffman/testdata/fixed_test.go"))
	fmt.Println("Both deterministic encoders agree on 316 code words and 228 reserved capacity words.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
