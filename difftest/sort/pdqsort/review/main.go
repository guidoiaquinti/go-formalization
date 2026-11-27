package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Variant, A, B, Limit int
	Values               []int
	Want                 string
}

func main() {
	ascending := make([]int, 13)
	for i := range ascending {
		ascending[i] = i
	}
	descending := make([]int, 50)
	for i := range descending {
		descending[i] = 50 - i
	}
	base := []fixture{
		{0, 0, 0, 0, nil, ""},
		{0, 0, 2, 0, []int{2, 1}, ""},
		{0, 0, 13, 0, ascending, ""},
		{0, 0, 13, 1, ascending, ""},
		{0, 0, 50, 1, descending, ""},
	}
	var input bytes.Buffer
	var cases []fixture
	for _, c := range base {
		for variant := range 2 {
			c.Variant = variant
			cases = append(cases, c)
			fmt.Fprintf(&input, "%d %d %d %d %d", variant, strconv.IntSize, c.A, c.B, c.Limit)
			for _, value := range c.Values {
				fmt.Fprintf(&input, " %d", value)
			}
			fmt.Fprintln(&input)
		}
	}
	rows := check.LeanRows("godiff-sort-pdqsort", len(cases), &input)
	for variant := range 2 {
		if rows[variant] != ":" || rows[2+variant] != "1,2:L1,0;S1,0" {
			log.Fatal("empty or insertion trace mismatch")
		}
		_, heap, _ := strings.Cut(rows[4+variant], ":")
		_, normal, _ := strings.Cut(rows[6+variant], ":")
		_, reverse, _ := strings.Cut(rows[8+variant], ":")
		if !strings.HasPrefix(heap, "L11,12;") || !strings.HasPrefix(normal, "L6,3;") ||
			!strings.HasPrefix(reverse, "L12,11;") || !strings.Contains(reverse, ";S0,49;S1,48;") {
			log.Fatal("missing heap fallback, median or descending ninther trace")
		}
		if strings.HasPrefix(heap, "L6,3;") || strings.HasPrefix(normal, "L11,12;") ||
			strings.HasPrefix(reverse, "L24,12;") {
			log.Fatal("fallback- or pivot-threshold mutant was accepted")
		}
	}
	file, err := os.CreateTemp("", "golean-pdqsort-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		c.Want = rows[i]
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_PDQ_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanPdqsort", "difftest/sort/pdqsort/testdata/pdqsort_test.go"))
	fmt.Println("Ten independent fixtures matched pinned Go and Lean, with six fallback and pivot-threshold witnesses.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
