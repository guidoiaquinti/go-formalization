package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func write(b *strings.Builder, r rune) (n int, err error, failure any) {
	defer func() { failure = recover() }()
	n, err = b.WriteRune(r)
	return
}

func main() {
	cases := []struct {
		r    rune
		want string
	}{
		{0, "\x00"}, {0x7f, "\x7f"}, {0x80, "\xc2\x80"}, {0x800, "\xe0\xa0\x80"},
		{0x10ffff, "\xf4\x8f\xbf\xbf"}, {-1, "\xef\xbf\xbd"},
		{0xd800, "\xef\xbf\xbd"}, {0x110000, "\xef\xbf\xbd"},
	}
	var oneByte, noReplace int
	for _, c := range cases {
		var b strings.Builder
		b.Grow(32)
		b.WriteString("\xffP")
		old := b.String()
		n, err, failure := write(&b, c.r)
		if n != len(c.want) || err != nil || failure != nil || b.String() != "\xffP"+c.want || old != "\xffP" {
			log.Fatalf("WriteRune(%#x): n=%d err=%v failure=%v final=%q old=%q", c.r, n, err, failure, b.String(), old)
		}
		if n != 1 {
			oneByte++
		}
		skipped := string(c.r)
		if c.r < 0 || c.r > 0x10ffff || c.r >= 0xd800 && c.r <= 0xdfff {
			skipped = ""
		}
		if skipped != c.want {
			noReplace++
		}
	}
	var initialized strings.Builder
	initialized.Grow(0)
	copied := initialized
	if _, _, failure := write(&copied, 'x'); failure != "strings: illegal use of non-zero Builder copied by value" ||
		copied.Len() != 0 || initialized.Len() != 0 {
		log.Fatalf("copied initialized-empty Builder: failure=%v lengths=%d/%d", failure, copied.Len(), initialized.Len())
	}
	var zero strings.Builder
	fresh := zero
	if n, err, failure := write(&fresh, 'x'); n != 1 || err != nil || failure != nil || fresh.String() != "x" || zero.Len() != 0 {
		log.Fatal("copy of zero Builder did not accept first write")
	}
	if _, _, failure := write(nil, 'x'); failure == nil {
		log.Fatal("nil receiver did not panic")
	} else if _, ok := failure.(runtime.Error); !ok {
		log.Fatalf("nil receiver recovered %T", failure)
	}
	if oneByte != 6 || noReplace != 3 {
		log.Fatalf("unexpected mutation counts: rune count=%d no replacement=%d", oneByte, noReplace)
	}
	fmt.Printf("Builder.WriteRune: %d fixtures pass, mutant failures rune count=%d no replacement=%d\n", len(cases)+3, oneByte, noReplace)
}
