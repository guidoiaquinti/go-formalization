package main

import (
	"fmt"
	"log"
	"math"
	"net/netip"
)

type fixture struct {
	input         netip.Addr
	request       int
	address       netip.Addr
	bits          int
	valid, single bool
}

func main() {
	zero := netip.Addr{}
	invalidBits := zero.Next()
	v4 := netip.MustParseAddr("192.0.2.15")
	v6 := netip.MustParseAddr("fe80::dead%rack")
	unzoned := netip.MustParseAddr("fe80::dead")
	mapped := netip.MustParseAddr("::ffff:192.0.2.15")
	cases := []fixture{
		{zero, 0, zero, -1, false, false},
		{invalidBits, 0, invalidBits, -1, false, false},
		{v4, 0, v4, 0, true, false},
		{v4, 24, v4, 24, true, false},
		{v4, 32, v4, 32, true, true},
		{v4, 33, v4, -1, false, false},
		{v6, 64, unzoned, 64, true, false},
		{v6, 128, unzoned, 128, true, true},
		{v6, math.MaxInt, unzoned, -1, false, false},
		{mapped, 32, mapped, 32, true, false},
		{mapped, 128, mapped, 128, true, true},
	}
	check := func(makePrefix func(netip.Addr, int) netip.Prefix, ok func(netip.Prefix, fixture) bool) int {
		fail := 0
		for _, c := range cases {
			if !ok(makePrefix(c.input, c.request), c) {
				fail++
			}
		}
		return fail
	}
	all := func(p netip.Prefix, c fixture) bool {
		return p.Addr() == c.address && p.Bits() == c.bits && p.IsValid() == c.valid && p.IsSingleIP() == c.single
	}
	assertCount("PrefixFrom and accessors", check(netip.PrefixFrom, all), 0, len(cases))
	assertCount("PrefixFrom masks host bits", check(func(a netip.Addr, n int) netip.Prefix {
		p := netip.PrefixFrom(a, n)
		if p.IsValid() {
			return p.Masked()
		}
		return p
	}, all), 4, len(cases))
	assertCount("Addr returns zero for invalid prefixes", check(netip.PrefixFrom, func(p netip.Prefix, c fixture) bool {
		a := netip.Addr{}
		if p.IsValid() {
			a = p.Addr()
		}
		return a == c.address
	}), 3, len(cases))
	assertCount("Bits rounds an invalid length to zero", check(netip.PrefixFrom, func(p netip.Prefix, c fixture) bool {
		return max(0, p.Bits()) == c.bits
	}), 4, len(cases))
	assertCount("IsValid rejects zero-length prefixes", check(netip.PrefixFrom, func(p netip.Prefix, c fixture) bool {
		return (p.Bits() > 0) == c.valid
	}), 1, len(cases))
	assertCount("IsSingleIP treats mapped IPv6 as IPv4", check(netip.PrefixFrom, func(p netip.Prefix, c fixture) bool {
		return (p.IsValid() && p.Bits() == p.Addr().Unmap().BitLen()) == c.single
	}), 2, len(cases))
}

func assertCount(name string, got, want, total int) {
	if got != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, got, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, got, total)
}
