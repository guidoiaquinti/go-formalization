package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/netip"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := [][16]byte{{}, {15: 1}, {10: 255, 11: 255, 12: 192, 13: 0, 14: 2, 15: 1}}
	var ones [16]byte
	for i := range ones {
		ones[i] = 255
		for bit := range 8 {
			var single [16]byte
			single[i] = 1 << bit
			edges = append(edges, single)
		}
	}
	edges = append(edges, ones)
	check.CompareSamples("godiff-addrfrom16", edges, func(r *rand.Rand, _ int) [16]byte {
		var b [16]byte
		binary.BigEndian.PutUint64(b[:8], r.Uint64())
		binary.BigEndian.PutUint64(b[8:], r.Uint64())
		return b
	}, func(b [16]byte) string {
		var fields [16]string
		for i, v := range b {
			fields[i] = strconv.Itoa(int(v))
		}
		return strings.Join(fields[:], " ")
	}, func(b [16]byte) string {
		ip := netip.AddrFrom16(b)
		actual := ip.As16()
		return fmt.Sprintf("%d %d %t [%s]", binary.BigEndian.Uint64(actual[:8]), binary.BigEndian.Uint64(actual[8:]), ip.Is6(), ip.Zone())
	}, func(b [16]byte) string {
		return fmt.Sprintf("%d %d true []", binary.BigEndian.Uint64(b[:8]), binary.BigEndian.Uint64(b[8:]))
	})
}
