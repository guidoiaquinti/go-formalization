package main

import (
	"cmp"
	"fmt"
	"log"
	"net/netip"
)

func main() {
	checkFields()
	checkCompare()
}

func checkFields() {
	zero := netip.Addr{}
	cases := []struct {
		addr  netip.Addr
		port  uint16
		valid bool
	}{
		{zero, 0, false},
		{zero.Next(), 65535, false},
		{netip.MustParseAddr("192.0.2.3"), 0, true},
		{netip.MustParseAddr("::ffff:192.0.2.3%zone"), 32768, true},
		{netip.MustParseAddr("fe80::1%scope"), 65535, true},
	}
	check := func(run func(netip.Addr, uint16, bool) bool) int {
		fail := 0
		for _, c := range cases {
			if !run(c.addr, c.port, c.valid) {
				fail++
			}
		}
		return fail
	}
	assertCount("AddrPortFrom and accessors", check(func(a netip.Addr, port uint16, valid bool) bool {
		value := netip.AddrPortFrom(a, port)
		return value.Addr() == a && value.Port() == port && value.IsValid() == valid
	}), 0, len(cases))
	assertCount("AddrPortFrom unmaps the supplied address", check(func(a netip.Addr, port uint16, _ bool) bool {
		return netip.AddrPortFrom(a.Unmap(), port).Addr() == a
	}), 1, len(cases))
	assertCount("Addr clears the stored zone", check(func(a netip.Addr, port uint16, _ bool) bool {
		return netip.AddrPortFrom(a, port).Addr().WithZone("") == a
	}), 2, len(cases))
	assertCount("Port clears the high bit", check(func(a netip.Addr, port uint16, _ bool) bool {
		return netip.AddrPortFrom(a, port).Port()&0x7fff == port
	}), 3, len(cases))
	assertCount("IsValid disallows port zero", check(func(a netip.Addr, port uint16, valid bool) bool {
		value := netip.AddrPortFrom(a, port)
		return (value.IsValid() && value.Port() != 0) == valid
	}), 1, len(cases))
}

func checkCompare() {
	makeAddrPort := func(a string, port uint16) netip.AddrPort {
		return netip.AddrPortFrom(netip.MustParseAddr(a), port)
	}
	zero := netip.Addr{}
	cases := []struct {
		a, b netip.AddrPort
		want int
	}{
		{makeAddrPort("192.0.2.3", 65535), makeAddrPort("::ffff:192.0.2.3", 0), -1},
		{makeAddrPort("fe80::1%10", 65535), makeAddrPort("fe80::1%2", 0), -1},
		{makeAddrPort("2001:db8:0:1::", 0), makeAddrPort("2001:db8::ffff:ffff:ffff:ffff", 65535), 1},
		{makeAddrPort("192.0.2.3", 32767), makeAddrPort("192.0.2.3", 32768), -1},
		{netip.AddrPortFrom(zero, 65535), netip.AddrPortFrom(zero.Next(), 0), -1},
		{makeAddrPort("fe80::1", 11), makeAddrPort("fe80::1%2", 11), -1},
		{makeAddrPort("fe80::1%zone", 11), makeAddrPort("fe80::1%zone", 11), 0},
	}
	for _, candidate := range []struct {
		name string
		run  func(netip.AddrPort, netip.AddrPort) int
		want int
	}{
		{"Compare", netip.AddrPort.Compare, 0},
		{"Compare checks ports first", func(a, b netip.AddrPort) int {
			if c := cmp.Compare(a.Port(), b.Port()); c != 0 {
				return c
			}
			return a.Addr().Compare(b.Addr())
		}, 4},
		{"Compare treats ports as signed", func(a, b netip.AddrPort) int {
			if c := a.Addr().Compare(b.Addr()); c != 0 {
				return c
			}
			return cmp.Compare(int16(a.Port()), int16(b.Port()))
		}, 1},
		{"Compare ignores address zones", func(a, b netip.AddrPort) int {
			aa := netip.AddrPortFrom(a.Addr().WithZone(""), a.Port())
			bb := netip.AddrPortFrom(b.Addr().WithZone(""), b.Port())
			return aa.Compare(bb)
		}, 2},
	} {
		fail := 0
		for _, c := range cases {
			if candidate.run(c.a, c.b) != c.want {
				fail++
			}
		}
		assertCount(candidate.name, fail, candidate.want, len(cases))
	}
}

func assertCount(name string, got, want, total int) {
	if got != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, got, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, got, total)
}
