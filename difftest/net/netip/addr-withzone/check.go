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
	family         int
	raw            [16]byte
	old, requested uint32
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
		{requested: 1},
		{family: 4, raw: mapped, requested: 1},
		{family: 6, requested: 1},
		{family: 6, raw: mapped, old: 1, requested: 2},
		{family: 6, raw: mapped, old: 1},
		{family: 6, raw: mapped, old: 1, requested: 1},
	}
	check.CompareSamples("godiff-addr-withzone", edges, func(r *rand.Rand, i int) input {
		var c input
		c.requested = r.Uint32()
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
				c.requested = 0
			}
			if i%5 == 3 {
				c.old = 0
			}
		}
		return c
	}, func(c input) string {
		return fmt.Sprintf("%d %d %d %d %d", c.family, binary.BigEndian.Uint64(c.raw[:8]), binary.BigEndian.Uint64(c.raw[8:]), c.old, c.requested)
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
		got := ip.WithZone(zone(c.requested))
		return observation(got.BitLen(), got.As16(), got.Zone())
	}, func(c input) string {
		switch c.family {
		case 4:
			return observation(32, c.raw, "")
		case 6:
			return observation(128, c.raw, zone(c.requested))
		default:
			return observation(0, c.raw, "")
		}
	})
}

func observation(width int, raw [16]byte, zone string) string {
	return fmt.Sprintf("%d %d %d [%s]", width, binary.BigEndian.Uint64(raw[:8]), binary.BigEndian.Uint64(raw[8:]), zone)
}
