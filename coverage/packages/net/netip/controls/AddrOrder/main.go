package main

import (
	"fmt"
	"log"
	"net/netip"
	"strconv"
)

func main() {
	v4 := netip.AddrFrom4([4]byte{192, 0, 2, 1})
	mapped := netip.AddrFrom16(v4.As16())
	v6 := netip.MustParseAddr("2001:db8::1")
	cases := []struct {
		name  string
		a, b  netip.Addr
		order int
	}{
		{"zero before IPv4", netip.Addr{}, v4, -1},
		{"IPv4 before mapped IPv6", v4, mapped, -1},
		{"equal zero receivers", netip.Addr{}, netip.Addr{}, 0},
		{"unscoped before scoped", v6, v6.WithZone("1"), -1},
		{"lexical order", v6.WithZone("10"), v6.WithZone("2"), -1},
		{"zone prefix", v6.WithZone("1"), v6.WithZone("10"), -1},
		{"same zone", v6.WithZone("2"), v6.WithZone("2"), 0},
		{"mapped zones", mapped.WithZone("10"), mapped.WithZone("2"), -1},
		{"zone after words", v6.WithZone("9"), netip.MustParseAddr("2001:db8::2").WithZone("1"), -1},
		{"unsigned high word", netip.MustParseAddr("7fff::ffff"), netip.MustParseAddr("8000::"), -1},
		{"unsigned low word", netip.MustParseAddr("2001:db8::7fff:ffff:ffff:ffff"), netip.MustParseAddr("2001:db8::8000:0:0:0"), -1},
	}
	var ignoreZones, numericZones, unmapped, nonStrict int
	for _, c := range cases {
		for _, reverse := range []bool{false, true} {
			a, b, want := c.a, c.b, c.order
			if reverse {
				a, b, want = b, a, -want
			}
			if got := a.Compare(b); got != want {
				log.Fatalf("%s, reversed=%t: Compare=%d want %d", c.name, reverse, got, want)
			}
			if got := a.Less(b); got != (want < 0) {
				log.Fatalf("%s, reversed=%t: Less=%t want %t", c.name, reverse, got, want < 0)
			}
			withoutZones := a.WithZone("").Compare(b.WithZone(""))
			if withoutZones != want {
				ignoreZones++
			}
			// Mutant parses tied IPv6 zone strings as numbers.
			numeric := withoutZones
			if numeric == 0 && a.Is6() {
				az, _ := strconv.Atoi(a.Zone())
				bz, _ := strconv.Atoi(b.Zone())
				switch {
				case az < bz:
					numeric = -1
				case az > bz:
					numeric = 1
				}
			}
			if numeric != want {
				numericZones++
			}
			if a.Unmap().Compare(b.Unmap()) != want {
				unmapped++
			}
			if (a.Compare(b) <= 0) != (want < 0) {
				nonStrict++
			}
		}
	}
	if ignoreZones != 8 || numericZones != 4 || unmapped != 4 || nonStrict != 4 {
		log.Fatalf("unexpected mutation failures: ignore zones=%d numeric zones=%d unmap=%d non-strict=%d",
			ignoreZones, numericZones, unmapped, nonStrict)
	}
	fmt.Printf("Addr order: %d directed fixtures passed, mutant failures ignore zones=%d numeric zones=%d unmap=%d non-strict=%d\n",
		2*len(cases), ignoreZones, numericZones, unmapped, nonStrict)
}
