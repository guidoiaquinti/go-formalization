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
	Family        int
	Hi, Lo        uint64
	Zone, Backing []byte
	Length        int
	Want          string
}

func main() {
	f, err := os.CreateTemp("", "golean-netip-addr-text-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	rng := rand.New(rand.NewSource(17))
	for profile := range 4 {
		var fixed []fixture
		for mask := range 256 {
			for _, spare := range []int{0, 1, 2, 3, 4, 6, 7, 8, 14, 15, 39, 40, 43, 48} {
				fixed = append(fixed, sample(rng, profile, mask, spare))
			}
		}
		compare(fmt.Sprintf("profile=%d fixed", profile), fixed, rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			blocks := 4
			if profile == 0 {
				blocks = 1
			}
			for block := range blocks {
				cases := make([]fixture, 5000)
				for i := range cases {
					cases[i] = sample(rng, profile, rng.Intn(256), rng.Intn(64))
				}
				compare(fmt.Sprintf("profile=%d seed=%d block=%d", profile, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_NETIP_ADDR_TEXT_CASES", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanAddrText", "difftest/net/netip/addr-text/testdata/addr_text_test.go"))
}

func sample(rng *rand.Rand, profile, mask, spare int) fixture {
	c := fixture{Length: rng.Intn(6)}
	c.Backing = make([]byte, c.Length+spare)
	rng.Read(c.Backing)
	if mask%11 == 0 && spare == 0 {
		c.Length, c.Backing = 0, nil
	}
	if profile == 0 {
		return c
	}
	if profile < 3 {
		c.Family = 4
		c.Lo = 0xffff00000000 | uint64(rng.Uint32())
		if profile == 2 {
			c.Family = 6
		}
	} else {
		c.Family = 6
		for i := range 8 {
			field := uint64(0)
			if mask&(1<<i) == 0 {
				field = uint64(rng.Intn(65535) + 1)
			}
			if i < 4 {
				c.Hi = c.Hi<<16 | field
			} else {
				c.Lo = c.Lo<<16 | field
			}
		}
		if c.Hi == 0 && c.Lo>>32 == 65535 {
			c.Hi = 1
		}
	}
	if c.Family == 6 {
		c.Zone = make([]byte, rng.Intn(18))
		rng.Read(c.Zone)
	}
	return c
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		missing := 0
		if c.Backing == nil {
			missing = 1
		}
		fmt.Fprintf(&wire, "%d %d %d %d %d %d", c.Family, c.Hi, c.Lo, len(c.Zone), c.Length, missing)
		for _, part := range [][]byte{c.Zone, c.Backing} {
			for _, v := range part {
				fmt.Fprintf(&wire, " %d", v)
			}
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-netip-addr-text", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("netip address text %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
