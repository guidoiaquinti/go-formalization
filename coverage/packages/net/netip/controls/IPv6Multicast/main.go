package main

import (
	"fmt"
	"log"
	"net/netip"
)

func matches(ip netip.Addr, want [16]byte) bool {
	return ip.As16() == want && ip.Zone() == "" && ip.Is6() && ip.IsLinkLocalMulticast()
}

func main() {
	nodes := netip.IPv6LinkLocalAllNodes()
	routers := netip.IPv6LinkLocalAllRouters()
	cases := []struct {
		label        string
		got, swapped netip.Addr
		last         byte
	}{
		{"IPv6LinkLocalAllNodes", nodes, routers, 1},
		{"IPv6LinkLocalAllRouters", routers, nodes, 2},
	}
	var swaps, scopes int
	for _, c := range cases {
		want := [16]byte{0: 255, 1: 2, 15: c.last}
		if !matches(c.got, want) {
			log.Fatalf("%s: unexpected address %s (%x)", c.label, c.got, c.got.As16())
		}
		if !matches(c.swapped, want) {
			swaps++
		}
		wrongScope := want
		wrongScope[1] = 1
		if !matches(netip.AddrFrom16(wrongScope), want) {
			scopes++
		}
	}
	if swaps != 2 || scopes != 2 {
		log.Fatalf("mutation failures: swapped endpoint=%d, wrong scope=%d", swaps, scopes)
	}
	fmt.Printf("IPv6 link-local multicast: %d fixtures pass, mutant failures swapped endpoint=%d wrong scope=%d\n",
		len(cases), swaps, scopes)
}
