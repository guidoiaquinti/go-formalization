package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/netip"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	ip := netip.IPv4Unspecified()
	b := ip.As16()
	got := fmt.Sprintf("%d %d %t [%s]", binary.BigEndian.Uint64(b[:8]), binary.BigEndian.Uint64(b[8:]), ip.Is6(), ip.Zone())
	model := check.LeanRows("godiff-ipv4unspecified", 1, nil)[0]
	const want = "0 281470681743360 false []"
	if got != want || model != want {
		log.Fatalf("IPv4Unspecified(): Go=%q Lean=%q want=%q", got, model, want)
	}
	fmt.Println("IPv4Unspecified(): Go and Lean agree on the unzoned IPv4 address 0.0.0.0.")
}
