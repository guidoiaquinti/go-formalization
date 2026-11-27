package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

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
	const all = math.MaxUint64
	cases := []fixture{
		{64, 65, false, 1, 0, 3, 7, "18446744073709551615,0|false|1,0|2,7|3,7|18446744073709551614,18446744073709551615|0,18446744073709551615|1,1|1,0:3,7|1,9223372036854775807|1,0"},
		{65, 64, false, all, all, 1, 2, "18446744073709551615,9223372036854775808|false|1,2|18446744073709551614,18446744073709551613|18446744073709551615,18446744073709551615|0,0|18446744073709551615,18446744073709551614|0,0|18446744073709551615,18446744073709551615:1,2|18446744073709551615,18446744073709551615|18446744073709551615,0"},
		{-1, 129, true, 0, 1, 2, 4, "panic|false|0,0|2,5|2,5|18446744073709551615,18446744073709551614|0,0|0,2|panic|panic|panic"},
		{129, 128, false, 0, 0, 11, 22, "panic|true|0,0|11,22|11,22|18446744073709551615,18446744073709551615|18446744073709551615,18446744073709551615|0,1|0,0:11,22|0,0|0,0"},
		{0, 0, false, 0, all, 4, 3, "0,0|false|0,3|4,18446744073709551612|4,18446744073709551615|18446744073709551615,0|0,18446744073709551614|1,0|0,18446744073709551615:4,3|18446744073709551615,18446744073709551615|0,0"},
	}
	pair := func(a, b uint64) string { return fmt.Sprintf("%d,%d", a, b) }
	reject := func(row, field int, mutation string) {
		if strings.Split(cases[row].Want, "|")[field] == mutation {
			panic("mutated observation survived")
		}
	}
	carry := cases[4]
	reject(4, 7, pair(carry.Hi, carry.Lo+1))
	borrow := cases[0]
	reject(0, 6, pair(borrow.Hi, borrow.Lo-1))
	reject(0, 8, pair(borrow.Lo, borrow.Hi)+":"+pair(borrow.VLo, borrow.VHi))
	reject(0, 10, pair(borrow.Hi, borrow.Lo|((uint64(1)<<63)-1)))
	reject(0, 0, pair(0, ^uint64(0)))
	reject(0, 1, fmt.Sprint(borrow.Lo == 0))
	cross := cases[1]
	reject(1, 2, pair(cross.Hi&cross.VLo, cross.Lo&cross.VHi))
	reject(0, 3, pair(borrow.Hi|borrow.VHi, borrow.Lo|borrow.VLo))
	reject(0, 4, pair(borrow.Hi^borrow.VHi, borrow.Lo^borrow.VLo))
	reject(0, 5, pair(^borrow.Hi, borrow.Lo))
	reject(4, 9, pair(carry.Hi, carry.Lo))
	fmt.Println("all eleven independent mutation witnesses rejected")

	file, err := os.CreateTemp("", "review-netip-u128-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	var protocol bytes.Buffer
	for _, c := range cases {
		missing := 0
		if c.Missing {
			missing = 1
		}
		fmt.Fprintln(&protocol, c.N, c.Bit, missing, c.Hi, c.Lo, c.VHi, c.VLo)
		must(encoder.Encode(c))
	}
	must(file.Close())
	for i, row := range check.LeanRows("godiff-netip-uint128", len(cases), &protocol) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("row %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_UINT128_CASES", file.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanUint128", "difftest/net/netip/uint128/testdata/uint128_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
