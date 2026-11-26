package main

import (
	"bytes"
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
}

func main() {
	prefix := [12]byte{10: 255, 11: 255}
	mapped := [16]byte{10: 255, 11: 255, 12: 192, 14: 2, 15: 1}
	edges := []receiver{{}, {4, mapped}, {6, mapped}, {family: 6}}
	for i := range prefix {
		for bit := range 8 {
			near := mapped
			near[i] ^= 1 << bit
			edges = append(edges, receiver{6, near})
		}
	}
	check.CompareSamples("godiff-addr-as4", edges, func(r *rand.Rand, i int) receiver {
		var raw [16]byte
		r.Read(raw[:])
		switch i % 4 {
		case 0:
			copy(raw[:], prefix[:])
		case 1:
			copy(raw[:], prefix[:])
			raw[r.Intn(12)] ^= 1 << r.Intn(8)
		case 2:
			copy(raw[:], prefix[:])
			return receiver{4, raw}
		}
		return receiver{6, raw}
	}, func(c receiver) string {
		return fmt.Sprintf("%d %d %d", c.family, binary.BigEndian.Uint64(c.raw[:8]), binary.BigEndian.Uint64(c.raw[8:]))
	}, func(c receiver) string {
		var ip netip.Addr
		switch c.family {
		case 4:
			ip = netip.AddrFrom4([4]byte(c.raw[12:]))
		case 6:
			ip = netip.AddrFrom16(c.raw)
		}
		return actual(ip)
	}, func(c receiver) string {
		if c.family == 4 || (c.family == 6 && bytes.Equal(c.raw[:12], prefix[:])) {
			return formatBytes([4]byte(c.raw[12:]))
		}
		if c.family == 0 {
			return "panic: As4 called on IP zero value"
		}
		return "panic: As4 called on IPv6 address"
	})
}

func actual(ip netip.Addr) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if msg, ok := p.(string); ok {
				result = "panic: " + msg
			} else {
				result = fmt.Sprintf("unexpected panic (%T): %v", p, p)
			}
		}
	}()
	return formatBytes(ip.As4())
}

func formatBytes(raw [4]byte) string {
	var fields [4]string
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields[:], " ")
}
