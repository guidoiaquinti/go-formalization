package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/netip"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	family int
	raw    [16]byte
	old    uint32
	bits   int
}

func zone(n uint32) string {
	if n == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(n), 10)
}

func main() {
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 1}
	edges := []input{
		{family: 4, raw: mapped, bits: 32}, {family: 4, raw: mapped, bits: 33}, {family: 6, raw: mapped, old: 1, bits: 128}, {family: 6, raw: mapped, old: 1, bits: 129}, {family: 6, raw: mapped, old: 1, bits: -1},
		{bits: 0},
		{family: 4, raw: mapped, bits: 1},
		{family: 6, bits: 1},
		{family: 6, raw: mapped, old: 1, bits: 2},
		{family: 6, raw: mapped, old: 1},
		{family: 6, raw: mapped, old: 1, bits: 1},
	}
	check.CompareSamples("godiff-prefixfrom", edges, func(r *rand.Rand, i int) input {
		var c input
		c.bits = r.Intn(132) - 2
		switch i % 5 {
		case 0:
		case 1:
			c.family = 4
			c.raw = mapped
			binary.BigEndian.PutUint32(c.raw[12:], r.Uint32())
		default:
			c.family = 6
			r.Read(c.raw[:])
			c.old = r.Uint32()
			if i%5 == 2 {
				c.bits = int(r.Uint64())
			}
			if i%5 == 3 {
				c.old = 0
			}
		}
		return c
	}, func(c input) string {
		return fmt.Sprintf("%d %d %d %d %d", c.family, binary.BigEndian.Uint64(c.raw[:8]), binary.BigEndian.Uint64(c.raw[8:]), c.old, c.bits)
	}, func(c input) string {
		var ip netip.Addr
		switch c.family {
		case 4:
			ip = netip.AddrFrom4([4]byte(c.raw[12:]))
		case 6:
			ip = netip.AddrFrom16(c.raw)
			if c.old != 0 {
				ip = netip.MustParseAddr(ip.String() + "%" + zone(c.old))
			}
		}
		p := netip.PrefixFrom(ip, c.bits)
		addr := p.Addr()
		return observation(addr.BitLen(), addr.As16(), addr.Zone(), p.Bits())
	}, func(c input) string {
		width := map[int]int{0: 0, 4: 32, 6: 128}[c.family]
		bits := -1
		if width != 0 && 0 <= c.bits && c.bits <= width {
			bits = c.bits
		}
		return observation(width, c.raw, "", bits)
	})
}

func observation(width int, raw [16]byte, zone string, bits int) string {
	return fmt.Sprintf("%d %d %d [%s] %d", width, binary.BigEndian.Uint64(raw[:8]), binary.BigEndian.Uint64(raw[8:]), zone, bits)
}
