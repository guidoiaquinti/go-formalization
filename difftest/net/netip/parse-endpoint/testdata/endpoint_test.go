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

func endpointBytes(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}

func endpointValue(a Addr, number int) string {
	if a.Is4() {
		raw := a.As4()
		return fmt.Sprintf("v4:%s::%d", endpointBytes(raw[:]), number)
	}
	if a.Is6() {
		raw := a.As16()
		return fmt.Sprintf("v6:%s:%s:%d", endpointBytes(raw[:]), endpointBytes([]byte(a.Zone())), number)
	}
	panic("invalid address on success")
}

func endpointFault(err error) string {
	switch e := err.(type) {
	case parseAddrError:
		return "address fault"
	case parsePrefixError:
		switch {
		case e.msg == "no '/'":
			return "no slash"
		case e.msg == "IPv6 zones cannot be present in a prefix":
			return "zone"
		case e.msg == "prefix length out of range":
			return "out of range"
		case strings.HasPrefix(e.msg, "bad bits after slash: "):
			return "bad bits"
		case strings.HasPrefix(e.msg, "ParseAddr("):
			return "address fault"
		}
	default:
		msg := err.Error()
		switch {
		case msg == "not an ip:port" || msg == "no IP" || msg == "no port" || msg == "missing ]":
			return msg
		case strings.HasPrefix(msg, "invalid port "):
			return "invalid port"
		case strings.HasSuffix(msg, ", square brackets can only be used with IPv6 addresses"):
			return "IPv4 bracketed"
		case strings.HasSuffix(msg, ", IPv6 addresses must be surrounded by square brackets"):
			return "IPv6 unbracketed"
		}
	}
	panic(fmt.Sprintf("unclassified error %T %q", err, err))
}

func endpointMust(parse func() string) (result string) {
	defer func() {
		if failure := recover(); failure != nil {
			result = "panic:" + endpointFault(failure.(error))
		}
	}()
	return parse()
}

func endpointPort(v AddrPort, err error, clear bool) string {
	if err != nil || clear {
		if v != (AddrPort{}) {
			panic("nonzero AddrPort on error or clear")
		}
		if err == nil {
			return "clear"
		}
		return "error:" + endpointFault(err)
	}
	return endpointValue(v.Addr(), int(v.Port()))
}

func endpointPrefix(v Prefix, err error, clear bool) string {
	if err != nil || clear {
		if v != (Prefix{}) {
			panic("nonzero Prefix on error or clear")
		}
		if err == nil {
			return "clear"
		}
		return "error:" + endpointFault(err)
	}
	return endpointValue(v.Addr(), v.Bits())
}

func TestGoLeanEndpoint(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_NETIP_ENDPOINT"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	for i := 0; ; i++ {
		var c struct {
			Raw                                []byte
			Want, ExpectedPort, ExpectedPrefix string
		}
		if err := decoder.Decode(&c); err != nil {
			if err == io.EOF && i != 0 {
				t.Logf("%d pinned rows per declaration", i)
				return
			}
			t.Fatal(err)
		}
		raw := string(c.Raw)
		port, portErr := ParseAddrPort(raw)
		prefix, prefixErr := ParsePrefix(raw)
		oldPort := AddrPortFrom(IPv6Loopback(), 1234)
		oldPrefix := PrefixFrom(IPv6Loopback(), 127)
		portText := oldPort.UnmarshalText(c.Raw)
		prefixText := oldPrefix.UnmarshalText(c.Raw)
		portWire := endpointPort(port, portErr, false)
		prefixWire := endpointPrefix(prefix, prefixErr, false)
		got := strings.Join([]string{
			portWire, endpointMust(func() string { return endpointPort(MustParseAddrPort(raw), nil, false) }),
			endpointPort(oldPort, portText, raw == ""),
			prefixWire, endpointMust(func() string { return endpointPrefix(MustParsePrefix(raw), nil, false) }),
			endpointPrefix(oldPrefix, prefixText, raw == ""),
		}, "|")
		if got != c.Want {
			t.Fatalf("%q: Go %q, Lean %q", raw, got, c.Want)
		}
		if c.ExpectedPort != "" && portWire != c.ExpectedPort || c.ExpectedPrefix != "" && prefixWire != c.ExpectedPrefix {
			t.Fatalf("%q differs from expected address bytes, zone or numeric value", raw)
		}
	}
}
