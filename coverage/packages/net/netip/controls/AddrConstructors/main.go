package main

import (
	"fmt"
	"log"
	"net/netip"
)

type outcome struct {
	octets            [16]byte
	valid, ipv4, ipv6 bool
	zone              string
}

func observe(addr netip.Addr) outcome {
	return outcome{addr.As16(), addr.IsValid(), addr.Is4(), addr.Is6(), addr.Zone()}
}

func main() {
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 17}
	native := [16]byte{0x20, 1, 0x0d, 0xb8, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	swapped := [16]byte{8, 9, 10, 11, 12, 13, 14, 15, 0x20, 1, 0x0d, 0xb8, 4, 5, 6, 7}
	var zero netip.Addr
	v4 := netip.AddrFrom4([4]byte{192, 0, 2, 17})
	mapped6 := netip.AddrFrom16(mapped)
	cases := []struct {
		label  string
		got    netip.Addr
		mutant netip.Addr
		want   outcome
	}{
		{"AddrFrom4", v4, mapped6, outcome{mapped, true, true, false, ""}},
		{"AddrFrom16 mapped", mapped6, mapped6.Unmap(), outcome{mapped, true, false, true, ""}},
		{"AddrFrom16 native", netip.AddrFrom16(native), netip.AddrFrom16(swapped), outcome{native, true, false, true, ""}},
		{"IPv4Unspecified", netip.IPv4Unspecified(), zero, outcome{[16]byte{10: 255, 11: 255}, true, true, false, ""}},
		{"IPv6Unspecified", netip.IPv6Unspecified(), zero, outcome{[16]byte{}, true, false, true, ""}},
		{"IPv6Loopback", netip.IPv6Loopback(), netip.AddrFrom16([16]byte{0: 1}), outcome{[16]byte{15: 1}, true, false, true, ""}},
	}
	count := 0
	for _, c := range cases {
		if got := observe(c.got); got != c.want {
			log.Fatalf("%s: got %+v, want %+v", c.label, got, c.want)
		}
		if observe(c.mutant) != c.want {
			count++
		}
	}
	if count != len(cases) {
		log.Fatalf("mutation failures=%d, want %d", count, len(cases))
	}
	fmt.Printf("address constructors: %d fixtures pass and reject %d mutations\n", len(cases), count)
}
