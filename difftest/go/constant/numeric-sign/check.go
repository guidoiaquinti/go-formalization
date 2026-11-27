package main

import (
	"bufio"
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type row struct {
	X, Y int64
	Mode int
	Want string
}

func oracle(c row) string {
	x, y := cmp.Compare(c.X, 0), cmp.Compare(c.Y, 0)
	dominant := func(a, b int) int {
		if a < 0 || b < 0 {
			return -1
		}
		return max(a, b)
	}
	z := dominant(x, y)
	alternate, combined := "1", fmt.Sprint(dominant(1, y))
	if c.Mode == 1 {
		alternate, combined = "panic:<nil> not numeric", "panic:<nil> not numeric"
	} else if c.Mode == 2 {
		alternate, combined = "panic:true not numeric", "panic:true not numeric"
	}
	return fmt.Sprintf("3|3|4|4|5|%d|%d|%d|%d|%d|%d|%s|%s", x, x, x, x, z,
		dominant(z, -x), alternate, combined)
}

func main() {
	file, err := os.CreateTemp("", "golean-constant-sign-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	writer := bufio.NewWriter(file)
	enc := json.NewEncoder(writer)
	compare := func(label string, cases []row) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.X, c.Y, c.Mode)
		}
		for i, result := range check.LeanRows("godiff-go-constant-numeric-sign", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if result != c.Want {
				log.Fatalf("%s row %d: Lean %q, want %q", label, i, result, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d observations\n", label, len(cases))
	}
	var fixed []row
	edges := []int64{math.MinInt64, -(1 << 32), -1, 0, 1, 1 << 32, math.MaxInt64}
	for _, x := range edges {
		for _, y := range edges {
			for mode := range 3 {
				fixed = append(fixed, row{x, y, mode, ""})
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]row, 20_000)
		for i := range cases {
			x, y := int64(r.Uint64()), int64(r.Uint64())
			if i%5 == 0 {
				x = edges[r.Intn(len(edges))]
			}
			if i%7 == 0 {
				y = edges[r.Intn(len(edges))]
			}
			cases[i] = row{x, y, r.Intn(3), ""}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_CONSTANT_SIGN_ROWS", file.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanNumericSign", "difftest/go/constant/numeric-sign/testdata/sign_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
