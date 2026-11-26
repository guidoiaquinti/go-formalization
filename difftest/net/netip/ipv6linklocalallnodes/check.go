package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/netip"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	ip := netip.IPv6LinkLocalAllNodes()
	b := ip.As16()
	got := fmt.Sprintf("%d %d %t [%s]", binary.BigEndian.Uint64(b[:8]), binary.BigEndian.Uint64(b[8:]), ip.Is6(), ip.Zone())
	model := check.LeanRows("godiff-ipv6linklocalallnodes", 1, nil)[0]
	want := fmt.Sprintf("%d 1 true []", uint64(0xff02000000000000))
	if got != want || model != want {
		log.Fatalf("IPv6LinkLocalAllNodes(): Go=%q Lean=%q want=%q", got, model, want)
	}
	fmt.Println("IPv6LinkLocalAllNodes(): Go and Lean agree on the unzoned IPv6 address ff02::1.")
}
