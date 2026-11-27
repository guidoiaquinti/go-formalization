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

type addr struct {
	family int
	raw    [16]byte
	zone   uint32
}

type pair struct{ a, b addr }

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
			edges = append(edges, pair{a, b})
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
	check.CompareSamples("godiff-addr-less", edges, func(r *rand.Rand, i int) pair {
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
		return pair{a, b}
	}, func(p pair) string { return p.a.input() + " " + p.b.input() },
		func(p pair) string { return strconv.FormatBool(p.a.goAddr().Less(p.b.goAddr())) },
		func(p pair) string {
			if p.a.family != p.b.family {
				return strconv.FormatBool(p.a.family < p.b.family)
			}
			if order := bytes.Compare(p.a.raw[:], p.b.raw[:]); order != 0 {
				return strconv.FormatBool(order < 0)
			}
			return strconv.FormatBool(p.a.family == 6 && p.a.zoneString() < p.b.zoneString())
		})
}
