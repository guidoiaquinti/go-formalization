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
	edges := [][4]byte{{}, {3: 1}, {192, 0, 2, 1}}
	var ones [4]byte
	for i := range ones {
		ones[i] = 255
		for bit := range 8 {
			var single [4]byte
			single[i] = 1 << bit
			edges = append(edges, single)
		}
	}
	edges = append(edges, ones)
	check.CompareSamples("godiff-addrfrom4", edges, func(r *rand.Rand, _ int) [4]byte {
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], r.Uint32())
		return b
	}, func(b [4]byte) string {
		var fields [4]string
		for i, v := range b {
			fields[i] = strconv.Itoa(int(v))
		}
		return strings.Join(fields[:], " ")
	}, func(b [4]byte) string {
		ip := netip.AddrFrom4(b)
		actual := ip.As16()
		return fmt.Sprintf("%d %d %t [%s]", binary.BigEndian.Uint64(actual[:8]), binary.BigEndian.Uint64(actual[8:]), ip.Is6(), ip.Zone())
	}, func(b [4]byte) string {
		return fmt.Sprintf("0 %d false []", uint64(0xffff00000000)+uint64(binary.BigEndian.Uint32(b[:])))
	})
}
