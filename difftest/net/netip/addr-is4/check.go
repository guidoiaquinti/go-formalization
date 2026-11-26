package main

import (
	"encoding/binary"
	"math/rand"
	"net/netip"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type receiver struct {
	addr   netip.Addr
	family int
}

func main() {
	edges := []receiver{
		{},
		{netip.IPv4Unspecified(), 4},
		{netip.IPv6Unspecified(), 6},
		{netip.MustParseAddr("::ffff:192.0.2.1"), 6},
		{netip.MustParseAddr("fe80::1%eth0"), 6},
	}
	check.CompareSamples("godiff-addr-is4", edges, func(r *rand.Rand, i int) receiver {
		if i%3 == 0 {
			var bytes [4]byte
			binary.BigEndian.PutUint32(bytes[:], r.Uint32())
			return receiver{netip.AddrFrom4(bytes), 4}
		}
		var bytes [16]byte
		binary.BigEndian.PutUint64(bytes[:8], r.Uint64())
		binary.BigEndian.PutUint64(bytes[8:], r.Uint64())
		addr := netip.AddrFrom16(bytes)
		if i%3 == 2 {
			addr = addr.WithZone(strconv.FormatUint(r.Uint64(), 16))
		}
		return receiver{addr, 6}
	}, func(c receiver) string {
		return strconv.Itoa(c.family)
	}, func(c receiver) string {
		return strconv.FormatBool(c.addr.Is4())
	}, func(c receiver) string {
		return strconv.FormatBool(c.family == 4)
	})
}
