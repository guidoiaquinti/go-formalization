package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

func main() {
	const mask = 0x0000ffffc0000200
	a := address{4, 0, mask + 129, false}
	b := address{4, 0, mask + 2, false}
	mapped := address{6, 0, a.Lo, false}
	zone := address{6, 0x20010db800000000, 0x123456, true}
	cases := []fixture{
		{a, b, 24, 24, fmt.Sprintf("4,0,%d,24|4,0,%d,24|false|true|true|1", mask, mask)},
		{mapped, b, 120, 24, fmt.Sprintf("6,0,%d,120|6,0,%d,120|false|false|false|1", mask, mask)},
		{a, address{}, 33, 0, "error:prefix length 33 too large for IPv4|0,0,0,-1|false|false|false|1"},
		{zone, zone, 64, 64, fmt.Sprintf("6,%d,0,64|6,%d,0,64|false|false|true|0", zone.Hi, zone.Hi)},
		{address{}, address{}, -1, 0, "error:negative Prefix bits|0,0,0,-1|true|false|false|0"},
	}
	var input bytes.Buffer
	var zoneTag = map[bool]int{true: 1}
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d %d %d %d %d %d\n",
			c.A.Family, c.A.Hi, c.A.Lo, zoneTag[c.A.Zone], c.N,
			c.B.Family, c.B.Hi, c.B.Lo, zoneTag[c.B.Zone], c.M)
	}
	rows := check.LeanRows("godiff-netip-prefix-ops", len(cases), &input)
	file, err := os.CreateTemp("", "golean-netip-prefix-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_PREFIX_OPS_CASES", file.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanPrefixOps", "difftest/net/netip/prefix-ops/testdata/prefix_ops_test.go"))
	if strings.HasSuffix(rows[0], "|0") || strings.Contains(rows[1], "|true|") ||
		strings.Contains(rows[2], "|true|") || !strings.Contains(rows[3], "|false|false|true|") {
		log.Fatal("incorrect compare tiebreak, mapped family, invalid zero, or zone mutant accepted")
	}
	fmt.Println("Five independent pinned native and Lean fixtures passed; four mutation witnesses rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
