package main

import (
	"fmt"
	"log"
	"net/netip"
	"os"
	"runtime"
	"strings"
)

const (
	global = 1 << iota
	private
	unspecified
)

var (
	private4 = []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("172.16.0.0/12"),
		netip.MustParsePrefix("192.168.0.0/16"),
	}
	private6 = netip.MustParsePrefix("fc00::/7")
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if runtime.Version() != "go"+strings.TrimSpace(string(pin)) {
		log.Fatalf("unexpected Go version %s", runtime.Version())
	}
	cases := []struct {
		address string
		flags   uint
	}{
		{"", 0},
		{"0.0.0.0", unspecified},
		{"::", unspecified},
		{"::%scope", global},
		{"::ffff:0.0.0.0%scope", 0},
		{"255.255.255.255", 0},
		{"::ffff:255.255.255.255", 0},
		{"127.1.2.3", 0},
		{"224.1.2.3", 0},
		{"169.254.1.2", 0},
		{"::1%scope", 0},
		{"ff02::1%scope", 0},
		{"febf::1%scope", 0},
		{"fec0::1", global},
		{"192.0.2.1", global},
		{"10.255.255.255", global | private},
		{"::ffff:10.1.2.3%scope", global | private},
		{"172.15.255.255", global},
		{"172.16.0.0", global | private},
		{"172.31.255.255", global | private},
		{"172.32.0.0", global},
		{"192.168.0.0", global | private},
		{"192.169.0.0", global},
		{"fc00::1", global | private},
		{"fdff::1%scope", global | private},
		{"fe00::1", global},
	}
	var failures, mutations [3]int
	for _, c := range cases {
		var ip netip.Addr
		if c.address != "" {
			ip = netip.MustParseAddr(c.address)
		}
		actual := [3]bool{ip.IsGlobalUnicast(), ip.IsPrivate(), ip.IsUnspecified()}
		unzoned := ip.WithZone("")
		privateWithoutUnmapping := private6.Contains(unzoned)
		for _, prefix := range private4 {
			privateWithoutUnmapping = privateWithoutUnmapping || prefix.Contains(unzoned)
		}
		wrong := [3]bool{
			unzoned.IsGlobalUnicast(),
			privateWithoutUnmapping,
			ip.Unmap().IsUnspecified(),
		}
		for i := range actual {
			want := c.flags&(1<<i) != 0
			if actual[i] != want {
				failures[i]++
			}
			if wrong[i] != want {
				mutations[i]++
			}
		}
	}
	if failures != [3]int{} || mutations != [3]int{1, 1, 1} {
		log.Fatalf("netip classifications: Go=%v mutants=%v", failures, mutations)
	}
	fmt.Printf("%s %s/%s: global, private, unspecified Go=%v mutants=%v across %d cases\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, failures, mutations, len(cases))
}
