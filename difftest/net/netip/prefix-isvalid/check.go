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
	rows := check.LeanRows("godiff-prefix-isvalid", 130, nil)
	for tag, row := range rows {
		if row != strconv.FormatBool(tag > 0) {
			log.Fatalf("tag %d: Lean=%q, expected=%t", tag, row, tag > 0)
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
	if got := (netip.Prefix{}).IsValid(); got {
		log.Fatalf("zero Prefix: Go=%t, expected=false", got)
	}
	calls := 1
	for _, addr := range addresses {
		compare := func(bits int) {
			want := false
			if addr.width != 0 && bits >= 0 && bits <= addr.width {
				want = true
			}
			if got := netip.PrefixFrom(addr.ip, bits).IsValid(); got != want {
				log.Fatalf("address %v requested /%d: Go=%t, expected=%t", addr.ip, bits, got, want)
			}
			calls++
		}
		for bits := -1; bits <= 129; bits++ {
			compare(bits)
		}
		compare(math.MinInt)
		compare(math.MaxInt)
	}
	fmt.Printf("Prefix.IsValid: %d Lean tags and %d Go cases agree.\n", len(rows), calls)
}
