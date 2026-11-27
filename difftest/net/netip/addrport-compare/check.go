package main

import (
	"bytes"
	"cmp"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/netip"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type addr struct {
	family int
	raw    [16]byte
	zone   uint32
}

type pair struct {
	a, b   addr
	ap, bp uint16
}

func (a addr) zoneString() string {
	if a.zone == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(a.zone), 10)
}

func (a addr) goAddr() netip.Addr {
	switch a.family {
	case 4:
		return netip.AddrFrom4([4]byte(a.raw[12:]))
	case 6:
		return netip.AddrFrom16(a.raw).WithZone(a.zoneString())
	default:
		return netip.Addr{}
	}
}

func (a addr) input() string {
	return fmt.Sprintf("%d %d %d %d", a.family, binary.BigEndian.Uint64(a.raw[:8]), binary.BigEndian.Uint64(a.raw[8:]), a.zone)
}

func main() {
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 1}
	addrs := []addr{{}, {4, mapped, 0}, {6, mapped, 0}, {family: 6}, {family: 6, zone: 2}, {family: 6, zone: 10}, {6, [16]byte{7: 1}, 0}, {6, [16]byte{8: 1}, 0}}
	var edges []pair
	for _, a := range addrs {
		for _, b := range addrs {
			for _, ports := range [][2]uint16{{0, 0}, {0, 65535}, {65535, 0}, {32767, 32768}} {
				edges = append(edges, pair{a, b, ports[0], ports[1]})
			}
		}
	}
	randomAddr := func(r *rand.Rand) addr {
		var a addr
		switch r.Intn(3) {
		case 0:
		case 1:
			a.family = 4
			a.raw = mapped
			binary.BigEndian.PutUint32(a.raw[12:], r.Uint32())
		case 2:
			a.family = 6
			r.Read(a.raw[:])
			a.zone = r.Uint32()
		}
		return a
	}
	check.CompareSamples("godiff-addrport-compare", edges, func(r *rand.Rand, i int) pair {
		a := randomAddr(r)
		b := a
		switch i % 4 {
		case 0:
			b = randomAddr(r)
		case 1:
			if a.family == 6 {
				b.zone = r.Uint32()
			}
		case 2:
			if a.family != 0 {
				b.raw[15] ^= 1
			}
		}
		ap := uint16(r.Uint32())
		bp := ap
		if i%4 != 3 {
			bp = uint16(r.Uint32())
		}
		return pair{a, b, ap, bp}
	}, func(p pair) string { return fmt.Sprintf("%s %d %s %d", p.a.input(), p.ap, p.b.input(), p.bp) },
		func(p pair) string {
			return strconv.Itoa(netip.AddrPortFrom(p.a.goAddr(), p.ap).Compare(netip.AddrPortFrom(p.b.goAddr(), p.bp)))
		},
		func(p pair) string {
			if p.a.family != p.b.family {
				return strconv.Itoa(cmp.Compare(p.a.family, p.b.family))
			}
			if order := bytes.Compare(p.a.raw[:], p.b.raw[:]); order != 0 {
				return strconv.Itoa(order)
			}
			if p.a.family == 6 {
				if order := strings.Compare(p.a.zoneString(), p.b.zoneString()); order != 0 {
					return strconv.Itoa(order)
				}
			}
			return strconv.Itoa(cmp.Compare(p.ap, p.bp))
		})
}
