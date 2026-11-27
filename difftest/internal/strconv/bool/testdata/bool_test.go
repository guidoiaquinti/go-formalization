package strconv_test

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	conv "internal/strconv"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanInternalBool(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_INTERNAL_STRCONV_BOOL_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := gob.NewDecoder(bufio.NewReader(f))
	parse := map[string]bool{
		"1": true, "t": true, "T": true, "true": true, "TRUE": true, "True": true,
		"0": false, "f": false, "F": false, "false": false, "FALSE": false, "False": false,
	}
	labels := map[int]string{1: "value out of range", 2: "invalid syntax", 3: "invalid base", 4: "invalid bit size"}
	for row := 0; ; row++ {
		var c struct {
			Bool, Nil   bool
			Tag, Length int
			Backing     []byte
			Raw, Want   string
		}
		if err := in.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d observations for all four pinned internal strconv declarations", row)
				return
			}
			t.Fatal(err)
		}
		suffix := "false"
		if c.Bool {
			suffix = "true"
		}
		value, known := parse[c.Raw]
		label, knownTag := labels[c.Tag]
		if !knownTag {
			label = "unknown error"
		}
		requireReuse := c.Length+len(suffix) <= len(c.Backing)
		visible := string(c.Backing[:c.Length]) + suffix
		old := bytes.Clone(c.Backing)
		if requireReuse {
			copy(old, visible)
		}
		ref := fmt.Sprintf("%t|%t|%s|%s|%s|%t|%s", value, !known, suffix, boolWire([]byte(visible)), boolWire(old), requireReuse, label)
		var backing []byte
		if !c.Nil {
			backing = make([]byte, len(c.Backing))
			copy(backing, c.Backing)
		}
		gotValue, failure := conv.ParseBool(c.Raw)
		if failure != nil && failure != conv.ErrSyntax {
			t.Fatalf("row %d: unexpected internal error type %T", row, failure)
		}
		appended := conv.AppendBool(backing[:c.Length:len(backing)], c.Bool)
		reused := len(backing) > 0 && &appended[0] == &backing[0]
		got := fmt.Sprintf("%t|%t|%s|%s|%s|%t|%s", gotValue, failure != nil, conv.FormatBool(c.Bool),
			boolWire(appended), boolWire(backing), reused, conv.Error(c.Tag).Error())
		if got != ref || got != c.Want {
			t.Fatalf("row %d: Go %q, Lean %q, independent oracle %q", row, got, c.Want, ref)
		}
	}
}

func boolWire(b []byte) string {
	fields := make([]string, len(b))
	for i, v := range b {
		fields[i] = strconv.Itoa(int(v))
	}
	return strings.Join(fields, ",")
}
