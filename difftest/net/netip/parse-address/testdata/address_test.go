package netip

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanParseAddress(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_NETIP_PARSE_ADDRESS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	known := 0
	for row := 0; ; row++ {
		var c struct {
			Input, Expected, Zone []byte
			Family                int
			Want                  string
		}
		if err := dec.Decode(&c); err != nil {
			if err == io.EOF && row > 0 {
				t.Logf("%d pinned rows per declaration, %d independent known-byte addresses", row, known)
				return
			}
			t.Fatal(err)
		}
		input := string(c.Input)
		v6, e6 := parseIPv6(input)
		a, ea := ParseAddr(input)
		dst := IPv6Loopback().WithZone("initial")
		mutErr := dst.UnmarshalText(c.Input)
		text := "clear"
		if len(c.Input) == 0 {
			if dst != (Addr{}) || mutErr != nil {
				t.Fatal("empty input did not clear address")
			}
		} else {
			text = parseObservation(dst, mutErr)
		}
		got := strings.Join([]string{parseObservation(v6, e6), parseObservation(a, ea), observeMust(input), text}, "|")
		if got != c.Want {
			t.Fatalf("row %d %q: Go %q want %q", row, input, got, c.Want)
		}
		if c.Family != 0 {
			known++
			expected := fmt.Sprintf("v%d:%s", c.Family, encodeAddressBytes(c.Expected))
			if c.Family == 6 {
				expected += ":" + encodeAddressBytes(c.Zone)
			}
			if parseObservation(a, ea) != expected || c.Family == 6 && parseObservation(v6, e6) != expected {
				t.Fatalf("row %d %q: address does not match independent construction bytes", row, input)
			}
		}
	}
}

func encodeAddressBytes(raw []byte) string {
	parts := make([]string, len(raw))
	for i, b := range raw {
		parts[i] = strconv.Itoa(int(b))
	}
	return strings.Join(parts, ",")
}

func observationError(err error) string {
	e, ok := err.(parseAddrError)
	if !ok {
		panic("unexpected address error type")
	}
	return "error:" + e.msg + ":" + encodeAddressBytes([]byte(e.in)) + ":" + encodeAddressBytes([]byte(e.at))
}

func parseObservation(a Addr, err error) string {
	if err != nil {
		if a != (Addr{}) {
			panic("nonzero address on failure")
		}
		return observationError(err)
	}
	if a.Is4() {
		raw := a.As4()
		if a.Zone() != "" {
			panic("IPv4 with a zone")
		}
		return "v4:" + encodeAddressBytes(raw[:])
	}
	if !a.Is6() {
		panic("invalid address on success")
	}
	raw := a.As16()
	return "v6:" + encodeAddressBytes(raw[:]) + ":" + encodeAddressBytes([]byte(a.Zone()))
}

func observeMust(input string) (result string) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(parseAddrError)
			if !ok {
				panic(failure)
			}
			result = "panic:" + observationError(err)
		}
	}()
	return parseObservation(MustParseAddr(input), nil)
}
