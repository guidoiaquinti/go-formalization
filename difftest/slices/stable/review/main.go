package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                             []int64
	Kind, A, M, B, Mode, FailAt, Mask int
	Want                              string
}

func main() {
	merged := []int64{77, 258, 514, 257, 513, 88}
	cases := []fixture{
		{merged, 0, 1, 3, 5, 1, -1, 0, "ok/77,257,258,513,514,88/|ok/77,258,257,514,513,88/257,514;513,258;257,258;513,514"},
		{merged, 0, 1, 3, 5, 1, 2, 0, "ok/77,257,258,513,514,88/|panic:2/77,258,257,514,513,88/257,514;513,258;257,258"},
		{[]int64{258, 257, 3, 77}, 1, 0, 0, 3, 1, -1, 0, "ok/3,257,258,77/|ok/3,258,257,77/257,258;3,257;3,258"},
		{[]int64{769, 514, 257, 77}, 1, 0, 0, 3, 1, 1, 0, "ok/257,514,769,77/|panic:1/514,769,257,77/514,769;257,769"},
		{[]int64{258, 77}, 1, 0, 0, 1, 0, -1, 0, "ok/258,77/|ok/258,77/"},
		{[]int64{258, 257, 77}, 1, 0, 0, 2, 0, -1, 0, "ok/257,258,77/|nil/258,257,77/"},
	}
	mutants := []string{
		"ok/77,257,258,513,514,88/|ok/77,257,258,513,514,88/257,514;513,258;257,258;513,514",
		"ok/77,257,258,513,514,88/|panic:2/77,258,514,257,513,88/257,514;513,258;257,258",
		"ok/3,257,258,77/|ok/3,257,258,77/257,258;3,257;3,258",
		"ok/257,514,769,77/|panic:1/769,514,257,77/514,769;257,769",
		"ok/258,77/|nil/258,77/",
		"ok/257,258,77/|ok/258,257,77/",
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d %d %d", c.Kind, c.A, c.M, c.B, c.Mode, c.FailAt, c.Mask)
		for _, v := range c.Words {
			fmt.Fprintf(&input, " %d", v)
		}
		fmt.Fprintln(&input)
	}
	rejected := 0
	for i, got := range check.LeanRows("godiff-slices-stable", len(cases), &input) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, cases[i].Want)
		}
		if got != mutants[i] {
			rejected++
		}
	}
	if rejected != len(mutants) {
		log.Fatalf("rejected %d of %d variants", rejected, len(mutants))
	}
	f, err := os.CreateTemp("", "golean-slices-stable-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_SLICES_STABLE_CASES", f.Name()))
	const base = "difftest/slices/stable/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanStable", base+"bridge_test.go", base+"stable_test.go"))
	fmt.Printf("Stable helper review: %d native and Lean fixtures, %d altered outputs rejected\n", len(cases), rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
