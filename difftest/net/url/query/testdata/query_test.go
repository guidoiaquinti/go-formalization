package url_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"runtime/metrics"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanURLQuery(t *testing.T) {
	t.Setenv("GODEBUG", "")
	file, err := os.Open(os.Getenv("GOLEAN_URL_QUERY_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	lastSetting := ""
	counter := []metrics.Sample{{Name: "/godebug/non-default-behavior/urlmaxqueryparams:events"}}
	readMetric := func() uint64 {
		metrics.Read(counter)
		return counter[0].Value.Uint64()
	}
	for i := 0; ; i++ {
		var c struct {
			Config        string
			Raw           []byte
			Params, Start int
			Present       bool
			Want          string
		}
		if err := decoder.Decode(&c); errors.Is(err, io.EOF) {
			if i == 0 {
				t.Fatal("missing query fixtures")
			}
			t.Logf("four pinned query entries matched %d rows", i)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		if c.Config != lastSetting {
			if err := os.Setenv("GODEBUG", "urlmaxqueryparams="+c.Config); err != nil {
				t.Fatal(err)
			}
			lastSetting = c.Config
		}
		m := url.Values(nil)
		if c.Start == 1 {
			m = url.Values{}
		}
		if c.Start == 2 {
			m = url.Values{"a": {"old"}, "z": nil}
		}
		before := readMetric()
		allowed := url.GoLeanParamsWithinMax(c.Params)
		delta := readMetric() - before
		if delta > 1 {
			t.Fatalf("metric increment: %d", delta)
		}
		query := string(c.Raw)
		private := catchQuery("assignment to entry in nil map", func() string {
			return queryReport(m, url.GoLeanParseQuery(m, query))
		})
		public := queryReport(url.ParseQuery(query))
		var u *url.URL
		if c.Present {
			u = &url.URL{RawQuery: query}
		}
		method := catchQuery("runtime error: invalid memory address or nil pointer dereference", func() string {
			return queryState(u.Query())
		})
		got := fmt.Sprintf("%t:%t|%s|%s|%s", allowed, delta == 1, private, public, method)
		if got != c.Want {
			t.Fatalf("row %d %+v: got %q want %q", i, c, got, c.Want)
		}
	}
}

func queryBytes(s string) string {
	fields := make([]string, len(s))
	for i := range fields {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return fmt.Sprintf("%d:%s", len(s), strings.Join(fields, ","))
}

func queryState(m url.Values) string {
	if m == nil {
		return "nil"
	}
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	var entries []string
	for _, key := range keys {
		words := m[key]
		payload := "nil"
		if words != nil {
			items := make([]string, len(words))
			for i, w := range words {
				items[i] = queryBytes(w)
			}
			payload = "list:" + strings.Join(items, ";")
		}
		entries = append(entries, queryBytes(key)+"="+payload)
	}
	return "map:" + strings.Join(entries, "/")
}

func queryReport(m url.Values, err error) string {
	status := "ok"
	if e, ok := err.(url.EscapeError); ok {
		status = "escape:" + queryBytes(string(e))
	} else if err != nil {
		switch err.Error() {
		case "number of URL query parameters exceeded limit":
			status = "limit"
		case "invalid semicolon separator in query":
			status = "semicolon"
		default:
			panic(err)
		}
	}
	return status + "/" + queryState(m)
}

func catchQuery(reason string, action func() string) (out string) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(runtime.Error)
			if !ok || err.Error() != reason {
				panic(p)
			}
			out = "panic"
		}
	}()
	return action()
}
