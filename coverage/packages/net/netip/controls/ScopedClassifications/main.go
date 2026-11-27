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
	loopback = 1 << iota
	multicast
	linkMulticast
	linkUnicast
	interfaceMulticast
)

var (
	loop4      = netip.MustParsePrefix("127.0.0.0/8")
	loop6      = netip.MustParseAddr("::1")
	multicast4 = netip.MustParsePrefix("224.0.0.0/8")
	multicast6 = netip.MustParsePrefix("ff00::/8")
	linkGroup4 = netip.MustParsePrefix("224.0.0.0/24")
	linkGroup6 = netip.MustParsePrefix("ff02::/16")
	link4      = netip.MustParsePrefix("169.254.0.0/16")
	link6      = netip.MustParsePrefix("fe80::/16")
	ifaceGroup = netip.MustParsePrefix("ff01::/16")
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
		{"127.255.255.255", loopback},
		{"::ffff:127.1.2.3%scope", loopback},
		{"::1%scope", loopback},
		{"126.255.255.255", 0},
		{"224.0.0.255", multicast | linkMulticast},
		{"::ffff:224.0.1.0%scope", multicast},
		{"239.255.255.255", multicast},
		{"240.0.0.0", 0},
		{"::ffff:169.254.1.2%scope", linkUnicast},
		{"169.253.255.255", 0},
		{"fe80::1", linkUnicast},
		{"febf::1%scope", linkUnicast},
		{"fec0::1", 0},
		{"ff02::1", multicast | linkMulticast},
		{"ff12::1%scope", multicast | linkMulticast},
		{"ff01::1", multicast | interfaceMulticast},
		{"ff31::1%scope", multicast | interfaceMulticast},
		{"ff03::1", multicast},
		{"2001:db8::1", 0},
	}
	var failures, mutations [5]int
	for _, c := range cases {
		var ip netip.Addr
		if c.address != "" {
			ip = netip.MustParseAddr(c.address)
		}
		actual := [5]bool{
			ip.IsLoopback(), ip.IsMulticast(), ip.IsLinkLocalMulticast(),
			ip.IsLinkLocalUnicast(), ip.IsInterfaceLocalMulticast(),
		}
		unzoned := ip.WithZone("")
		unmapped := unzoned.Unmap()
		wrong := [5]bool{
			loop4.Contains(unzoned) || unzoned == loop6,
			multicast4.Contains(unmapped) || multicast6.Contains(unmapped),
			linkGroup4.Contains(unmapped) || linkGroup6.Contains(unmapped),
			link4.Contains(unmapped) || link6.Contains(unmapped),
			ifaceGroup.Contains(unzoned),
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
	if failures != [5]int{} || mutations != [5]int{1, 1, 1, 1, 1} {
		log.Fatalf("netip classifications: Go=%v mutants=%v", failures, mutations)
	}
	fmt.Printf("%s %s/%s: scoped address classifications Go=%v mutants=%v across %d cases\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, failures, mutations, len(cases))
}
