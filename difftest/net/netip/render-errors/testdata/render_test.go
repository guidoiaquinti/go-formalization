package netip

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func renderWire(s string) string {
	out := make([]string, len(s))
	for i := range len(s) {
		out[i] = strconv.Itoa(int(s[i]))
	}
	return strings.Join(out, ",")
}

func TestGoLeanRenderErrors(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_NETIP_RENDER_ERRORS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	for i := 0; ; i++ {
		var row struct {
			Input, Message, Suffix []byte
			Want                   string
		}
		if err := decoder.Decode(&row); err != nil {
			if err == io.EOF && i != 0 {
				t.Logf("%d pinned rows per declaration", i)
				return
			}
			t.Fatal(err)
		}
		input, msg, suffix := string(row.Input), string(row.Message), string(row.Suffix)
		address := parseAddrError{in: input, msg: msg, at: suffix}
		prefix := parsePrefixError{in: input, msg: msg}
		got := renderWire(address.Error()) + "|" + renderWire(prefix.Error())
		if got != row.Want {
			t.Fatalf("row %d (%q, %q, %q): Go %q, Lean %q", i, input, msg, suffix, got, row.Want)
		}
	}
}
