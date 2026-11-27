package netip

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func goLeanAddrTextCSV(raw []byte) string {
	fields := make([]string, len(raw))
	for i, v := range raw {
		fields[i] = strconv.Itoa(int(v))
	}
	return strings.Join(fields, ",")
}
func goLeanAddrTextAppend(reused bool, out, old []byte) string {
	return fmt.Sprintf("%t:%t:%s:%s", reused, out == nil, goLeanAddrTextCSV(out), goLeanAddrTextCSV(old))
}
func goLeanAddrTextReference(family int, hi, lo uint64, zone []byte) (string, string, []string, int) {
	if family == 0 {
		return "invalid IP", "invalid IP", nil, 0
	}
	var words [8]uint16
	var padded []string
	for i := range words {
		word := hi
		if i >= 4 {
			word = lo
		}
		words[i] = uint16(word >> (16 * (3 - i%4)))
		padded = append(padded, fmt.Sprintf("%04x", words[i]))
	}
	zoneText := ""
	if len(zone) > 0 {
		zoneText = "%" + string(zone)
	}
	if family == 4 || hi == 0 && lo>>32 == 65535 {
		ipv4 := fmt.Sprintf("%d.%d.%d.%d", byte(lo>>24), byte(lo>>16), byte(lo>>8), byte(lo))
		chunks := make([]string, 0, len(ipv4)+3)
		profile := 1
		if family == 6 {
			chunks = append(chunks, "::ffff:")
			profile = 2
		}
		for i := range len(ipv4) {
			chunks = append(chunks, ipv4[i:i+1])
		}
		if zoneText != "" {
			chunks = append(chunks, "%", string(zone))
		}
		expanded := ipv4
		if family == 6 {
			expanded = strings.Join(padded, ":") + zoneText
		}
		return strings.Join(chunks, ""), expanded, chunks, profile
	}
	bestStart, bestEnd := 0, 0
	for hi := 8; hi >= 2; hi-- {
		for lo := hi - 2; lo >= 0; lo-- {
			if hi-lo < bestEnd-bestStart {
				continue
			}
			if !slices.ContainsFunc(words[lo:hi], func(w uint16) bool { return w != 0 }) {
				if hi-lo > bestEnd-bestStart || lo < bestStart {
					bestStart, bestEnd = lo, hi
				}
			}
		}
	}
	var chunks []string
	for i := 0; i < 8; i++ {
		if i == bestStart && bestEnd > bestStart {
			chunks = append(chunks, "::")
			i = bestEnd
			if i == 8 {
				break
			}
		} else if i > 0 {
			chunks = append(chunks, ":")
		}
		hex := strconv.FormatUint(uint64(words[i]), 16)
		for j := range len(hex) {
			chunks = append(chunks, hex[j:j+1])
		}
	}
	if zoneText != "" {
		chunks = append(chunks, "%", string(zone))
	}
	return strings.Join(chunks, ""), strings.Join(padded, ":") + zoneText, chunks, 3
}
func TestGoLeanAddrText(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_NETIP_ADDR_TEXT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var counts [4]int
	for scanner.Scan() {
		var c struct {
			Family        int
			Hi, Lo        uint64
			Zone, Backing []byte
			Length        int
			Want          string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		var ip Addr
		switch c.Family {
		case 4:
			var raw [4]byte
			binary.BigEndian.PutUint32(raw[:], uint32(c.Lo))
			ip = AddrFrom4(raw)
		case 6:
			var raw [16]byte
			binary.BigEndian.PutUint64(raw[:8], c.Hi)
			binary.BigEndian.PutUint64(raw[8:], c.Lo)
			ip = AddrFrom16(raw).WithZone(string(c.Zone))
		}
		text, expanded, parts, profile := goLeanAddrTextReference(c.Family, c.Hi, c.Lo, c.Zone)
		encoded := strings.Join(parts, "")
		native := func(fn func([]byte) []byte) string {
			old := slices.Clone(c.Backing)
			out := fn(old[:c.Length:len(old)])
			same := len(out) == 0 || len(old) > 0 && &out[0] == &old[0]
			return goLeanAddrTextAppend(same, out, old)
		}
		old := slices.Clone(c.Backing)
		out := slices.Clone(c.Backing[:c.Length])
		reuse, cursor := true, c.Length
		for _, part := range parts {
			out = append(out, part...)
			if reuse && cursor+len(part) <= len(old) {
				copy(old[cursor:], part)
				cursor += len(part)
			} else {
				reuse = false
			}
		}
		expectedAppend := goLeanAddrTextAppend(reuse, out, old)
		direct, wantDirect := "skip", "skip"
		if profile > 0 {
			wantDirect = expectedAppend
			switch profile {
			case 1:
				direct = native(ip.appendTo4)
			case 2:
				direct = native(ip.appendTo4In6)
			case 3:
				direct = native(ip.appendTo6)
			}
		}
		marshaled, err := ip.MarshalText()
		if err != nil || marshaled == nil {
			t.Fatalf("unexpected MarshalText %s %v", marshaled, err)
		}
		actual := strings.Join([]string{
			goLeanAddrTextCSV([]byte(ip.String())), goLeanAddrTextCSV([]byte(ip.StringExpanded())), goLeanAddrTextCSV(marshaled),
			native(ip.AppendTo), native(func(b []byte) []byte {
				ret, err := ip.AppendText(b)
				if err != nil {
					t.Fatal(err)
				}
				return ret
			}), direct,
		}, "|")
		want := strings.Join([]string{
			goLeanAddrTextCSV([]byte(text)), goLeanAddrTextCSV([]byte(expanded)), goLeanAddrTextCSV([]byte(encoded)),
			expectedAppend, expectedAppend, wantDirect,
		}, "|")
		if actual != want || actual != c.Want {
			t.Fatalf("%+v: native %s oracle %s", c, actual, want)
		}
		counts[profile]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("native IP text rows by profile: %v", counts)
}
