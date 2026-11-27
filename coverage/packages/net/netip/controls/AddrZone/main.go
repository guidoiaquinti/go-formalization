package main

import (
	"fmt"
	"log"
	"net/netip"
)

func main() {
	mapped := netip.MustParseAddr("::ffff:192.0.2.8%rack:α")
	v6 := netip.MustParseAddr("fe80::1%edge 2")
	v4 := netip.MustParseAddr("192.0.2.8")
	zero := netip.Addr{}
	zoneCases := []struct {
		input netip.Addr
		want  string
	}{
		{mapped, "rack:α"}, {v6, "edge 2"}, {v4, ""}, {zero, ""}, {netip.IPv6Unspecified(), ""},
	}
	count := func(run func(netip.Addr) string) int {
		fail := 0
		for _, c := range zoneCases {
			if run(c.input) != c.want {
				fail++
			}
		}
		return fail
	}
	assertCount("Zone", count(netip.Addr.Zone), 0, len(zoneCases))
	assertCount("Zone drops mapped address zone", count(func(a netip.Addr) string { return a.Unmap().Zone() }), 1, len(zoneCases))
	assertCount("Zone always returns empty", count(func(netip.Addr) string { return "" }), 2, len(zoneCases))

	byteZone := "raw\xff%zone"
	withZoneCases := []struct {
		input   netip.Addr
		request string
		want    string
	}{
		{mapped, "", ""},
		{mapped, byteZone, byteZone},
		{v6, "new:é", "new:é"},
		{v6, "edge 2", "edge 2"},
		{v4, "ignored", ""},
		{zero, "ignored", ""},
	}
	checkZone := func(run func(netip.Addr, string) netip.Addr) int {
		fail := 0
		for _, c := range withZoneCases {
			got := run(c.input, c.request)
			if got.BitLen() != c.input.BitLen() || got.As16() != c.input.As16() || got.Zone() != c.want {
				fail++
			}
		}
		return fail
	}
	assertCount("WithZone", checkZone(netip.Addr.WithZone), 0, len(withZoneCases))
	assertCount("WithZone keeps the old zone on an empty request", checkZone(func(a netip.Addr, z string) netip.Addr {
		if z == "" {
			return a
		}
		return a.WithZone(z)
	}), 1, len(withZoneCases))
	assertCount("WithZone converts mapped IPv6 to IPv4", checkZone(func(a netip.Addr, z string) netip.Addr {
		return a.Unmap().WithZone(z)
	}), 2, len(withZoneCases))

	near := netip.MustParseAddr("::fffe:c000:208%keep")
	unmapCases := []struct{ input, want netip.Addr }{
		{mapped, v4}, {mapped.WithZone(""), v4}, {near, near}, {v6, v6}, {v4, v4}, {zero, zero},
	}
	checkUnmap := func(run func(netip.Addr) netip.Addr) int {
		fail := 0
		for _, c := range unmapCases {
			got := run(c.input)
			if got.BitLen() != c.want.BitLen() || got.As16() != c.want.As16() || got.Zone() != c.want.Zone() {
				fail++
			}
		}
		return fail
	}
	assertCount("Unmap", checkUnmap(netip.Addr.Unmap), 0, len(unmapCases))
	assertCount("Unmap leaves all addresses unchanged", checkUnmap(func(a netip.Addr) netip.Addr { return a }), 2, len(unmapCases))
	assertCount("Unmap ignores the final mapped-prefix bit", checkUnmap(func(a netip.Addr) netip.Addr {
		raw := a.As16()
		if a.Is6() && raw[10] == 255 && raw[11]&0xfe == 0xfe {
			return netip.AddrFrom4([4]byte(raw[12:]))
		}
		return a.Unmap()
	}), 1, len(unmapCases))
}

func assertCount(name string, got, want, total int) {
	if got != want {
		log.Fatalf("%s disagreed on %d fixtures, want %d", name, got, want)
	}
	fmt.Printf("%s disagreed on %d of %d fixtures\n", name, got, total)
}
