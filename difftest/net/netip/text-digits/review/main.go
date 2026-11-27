package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Value   uint16
	Length  int
	Backing []byte
	Want    string
}

func main() {
	cases := []fixture{
		{0xabff, 1, []byte{7, 85, 86}, "false:7,50,53,53:7,50,53|false:7,97,98,102,102:7,97,98|false:7,97,98,102,102:7,85,86"},
		{0xabff, 1, []byte{7, 85, 86, 87, 88}, "true:7,50,53,53:7,50,53,53,88|true:7,97,98,102,102:7,97,98,102,102|true:7,97,98,102,102:7,97,98,102,102"},
		{10, 0, []byte{85}, "false:49,48:49|true:97:97|false:48,48,48,97:85"},
		{0, 0, nil, "false:48:|false:48:|false:48,48,48,48:"},
		{256, 0, []byte{}, "false:48:|false:49,48,48:|false:48,49,48,48:"},
	}
	fields := strings.Split(cases[0].Want, "|")
	if fields[0] == "false:7,50,53,53:7,85,86" ||
		fields[1] == "false:7,97,98,102,102:7,85,86" ||
		fields[2] == "false:7,97,98,102,102:7,97,98" {
		panic("atomicity mutation survived")
	}
	if strings.Split(cases[2].Want, "|")[2] == "false:97:85" {
		panic("padding mutation survived")
	}
	fmt.Println("detected atomizing decimal and hex, splitting padded append, and omitted padding")

	f, err := os.CreateTemp("", "review-netip-digits-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	var protocol bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&protocol, c.Value, " ", c.Length)
		for _, b := range c.Backing {
			fmt.Fprint(&protocol, " ", b)
		}
		fmt.Fprintln(&protocol)
		must(encoder.Encode(c))
	}
	must(f.Close())
	for i, row := range check.LeanRows("godiff-netip-text-digits", len(cases), &protocol) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_NETIP_DIGITS_CASES", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanTextDigits", "difftest/net/netip/text-digits/testdata/text_digits_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
