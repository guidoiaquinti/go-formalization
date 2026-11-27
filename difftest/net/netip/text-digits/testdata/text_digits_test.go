package netip

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func goLeanDigitCSV(raw []byte) string {
	out := make([]string, len(raw))
	for i, v := range raw {
		out[i] = strconv.Itoa(int(v))
	}
	return strings.Join(out, ",")
}
func goLeanDigitWire(reused bool, out, backing []byte) string {
	return fmt.Sprintf("%t:%s:%s", reused, goLeanDigitCSV(out), goLeanDigitCSV(backing))
}
func TestGoLeanTextDigits(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_NETIP_DIGITS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		var c struct {
			Value   uint16
			Length  int
			Backing []byte
			Want    string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		var got []string
		hex := strconv.FormatUint(uint64(c.Value), 16)
		padded := strings.Repeat("0", 4-len(hex)) + hex
		parts := []string{strconv.FormatUint(uint64(byte(c.Value)), 10), hex, padded}
		for i, appendDigits := range []func([]byte) []byte{
			func(b []byte) []byte { return appendDecimal(b, byte(c.Value)) },
			func(b []byte) []byte { return appendHex(b, c.Value) },
			func(b []byte) []byte { return appendHexPad(b, c.Value) },
		} {
			var backing, container []byte
			if c.Backing != nil {
				container = append([]byte{165}, c.Backing...)
				container = append(container, 166)
				backing = container[1 : len(container)-1 : len(container)-1]
			}
			out := appendDigits(backing[:c.Length])
			alias := len(backing) > 0 && &out[0] == &backing[0]
			actual := goLeanDigitWire(alias, out, backing)
			wantBacking := slices.Clone(c.Backing)
			reuse := len(parts[i]) <= len(wantBacking)-c.Length
			if i != 2 || reuse {
				copy(wantBacking[c.Length:], parts[i])
			}
			wantOut := append(slices.Clone(c.Backing[:c.Length]), parts[i]...)
			if want := goLeanDigitWire(reuse, wantOut, wantBacking); actual != want {
				t.Fatalf("fixture %+v method %d native %s want %s", c, i, actual, want)
			}
			if container != nil && (container[0] != 165 || container[len(container)-1] != 166) {
				t.Fatal("append changed surrounding sentinel bytes")
			}
			got = append(got, actual)
		}
		if observed := strings.Join(got, "|"); observed != c.Want {
			t.Fatalf("fixture %+v: native %s", c, observed)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("native numeric text observations per declaration: %d", count)
}
