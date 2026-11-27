package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	const infinity = uint64(0x7ff0000000000000)
	cases := []struct {
		word uint64
		exp  int
		want uint64
	}{
		{0x3ff0000000000000, -1075, 0},
		{0x3ff8000000000000, -1075, 1},
		{0x3ff8000000000000, -1074, 2},
		{0x3ff4000000000000, -1074, 1},
		{0xbff0000000000000, -1075, 1 << 63},
		{1, 1074, 0x3ff0000000000000},
		{0x4000000000000000, math.MaxInt, 0},
		{0x3fe0000000000000, math.MinInt, infinity},
		{1 << 63, math.MaxInt, 1 << 63},
		{0xfff8000000000042, math.MinInt, 0xfff8000000000042},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&rows, c.word, c.exp, c.want)
	}
	var noWrapFailures int
	for _, c := range cases[6:8] {
		input := math.Float64frombits(c.word)
		total := new(big.Int).Add(big.NewInt(int64(c.exp)), big.NewInt(int64(math.Ilogb(input))))
		var mutant uint64
		switch {
		case total.Cmp(big.NewInt(1023)) > 0:
			mutant = infinity
		case total.Cmp(big.NewInt(-1075)) < 0:
		default:
			log.Fatalf("unexpected unrestricted exponent %s", total)
		}
		if mutant != c.want {
			noWrapFailures++
		}
	}
	if noWrapFailures != 2 {
		log.Fatalf("unrestricted-scaling mutant failed %d cases, want two", noWrapFailures)
	}

	dir, err := os.MkdirTemp("", "golean-ldexp-generic-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_LDEXP_GENERIC_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/ldexp-generic/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanLdexpGeneric", fixtures+"bridge_test.go", fixtures+"ldexp_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s/%s: ldexp passed 10 fixtures; unrestricted scaling disagreed at both native exponent limits\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
