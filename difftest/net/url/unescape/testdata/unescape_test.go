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

func TestGoLeanUnescape(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_URL_UNESCAPE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	wire := func(output string, err error) string {
		raw, status := output, "ok"
		switch e := err.(type) {
		case nil:
		case url.EscapeError:
			status, raw = "escape", string(e)
		case url.InvalidHostError:
			status, raw = "host", string(e)
		default:
			t.Fatalf("unexpected error %T", err)
		}
		if err != nil && output != "" {
			t.Fatal("unescape returned partial output")
		}
		fields := make([]string, len(raw))
		for i := range fields {
			fields[i] = strconv.Itoa(int(raw[i]))
		}
		return fmt.Sprintf("%s:%s", status, strings.Join(fields, ","))
	}
	for i := 0; ; i++ {
		var row struct {
			Raw  []byte
			Mode byte
			Want string
		}
		if err := decoder.Decode(&row); errors.Is(err, io.EOF) {
			if i == 0 {
				t.Fatal("missing unescape fixtures")
			}
			t.Logf("three pinned unescape entries matched %d rows", i)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		raw := string(row.Raw)
		got := strings.Join([]string{wire(url.GoLeanUnescape(raw, row.Mode)), wire(url.QueryUnescape(raw)), wire(url.PathUnescape(raw))}, "|")
		if got != row.Want {
			t.Fatalf("row %d %+v: got %q want %q", i, row, got, row.Want)
		}
	}
}
