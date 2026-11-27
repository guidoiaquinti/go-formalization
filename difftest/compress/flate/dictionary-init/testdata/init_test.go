package flate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func dictionaryBytes(cells []byte) string {
	result := make([]string, len(cells))
	for i, b := range cells {
		result[i] = strconv.Itoa(int(b))
	}
	return strings.Join(result, ",")
}

func TestGoLeanDictionaryInit(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_DICT_INIT"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := json.NewDecoder(f)
	for i := 0; ; i++ {
		var c struct {
			Old, Preset                 []byte
			Length, Read, Written, Size int
			Full                        bool
			Want                        string
		}
		err := input.Decode(&c)
		if err == io.EOF && i > 0 {
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		original := slices.Clone(c.Preset)
		d := dictDecoder{hist: c.Old[:c.Length:len(c.Old)], wrPos: c.Written, rdPos: c.Read, full: c.Full}
		d.init(c.Size, c.Preset)
		next := d.hist[:cap(d.hist)]
		reused := len(c.Old) == 0 && cap(d.hist) == 0
		if len(c.Old) > 0 && len(next) > 0 {
			reused = &c.Old[0] == &next[0]
		}
		got := fmt.Sprintf("%t|%d|%d|%d|%d|%t|%t|%s|%s", d.hist != nil, len(d.hist), cap(d.hist), d.rdPos, d.wrPos, d.full,
			reused, dictionaryBytes(next), dictionaryBytes(c.Old))
		if got != c.Want || !slices.Equal(original, c.Preset) {
			t.Fatalf("row %d: Go %q want %q", i, got, c.Want)
		}
	}
}
