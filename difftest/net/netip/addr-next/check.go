package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/netip"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type receiver struct {
	family int
	raw    [16]byte
	zone   uint32
}

func main() {
	prefix := [12]byte{10: 255, 11: 255}
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 1}
	edges := []receiver{{}, {4, mapped, 0}, {6, mapped, 1}, {family: 6}}
	for _, text := range []string{"255.255.255.254", "255.255.255.255", "::ffff:255.255.255.255", "::ffff:0.0.0.0", "::ffff:ffff:ffff:ffff", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"} {
		ip := netip.MustParseAddr(text)
		if ip.Is4() {
			edges = append(edges, receiver{4, ip.As16(), 0})
		} else {
			edges = append(edges, receiver{6, ip.As16(), 0}, receiver{6, ip.As16(), 1})
		}
	}
	check.CompareSamples("godiff-addr-next", edges, func(r *rand.Rand, i int) receiver {
		var raw [16]byte
		r.Read(raw[:])
		switch i % 4 {
		case 0:
			copy(raw[:], prefix[:])
		case 1:
			copy(raw[:], prefix[:])
			raw[r.Intn(12)] ^= 1 << r.Intn(8)
		case 2:
			copy(raw[:], prefix[:])
			return receiver{4, raw, 0}
		}
		return receiver{6, raw, r.Uint32()}
	}, func(c receiver) string {
		return fmt.Sprintf("%d %d %d %d", c.family, binary.BigEndian.Uint64(c.raw[:8]), binary.BigEndian.Uint64(c.raw[8:]), c.zone)
	}, func(c receiver) string {
		var ip netip.Addr
		switch c.family {
		case 4:
			ip = netip.AddrFrom4([4]byte(c.raw[12:]))
		case 6:
			ip = netip.AddrFrom16(c.raw)
		}
		if c.zone != 0 {
			ip = ip.WithZone(strconv.FormatUint(uint64(c.zone), 10))
		}
		got := ip.Next()
		return observation(got.BitLen(), got.As16(), got.Zone())
	}, func(c receiver) string {
		raw := c.raw
		carry := byte(1)
		start := 0
		if c.family == 4 {
			start = 12
		}
		for i := 15; i >= start && carry != 0; i-- {
			raw[i]++
			if raw[i] != 0 {
				carry = 0
			}
		}
		if carry != 0 {
			return observation(0, [16]byte{}, "")
		}
		width := map[int]int{0: 0, 4: 32, 6: 128}[c.family]
		zone := ""
		if c.zone != 0 {
			zone = strconv.FormatUint(uint64(c.zone), 10)
		}
		return observation(width, raw, zone)
	})
}

func observation(width int, raw [16]byte, zone string) string {
	return fmt.Sprintf("%d %d %d [%s]", width, binary.BigEndian.Uint64(raw[:8]), binary.BigEndian.Uint64(raw[8:]), zone)
}
