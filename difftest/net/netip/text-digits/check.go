package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Value   uint16
	Length  int
	Backing []byte
	Want    string
}

func main() {
	f, err := os.CreateTemp("", "golean-netip-digits-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	for spare := range 5 {
		for start := 0; start < 65536; start += 4096 {
			var cases []fixture
			for value := start; value < start+4096; value++ {
				length := value % 4
				b := bytes.Repeat([]byte{byte(value>>8) ^ 0xff}, length+spare)
				cases = append(cases, fixture{uint16(value), length, b, ""})
			}
			compare(cases, rows)
		}
	}
	compare([]fixture{{Value: 0}, {Value: 255}, {Value: 65535}}, rows)
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_NETIP_DIGITS_CASES", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanTextDigits", "difftest/net/netip/text-digits/testdata/text_digits_test.go"))
}

func compare(cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d", c.Value, c.Length)
		for _, b := range c.Backing {
			fmt.Fprintf(&wire, " %d", b)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-netip-text-digits", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
