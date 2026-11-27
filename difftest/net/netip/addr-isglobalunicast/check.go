package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
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
	edges := []receiver{{}, {family: 6}, {family: 6, zone: 1}}
	for _, text := range []string{"0.0.0.0", "0.0.0.1", "255.255.255.255", "255.255.255.254", "127.0.0.1", "10.0.0.1", "169.254.0.1", "224.0.0.1", "::1", "fc00::1", "fe80::1", "febf::1", "fec0::1", "ff02::1"} {
		ip := netip.MustParseAddr(text)
		if ip.Is4() {
			edges = append(edges, receiver{4, ip.As16(), 0}, receiver{6, ip.As16(), 1})
		} else {
			edges = append(edges, receiver{6, ip.As16(), 0}, receiver{6, ip.As16(), 1})
		}
	}
	check.CompareSamples("godiff-addr-isglobalunicast", edges, func(r *rand.Rand, i int) receiver {
		var raw [16]byte
		r.Read(raw[:])
		if i%3 != 0 {
			raw = netip.AddrFrom4([4]byte(raw[12:])).As16()
			if i%3 == 1 {
				return receiver{4, raw, 0}
			}
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
			ip = netip.AddrFrom16(c.raw).WithZone(c.zoneString())
		}
		return strconv.FormatBool(ip.IsGlobalUnicast())
	}, func(c receiver) string {
		want := c.family != 0 && net.IP(c.raw[:]).IsGlobalUnicast()
		// net.IP has no zone. A scoped IPv6 zero address differs from the unzoned constant.
		if c.family == 6 && c.raw == [16]byte{} && c.zone != 0 {
			want = true
		}
		return strconv.FormatBool(want)
	})
}
