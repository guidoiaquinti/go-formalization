package main

import (
	"fmt"
	"log"
	"net/netip"
)

type fixture struct {
	input string
	want  string
}

func main() {
	next := []fixture{
		{"", "invalid-plus-one"},
		{"255.255.255.254", "255.255.255.255"},
		{"255.255.255.255", ""},
		{"::ffff:255.255.255.255%eth0", "::1:0:0:0%eth0"},
		{"2001:db8::ffff:ffff:ffff:ffff%rack", "2001:db8:0:1::%rack"},
		{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff%zone", ""},
	}
	prev := []fixture{
		{"", ""},
		{"0.0.0.0", ""},
		{"0.0.0.1", "0.0.0.0"},
		{"::ffff:0.0.0.0%eth0", "::fffe:ffff:ffff%eth0"},
		{"2001:db8:0:1::%rack", "2001:db8:0:0:ffff:ffff:ffff:ffff%rack"},
		{"::%zone", ""},
	}
	check("Next", netip.Addr.Next, next, 0)
	check("Next ignores invalid input", func(a netip.Addr) netip.Addr {
		if !a.IsValid() {
			return a
		}
		return a.Next()
	}, next, 1)
	check("Next loses carry into the high limb", dropHighChange(netip.Addr.Next), next, 1)
	check("Next drops the zone", dropZone(netip.Addr.Next), next, 2)
	check("Prev", netip.Addr.Prev, prev, 0)
	check("Prev loses borrow from the high limb", dropHighChange(netip.Addr.Prev), prev, 1)
	check("Prev drops the zone", dropZone(netip.Addr.Prev), prev, 2)
	check("Prev treats mapped addresses as IPv4", func(a netip.Addr) netip.Addr {
		return a.Unmap().Prev()
	}, prev, 1)
	fmt.Printf("zero Addr after Next: BitLen=%d As16=%x equalZero=%t\n",
		(netip.Addr{}).Next().BitLen(), (netip.Addr{}).Next().As16(), (netip.Addr{}).Next() == (netip.Addr{}))
}

func check(name string, run func(netip.Addr) netip.Addr, tests []fixture, wantCount int) {
	failures := 0
	for _, test := range tests {
		actual := run(addr(test.input))
		var good bool
		if test.want == "invalid-plus-one" {
			good = !actual.IsValid() && actual.BitLen() == 0 && actual.Zone() == "" &&
				actual.As16() == ([16]byte{15: 1}) && actual != (netip.Addr{})
		} else {
			want := addr(test.want)
			good = actual == want && actual.As16() == want.As16() && actual.BitLen() == want.BitLen() && actual.Zone() == want.Zone()
		}
		if !good {
			failures++
			if wantCount == 0 {
				log.Fatalf("%s(%q) = (%d,%x,%q), want %q", name, test.input, actual.BitLen(), actual.As16(), actual.Zone(), test.want)
			}
		}
	}
	if failures != wantCount {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, failures, wantCount)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, failures, len(tests))
}

func addr(text string) netip.Addr {
	if text == "" {
		return netip.Addr{}
	}
	return netip.MustParseAddr(text)
}

func dropHighChange(next func(netip.Addr) netip.Addr) func(netip.Addr) netip.Addr {
	return func(a netip.Addr) netip.Addr {
		got := next(a)
		if !a.Is6() || !got.Is6() {
			return got
		}
		original, changed := a.As16(), got.As16()
		copy(changed[:8], original[:8])
		return netip.AddrFrom16(changed).WithZone(got.Zone())
	}
}

func dropZone(next func(netip.Addr) netip.Addr) func(netip.Addr) netip.Addr {
	return func(a netip.Addr) netip.Addr { return next(a).WithZone("") }
}
