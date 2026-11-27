package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	N                int
	Bit              uint8
	Missing          bool
	Hi, Lo, VHi, VLo uint64
	Want             string
}

func main() {
	f, err := os.CreateTemp("", "golean-uint128-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	compare("fixed", boundaries(), rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				c := fixture{rng.Intn(258) - 1, uint8(rng.Uint32()), i%11 == 0, rng.Uint64(), rng.Uint64(), rng.Uint64(), rng.Uint64(), ""}
				switch i % 13 {
				case 0:
					c.Lo = 0
				case 1:
					c.Lo = math.MaxUint64
				case 2:
					c.Hi, c.Lo = 0, 0
				case 3:
					c.Hi, c.Lo = math.MaxUint64, math.MaxUint64
				}
				cases[i] = c
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_UINT128_CASES", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanUint128", "difftest/net/netip/uint128/testdata/uint128_test.go"))
}
func boundaries() []fixture {
	var cases []fixture
	vals := [][2]uint64{{0, 0}, {0, 1}, {1, 0}, {1, math.MaxUint64}, {math.MaxUint64, 0}, {math.MaxUint64, math.MaxUint64}, {1 << 63, 1 << 63}}
	for i := range 256 {
		for j, u := range vals {
			v := vals[(i+j+1)%len(vals)]
			cases = append(cases, fixture{i, uint8(i), i%2 == 0, u[0], u[1], v[0], v[1], ""})
		}
	}
	for _, n := range []int{math.MinInt, -2, -1, math.MaxInt} {
		cases = append(cases, fixture{N: n, Hi: math.MaxUint64, Lo: 1, VHi: 3, VLo: 4})
	}
	return cases
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		missing := 0
		if c.Missing {
			missing = 1
		}
		fmt.Fprintf(&wire, "%d %d %d %d %d %d %d\n", c.N, c.Bit, missing, c.Hi, c.Lo, c.VHi, c.VLo)
	}
	for i, got := range check.LeanRows("godiff-netip-uint128", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("netip uint128 %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
