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
	mapped := [16]byte{10: 255, 11: 255, 12: 127, 14: 2, 15: 1}
	edges := []receiver{{}, {4, mapped}, {6, mapped}, {family: 6}, {6, [16]byte{15: 1}}}
	for _, head := range []byte{126, 128} {
		near := mapped
		near[12] = head
		edges = append(edges, receiver{4, near}, receiver{6, near})
	}
	for i := range 16 {
		for bit := range 8 {
			near := [16]byte{15: 1}
			near[i] ^= 1 << bit
			edges = append(edges, receiver{6, near})
		}
	}
	for i := range prefix {
		for bit := range 8 {
			near := mapped
			near[i] ^= 1 << bit
			edges = append(edges, receiver{6, near})
		}
	}
	check.CompareSamples("godiff-addr-isloopback", edges, func(r *rand.Rand, i int) receiver {
		var raw [16]byte
		r.Read(raw[:])
		switch i % 4 {
		case 0:
			copy(raw[:], prefix[:])
			raw[12] = 127
		case 1:
			copy(raw[:], prefix[:])
			raw[12] = 127
			raw[r.Intn(12)] ^= 1 << r.Intn(8)
		case 2:
			copy(raw[:], prefix[:])
			raw[12] = 127
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
		return strconv.FormatBool(ip.IsLoopback())
	}, func(c receiver) string {
		ipv4 := c.family == 4 || (c.family == 6 && bytes.Equal(c.raw[:12], prefix[:]))
		return strconv.FormatBool((ipv4 && c.raw[12] == 127) || (c.family == 6 && c.raw == [16]byte{15: 1}))
	})
}
