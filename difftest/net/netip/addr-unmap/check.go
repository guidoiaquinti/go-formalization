package main

import (
	"bytes"
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
	edges := []receiver{{}, {4, mapped, 0}, {6, mapped, 0}, {6, mapped, 1}, {family: 6}}
	for i := range prefix {
		for bit := range 8 {
			near := mapped
			near[i] ^= 1 << bit
			edges = append(edges, receiver{6, near, 1})
		}
	}
	check.CompareSamples("godiff-addr-unmap", edges, func(r *rand.Rand, i int) receiver {
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
		got := ip.Unmap()
		return observation(got.BitLen(), got.As16(), got.Zone())
	}, func(c receiver) string {
		switch c.family {
		case 0:
			return observation(0, c.raw, "")
		case 4:
			return observation(32, c.raw, "")
		}
		if bytes.Equal(c.raw[:12], prefix[:]) {
			return observation(32, c.raw, "")
		}
		zone := ""
		if c.zone != 0 {
			zone = strconv.FormatUint(uint64(c.zone), 10)
		}
		return observation(128, c.raw, zone)
	})
}

func observation(width int, raw [16]byte, zone string) string {
	return fmt.Sprintf("%d %d %d [%s]", width, binary.BigEndian.Uint64(raw[:8]), binary.BigEndian.Uint64(raw[8:]), zone)
}
