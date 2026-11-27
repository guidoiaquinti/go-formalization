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

func TestGoLeanPathFragment(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_URL_PATH_FRAGMENT_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for i := 0; ; i++ {
		var c struct {
			Present              bool
			Mode                 byte
			Decoded, Raw, Update []byte
			Want                 string
		}
		if err := decoder.Decode(&c); errors.Is(err, io.EOF) {
			if i == 0 {
				t.Fatal("missing path and fragment fixtures")
			}
			t.Logf("all five pinned URL field declarations matched %d rows", i)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		decoded, raw, update := string(c.Decoded), string(c.Raw), string(c.Update)
		state := func() *url.URL {
			if !c.Present {
				return nil
			}
			return &url.URL{Scheme: "s", Opaque: "o", User: url.UserPassword("u", "p"), Host: "h",
				Path: decoded, RawPath: raw, Fragment: decoded, RawFragment: raw, RawQuery: "a=b", ForceQuery: true, OmitHost: true}
		}
		p := state()
		f := state()
		unchanged := func(u *url.URL, read func(*url.URL) string) string {
			return nativeNoNil(func() string {
				var before url.URL
				if u != nil {
					before = *u
				}
				output := read(u)
				if *u != before {
					t.Fatal("escaped accessor changed URL fields")
				}
				return "ok:" + nativePartBytes(output)
			})
		}
		ep := unchanged(p, (*url.URL).EscapedPath)
		ef := unchanged(f, (*url.URL).EscapedFragment)
		sp := nativeNoNil(func() string {
			var before url.URL
			if p != nil {
				before = *p
			}
			err := url.GoLeanSetPath(p, update)
			if p == nil {
				return nativeSetterStatus(err) + "/nil"
			}
			before.Path, before.RawPath = p.Path, p.RawPath
			if *p != before {
				t.Fatal("setPath changed unrelated fields")
			}
			return nativeSetterStatus(err) + "/" + nativePartBytes(p.Path) + "/" + nativePartBytes(p.RawPath)
		})
		sf := nativeNoNil(func() string {
			var before url.URL
			if f != nil {
				before = *f
			}
			err := url.GoLeanSetFragment(f, update)
			if f == nil {
				return nativeSetterStatus(err) + "/nil"
			}
			before.Fragment, before.RawFragment = f.Fragment, f.RawFragment
			if *f != before {
				t.Fatal("setFragment changed unrelated fields")
			}
			return nativeSetterStatus(err) + "/" + nativePartBytes(f.Fragment) + "/" + nativePartBytes(f.RawFragment)
		})
		got := strings.Join([]string{strconv.FormatBool(url.GoLeanValidEncoded(raw, c.Mode)), ep, ef, sp, sf}, "|")
		if got != c.Want {
			t.Fatalf("row %d mode %d (%q,%q,%q) got %q want %q", i, c.Mode, decoded, raw, update, got, c.Want)
		}
	}
}

func nativePartBytes(raw string) string {
	fields := make([]string, len(raw))
	for i := range fields {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(fields, ",")
}

func nativeSetterStatus(err error) string {
	switch e := err.(type) {
	case nil:
		return "ok"
	case url.EscapeError:
		return "escape:" + nativePartBytes(string(e))
	default:
		panic(fmt.Sprintf("unexpected setter error %T", err))
	}
}

func nativeNoNil(action func() string) (result string) {
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
