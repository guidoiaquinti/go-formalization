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
}

func main() {
	prefix := [12]byte{10: 255, 11: 255}
	mapped := [16]byte{10: 255, 11: 255, 12: 172, 13: 16, 15: 1}
	edges := []receiver{{}}
	for _, pair := range [][2]byte{{10, 0}, {172, 16}, {172, 31}, {192, 168}, {9, 255}, {11, 0}, {172, 15}, {172, 32}, {192, 167}, {192, 169}} {
		raw := mapped
		copy(raw[12:14], pair[:])
		edges = append(edges, receiver{4, raw}, receiver{6, raw})
	}
	for octet := range 256 {
		edges = append(edges, receiver{6, [16]byte{0: byte(octet), 15: 1}})
	}
	check.CompareSamples("godiff-addr-isprivate", edges, func(r *rand.Rand, i int) receiver {
		var raw [16]byte
		r.Read(raw[:])
		switch i % 4 {
		case 0:
			copy(raw[:], prefix[:])
			copy(raw[12:14], mapped[12:14])
		case 1:
			copy(raw[:], prefix[:])
			copy(raw[12:14], mapped[12:14])
			raw[r.Intn(12)] ^= 1 << r.Intn(8)
		case 2:
			copy(raw[:], prefix[:])
			copy(raw[12:14], mapped[12:14])
			return receiver{4, raw}
		}
		return receiver{6, raw}
	}, func(c receiver) string {
		return fmt.Sprintf("%d %d %d 0", c.family, binary.BigEndian.Uint64(c.raw[:8]), binary.BigEndian.Uint64(c.raw[8:]))
	}, func(c receiver) string {
		var ip netip.Addr
		switch c.family {
		case 4:
			ip = netip.AddrFrom4([4]byte(c.raw[12:]))
		case 6:
			ip = netip.AddrFrom16(c.raw)
		}
		return strconv.FormatBool(ip.IsPrivate())
	}, func(c receiver) string {
		ipv4 := c.family == 4 || (c.family == 6 && bytes.Equal(c.raw[:12], prefix[:]))
		a, b := c.raw[12], c.raw[13]
		return strconv.FormatBool((ipv4 && (a == 10 || a == 172 && b >= 16 && b <= 31 || a == 192 && b == 168)) ||
			(c.family == 6 && (c.raw[0] == 252 || c.raw[0] == 253)))
	})
}
