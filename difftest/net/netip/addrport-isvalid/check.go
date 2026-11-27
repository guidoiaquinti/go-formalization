package main

import (
	"fmt"
	"log"
	"net/netip"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		valid     bool
		addresses []netip.Addr
	}{
		{false, []netip.Addr{{}}},
		{true, []netip.Addr{netip.IPv4Unspecified(), netip.MustParseAddr("192.0.2.1")}},
		{true, []netip.Addr{netip.IPv6Unspecified(), netip.MustParseAddr("::ffff:192.0.2.1"), netip.MustParseAddr("fe80::1%eth0")}},
	}
	const ports = 1 << 16
	rows := check.LeanRows("godiff-addrport-isvalid", ports*len(cases), nil)
	calls := 0
	for i, tc := range cases {
		for port := range ports {
			row := rows[i*ports+port]
			if row != strconv.FormatBool(tc.valid) {
				log.Fatalf("class %d port %d: Lean=%q, expected=%t", i, port, row, tc.valid)
			}
			for _, ip := range tc.addresses {
				if got := netip.AddrPortFrom(ip, uint16(port)).IsValid(); got != tc.valid {
					log.Fatalf("address %v port %d: Go=%t, expected=%t", ip, port, got, tc.valid)
				}
				calls++
			}
		}
	}
	fmt.Printf("AddrPort.IsValid: %d Lean cases and %d Go cases agree.\n", len(rows), calls)
}
