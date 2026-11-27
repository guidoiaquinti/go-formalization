package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	boundary := new(big.Int).Lsh(big.NewInt(1), 4095)
	below := new(big.Int).Sub(boundary, big.NewInt(1)).String()
	at := boundary.String()
	fixed := []struct {
		Word                           int64
		Big                            string
		Present, Fits32, Fits63, Small bool
	}{
		{math.MinInt64, "-18446744073709551617", true, false, false, true},
		{math.MaxInt64, "18446744073709551617", true, false, false, true},
		{-2147483649, "0", true, false, true, true},
		{-2147483648, "-1", true, true, true, true},
		{2147483647, "1", true, true, true, true},
		{2147483648, "64", false, false, true, true},
		{-4611686018427387905, below, true, false, false, true},
		{-4611686018427387904, "-" + at, true, false, true, false},
		{4611686018427387903, at, true, false, true, false},
		{4611686018427387904, "-" + below, true, false, false, true},
		{0, at, false, true, true, false},
	}
	f, err := os.CreateTemp("", "golean-constant-render-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	wants := make([]string, 0, len(fixed))
	for _, c := range fixed {
		present := 0
		large := "<nil>"
		if c.Present {
			present = 1
			large = c.Big
		}
		fmt.Fprintln(&input, c.Word, c.Big, present)
		word := strconv.FormatInt(c.Word, 10)
		want := fmt.Sprintf("%s|%s|%s|%s|%t|%t|%t", word, word, large, large, c.Fits32, c.Fits63, c.Small)
		wants = append(wants, want)
		must(enc.Encode(struct {
			Word    int64
			Big     string
			Present bool
			Want    string
		}{c.Word, c.Big, c.Present, want}))
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-go-constant-integer-render", len(fixed), &input) {
		if got != wants[i] {
			log.Fatalf("fixture %d: Lean differs from exact control", i)
		}
	}
	must(os.Setenv("GOLEAN_CONSTANT_INTEGER_RENDER", f.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanIntegerRender", "difftest/go/constant/integer-render/testdata/render_test.go"))
	fmt.Printf("%d pinned Go and Lean controls for nil rendering and signed and magnitude width thresholds\n", len(fixed))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
