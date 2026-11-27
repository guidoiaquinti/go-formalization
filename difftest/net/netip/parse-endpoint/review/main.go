package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Raw                          []byte
	Want                         string
	ExpectedPort, ExpectedPrefix string
}

func outcomes(port, prefix string, empty bool) string {
	transcript := func(observation string) []string {
		if !strings.HasPrefix(observation, "error:") {
			return []string{observation, observation, observation}
		}
		err := strings.TrimPrefix(observation, "error:")
		text := observation
		if empty {
			text = "clear"
		}
		return []string{observation, "panic:" + err, text}
	}
	return strings.Join(append(transcript(port), transcript(prefix)...), "|")
}

func main() {
	const noSlash = "error:no slash"
	const noEndpoint = "error:not an ip:port"
	const v4 = "v4:192,0,2,129::"
	large := "error:out of range"
	if strconv.IntSize == 32 {
		large = "error:bad bits"
	}
	cases := []fixture{
		{nil, outcomes(noEndpoint, noSlash, true), "", ""},
		{[]byte("192.0.2.129:00080"), outcomes(v4+"80", noSlash, false), v4 + "80", ""},
		{[]byte("[192.0.2.129]:7"), outcomes("error:IPv4 bracketed", noSlash, false), "", ""},
		{[]byte("::1:7"), outcomes("error:IPv6 unbracketed", noSlash, false), "", ""},
		{[]byte("bad:65536"), outcomes("error:invalid port", noSlash, false), "", ""},
		{[]byte("192.0.2.129/24"), outcomes(noEndpoint, v4+"24", false), "", v4 + "24"},
		{[]byte("192.0.2.129/024"), outcomes(noEndpoint, "error:bad bits", false), "", ""},
		{[]byte("::1%zone/64"), outcomes("error:invalid port", "error:zone", false), "", ""},
		{[]byte("192.0.2.129/2147483648"), outcomes(noEndpoint, large, false), "", ""},
	}
	f, err := os.CreateTemp("", "golean-netip-endpoint-review-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, strconv.IntSize)
		for _, b := range c.Raw {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(f.Close())
	var rejected [8]int
	for i, got := range check.LeanRows("godiff-netip-parse-endpoint", len(cases), &input) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d %q: Lean %q want %q", i, c.Raw, got, c.Want)
		}
		parts := strings.Split(got, "|")
		wrong := [8]string{parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[0], parts[3]}
		switch string(c.Raw) {
		case "bad:65536":
			wrong[0] = "error:address fault"
		case "":
			wrong[1] = parts[0]
			wrong[2] = "error:not an ip:port"
			wrong[4] = parts[3]
			wrong[5] = "error:no slash"
		case "192.0.2.129/24":
			wrong[3] = "v4:192,0,2,0::24"
		case "[192.0.2.129]:7":
			wrong[6] = v4 + "7"
		case "192.0.2.129/024":
			wrong[7] = v4 + "24"
		}
		for j, value := range wrong {
			indices := [8]int{0, 1, 2, 3, 4, 5, 0, 3}
			if value != parts[indices[j]] {
				rejected[j]++
			}
		}
	}
	if want := [8]int{1, 1, 1, 1, 1, 1, 1, 1}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	rawWide := cases[len(cases)-1].Raw
	var widths bytes.Buffer
	for _, width := range []int{32, 64} {
		fmt.Fprint(&widths, width)
		for _, b := range rawWide {
			fmt.Fprint(&widths, " ", b)
		}
		fmt.Fprintln(&widths)
	}
	widthWant := []string{outcomes(noEndpoint, "error:bad bits", false),
		outcomes(noEndpoint, "error:out of range", false)}
	for i, got := range check.LeanRows("godiff-netip-parse-endpoint", len(widthWant), &widths) {
		if got != widthWant[i] {
			log.Fatalf("prefix profile %d: Lean %q want %q", i, got, widthWant[i])
		}
	}
	must(os.Setenv("GOLEAN_NETIP_ENDPOINT", f.Name()))
	must(check.TestInPackage("net/netip", "TestGoLeanEndpoint",
		"difftest/net/netip/parse-endpoint/testdata/endpoint_test.go"))
	fmt.Printf("netip endpoint review: %d native and Lean cases, mutation counts %v\n", len(cases), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
