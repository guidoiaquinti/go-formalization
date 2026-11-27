package netip

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func goLeanAddrState(ip Addr) string {
	kind := 0
	switch {
	case ip.Is4():
		kind = 4
	case ip.Is6():
		kind = 6
	}
	fields := make([]string, len(ip.Zone()))
	for j := range fields {
		fields[j] = strconv.Itoa(int(ip.Zone()[j]))
	}
	return fmt.Sprintf("%d:%d:%d:%s", kind, ip.addr.hi, ip.addr.lo, strings.Join(fields, ","))
}

func TestGoLeanAddrFields(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_ADDR_FIELDS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		var c struct {
			Family    int
			Hi, Lo    uint64
			Index     uint8
			Zone, Raw []byte
			Want      string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		var address Addr
		var rawIP [16]byte
		binary.BigEndian.PutUint64(rawIP[:8], c.Hi)
		binary.BigEndian.PutUint64(rawIP[8:], c.Lo)
		switch c.Family {
		case 4:
			address = AddrFrom4([4]byte(rawIP[12:]))
		case 6:
			address = AddrFrom16(rawIP).WithZone(string(c.Zone))
		}
		raw := c.Raw
		if raw != nil {
			raw = append([]byte{161}, raw...)
			raw = append(raw, 162)
			raw = raw[1 : len(raw)-1]
		}
		original := slices.Clone(raw)
		made, ok := AddrFromSlice(raw)
		if !slices.Equal(raw, original) {
			t.Fatal("constructor changed source")
		}
		before := goLeanAddrState(made)
		for j := range raw {
			raw[j] ^= 255
		}
		if goLeanAddrState(made) != before {
			t.Fatal("constructor aliases input")
		}
		if !ok && made != (Addr{}) {
			t.Fatal("failed conversion returned nonzero address")
		}
		actualFrom := "none"
		if ok {
			actualFrom = before
		}

		i := int(c.Index)
		shift := (280 - 8*(i%32)) % 256
		v4 := byte(0)
		if shift < 64 {
			v4 = rawIP[15-shift/8]
		}
		v6, w6 := rawIP[i%16], binary.BigEndian.Uint16(rawIP[2*(i%8):])
		from := "none"
		switch len(c.Raw) {
		case 4:
			from = fmt.Sprintf("4:0:%d:", uint64(0xffff00000000)|uint64(binary.BigEndian.Uint32(c.Raw)))
		case 16:
			from = fmt.Sprintf("6:%d:%d:", binary.BigEndian.Uint64(c.Raw), binary.BigEndian.Uint64(c.Raw[8:]))
		}
		want := fmt.Sprintf("%d|%d|%d|%t|%t|%d:%d:%d:|%s", v4, v6, w6, c.Family == 0, c.Family == 6 && len(c.Zone) > 0, c.Family, c.Hi, c.Lo, from)
		observed := fmt.Sprintf("%d|%d|%d|%t|%t|%s|%s", address.v4(c.Index), address.v6(c.Index), address.v6u16(c.Index),
			address.isZero(), address.hasZone(), goLeanAddrState(address.withoutZone()), actualFrom)
		if want != c.Want || observed != want {
			t.Fatalf("%+v: native %s, oracle %s", c, observed, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("native rows per address field declaration: %d", count)
}
