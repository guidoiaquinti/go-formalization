package url_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanEscape(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_URL_ESCAPE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	wire := func(raw string) string {
		fields := make([]string, len(raw))
		for i := range len(raw) {
			fields[i] = strconv.Itoa(int(raw[i]))
		}
		return strings.Join(fields, ",")
	}
	for i := 0; ; i++ {
		var c struct {
			Raw  []byte
			Mode byte
			Want string
		}
		if err := decoder.Decode(&c); errors.Is(err, io.EOF) {
			if i == 0 {
				t.Fatal("missing escape fixtures")
			}
			t.Logf("all six pinned escape entries matched %d rows", i)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		fields := make([]string, len(c.Raw))
		for j, b := range c.Raw {
			escape, hex, unhex := url.GoLeanEscapeClassification(b, c.Mode)
			fields[j] = fmt.Sprintf("%t:%t:%d", escape, hex, unhex)
		}
		raw := string(c.Raw)
		got := strings.Join([]string{strings.Join(fields, ","), wire(url.GoLeanEscape(raw, c.Mode)), wire(url.QueryEscape(raw)), wire(url.PathEscape(raw))}, "|")
		if got != c.Want {
			t.Fatalf("row %d (%+v): got %q want %q", i, c, got, c.Want)
		}
	}
}
