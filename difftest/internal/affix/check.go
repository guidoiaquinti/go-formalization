// Package affix exercises matching byte prefixes and suffixes.
package affix

import (
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type pair struct {
	text, affix string
}

func Compare(binary string, suffix bool, actual func(string, string) bool) {
	compare(binary, suffix, func(p pair) string {
		return strconv.FormatBool(actual(p.text, p.affix))
	}, func(p pair) string {
		_, found := boundary(p, suffix)
		return strconv.FormatBool(found)
	})
}

func CompareTrim(binary string, suffix bool, actual func(string, string) string) {
	compare(binary, suffix, func(p pair) string {
		return string(appendBytes([]byte("bytes"), actual(p.text, p.affix)))
	}, func(p pair) string {
		n, found := boundary(p, suffix)
		out := p.text
		if found {
			if suffix {
				out = p.text[:n]
			} else {
				out = p.text[n:]
			}
		}
		return string(appendBytes([]byte("bytes"), out))
	})
}

func compare(binary string, suffix bool, actual, expected func(pair) string) {
	values := []string{"", "a", "ab", "abc", "aaa", "a\x00", "\x00", "\xff", "\xc0\x80", "é", "🔧"}
	var edges []pair
	for _, text := range values {
		for _, affix := range values {
			edges = append(edges, pair{text, affix})
		}
	}
	for b := range 256 {
		s := string([]byte{byte(b)})
		text := s + "!"
		if suffix {
			text = "!" + s
		}
		edges = append(edges, pair{text, s})
	}
	check.CompareSamples(binary, edges, func(r *rand.Rand, i int) pair {
		data := make([]byte, r.Intn(33))
		r.Read(data)
		length := r.Intn(len(data) + 1)
		offset := 0
		if suffix {
			offset = len(data) - length
		}
		affix := append([]byte(nil), data[offset:offset+length]...)
		switch i % 4 {
		case 1:
			affix = append(affix, byte(r.Intn(256)))
		case 2:
			if length > 0 {
				affix[r.Intn(length)] ^= 1
			}
		case 3:
			r.Read(affix)
		}
		return pair{string(data), string(affix)}
	}, encode, actual, expected)
}

// Look for an exact byte concatenation at the chosen end.
func boundary(p pair, suffix bool) (int, bool) {
	for n := range len(p.text) + 1 {
		if !suffix && p.affix+p.text[n:] == p.text {
			return n, true
		}
		if suffix && p.text[:n]+p.affix == p.text {
			return n, true
		}
	}
	return 0, false
}

func encode(p pair) string {
	out := strconv.AppendInt(nil, int64(len(p.affix)), 10)
	return string(appendBytes(appendBytes(out, p.affix), p.text))
}

func appendBytes(out []byte, text string) []byte {
	for i := range len(text) {
		out = strconv.AppendUint(append(out, ' '), uint64(text[i]), 10)
	}
	return out
}
