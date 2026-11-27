package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net/netip"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	family        int
	hi, lo        uint64
	port          uint16
	bits, length  int
	zone, backing []byte
}

func address(c fixture) netip.Addr {
	switch c.family {
	case 4:
		var a [4]byte
		binary.BigEndian.PutUint32(a[:], uint32(c.lo))
		return netip.AddrFrom4(a)
	case 6:
		var a [16]byte
		binary.BigEndian.PutUint64(a[:8], c.hi)
		binary.BigEndian.PutUint64(a[8:], c.lo)
		return netip.AddrFrom16(a).WithZone(string(c.zone))
	default:
		return netip.Addr{}
	}
}

func csv(raw []byte) string {
	out := make([]string, len(raw))
	for i, v := range raw {
		out[i] = strconv.Itoa(int(v))
	}
	return strings.Join(out, ",")
}
func byteWire(reused bool, out, old []byte) string {
	return fmt.Sprintf("%t:%t:%s:%s", reused, out == nil, csv(out), csv(old))
}
func observedAppend(c fixture, fn func([]byte) []byte) string {
	old := slices.Clone(c.backing)
	out := fn(old[:c.length:len(old)])
	reuse := len(out) == 0 || len(old) > 0 && &out[0] == &old[0]
	return byteWire(reuse, out, old)
}
func appendOracle(c fixture, chunks []string) string {
	old := slices.Clone(c.backing)
	out := slices.Clone(c.backing[:c.length])
	cursor, reuse := c.length, true
	for _, chunk := range chunks {
		out = append(out, chunk...)
		if reuse && len(chunk) <= len(old)-cursor {
			copy(old[cursor:], chunk)
			cursor += len(chunk)
		} else {
			reuse = false
		}
	}
	return byteWire(reuse, out, old)
}

func characters(s string) []string {
	out := make([]string, len(s))
	for i := range out {
		out[i] = s[i : i+1]
	}
	return out
}

// Use the already observed public Addr text as the address oracle. Recover
// each source append stage from its separators and the raw zone suffix.
func hostChunks(ip netip.Addr) []string {
	if !ip.IsValid() {
		return nil
	}
	host := ip.WithZone("").String()
	var out []string
	if ip.Is4In6() {
		out = append(out, "::ffff:")
		host = strings.TrimPrefix(host, "::ffff:")
	}
	for i := 0; i < len(host); i++ {
		if host[i] == ':' && i+1 < len(host) && host[i+1] == ':' {
			out = append(out, "::")
			i++
		} else {
			out = append(out, host[i:i+1])
		}
	}
	if zone := ip.Zone(); zone != "" {
		out = append(out, "%", zone)
	}
	return out
}

func observe(c fixture) (string, string) {
	ip := address(c)
	p, q := netip.AddrPortFrom(ip, c.port), netip.PrefixFrom(ip, c.bits)
	portParts := hostChunks(ip)
	portString := "invalid AddrPort"
	if ip.IsValid() {
		if ip.Is6() {
			portParts = append([]string{"["}, portParts...)
			portParts = append(portParts, "]")
		}
		portParts = append(portParts, ":", strconv.FormatUint(uint64(c.port), 10))
		portString = strings.Join(portParts, "")
	}
	prefixString := "invalid Prefix"
	var prefixParts []string
	if ip.IsValid() {
		if c.bits < 0 {
			prefixParts = []string{"invalid Prefix"}
		} else {
			prefixParts = append(hostChunks(ip.WithZone("")), "/")
			prefixParts = append(prefixParts, characters(strconv.Itoa(c.bits))...)
			prefixString = strings.Join(prefixParts, "")
		}
	}
	portAppend, prefixAppend := appendOracle(c, portParts), appendOracle(c, prefixParts)
	want := []string{csv([]byte(portString)), csv([]byte(strings.Join(portParts, ""))), portAppend, portAppend,
		csv([]byte(prefixString)), csv([]byte(strings.Join(prefixParts, ""))), prefixAppend, prefixAppend}
	marshal := func(fn func() ([]byte, error)) string {
		ret, err := fn()
		if err != nil || ret == nil {
			log.Fatalf("unexpected marshal result %v, %v", ret, err)
		}
		return csv(ret)
	}
	appendText := func(fn func([]byte) ([]byte, error)) string {
		return observedAppend(c, func(b []byte) []byte {
			out, err := fn(b)
			if err != nil {
				log.Fatal(err)
			}
			return out
		})
	}
	got := []string{csv([]byte(p.String())), marshal(p.MarshalText), observedAppend(c, p.AppendTo), appendText(p.AppendText),
		csv([]byte(q.String())), marshal(q.MarshalText), observedAppend(c, q.AppendTo), appendText(q.AppendText)}
	return strings.Join(got, "|"), strings.Join(want, "|")
}

func compare(label string, cases []fixture) {
	var input bytes.Buffer
	for _, c := range cases {
		missing := 0
		if c.backing == nil {
			missing = 1
		}
		fmt.Fprintf(&input, "%d %d %d %d %d %d %d %d", c.family, c.hi, c.lo, c.port, c.bits, len(c.zone), c.length, missing)
		for _, raw := range [][]byte{c.zone, c.backing} {
			for _, v := range raw {
				fmt.Fprintf(&input, " %d", v)
			}
		}
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-netip-port-prefix-text", len(cases), &input) {
		c := cases[i]
		native, want := observe(c)
		if got != want || native != want {
			log.Fatalf("%s row %d %+v: native %q Lean %q oracle %q", label, i, c, native, got, want)
		}
	}
	fmt.Printf("netip port and prefix text %s: %d matched.\n", label, len(cases))
}

func normal(c fixture) fixture {
	switch c.family {
	case 0:
		c.hi, c.lo, c.bits, c.zone = 0, 0, -1, nil
	case 4:
		c.hi, c.lo, c.zone = 0, 0xffff00000000|(c.lo&0xffffffff), nil
		if c.bits > 32 {
			c.bits = -1
		}
	case 6:
		if c.bits > 128 {
			c.bits = -1
		}
	}
	return c
}

func main() {
	var fixed []fixture
	ips := []fixture{{}, {family: 4, lo: 0x01020004}, {family: 6, lo: 0xffffaabbccdd, zone: []byte{0, 255, 37}},
		{family: 6, hi: 0x20010db800000001, lo: 0x0000000000010001, zone: []byte("zone")}}
	for _, c := range ips {
		for _, port := range []uint16{0, 1, 9, 10, 99, 100, 999, 1000, 9999, 10000, 65535} {
			for _, bits := range []int{-1, 0, 1, 9, 10, 24, 32, 33, 99, 100, 127, 128, 129} {
				for _, spare := range []int{0, 1, 2, 3, 4, 5, 8, 13, 24, 39, 51} {
					c.port, c.bits, c.length = port, bits, 2
					c.backing = bytes.Repeat([]byte{165}, c.length+spare)
					fixed = append(fixed, normal(c))
				}
			}
		}
	}
	for _, c := range ips {
		fixed = append(fixed, normal(c))
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				c := fixture{family: []int{0, 4, 6, 6}[rng.Intn(4)], hi: rng.Uint64(), lo: rng.Uint64(), port: uint16(rng.Uint32()), bits: rng.Intn(130) - 1, length: rng.Intn(6)}
				c.backing = make([]byte, c.length+rng.Intn(58))
				c.zone = make([]byte, rng.Intn(13))
				rng.Read(c.backing)
				rng.Read(c.zone)
				if i%5 == 0 {
					c.hi, c.lo = 0, 0xffff00000000|(c.lo&0xffffffff)
				}
				if i%11 == 0 {
					c.length, c.backing = 0, nil
				}
				cases[i] = normal(c)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases)
		}
	}
}
