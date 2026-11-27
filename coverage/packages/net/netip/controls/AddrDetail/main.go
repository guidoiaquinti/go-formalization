package main

import (
	"fmt"
	"net/netip"
)

func main() {
	cases := []struct {
		addr netip.Addr
		bits int
	}{
		{netip.Addr{}, 0},
		{netip.IPv4Unspecified(), 32},
		{netip.IPv6Unspecified(), 128},
		{netip.MustParseAddr("::ffff:192.0.2.1"), 128},
		{netip.MustParseAddr("fe80::1%eth0"), 128},
		{netip.MustParseAddr("192.0.2.1"), 32},
	}
	var real, unmapped, unspecified int
	for _, c := range cases {
		valid := c.bits != 0
		if c.addr.BitLen() != c.bits || c.addr.IsValid() != valid {
			real++
		}
		if c.addr.Unmap().BitLen() != c.bits {
			unmapped++
		}
		if (c.addr.IsValid() && !c.addr.IsUnspecified()) != valid {
			unspecified++
		}
	}
	fmt.Printf("Addr detail: Go=%d/%d unmap-first=%d/%d reject-unspecified=%d/%d\n",
		real, len(cases), unmapped, len(cases), unspecified, len(cases))
	if real != 0 || unmapped != 1 || unspecified != 2 {
		panic("unexpected Addr detail control counts")
	}
}
