package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Family    int
	Hi, Lo    uint64
	Index     uint8
	Zone, Raw []byte
	Want      string
}

func main() {
	const high, low = 0x0102030405060708, 0x1112131415161718
	const v4Mapped = 0x0000ffffc0000201
	v4 := []byte{192, 0, 2, 1}
	mapped := append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255}, v4...)
	zone := []byte{255, 97}
	cases := []fixture{
		{0, 0, 0, 0, nil, []byte{1, 2, 3}, "0|0|0|true|false|0:0:0:|none"},
		{4, 0, v4Mapped, 32, nil, v4, fmt.Sprintf("192|0|0|false|false|4:0:%d:|4:0:%d:", v4Mapped, v4Mapped)},
		{6, high, low, 8, zone, mapped, fmt.Sprintf("0|17|258|false|true|6:%d:%d:|6:0:%d:", high, low, v4Mapped)},
		{6, high, low, 15, nil, nil, fmt.Sprintf("0|24|5912|false|false|6:%d:%d:|none", high, low)},
		{6, 0, 0, 0, nil, make([]byte, 16), "0|0|0|false|false|6:0:0:|6:0:0:"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d", c.Family, c.Hi, c.Lo, c.Index, len(c.Zone))
		for _, group := range [][]byte{c.Zone, c.Raw} {
			for _, b := range group {
				fmt.Fprintf(&input, " %d", b)
			}
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-netip-addr-fields", len(cases), &input)
	file, err := os.CreateTemp("", "golean-netip-addr-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_ADDR_FIELDS_CASES", file.Name()))
	const data = "difftest/net/netip/addr-fields/testdata/"
	must(check.TestInPackage("net/netip", "TestGoLeanAddrFields", data+"addr_fields_test.go"))
	if strings.HasPrefix(rows[1], "0|") || strings.HasSuffix(rows[2], fmt.Sprintf("|4:0:%d:", v4Mapped)) ||
		strings.Contains(rows[2], ":255,97|") || strings.Contains(rows[4], "|true|") {
		log.Fatal("address-field mutation accepted")
	}
	fmt.Println("Five independent native and Lean address rows passed with four mutation witnesses.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
