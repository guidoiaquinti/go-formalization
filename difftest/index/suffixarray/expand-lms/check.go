package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Mode, Bins int
	Text       []byte
	Prefix     []int
	Want       string
}

func main() {
	f, err := os.CreateTemp("", "golean-suffix-expand-")
	must(err)
	defer os.Remove(f.Name())
	out := bufio.NewWriter(f)
	enc := gob.NewEncoder(out)
	var fixed []fixture
	for size := 3; size <= 7; size++ {
		limit := 1
		for range size {
			limit *= 3
		}
		for code := range limit {
			text := make([]byte, size)
			q := code
			for i := range text {
				text[i] = []byte{0, 1, 255}[q%3]
				q /= 3
			}
			prefix := lms(text)
			if len(prefix) == 0 {
				continue
			}
			for mode := range 3 {
				fixed = append(fixed, fixture{mode, bins(text), text, prefix, ""})
			}
		}
	}
	compare("fixed", fixed, enc)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				text := make([]byte, rng.Intn(42)+3)
				rng.Read(text)
				if i%3 == 0 {
					for j := range text {
						text[j] %= 4
					}
				}
				prefix := lms(text)
				if len(prefix) == 0 {
					text = []byte{2, 1, 3}
					prefix = []int{1}
				}
				rows[i] = fixture{rng.Intn(3), bins(text), text, prefix, ""}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows, enc)
		}
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_SUFFIX_EXPAND_ROWS", f.Name()))
	must(check.TestInPackage("index/suffixarray", "TestGoLeanExpandLms", "difftest/index/suffixarray/expand-lms/testdata/expand_test.go"))
}

func bins(text []byte) int {
	return 2*int(slices.Max(text)) + 1
}

func lms(text []byte) []int {
	// Lexicographic order of neighboring suffixes independently identifies
	// S-type positions. A transition from L to S starts an LMS suffix.
	isS := make([]bool, len(text))
	for i := 0; i < len(text)-1; i++ {
		isS[i] = slices.Compare(text[i:], text[i+1:]) < 0
	}
	var out []int
	for i := 1; i < len(isS); i++ {
		if !isS[i-1] && isS[i] {
			out = append(out, i)
		}
	}
	slices.SortFunc(out, func(i, j int) int { return slices.Compare(text[i:], text[j:]) })
	return out
}

func compare(label string, rows []fixture, enc *gob.Encoder) {
	var input bytes.Buffer
	for _, r := range rows {
		fmt.Fprint(&input, r.Mode, " ", r.Bins, " ", len(r.Text))
		for _, b := range r.Text {
			fmt.Fprint(&input, " ", b)
		}
		for _, i := range r.Prefix {
			fmt.Fprint(&input, " ", i)
		}
		fmt.Fprintln(&input)
	}
	for i, want := range check.LeanRows("godiff-suffixarray-expand-lms", len(rows), &input) {
		r := rows[i]
		r.Want = want
		must(enc.Encode(r))
	}
	fmt.Printf("suffix expansion %s: %d Lean rows\n", label, len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
