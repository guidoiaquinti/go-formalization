package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Family    int
	Hi, Lo    uint64
	Index     uint8
	Zone, Raw []byte
	Want      string
}

func canonical(c fixture) fixture {
	switch c.Family {
	case 0:
		c.Hi, c.Lo, c.Zone = 0, 0, nil
	case 4:
		c.Hi, c.Lo, c.Zone = 0, 0xffff00000000|(c.Lo&0xffffffff), nil
	}
	return c
}

func main() {
	f, err := os.CreateTemp("", "golean-addr-fields-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	var fixed []fixture
	for i := range 256 {
		for _, family := range []int{0, 4, 6} {
			for _, zone := range [][]byte{nil, []byte("link"), {0, 255, 37}} {
				raw := make([]byte, i%24)
				for j := range raw {
					raw[j] = byte(i + j)
				}
				fixed = append(fixed, canonical(fixture{family, 0x0123456789abcdef, 0xfedcba9876543210, uint8(i), zone, raw, ""}))
			}
		}
	}
	for _, raw := range [][]byte{nil, {}, {0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0}} {
		fixed = append(fixed, fixture{Raw: raw})
	}
	compare("fixed", fixed, rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				raw := make([]byte, []int{4, 16, rng.Intn(44)}[i%3])
				zone := make([]byte, rng.Intn(12))
				rng.Read(raw)
				rng.Read(zone)
				c := fixture{[]int{0, 4, 6}[rng.Intn(3)], rng.Uint64(), rng.Uint64(), uint8(rng.Uint32()), zone, raw, ""}
				if i%5 == 0 {
					c.Hi, c.Lo = 0, 0xffff00000000|(c.Lo&0xffffffff)
				}
				cases[i] = canonical(c)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_ADDR_FIELDS_CASES", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanAddrFields", "difftest/net/netip/addr-fields/testdata/addr_fields_test.go"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Family, c.Hi, c.Lo, c.Index, len(c.Zone))
		for _, part := range [][]byte{c.Zone, c.Raw} {
			for _, b := range part {
				fmt.Fprintf(&wire, " %d", b)
			}
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-netip-addr-fields", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("netip Addr fields %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
