package main

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

func replaceEveryInvalidByte(input, replacement string) string {
	var output strings.Builder
	for offset, r := range input {
		if r == utf8.RuneError && !strings.HasPrefix(input[offset:], "�") {
			output.WriteString(replacement)
		} else {
			output.WriteRune(r)
		}
	}
	return output.String()
}

func replaceValidRuneError(input, replacement string) string {
	const sentinel = "\xff"
	input = strings.ReplaceAll(input, "�", sentinel)
	return strings.ToValidUTF8(input, replacement)
}

func main() {
	fixtures := []struct{ input, replacement, want string }{
		{"", "!", ""},
		{"A�é🚀", "!", "A�é🚀"},
		{"\xff\xfe", "?", "?"},
		{"\xffX\xfe", "?", "?X?"},
		{"\xff�\xfe", "?", "?�?"},
		{"\xe2\x82", "替", "替"},
		{"\xed\xa0\x80A", "?", "?A"},
		{"\xff\x80A\xfe", "", "A"},
		{"X\xc2\xa0\xff", "\xff\x00", "X\xc2\xa0\xff\x00"},
		{"\xc2\xc2\xa0\xc2", "?", "?\xc2\xa0?"},
		{"\x00X", "?", "\x00X"},
		{"\xff", "\xff\xe2\x82", "\xff\xe2\x82"},
	}
	for _, c := range []struct {
		name    string
		mapUTF8 func(string, string) string
		want    int
	}{
		{"Go", strings.ToValidUTF8, 0},
		{"Replace every malformed byte", replaceEveryInvalidByte, 3},
		{"Treat literal RuneError as malformed", replaceValidRuneError, 2},
	} {
		bad := 0
		for _, f := range fixtures {
			if got := c.mapUTF8(f.input, f.replacement); got != f.want {
				bad++
				if c.want == 0 {
					log.Fatalf("ToValidUTF8(%x, %x) = %x, want %x", f.input, f.replacement, got, f.want)
				}
			}
		}
		fmt.Printf("%s disagreed on %d of %d fixtures\n", c.name, bad, len(fixtures))
		if bad != c.want {
			log.Fatalf("%s wanted %d differences", c.name, c.want)
		}
	}
}
