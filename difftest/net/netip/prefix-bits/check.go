package main

import (
	"fmt"
	"log"
	"math"
	"net/netip"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	rows := check.LeanRows("godiff-prefix-bits", 130, nil)
	for tag, row := range rows {
		if row != strconv.Itoa(tag-1) {
			log.Fatalf("tag %d: Lean=%q, expected=%d", tag, row, tag-1)
		}
	}
	addresses := []struct {
		ip    netip.Addr
		width int
	}{
		{},
		{netip.IPv4Unspecified(), 32},
		{netip.IPv6Unspecified(), 128},
		{netip.MustParseAddr("::ffff:192.0.2.1"), 128},
		{netip.MustParseAddr("fe80::1%eth0"), 128},
	}
	if got := (netip.Prefix{}).Bits(); got != -1 {
		log.Fatalf("zero Prefix: Go=%d, expected=-1", got)
	}
	calls := 1
	for _, addr := range addresses {
		compare := func(bits int) {
			want := -1
			if addr.width != 0 && bits >= 0 && bits <= addr.width {
				want = bits
			}
			if got := netip.PrefixFrom(addr.ip, bits).Bits(); got != want {
				log.Fatalf("address %v requested /%d: Go=%d, expected=%d", addr.ip, bits, got, want)
			}
			calls++
		}
		for bits := -1; bits <= 129; bits++ {
			compare(bits)
		}
		compare(math.MinInt)
		compare(math.MaxInt)
	}
	fmt.Printf("Prefix.Bits: %d Lean tags and %d Go cases agree.\n", len(rows), calls)
}
