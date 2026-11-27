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
	mapped := [16]byte{10: 255, 11: 255}
	edges := []receiver{{}, {4, mapped, 0}, {family: 6}, {family: 6, zone: 1}, {6, mapped, 0}}
	for i := range 16 {
		for bit := range 8 {
			var raw [16]byte
			raw[i] = 1 << bit
			edges = append(edges, receiver{6, raw, 0})
		}
	}
	check.CompareSamples("godiff-addr-isunspecified", edges, func(r *rand.Rand, i int) receiver {
		var c receiver
		switch i % 6 {
		case 0:
		case 1:
			c.family = 4
			c.raw = mapped
		case 2:
			c.family = 4
			c.raw = mapped
			binary.BigEndian.PutUint32(c.raw[12:], r.Uint32())
		case 3:
			c.family = 6
		default:
			c.family = 6
			r.Read(c.raw[:])
			if i%6 == 5 {
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
		return strconv.FormatBool(ip.IsUnspecified())
	}, func(c receiver) string {
		return strconv.FormatBool((c.family == 4 && c.raw == mapped) || (c.family == 6 && c.raw == [16]byte{} && c.zone == 0))
	})
}
