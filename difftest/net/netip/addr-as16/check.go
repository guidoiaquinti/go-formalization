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

type receiver struct {
	family int
	raw    [16]byte
	scoped bool
}

func main() {
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 1}
	edges := []receiver{{}, {4, mapped, false}, {6, mapped, false}, {6, mapped, true}, {family: 6}}
	for i := range 16 {
		for bit := range 8 {
			var single [16]byte
			single[i] = 1 << bit
			edges = append(edges, receiver{6, single, false})
		}
	}
	check.CompareSamples("godiff-addr-as16", edges, func(r *rand.Rand, i int) receiver {
		var raw [16]byte
		r.Read(raw[:])
		if i%3 == 0 {
			copy(raw[:12], mapped[:12])
			return receiver{4, raw, false}
		}
		return receiver{6, raw, i%3 == 1}
	}, func(c receiver) string {
		return fmt.Sprintf("%d %d", binary.BigEndian.Uint64(c.raw[:8]), binary.BigEndian.Uint64(c.raw[8:]))
	}, func(c receiver) string {
		var ip netip.Addr
		switch c.family {
		case 4:
			ip = netip.AddrFrom4([4]byte(c.raw[12:]))
		case 6:
			ip = netip.AddrFrom16(c.raw)
		}
		if c.scoped {
			ip = ip.WithZone("eth0")
		}
		return formatBytes(ip.As16())
	}, func(c receiver) string {
		return formatBytes(c.raw)
	})
}

func formatBytes(raw [16]byte) string {
	var fields [16]string
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields[:], " ")
}
