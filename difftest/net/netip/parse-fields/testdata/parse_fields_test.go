package netip

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanParseFields(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_NETIP_PARSE_FIELDS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := bufio.NewScanner(f)
	rows := 0
	for in.Scan() {
		left, want, ok := strings.Cut(in.Text(), "\t")
		if !ok {
			t.Fatal("missing Lean columns")
		}
		var lo, hi int
		var start [4]byte
		var text string
		if _, err := fmt.Sscanf(left, "%d %d %d %d %d %d %q", &lo, &hi, &start[0], &start[1], &start[2], &start[3], &text); err != nil {
			t.Fatal(err)
		}
		cells := start
		part := parseIPv4Fields(text, lo, hi, cells[:])
		a, whole := parseIPv4(text)
		wholeStatus := errorWire(whole)
		if whole == nil {
			if !a.Is4() || a.Zone() != "" {
				t.Fatal("unexpected address family or zone")
			}
			raw := a.As4()
			wholeStatus += ":" + bytesWire(raw[:])
		} else if a != (Addr{}) {
			t.Fatal("nonzero failure result")
		}
		ip, port, bracketed, err := splitAddrPort(text)
		split := splitWire(ip, port, bracketed, err)
		if err != nil && (ip != "" || port != "" || bracketed) {
			t.Fatal("unexpected partial split on error")
		}
		got := strings.Join([]string{bytesWire(cells[:]), errorWire(part), wholeStatus, split}, "|")
		refCells, refPart := fieldsOracle(text, lo, hi, start)
		refWhole, refErr := fieldsOracle(text, 0, len(text), [4]byte{})
		wholeWant := errorWire(refErr)
		if refErr == nil {
			wholeWant += ":" + bytesWire(refWhole[:])
		}
		oracle := strings.Join([]string{bytesWire(refCells[:]), errorWire(refPart), wholeWant, splitOracle(text)}, "|")
		if got != want || got != oracle {
			t.Fatalf("row %d: native %q Lean %q oracle %q", rows, got, want, oracle)
		}
		rows++
	}
	if err := in.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("no rows")
	}
	t.Logf("%d pinned native cases per declaration", rows)
}

func bytesWire(raw []byte) string {
	s := make([]string, len(raw))
	for i, b := range raw {
		s[i] = strconv.Itoa(int(b))
	}
	return strings.Join(s, ",")
}

func errorWire(err error) string {
	if err == nil {
		return "ok"
	}
	e, ok := err.(parseAddrError)
	if !ok {
		panic("unexpected parser error type")
	}
	return e.msg + ":" + bytesWire([]byte(e.in)) + ":" + bytesWire([]byte(e.at))
}

func splitWire(ip, port string, bracketed bool, err error) string {
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%s:%s:%t", bytesWire([]byte(ip)), bytesWire([]byte(port)), bracketed)
}

func fieldsOracle(text string, lo, hi int, cells [4]byte) ([4]byte, error) {
	s := text[lo:hi]
	segments := strings.Split(s, ".")
	offset := 0
	for i, octet := range segments {
		start := offset
		fault := func(reason, at string) ([4]byte, error) { return cells, parseAddrError{in: text, msg: reason, at: at} }
		for j, c := range []byte(octet) {
			if c < '0' || c > '9' {
				return fault("unexpected character", s[start+j:])
			}
			if j == 1 && octet[0] == '0' {
				return fault("IPv4 field has octet with leading zero", "")
			}
			if j == 3 || j == 2 && octet[:3] > "255" {
				return fault("IPv4 field has value >255", "")
			}
		}
		offset += len(octet)
		if i < len(segments)-1 {
			if octet == "" || offset == len(s)-1 {
				return fault("IPv4 field must have at least one digit", s[offset:])
			}
			if i == 3 {
				return fault("IPv4 address too long", "")
			}
		} else if i < 3 {
			return fault("IPv4 address too short", "")
		}
		value, err := strconv.Atoi(octet)
		if err != nil {
			panic(err)
		}
		cells[i] = byte(value)
		offset++
	}
	return cells, nil
}

func splitOracle(text string) string {
	pieces := strings.Split(text, ":")
	if len(pieces) < 2 {
		return "not an ip:port"
	}
	host := strings.Join(pieces[:len(pieces)-1], ":")
	port := pieces[len(pieces)-1]
	if host == "" {
		return "no IP"
	}
	if port == "" {
		return "no port"
	}
	unwrapped, bracketed := strings.CutPrefix(host, "[")
	if bracketed {
		var closed bool
		host, closed = strings.CutSuffix(unwrapped, "]")
		if !closed {
			return "missing ]"
		}
	}
	return splitWire(host, port, bracketed, nil)
}
