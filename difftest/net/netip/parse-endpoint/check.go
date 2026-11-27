package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Raw                                []byte
	Want, ExpectedPort, ExpectedPrefix string
}

func main() {
	f, err := os.CreateTemp("", "golean-netip-endpoint-*")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, strconv.IntSize)
			for _, b := range c.Raw {
				fmt.Fprint(&input, " ", b)
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-netip-parse-endpoint", len(cases), &input) {
			c := cases[i]
			c.Want = row
			must(encoder.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows recorded\n", label, len(cases))
	}
	hosts := []string{"", "192.0.2.1", "[192.0.2.1]", "::", "[::]", "fe80::1%eth", "[fe80::1%eth]",
		"[::ffff:192.0.2.1]", "2001:db8::1", "[2001:db8::1", "192.0.2.256", "x", "%zone", "fe80::1%a/b"}
	values := []string{"", "0", "00", "001", "1", "9", "32", "33", "128", "129", "65535", "65536", "-1", "+1",
		"2147483647", "2147483648", "9223372036854775807", "9223372036854775808", "1_2", "1\xff"}
	var fixed []fixture
	for _, h := range hosts {
		fixed = append(fixed, fixture{Raw: []byte(h)})
		for _, n := range values {
			for _, sep := range []string{":", "/"} {
				fixed = append(fixed, fixture{Raw: []byte(h + sep + n)})
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20000)
		for i := range cases {
			last := rng.Intn(256)
			limit := rng.Intn(33)
			port := rng.Intn(65536)
			switch i % 5 {
			case 0:
				cases[i] = fixture{Raw: fmt.Appendf(nil, "192.0.2.%d:%05d", last, port),
					ExpectedPort: fmt.Sprintf("v4:192,0,2,%d::%d", last, port)}
			case 1:
				cases[i] = fixture{Raw: fmt.Appendf(nil, "192.0.2.%d/%d", last, limit),
					ExpectedPrefix: fmt.Sprintf("v4:192,0,2,%d::%d", last, limit)}
			case 2:
				cases[i] = fixture{Raw: fmt.Appendf(nil, "[::%x%%z]:%d", last, port),
					ExpectedPort: fmt.Sprintf("v6:0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,%d:122:%d", last, port)}
			case 3:
				cases[i] = fixture{Raw: fmt.Appendf(nil, "2001:db8::%x/%d", last, limit+96),
					ExpectedPrefix: fmt.Sprintf("v6:32,1,13,184,0,0,0,0,0,0,0,0,0,0,0,%d::%d", last, limit+96)}
			case 4:
				h, n := hosts[rng.Intn(len(hosts))], values[rng.Intn(len(values))]
				sep := []string{":", "/", "//", ""}[rng.Intn(4)]
				cases[i] = fixture{Raw: []byte(h + sep + n)}
			}
		}
		compare(fmt.Sprint("seed ", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_NETIP_ENDPOINT", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanEndpoint", "difftest/net/netip/parse-endpoint/testdata/endpoint_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
