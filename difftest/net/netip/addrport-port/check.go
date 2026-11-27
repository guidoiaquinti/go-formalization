package main

import (
	"fmt"
	"log"
	"net/netip"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	addresses := []netip.Addr{
		{},
		netip.MustParseAddr("192.0.2.1"),
		netip.IPv6Loopback(),
		netip.MustParseAddr("::ffff:192.0.2.1"),
		netip.MustParseAddr("fe80::1%eth0"),
	}
	if got := (netip.AddrPort{}).Port(); got != 0 {
		log.Fatalf("zero AddrPort: Go=%d, expected=0", got)
	}
	const count = 1 << 16
	for port, row := range check.LeanRows("godiff-addrport-port", count, nil) {
		if want := strconv.Itoa(port); row != want {
			log.Fatalf("port %d: Lean=%q, expected=%q", port, row, want)
		}
		for _, ip := range addresses {
			got := netip.AddrPortFrom(ip, uint16(port)).Port()
			if got != uint16(port) {
				log.Fatalf("address %v, port %d: Go=%d, expected=%d", ip, port, got, port)
			}
		}
	}
	fmt.Printf("AddrPort.Port: %d Lean ports and %d Go receivers agree with the specification.\n", count, count*len(addresses)+1)
}
