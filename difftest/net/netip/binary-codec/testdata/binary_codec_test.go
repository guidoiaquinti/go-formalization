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

type goLeanBinaryIP struct {
	family int
	hi, lo uint64
	zone   string
}

func (i goLeanBinaryIP) wire() string {
	return fmt.Sprintf("%d:%d:%d:%s", i.family, i.hi, i.lo, goLeanBinaryCSV([]byte(i.zone)))
}
func (i goLeanBinaryIP) chunks() [][]byte {
	switch i.family {
	case 4:
		return [][]byte{binary.BigEndian.AppendUint32(nil, uint32(i.lo))}
	case 6:
		return [][]byte{binary.BigEndian.AppendUint64(nil, i.hi), binary.BigEndian.AppendUint64(nil, i.lo), []byte(i.zone)}
	default:
		return nil
	}
}
func (i goLeanBinaryIP) native() Addr {
	switch i.family {
	case 4:
		var ipv4 [4]byte
		binary.BigEndian.PutUint32(ipv4[:], uint32(i.lo))
		return AddrFrom4(ipv4)
	case 6:
		var ipv6 [16]byte
		binary.BigEndian.PutUint64(ipv6[:8], i.hi)
		binary.BigEndian.PutUint64(ipv6[8:], i.lo)
		return AddrFrom16(ipv6).WithZone(i.zone)
	default:
		return Addr{}
	}
}
func goLeanGetBinaryIP(ip Addr) goLeanBinaryIP {
	f := 0
	if ip.Is4() {
		f = 4
	}
	if ip.Is6() {
		f = 6
	}
	return goLeanBinaryIP{f, ip.addr.hi, ip.addr.lo, ip.Zone()}
}
func goLeanDecodeBinaryIP(raw []byte) (goLeanBinaryIP, bool) {
	switch len(raw) {
	case 0:
		return goLeanBinaryIP{}, true
	case 4:
		return goLeanBinaryIP{family: 4, lo: 0xffff00000000 | uint64(binary.BigEndian.Uint32(raw))}, true
	}
	if len(raw) < 16 {
		return goLeanBinaryIP{}, false
	}
	return goLeanBinaryIP{6, binary.BigEndian.Uint64(raw), binary.BigEndian.Uint64(raw[8:]), string(raw[16:])}, true
}
func goLeanBinaryCSV(raw []byte) string {
	fields := make([]string, len(raw))
	for i, v := range raw {
		fields[i] = strconv.Itoa(int(v))
	}
	return strings.Join(fields, ",")
}
func goLeanBinaryAppendWire(reused bool, output, original []byte) string {
	return fmt.Sprintf("%t:%t:%s:%s", reused, output == nil, goLeanBinaryCSV(output), goLeanBinaryCSV(original))
}
func goLeanNativeBinaryAppend(t *testing.T, backing []byte, n int, fn func([]byte) ([]byte, error)) string {
	t.Helper()
	cells := slices.Clone(backing)
	result, err := fn(cells[:n:len(cells)])
	if err != nil {
		t.Fatal(err)
	}
	reused := len(result) == 0 || len(cells) > 0 && &result[0] == &cells[0]
	return goLeanBinaryAppendWire(reused, result, cells)
}
func goLeanBinaryAppendOracle(backing []byte, n int, parts [][]byte) string {
	old := slices.Clone(backing)
	output := slices.Clone(backing[:n])
	cursor := n
	reused := true
	for _, part := range parts {
		output = append(output, part...)
		if reused && len(part) <= len(old)-cursor {
			copy(old[cursor:], part)
			cursor += len(part)
		} else {
			reused = false
		}
	}
	return goLeanBinaryAppendWire(reused, output, old)
}
func goLeanBinaryMarshal(t *testing.T, fn func() ([]byte, error)) string {
	t.Helper()
	out, err := fn()
	if err != nil || out == nil || len(out) != cap(out) {
		t.Fatalf("unexpected marshal geometry: %v len=%d cap=%d err=%v", out, len(out), cap(out), err)
	}
	return goLeanBinaryCSV(out)
}
func goLeanBinaryError(t *testing.T, err error) bool {
	t.Helper()
	if err == nil {
		return true
	}
	if err.Error() != "unexpected slice size" {
		t.Fatalf("unexpected decode error: %v", err)
	}
	return false
}

func TestGoLeanBinaryCodec(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_NETIP_BINARY_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		var c struct {
			Family                 int
			Hi, Lo                 uint64
			Port                   uint16
			Bits, Length           int
			Zone, Backing, Payload []byte
			Want                   string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		input := goLeanBinaryIP{c.Family, c.Hi, c.Lo, string(c.Zone)}
		ip := input.native()
		port := AddrPortFrom(ip, c.Port)
		prefix := PrefixFrom(ip, c.Bits)
		addressParts := input.chunks()
		portParts := append(input.chunks(), binary.LittleEndian.AppendUint16(nil, c.Port))
		stripped := input
		stripped.zone = ""
		prefixParts := append(stripped.chunks(), []byte{byte(c.Bits)})
		a, okA := goLeanDecodeBinaryIP(c.Payload)
		if !okA {
			a = input
		}
		p, portNum, okP := input, c.Port, false
		if n := len(c.Payload); n >= 2 {
			if ip, ok := goLeanDecodeBinaryIP(c.Payload[:n-2]); ok {
				p, portNum, okP = ip, binary.LittleEndian.Uint16(c.Payload[n-2:]), true
			}
		}
		q, bits, okQ := stripped, c.Bits, false
		if n := len(c.Payload); n >= 1 {
			if ip, ok := goLeanDecodeBinaryIP(c.Payload[:n-1]); ok {
				q, bits, okQ = ip, int(c.Payload[n-1]), true
				q.zone = ""
				if q.family == 0 || q.family == 4 && bits > 32 || q.family == 6 && bits > 128 {
					bits = -1
				}
			}
		}
		want := strings.Join([]string{
			strconv.Itoa(len(slices.Concat(addressParts...))), goLeanBinaryCSV(slices.Concat(addressParts...)), goLeanBinaryAppendOracle(c.Backing, c.Length, addressParts),
			fmt.Sprintf("%t:%s", okA, a.wire()), goLeanBinaryCSV(slices.Concat(portParts...)), goLeanBinaryAppendOracle(c.Backing, c.Length, portParts),
			fmt.Sprintf("%t:%s:%d", okP, p.wire(), portNum), goLeanBinaryCSV(slices.Concat(prefixParts...)), goLeanBinaryAppendOracle(c.Backing, c.Length, prefixParts),
			fmt.Sprintf("%t:%s:%d", okQ, q.wire(), bits),
		}, "|")
		gotA, gotP, gotQ := ip, port, prefix
		payload := slices.Clone(c.Payload)
		hasA, hasP, hasQ := goLeanBinaryError(t, gotA.UnmarshalBinary(payload)), goLeanBinaryError(t, gotP.UnmarshalBinary(payload)), goLeanBinaryError(t, gotQ.UnmarshalBinary(payload))
		if !slices.Equal(payload, c.Payload) {
			t.Fatal("input mutated before return")
		}
		// Check that decoded zones retain their own content after the payload changes.
		clear(payload)
		got := strings.Join([]string{
			strconv.Itoa(ip.marshalBinarySize()), goLeanBinaryMarshal(t, ip.MarshalBinary), goLeanNativeBinaryAppend(t, c.Backing, c.Length, ip.AppendBinary),
			fmt.Sprintf("%t:%s", hasA, goLeanGetBinaryIP(gotA).wire()), goLeanBinaryMarshal(t, port.MarshalBinary), goLeanNativeBinaryAppend(t, c.Backing, c.Length, port.AppendBinary),
			fmt.Sprintf("%t:%s:%d", hasP, goLeanGetBinaryIP(gotP.Addr()).wire(), gotP.Port()), goLeanBinaryMarshal(t, prefix.MarshalBinary), goLeanNativeBinaryAppend(t, c.Backing, c.Length, prefix.AppendBinary),
			fmt.Sprintf("%t:%s:%d", hasQ, goLeanGetBinaryIP(gotQ.Addr()).wire(), gotQ.Bits()),
		}, "|")
		if got != want || got != c.Want {
			t.Fatalf("%+v\nnative %s\noracle %s", c, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("netip binary codec: %d native rows per declaration", count)
}
