package url_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanURLLexical(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_URL_LEXICAL_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	wire := func(raw string) string {
		fields := make([]string, len(raw))
		for i := range fields {
			fields[i] = strconv.Itoa(int(raw[i]))
		}
		return fmt.Sprintf("%d:%s", len(raw), strings.Join(fields, ","))
	}
	for i := 0; ; i++ {
		var c struct {
			Present bool
			Raw     []byte
			Want    string
		}
		if err := decoder.Decode(&c); errors.Is(err, io.EOF) {
			if i == 0 {
				t.Fatal("missing lexical fixtures")
			}
			t.Logf("seven pinned URL lexical declarations matched %d rows", i)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		raw := string(c.Raw)
		host, port := url.GoLeanSplitHostPort(raw)
		prefix, suffix, schemeErr := url.GoLeanGetScheme(raw)
		scheme := wire(prefix) + "/" + wire(suffix)
		if schemeErr != nil {
			if schemeErr.Error() != "missing protocol scheme" || prefix != "" || suffix != "" {
				t.Fatalf("unexpected scheme failure %q %q %v", prefix, suffix, schemeErr)
			}
			scheme = "missing"
		}
		var u *url.URL
		if c.Present {
			u = &url.URL{Host: raw, Scheme: raw}
		}
		h := catchNil(func() string { return wire(u.Hostname()) })
		p := catchNil(func() string { return wire(u.Port()) })
		a := catchNil(func() string { return strconv.FormatBool(u.IsAbs()) })
		got := strings.Join([]string{strconv.FormatBool(url.GoLeanValidOptionalPort(raw)), wire(host) + "/" + wire(port),
			h, p, a, scheme, strconv.FormatBool(url.GoLeanContainsCTL(raw))}, "|")
		if got != c.Want {
			t.Fatalf("row %d input %q got %q want %q", i, raw, got, c.Want)
		}
	}
}

func catchNil(action func() string) (result string) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(p)
			}
			result = "panic"
		}
	}()
	return action()
}
