package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	bits := math.Float64bits
	type fixture struct {
		word uint64
		want [2]uint64
	}
	same := func(word, want uint64) fixture {
		return fixture{word, [2]uint64{want, want}}
	}
	cases := []fixture{
		same(0, sign|inf), same(sign, sign|inf),
		same(inf, inf), same(sign|inf, inf+1),
		same(inf+1, inf+1), same(sign|inf|2, inf+1),
		same(bits(-1), inf+1), same(bits(1), 0),
		{1, [2]uint64{13867497231575896609, 13867767316136018803}},
		{2, [2]uint64{13867497231575896609, 13867762020360330133}},
		{3, [2]uint64{13867497231575896609, 13867758922530140030}},
		same(bits(10), bits(1)), same(bits(1e15), bits(math.Nextafter(15, 0))),
		same(bits(1000), bits(3)),
		same(0x3ff0000000000001, 4367307993162114317),
		same(0x3fefffffffffffff, 13586176430389519630),
	}
	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	profile := 0
	switch runtime.GOARCH {
	case "arm64":
		profile = 1
	case "amd64":
	default:
		log.Fatalf("unreviewed architecture %s", runtime.GOARCH)
	}
	var inputs, fixtureRows bytes.Buffer
	var swappedFailures, powerShortcutFailures int
	for _, c := range cases {
		x := math.Float64frombits(c.word)
		got := observe(c.want[profile])
		fmt.Fprintln(&fixtureRows, c.word, got, 0, true)
		if product := observe(bits(math.Log(x) * float64(1/math.Ln10))); got != product {
			log.Fatalf("Log10(%#x)=%s want rounded product %s", c.word, got, product)
		}
		if got != observe(c.want[1-profile]) {
			swappedFailures++
		}
		if c.word == bits(1e15) && got != observe(bits(15)) {
			powerShortcutFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-log10", 2*len(cases), &inputs) {
		if want := observe(cases[i/2].want[i%2]); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/2, i%2, row, want)
		}
	}
	if swappedFailures != 3 || powerShortcutFailures != 1 {
		log.Fatalf("exchanging profiles failed %d; exact-power shortcut failed %d", swappedFailures, powerShortcutFailures)
	}
	dir, err := os.MkdirTemp("", "golean-private-log10-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtureRows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_LOG10_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/log10/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanLog10", fixtures+"bridge_test.go", fixtures+"log10_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Private log10 passed %d fixtures; exchanging profiles failed %d; exact-power shortcut failed %d\n",
		len(cases), swappedFailures, powerShortcutFailures)
}
