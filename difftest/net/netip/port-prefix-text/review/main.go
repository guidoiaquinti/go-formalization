package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net/netip"
	"os/exec"
	"strconv"
	"strings"
)

type fixture struct {
	family                              int
	hi, lo                              uint64
	port                                uint16
	bits                                int
	zone                                string
	spare                               int
	nilBuffer                           bool
	portString, portText, portOld       string
	prefixString, prefixText, prefixOld string
}

func main() {
	cases := []fixture{
		{family: 0, port: 65535, bits: -1, nilBuffer: true, portString: "invalid AddrPort", prefixString: "invalid Prefix"},
		{family: 0, bits: -1, portString: "invalid AddrPort", prefixString: "invalid Prefix"},
		{family: 4, lo: 0xc0000201, port: 0, bits: 24, spare: 10,
			portString: "192.0.2.1:0", portText: "192.0.2.1:0", portOld: "192.0.2.1:",
			prefixString: "192.0.2.1/24", prefixText: "192.0.2.1/24", prefixOld: "192.0.2.1/"},
		{family: 4, lo: 0xc0000201, port: 7, bits: -1, spare: 8,
			portString: "192.0.2.1:7", portText: "192.0.2.1:7", portOld: "192.0.2.",
			prefixString: "invalid Prefix", prefixText: "invalid Prefix"},
		{family: 6, port: 65535, bits: 128, spare: 5,
			portString: "[::]:65535", portText: "[::]:65535", portOld: "[::]:",
			prefixString: "::/128", prefixText: "::/128", prefixOld: "::/12"},
		{family: 6, lo: 0x0000ffffc0000201, port: 42, bits: 96, zone: "z\x00\xff", spare: 6,
			portString: "[::ffff:192.0.2.1%z\x00\xff]:42", portText: "[::ffff:192.0.2.1%z\x00\xff]:42", portOld: "[",
			prefixString: "::ffff:192.0.2.1/96", prefixText: "::ffff:192.0.2.1/96"},
		{family: 6, hi: 0x2001000000000001, lo: 0x0000000000010001, port: 23, bits: 64, zone: "a", spare: 40,
			portString: "[2001::1:0:0:1:1%a]:23", portText: "[2001::1:0:0:1:1%a]:23", portOld: "[2001::1:0:0:1:1%a]:23",
			prefixString: "2001::1:0:0:1:1/64", prefixText: "2001::1:0:0:1:1/64", prefixOld: "2001::1:0:0:1:1/64"},
	}
	var rows bytes.Buffer
	var want []string
	var counts [8]int
	for i, c := range cases {
		length, missing := 1, 0
		if c.nilBuffer {
			length, missing = 0, 1
		}
		fmt.Fprint(&rows, c.family, " ", c.hi, " ", c.lo, " ", c.port, " ", c.bits, " ", len(c.zone), " ", length, " ", missing)
		for j := range len(c.zone) {
			fmt.Fprint(&rows, " ", c.zone[j])
		}
		for range length + c.spare {
			fmt.Fprint(&rows, " 88")
		}
		fmt.Fprintln(&rows)
		fixed := expected(c)
		if got := native(c); got != strings.Join(fixed, "|") {
			log.Fatalf("native case %d got %q want %q", i, got, fixed)
		}
		want = append(want, strings.Join(fixed, "|"))
		for slot, alt := range mutations(c, fixed) {
			if alt != fixed[slot] {
				counts[slot]++
			}
		}
	}
	wantCounts := [8]int{3, 2, 4, 4, 2, 2, 2, 2}
	if counts != wantCounts {
		log.Fatalf("mutant counts %v want %v", counts, wantCounts)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-netip-port-prefix-text")
	cmd.Stdin = &rows
	source, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(source), "\n"), "\n")
	if len(got) != len(want) {
		log.Fatalf("got %d Lean rows, want %d", len(got), len(want))
	}
	for i, row := range got {
		if row != want[i] {
			log.Fatalf("Lean row %d got %q want %q", i, row, want[i])
		}
	}
	fmt.Printf("netip port/prefix text review: %d pinned native/Lean rows with mutation counts %v\n", len(cases), counts)
}

func expected(c fixture) []string {
	initial := "X"
	if c.nilBuffer {
		initial = ""
	}
	wire := func(text, old string) string {
		return fmt.Sprintf("%t:%t:%s:%s", len(text) <= c.spare, c.nilBuffer && text == "",
			csv(initial+text), csv(initial+old+strings.Repeat("X", c.spare-len(old))))
	}
	pa, pb := wire(c.portText, c.portOld), wire(c.prefixText, c.prefixOld)
	return []string{csv(c.portString), csv(c.portText), pa, pa, csv(c.prefixString), csv(c.prefixText), pb, pb}
}

func mutations(c fixture, want []string) []string {
	out := append([]string(nil), want...)
	if c.family == 6 {
		out[0] = csv(strings.NewReplacer("[", "", "]", "").Replace(c.portString))
	}
	if c.family == 0 {
		out[1] = csv(c.portString)
		out[5] = csv(c.prefixString)
	}
	if c.bits >= 0 && c.zone != "" {
		out[4] = csv(strings.Replace(c.prefixString, "/", "%"+c.zone+"/", 1))
	}
	initial := "X"
	if c.nilBuffer {
		initial = ""
	}
	atomic := func(text string) string {
		reuse := len(text) <= c.spare
		wrote := ""
		if reuse {
			wrote = text
		}
		return fmt.Sprintf("%t:%t:%s:%s", reuse, c.nilBuffer && text == "", csv(initial+text),
			csv(initial+wrote+strings.Repeat("X", c.spare-len(wrote))))
	}
	out[2], out[3] = atomic(c.portText), atomic(c.portText)
	out[6], out[7] = atomic(c.prefixText), atomic(c.prefixText)
	return out
}

func native(c fixture) string {
	var ip netip.Addr
	switch c.family {
	case 4:
		var raw [4]byte
		binary.BigEndian.PutUint32(raw[:], uint32(c.lo))
		ip = netip.AddrFrom4(raw)
	case 6:
		var raw [16]byte
		binary.BigEndian.PutUint64(raw[:8], c.hi)
		binary.BigEndian.PutUint64(raw[8:], c.lo)
		ip = netip.AddrFrom16(raw).WithZone(c.zone)
	}
	port, prefix := netip.AddrPortFrom(ip, c.port), netip.PrefixFrom(ip, c.bits)
	check := func(b []byte, err error) []byte {
		if err != nil || b == nil {
			log.Fatalf("marshal result %x %v", b, err)
		}
		return b
	}
	pt, qt := check(port.MarshalText()), check(prefix.MarshalText())
	observe := func(call func([]byte) ([]byte, error)) string {
		var old []byte
		length := 0
		if !c.nilBuffer {
			old = bytes.Repeat([]byte("X"), c.spare+1)
			length = 1
		}
		res, err := call(old[:length:len(old)])
		if err != nil {
			log.Fatal(err)
		}
		reused := len(res) == 0 || len(old) > 0 && &res[0] == &old[0]
		return fmt.Sprintf("%t:%t:%s:%s", reused, res == nil, csv(string(res)), csv(string(old)))
	}
	p1 := observe(func(b []byte) ([]byte, error) { return port.AppendTo(b), nil })
	q1 := observe(func(b []byte) ([]byte, error) { return prefix.AppendTo(b), nil })
	return strings.Join([]string{csv(port.String()), csv(string(pt)), p1, observe(port.AppendText),
		csv(prefix.String()), csv(string(qt)), q1, observe(prefix.AppendText)}, "|")
}

func csv(raw string) string {
	words := make([]string, len(raw))
	for i := range words {
		words[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(words, ",")
}
