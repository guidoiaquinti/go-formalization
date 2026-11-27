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
	Family                 int
	Hi, Lo                 uint64
	Port                   uint16
	Bits, Length           int
	Zone, Backing, Payload []byte
	Want                   string
}

func canonical(c fixture) fixture {
	width := 128
	switch c.Family {
	case 0:
		c.Hi, c.Lo, c.Zone, width = 0, 0, nil, 0
	case 4:
		c.Hi, c.Lo, c.Zone, width = 0, 0xffff00000000|(c.Lo&0xffffffff), nil, 32
	}
	if width == 0 || c.Bits < 0 || c.Bits > width {
		c.Bits = -1
	}
	return c
}
func main() {
	f, err := os.CreateTemp("", "golean-netip-binary-codec-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	var fixed []fixture
	for _, family := range []int{0, 4, 6} {
		for _, zone := range [][]byte{nil, {0, 255, 37}} {
			for _, spare := range []int{0, 1, 2, 3, 4, 7, 8, 9, 15, 16, 17, 18, 19, 20, 21, 22} {
				for n := range 26 {
					raw := bytes.Repeat([]byte{123}, n)
					if n > 0 {
						raw[n-1] = 255
					}
					fixed = append(fixed, canonical(fixture{family, 0x0123456789abcdef, 0xffff01020304, 0xabcd, 33, 2, zone, bytes.Repeat([]byte{78}, 2+spare), raw, ""}))
				}
			}
		}
	}
	for _, family := range []int{0, 4, 6} {
		fixed = append(fixed, canonical(fixture{Family: family}), canonical(fixture{Family: family, Backing: []byte{}}))
	}
	compare("fixed", fixed, rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				length := rng.Intn(13)
				backing := make([]byte, length+rng.Intn(37))
				zone := make([]byte, rng.Intn(13))
				payload := make([]byte, rng.Intn(43))
				rng.Read(backing)
				rng.Read(zone)
				rng.Read(payload)
				if i%11 == 0 {
					length, backing = 0, nil
				}
				c := fixture{[]int{0, 4, 6}[rng.Intn(3)], rng.Uint64(), rng.Uint64(), uint16(rng.Uint32()), rng.Intn(132) - 1, length, zone, backing, payload, ""}
				if i%7 == 0 {
					c.Hi, c.Lo = 0, 0xffff00000000|(c.Lo&0xffffffff)
				}
				cases[i] = canonical(c)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_NETIP_BINARY_CASES", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanBinaryCodec", "difftest/net/netip/binary-codec/testdata/binary_codec_test.go"))
}
func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		nilTag := 0
		if c.Backing == nil {
			nilTag = 1
		}
		fmt.Fprintf(&wire, "%d %d %d %d %d %d %d %d %d", c.Family, c.Hi, c.Lo, c.Port, c.Bits, len(c.Zone), c.Length, len(c.Backing), nilTag)
		for _, part := range [][]byte{c.Zone, c.Backing, c.Payload} {
			for _, v := range part {
				fmt.Fprintf(&wire, " %d", v)
			}
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-netip-binary-codec", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("netip binary codec %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
