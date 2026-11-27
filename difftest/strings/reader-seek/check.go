package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text                   string
	cursor, offset         int64
	whence                 int
	primeRune, nilReceiver bool
}

func main() {
	edges := []input{{nilReceiver: true}, {nilReceiver: true, whence: -1}}
	for _, text := range []string{"", "abc", "猫x", "\xff\x00"} {
		for _, cursor := range []int64{0, 1, math.MaxInt64} {
			for _, offset := range []int64{math.MinInt64, -2, -1, 0, 1, math.MaxInt64} {
				for _, whence := range []int{-1, io.SeekStart, io.SeekCurrent, io.SeekEnd, 3} {
					edges = append(edges, input{text: text, cursor: cursor, offset: offset, whence: whence, primeRune: true})
				}
			}
		}
	}
	check.CompareSamples("godiff-strings-reader-seek", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	cursor := int64(rng.Intn(len(raw) + 2))
	if i%4 == 0 {
		cursor = rng.Int63()
	}
	offset := int64(rng.Intn(129) - 64)
	if i%3 == 0 {
		offset = int64(rng.Uint64())
	}
	whence := rng.Intn(5) - 1
	if i%16 == 0 {
		whence = int(rng.Uint64())
	}
	return input{text: string(raw), cursor: cursor, offset: offset, whence: whence, primeRune: i%2 == 1}
}

func start(c input) (cursor, previous int64) {
	cursor, previous = c.cursor, -1
	if c.primeRune && cursor < int64(len(c.text)) {
		_, width := utf8.DecodeRuneInString(c.text[cursor:])
		previous, cursor = cursor, cursor+int64(width)
	}
	return
}

func encode(c input) string {
	tag := 1
	if c.nilReceiver {
		tag = 0
	}
	cursor, previous := start(c)
	row := fmt.Appendf(nil, "%d %d %d %d %d", tag, cursor, previous, c.offset, c.whence)
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func describe(pos int64, status string, cursor, previous int64) string {
	return fmt.Sprintf("%d %s %d %d", pos, status, cursor, previous)
}

func actual(c input) (row string) {
	var r *strings.Reader
	if c.nilReceiver {
		defer func() {
			if p := recover(); p != nil {
				if _, ok := p.(runtime.Error); !ok {
					panic(p)
				}
				row = "panic"
			}
		}()
	} else {
		r = strings.NewReader(c.text)
		if pos, err := r.Seek(c.cursor, io.SeekStart); err != nil || pos != c.cursor {
			log.Fatalf("positioning reader at %d returned %d, %v", c.cursor, pos, err)
		}
		if c.primeRune {
			if _, _, err := r.ReadRune(); err != nil && err != io.EOF {
				log.Fatal(err)
			}
		}
	}
	pos, err := r.Seek(c.offset, c.whence)
	status := "ok"
	if err != nil {
		switch err.Error() {
		case "strings.Reader.Seek: invalid whence":
			status = "whence"
		case "strings.Reader.Seek: negative position":
			status = "negative"
		default:
			log.Fatalf("unexpected Seek error: %v", err)
		}
	}
	fields := reflect.ValueOf(r).Elem()
	if text := fields.FieldByName("s").String(); text != c.text {
		log.Fatalf("Seek changed backing string from %q to %q", c.text, text)
	}
	return describe(pos, status, fields.FieldByName("i").Int(), fields.FieldByName("prevRune").Int())
}

func expected(c input) string {
	if c.nilReceiver {
		return "panic"
	}
	cursor, _ := start(c)
	var base int64
	switch c.whence {
	case io.SeekStart:
	case io.SeekCurrent:
		base = cursor
	case io.SeekEnd:
		base = int64(len(c.text))
	default:
		return describe(0, "whence", cursor, -1)
	}
	absolute := new(big.Int).Add(big.NewInt(base), big.NewInt(c.offset))
	if absolute.Cmp(big.NewInt(math.MaxInt64)) > 0 || absolute.Sign() < 0 {
		// Base is nonnegative and both operands fit in int64. Positive
		// overflow wraps to a negative result and negative underflow is impossible.
		return describe(0, "negative", cursor, -1)
	}
	pos := absolute.Int64()
	return describe(pos, "ok", pos, -1)
}
