package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net/netip"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	ip    netip.Addr
	width int
	bits  int
}

func main() {
	addresses := []input{
		{},
		{ip: netip.IPv4Unspecified(), width: 32},
		{ip: netip.IPv6Unspecified(), width: 128},
		{ip: netip.MustParseAddr("::ffff:192.0.2.1"), width: 128},
		{ip: netip.MustParseAddr("fe80::1%eth0"), width: 128},
	}

	var cases []input
	for _, c := range addresses {
		for bits := -1; bits <= 129; bits++ {
			c.bits = bits
			cases = append(cases, c)
		}
		c.bits = math.MinInt
		cases = append(cases, c)
		c.bits = math.MaxInt
		cases = append(cases, c)
	}
	var stdin bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&stdin, c.width, c.bits)
	}
	if (netip.Prefix{}).IsSingleIP() {
		log.Fatal("zero Prefix reports a single IP")
	}
	for i, row := range check.LeanRows("godiff-prefix-issingleip", len(cases), &stdin) {
		c := cases[i]
		want := c.width != 0 && c.bits == c.width
		if got := netip.PrefixFrom(c.ip, c.bits).IsSingleIP(); got != want || row != strconv.FormatBool(want) {
			log.Fatalf("address %v requested /%d: Go=%t, Lean=%q, expected=%t", c.ip, c.bits, got, row, want)
		}
	}
	fmt.Printf("Prefix.IsSingleIP: %d paired cases and the Go zero Prefix agree.\n", len(cases))
}
