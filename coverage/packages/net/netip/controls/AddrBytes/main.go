package main

import (
	"bytes"
	"fmt"
	"log"
	"net/netip"
	"slices"
)

type as4Result struct {
	octets     [4]byte
	panicValue any
}

func as4(ip netip.Addr) (result as4Result) {
	defer func() { result.panicValue = recover() }()
	result.octets = ip.As4()
	return result
}

func main() {
	v4 := [4]byte{192, 0, 2, 17}
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 17}
	native6 := [16]byte{0x20, 1, 0x0d, 0xb8, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	nearMapped := mapped
	nearMapped[7] = 1
	zeroIPv4 := [16]byte{10: 255, 11: 255}
	zero := [16]byte{}
	ok := as4Result{octets: v4}
	ipv6Panic := as4Result{panicValue: "As4 called on IPv6 address"}
	cases := []struct {
		label     string
		addr      netip.Addr
		want4     as4Result
		want16    [16]byte
		wantSlice []byte
	}{
		{"invalid", netip.Addr{}, as4Result{panicValue: "As4 called on IP zero value"}, zero, nil},
		{"IPv4 unspecified", netip.AddrFrom4([4]byte{}), as4Result{}, zeroIPv4, []byte{0, 0, 0, 0}},
		{"IPv4", netip.AddrFrom4(v4), ok, mapped, v4[:]},
		{"mapped", netip.AddrFrom16(mapped), ok, mapped, mapped[:]},
		{"scoped mapped", netip.AddrFrom16(mapped).WithZone("eth0"), ok, mapped, mapped[:]},
		{"IPv6 unspecified", netip.AddrFrom16(zero), ipv6Panic, zero, zero[:]},
		{"scoped native IPv6", netip.AddrFrom16(native6).WithZone("eth0"), ipv6Panic, native6, native6[:]},
		{"near-mapped", netip.AddrFrom16(nearMapped), ipv6Panic, nearMapped, nearMapped[:]},
	}
	var mappedRejected, reversedWords, unmappedSlice int
	for _, c := range cases {
		got4, got16, gotSlice := as4(c.addr), c.addr.As16(), c.addr.AsSlice()
		if got4 != c.want4 || got16 != c.want16 ||
			(gotSlice == nil) != (c.wantSlice == nil) || !bytes.Equal(gotSlice, c.wantSlice) {
			log.Fatalf("%s: As4=%+v, As16=%x, AsSlice=%#v", c.label, got4, got16, gotSlice)
		}

		// Mutant: accept only native IPv4 while preserving invalid-address panic.
		if c.addr.Is6() && ipv6Panic != c.want4 {
			mappedRejected++
		}
		// Mutant: write each 64-bit word in little-endian order.
		slices.Reverse(got16[:8])
		slices.Reverse(got16[8:])
		if got16 != c.want16 {
			reversedWords++
		}
		// Mutant: normalize mapped IPv6 to IPv4 before selecting the slice width.
		unmapped := c.addr.Unmap().AsSlice()
		if (unmapped == nil) != (c.wantSlice == nil) || !bytes.Equal(unmapped, c.wantSlice) {
			unmappedSlice++
		}
	}
	if mappedRejected != 2 || reversedWords != 6 || unmappedSlice != 2 {
		log.Fatalf("unexpected mutation failure counts: As4=%d As16=%d AsSlice=%d", mappedRejected, reversedWords, unmappedSlice)
	}
	fmt.Printf("Addr byte conversions: %d fixtures pass, mutant failures As4=%d As16=%d AsSlice=%d\n",
		len(cases), mappedRejected, reversedWords, unmappedSlice)
}
