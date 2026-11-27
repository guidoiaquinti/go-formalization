package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tag     int
	Integer string
	Bytes   []byte
	Want    string
}

func main() {
	f, err := os.CreateTemp("", "golean-constant-integer-")
	must(err)
	defer os.Remove(f.Name())
	out := bufio.NewWriter(f)
	rows := json.NewEncoder(out)
	var fixed []fixture
	for _, bits := range []int{0, 7, 8, 15, 16, 31, 32, 62, 63, 64, 65, 95, 127, 128, 129, 191, 255, 256} {
		pivot := new(big.Int).Lsh(big.NewInt(1), uint(bits))
		for _, sign := range []int64{1, -1} {
			for _, delta := range []int64{-1, 0, 1} {
				x := new(big.Int).Mul(pivot, big.NewInt(sign))
				x.Add(x, big.NewInt(delta))
				for tag := range 4 {
					for _, raw := range [][]byte{nil, {}, {0}, {1, 0}, {0, 0, 0, 0, 0, 0, 0, 128}, {255, 255, 255, 255, 255, 255, 255, 255}, {0, 0, 0, 0, 0, 0, 0, 0, 1}} {
						fixed = append(fixed, fixture{tag, x.String(), raw, ""})
					}
				}
			}
		}
	}
	compare("fixed", fixed, rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				bits := make([]byte, rng.Intn(40))
				rng.Read(bits)
				n := new(big.Int).SetBytes(bits)
				if i%2 == 0 {
					n.Neg(n)
				}
				raw := make([]byte, rng.Intn(40))
				rng.Read(raw)
				if i%5 == 0 && len(raw) > 0 {
					raw[len(raw)-1] = 0
				}
				tag := i % 4
				cases[i] = fixture{tag, n.String(), raw, ""}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_CONSTANT_INTEGER_CASES", f.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanConstantInteger", "difftest/go/constant/integer/testdata/integer_test.go"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %s", c.Tag, c.Integer)
		for _, b := range c.Bytes {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-go-constant-integer", len(cases), &input) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("constant integer %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
