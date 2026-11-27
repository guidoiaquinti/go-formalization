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
	mapped := [16]byte{10: 255, 11: 255, 12: 169, 13: 254, 15: 1}
	edges := []receiver{{}}
	for octet := range 256 {
		for pos := 12; pos <= 13; pos++ {
			v4 := mapped
			v4[pos] = byte(octet)
			edges = append(edges, receiver{4, v4}, receiver{6, v4})
		}
		edges = append(edges, receiver{6, [16]byte{0: 254, 1: byte(octet), 15: 1}},
			receiver{6, [16]byte{0: byte(octet), 1: 128, 15: 1}})
	}
	check.CompareSamples("godiff-addr-islinklocalunicast", edges, func(r *rand.Rand, i int) receiver {
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
		return strconv.FormatBool(ip.IsLinkLocalUnicast())
	}, func(c receiver) string {
		ipv4 := c.family == 4 || (c.family == 6 && bytes.Equal(c.raw[:12], prefix[:]))
		return strconv.FormatBool((ipv4 && c.raw[12] == 169 && c.raw[13] == 254) || (c.family == 6 && c.raw[0] == 254 && c.raw[1] >= 128 && c.raw[1] < 192))
	})
}
