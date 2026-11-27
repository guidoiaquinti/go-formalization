package main

import (
	"bytes"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const spaces = "\t\n\v\f\r \u0085\u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000"

func main() {
	if fields := bytes.Fields([]byte{}); fields == nil || len(fields) != 0 {
		log.Fatal("bytes.Fields on a nonnil empty slice must return a nonnil empty result")
	}
	edges := []string{"", " \n ", "word", "  one    two  three   ", "one\u200btwo", "a\xff\xc2 b\uFFFD", spaces}
	for b := range 256 {
		edges = append(edges, " x"+string([]byte{byte(b)})+"y \t")
	}
	for _, space := range spaces {
		for _, r := range []rune{space - 1, space, space + 1} {
			edges = append(edges, string(r), " x"+string(r)+"y"+string(r)+string(r)+"z ")
		}
	}
	check.CompareSamples("godiff-bytes-fields", edges, sample, encodeInput, actual, expected)
}

func sample(r *rand.Rand, i int) string {
	var s strings.Builder
	whitespace := []rune(spaces)
	for range r.Intn(97) {
		if i%2 == 0 {
			if r.Intn(3) == 0 {
				s.WriteByte(" \t\n\v\f\r"[r.Intn(6)])
			} else {
				s.WriteByte(byte(r.Intn(128)))
			}
		} else {
			switch r.Intn(3) {
			case 0:
				s.WriteRune(whitespace[r.Intn(len(whitespace))])
			case 1:
				s.WriteByte(byte(r.Intn(256)))
			case 2:
				s.WriteRune(rune(r.Intn(0x110000)))
			}
		}
	}
	return s.String()
}

func encodeInput(s string) string {
	row := []byte{'0'}
	for i := range len(s) {
		row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
	}
	return string(row)
}

type fieldView struct {
	text     string
	capacity int
}

func describe(fields []fieldView) string {
	if fields == nil {
		return "nil"
	}
	encoded := make([]string, len(fields))
	for i, field := range fields {
		bytes := make([]string, len(field.text))
		for j := range len(field.text) {
			bytes[j] = strconv.Itoa(int(field.text[j]))
		}
		encoded[i] = "f" + strconv.Itoa(field.capacity) + ":" + strings.Join(bytes, ",")
	}
	return "slice" + strings.Join(encoded, ";")
}

func actual(s string) string {
	var src []byte
	if s != "" {
		src = make([]byte, len(s), len(s)+4)
		copy(src, s)
	}
	fields := bytes.Fields(src)
	if fields == nil {
		return "nil"
	}
	views := make([]fieldView, len(fields))
	for i, field := range fields {
		if len(field) == 0 {
			log.Fatalf("bytes.Fields(%q) produced an empty field", s)
		}
		views[i] = fieldView{string(field), cap(field)}
	}
	return describe(views)
}

func expected(s string) string {
	var fields []string
	start := -1
	for offset, r := range s {
		if unicode.IsSpace(r) {
			if start >= 0 {
				fields = append(fields, s[start:offset])
				start = -1
			}
		} else if start < 0 {
			start = offset
		}
	}
	if start >= 0 {
		fields = append(fields, s[start:])
	}
	views := make([]fieldView, len(fields))
	for i, field := range fields {
		views[i] = fieldView{field, len(field)}
	}
	return describe(views)
}
