package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/netip"
	"os"
	"runtime"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}
	cases := []struct {
		name string
		addr netip.Addr
		want [3]bool
	}{
		{"zero", netip.Addr{}, [3]bool{}},
		{"IPv4 unspecified", netip.IPv4Unspecified(), [3]bool{true, false, false}},
		{"IPv6 unspecified", netip.IPv6Unspecified(), [3]bool{false, true, false}},
		{"mapped IPv6", netip.MustParseAddr("::ffff:192.0.2.1"), [3]bool{false, true, true}},
		{"mapped IPv6 with zone", netip.MustParseAddr("::ffff:192.0.2.1%scope"), [3]bool{false, true, true}},
		{"nonzero high word", netip.MustParseAddr("1::ffff:192.0.2.1"), [3]bool{false, true, false}},
		{"different marker", netip.MustParseAddr("::fffe:192.0.2.1"), [3]bool{false, true, false}},
	}
	var failures [3]int
	for _, test := range cases {
		ip := test.addr
		got := [3]bool{ip.Is4(), ip.Is6(), ip.Is4In6()}
		if got != test.want {
			log.Fatalf("%s produced %v, want %v", test.name, got, test.want)
		}
		raw := ip.As16()
		mutant := [3]bool{
			ip.Unmap().Is4(),
			!ip.Is4(),
			ip.Is6() && binary.BigEndian.Uint32(raw[8:12]) == 0xffff,
		}
		for i := range mutant {
			if mutant[i] != test.want[i] {
				failures[i]++
			}
		}
	}
	if want := [3]int{2, 1, 1}; failures != want {
		log.Fatalf("mutant failure counts %v, want %v", failures, want)
	}
	fmt.Printf("%s %s/%s: %d receivers agree, mutant failures for Is4, Is6, Is4In6 = %v\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, len(cases), failures)
}
