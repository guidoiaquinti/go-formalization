package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	X    int64
	Big  string
	Tag  int
	Want string
}

func oracle(c fixture) string {
	real := strconv.Itoa(c.Tag)
	if c.Tag == 0 {
		real = "nil"
	} else if c.Tag == 1 {
		real = "typedNil"
	}
	return fmt.Sprintf("0|0/1|0/512|%d|%d/1|%d/512|%s/1|%s:0", c.X, c.X, c.X, c.Big, real)
}

func main() {
	file, err := os.CreateTemp("", "golean-constant-conversion-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	writer := bufio.NewWriter(file)
	enc := json.NewEncoder(writer)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.X, c.Big, c.Tag)
		}
		for i, result := range check.LeanRows("godiff-go-constant-numeric-conversions", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if result != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, result, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d observations\n", label, len(cases))
	}
	var fixed []fixture
	for _, x := range []int64{math.MinInt64, -1, 0, 1, 1 << 53, 1<<53 + 1, math.MaxInt64} {
		for _, shift := range []uint{0, 63, 512, 600} {
			for tag := range 5 {
				fixed = append(fixed, fixture{x, new(big.Int).Lsh(big.NewInt(x), shift).String(), tag, ""})
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			x := int64(r.Uint64())
			cases[i] = fixture{x, new(big.Int).Lsh(big.NewInt(x), uint(r.Intn(650))).String(), r.Intn(5), ""}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_CONSTANT_CONVERSION_ROWS", file.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanNumericConversions", "difftest/go/constant/numeric-conversions/testdata/conversion_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
