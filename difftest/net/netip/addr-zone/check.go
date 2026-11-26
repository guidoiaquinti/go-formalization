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

func (c receiver) zoneString() string {
	if c.zone == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(c.zone), 10)
}

func main() {
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 1}
	edges := []receiver{{}, {4, mapped, 0}, {family: 6}, {6, mapped, 1}, {6, [16]byte{0: 0xfe, 1: 0x80, 15: 1}, 2}}
	check.CompareSamples("godiff-addr-zone", edges, func(r *rand.Rand, i int) receiver {
		var c receiver
		switch i % 4 {
		case 0:
		case 1:
			c.family = 4
			c.raw = mapped
			binary.BigEndian.PutUint32(c.raw[12:], r.Uint32())
		default:
			c.family = 6
			r.Read(c.raw[:])
			if i%4 == 3 {
				c.zone = r.Uint32()
			}
		}
		return c
	}, func(c receiver) string {
		return fmt.Sprintf("%d %d %d %d", c.family, binary.BigEndian.Uint64(c.raw[:8]), binary.BigEndian.Uint64(c.raw[8:]), c.zone)
	}, func(c receiver) string {
		var ip netip.Addr
		switch c.family {
		case 4:
			ip = netip.AddrFrom4([4]byte(c.raw[12:]))
		case 6:
			ip = netip.AddrFrom16(c.raw).WithZone(c.zoneString())
		}
		return "[" + ip.Zone() + "]"
	}, func(c receiver) string {
		return "[" + c.zoneString() + "]"
	})
}
