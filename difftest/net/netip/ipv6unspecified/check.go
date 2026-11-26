package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/netip"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	ip := netip.IPv6Unspecified()
	b := ip.As16()
	got := fmt.Sprintf("%d %d %t [%s]", binary.BigEndian.Uint64(b[:8]), binary.BigEndian.Uint64(b[8:]), ip.Is6(), ip.Zone())
	model := check.LeanRows("godiff-ipv6unspecified", 1, nil)[0]
	const want = "0 0 true []"
	if got != want || model != want {
		log.Fatalf("IPv6Unspecified(): Go=%q Lean=%q want=%q", got, model, want)
	}
	fmt.Println("IPv6Unspecified(): Go and Lean agree on the unzoned IPv6 address ::.")
}
