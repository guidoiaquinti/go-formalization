package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Input, Expected, Zone []byte
	Family                int
	Want                  string
}

func main() {
	f, err := os.CreateTemp("", "golean-netip-address-")
	must(err)
	defer os.Remove(f.Name())
	out := bufio.NewWriter(f)
	encoder := json.NewEncoder(out)
	boundaries := []string{"", "0", ".1.2.3", "127.0.0.1", "::", "::1", "2001:db8:0:1:2:3:4:5", "2001:DB8::AbCD",
		"fe80::1%lo0", "fe80::1%\x00\xff%zone", "::%", "%zone", ":", "1:", ":1", "1::", ":::", "::1::2",
		"1:2:3:4:5:6:7:8:", "1:2:3:4:5:6:7:8:9", "1:2:3:4:5:6:7::8", "12345::", "1:g", "1:2-",
		"::ffff:192.0.2.1", "::ffff:192.0.002.1", "1:2:3:4:5:6:192.0.2.1",
		"1:2:3:4:5:192.0.2.1", "::1:2:3:4:5:6:7:192.0.2.1", "::1:2:3:4:5:6:192.0.2.1",
		"::ffff:192.0.2.1%zone", "1.2.3.4:7", "1.2.3.4%zone", "x\xff", "\x00:01"}
	var fixed []fixture
	for _, raw := range boundaries {
		for end := 0; end <= len(raw); end++ {
			fixed = append(fixed, fixture{Input: []byte(raw[:end])})
		}
	}
	compare("fixed", fixed, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				switch i % 5 {
				case 0:
					raw := make([]byte, 4)
					rng.Read(raw)
					rows[i] = fixture{Input: fmt.Appendf(nil, "%d.%d.%d.%d", raw[0], raw[1], raw[2], raw[3]), Expected: raw, Family: 4}
				case 1, 2, 3:
					rows[i] = validV6(rng, i%5)
				case 4:
					raw := make([]byte, rng.Intn(60))
					rng.Read(raw)
					if i%2 == 0 {
						const alphabet = "abcdefABD1234567890::.:..%%"
						for j := range raw {
							raw[j] = alphabet[int(raw[j])%len(alphabet)]
						}
					}
					rows[i] = fixture{Input: raw}
				}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows, encoder)
		}
	}
	must(out.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_NETIP_PARSE_ADDRESS", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanParseAddress", "difftest/net/netip/parse-address/testdata/address_test.go"))
}

func validV6(rng *rand.Rand, style int) fixture {
	raw := make([]byte, 16)
	rng.Read(raw)
	groups := make([]string, 8)
	for i := range groups {
		groups[i] = fmt.Sprintf("%x", int(raw[2*i])*256+int(raw[2*i+1]))
	}
	var text string
	switch style {
	case 1:
		text = strings.Join(groups, ":")
	case 2:
		start := rng.Intn(8)
		stop := start + 1 + rng.Intn(8-start)
		clear(raw[2*start : 2*stop])
		text = strings.Join(groups[:start], ":") + "::" + strings.Join(groups[stop:], ":")
	case 3:
		clear(raw[:10])
		raw[10], raw[11] = 255, 255
		text = fmt.Sprintf("::ffff:%d.%d.%d.%d", raw[12], raw[13], raw[14], raw[15])
	}
	var zone []byte
	if rng.Intn(2) == 0 {
		zone = make([]byte, rng.Intn(8)+1)
		rng.Read(zone)
		text += "%" + string(zone)
	}
	return fixture{Input: []byte(text), Expected: raw, Zone: zone, Family: 6}
}

func compare(label string, rows []fixture, encoder *json.Encoder) {
	var input bytes.Buffer
	for _, r := range rows {
		fmt.Fprint(&input, len(r.Input))
		for _, b := range r.Input {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	for i, result := range check.LeanRows("godiff-netip-parse-address", len(rows), &input) {
		r := rows[i]
		r.Want = result
		must(encoder.Encode(r))
	}
	fmt.Printf("netip address parsing %s: %d Lean rows\n", label, len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
