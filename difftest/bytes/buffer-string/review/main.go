package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	tag, length, offset int
	raw                 []byte
	want                string
}

func main() {
	cases := []fixture{
		{0, 0, 0, nil, "60,110,105,108,62"},
		{1, 0, 0, nil, "::"},
		{2, 3, 1, []byte{8, 10, 20, 40}, "10,20:11,21:165,165,165,165"},
		{2, 2, 2, []byte{1, 2, 3}, "::165,165,165"},
		{2, 2, 0, []byte{255, 0}, "255,0:254,1:165,165"},
		{2, 0, 0, []byte{12, 14}, "::165,165"},
	}
	var rows bytes.Buffer
	aliasFails, emptyNilFails := 0, 0
	for i, c := range cases {
		fmt.Fprint(&rows, c.tag)
		if c.tag != 0 {
			fmt.Fprint(&rows, " ", c.length, " ", c.offset)
			for _, raw := range c.raw {
				fmt.Fprint(&rows, " ", raw)
			}
		}
		fmt.Fprintln(&rows)
		got, alias, emptyNil := native(c)
		if got != c.want {
			log.Fatalf("Go fixture %d: got %q, want %q", i, got, c.want)
		}
		if alias != c.want {
			aliasFails++
		}
		if emptyNil != c.want {
			emptyNilFails++
		}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-string", len(cases), &rows) {
		if got != cases[i].want {
			log.Fatalf("Lean fixture %d: got %q, want %q", i, got, cases[i].want)
		}
	}
	if aliasFails != 2 || emptyNilFails != 1 {
		log.Fatalf("mutation failures alias=%d empty nil=%d, want 2 and 1", aliasFails, emptyNilFails)
	}
	fmt.Printf("bytes.Buffer.String passed %d fixed Go and Lean fixtures; alias and empty-nil mutants failed %d and %d.\n",
		len(cases), aliasFails, emptyNilFails)
}

func native(c fixture) (got, alias, emptyNil string) {
	if c.tag == 0 {
		var b *bytes.Buffer
		out := wire(b.String())
		return out, out, ""
	}
	backing := bytes.Clone(c.raw)
	var b bytes.Buffer
	if c.tag == 2 {
		b = *bytes.NewBuffer(backing[:c.length:len(backing)])
		b.Next(c.offset)
	}
	before := b.String()
	view := b.Bytes()
	for i := range view {
		view[i] ^= 1
	}
	changed := b.String()
	b.Reset()
	b.Write(bytes.Repeat([]byte{165}, len(backing)))
	last := ":" + wire(changed) + ":" + wire(b.String())
	got = wire(before) + last
	alias = wire(string(backing[c.offset:c.length])) + last
	return got, alias, got
}

func wire(s string) string {
	fields := make([]string, len(s))
	for i := range len(s) {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return strings.Join(fields, ",")
}
