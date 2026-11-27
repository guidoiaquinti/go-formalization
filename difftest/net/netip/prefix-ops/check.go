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

type address struct {
	Family int
	Hi, Lo uint64
	Zone   bool
}
type fixture struct {
	A, B address
	N, M int
	Want string
}

func canonical(a address) address {
	if a.Family == 0 {
		return address{}
	}
	if a.Family == 4 {
		return address{4, 0, 0xffff00000000 | (a.Lo & 0xffffffff), false}
	}
	return a
}

func boundaries() []fixture {
	var out []fixture
	for _, a := range []address{{}, {Family: 4}, {4, 0, 0xc00002fe, false}, {Family: 6}, {6, 0, 0xffffc00002fe, true}, {6, 0x20010db8f7ffff00, 0x765432100000ff01, false}} {
		a = canonical(a)
		for _, b := range []address{{}, a, canonical(address{a.Family, a.Hi, a.Lo ^ 1, !a.Zone}),
			canonical(address{a.Family, a.Hi ^ 1<<63, a.Lo ^ 1<<31, false}), canonical(address{4, 0, a.Lo, false}), canonical(address{6, a.Hi, a.Lo, false})} {
			for _, n := range []int{math.MinInt, -1, 0, 1, 8, 16, 31, 32, 33, 64, 65, 95, 96, 97, 127, 128, 129, math.MaxInt} {
				for _, m := range []int{-1, 0, 24, 31, 32, 64, 95, 120, 128, 129} {
					out = append(out, fixture{a, b, n, m, ""})
				}
			}
		}
	}
	return out
}

func main() {
	f, err := os.CreateTemp("", "golean-prefix-ops-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	compare("fixed", boundaries(), rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		randomAddress := func() address {
			family := []int{0, 4, 6, 6, 6}[rng.Intn(5)]
			hi, lo := rng.Uint64(), rng.Uint64()
			if rng.Intn(5) == 0 {
				hi, lo = 0, 0xffff00000000|lo&0xffffffff
			}
			return canonical(address{family, hi, lo, rng.Intn(4) == 0})
		}
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				a, b := randomAddress(), randomAddress()
				if i%2 == 0 {
					b = a
					if i%3 != 0 {
						b.Lo ^= 1 << rng.Intn(32)
					}
					if i%5 == 0 {
						b.Zone = !b.Zone
					}
					b = canonical(b)
				}
				width := func(kind int) int {
					if kind == 4 {
						return 32
					}
					return 128
				}
				n, m := rng.Intn(width(a.Family)+7)-3, rng.Intn(width(b.Family)+7)-3
				if i%11 == 0 {
					m = n
				}
				cases[i] = fixture{a, b, n, m, ""}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_PREFIX_OPS_CASES", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanPrefixOps", "difftest/net/netip/prefix-ops/testdata/prefix_ops_test.go"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		for _, a := range []struct {
			address
			bits int
		}{{c.A, c.N}, {c.B, c.M}} {
			zone := 0
			if a.Zone {
				zone = 1
			}
			fmt.Fprintf(&wire, "%d %d %d %d %d ", a.Family, a.Hi, a.Lo, zone, a.bits)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-netip-prefix-ops", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("netip prefix ops %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
