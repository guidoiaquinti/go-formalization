package strings_test

import (
	"slices"
	"strings"
	"testing"
)

func observe(fn func(string, []byte, bool) string, input string, table []byte, nilTable bool) (out string, panicked bool) {
	defer func() { panicked = recover() != nil }()
	return fn(input, table, nilTable), false
}

func asciiOnly(input string, table []byte, nilTable bool) string {
	if nilTable && input != "" {
		panic("nil table")
	}
	for len(input) != 0 && input[len(input)-1] < 128 && slices.Contains(table, input[len(input)-1]) {
		input = input[:len(input)-1]
	}
	return input
}

func TestGoLeanTrimRightASCIIReview(t *testing.T) {
	fixtures := []struct {
		input    string
		table    []byte
		nilTable bool
		want     string
		panic    bool
	}{
		{"", nil, true, "", false},
		{"\xff", nil, true, "", true},
		{"a", nil, true, "", true},
		{"\x00\xff", nil, false, "\x00\xff", false},
		{"aaXaa\x00", []byte{0, 'a'}, false, "aaX", false},
		{"aaaa", []byte{'a'}, false, "", false},
		{"\xffé\xff", []byte{0xff, 0xa9}, false, "\xff\xc3", false},
		{"\x80x\xe2\x80", []byte{0x80, 0xe2}, false, "\x80x", false},
		{"é", []byte{0xc3}, false, "é", false},
		{"\xff", []byte{0xff}, false, "", false},
	}
	for _, c := range []struct {
		name string
		trim func(string, []byte, bool) string
		want int
	}{
		{"Go", strings.GoLeanTrimRightASCII, 0},
		{"Restrict table to ASCII", asciiOnly, 3},
		{"Nil table acts like empty table", func(input string, table []byte, _ bool) string {
			return strings.GoLeanTrimRightASCII(input, table, false)
		}, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			got, panicked := observe(c.trim, f.input, f.table, f.nilTable)
			if got != f.want || panicked != f.panic {
				bad++
				if c.want == 0 {
					t.Errorf("trimRightASCII(%x,%x,nil=%t) gave %x panic=%t, want %x panic=%t",
						f.input, f.table, f.nilTable, got, panicked, f.want, f.panic)
				}
			}
		}
		t.Logf("%s disagreed on %d of %d", c.name, bad, len(fixtures))
		if bad != c.want {
			t.Errorf("%s wanted %d differences", c.name, c.want)
		}
	}
}
