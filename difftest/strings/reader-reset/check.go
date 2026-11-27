package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	old, replacement string
	offset           int64
	readRune         bool
	nilReceiver      bool
	zeroValue        bool
}

type state struct {
	cursor, previous int64
}

func main() {
	edges := []input{{nilReceiver: true}, {nilReceiver: true, replacement: "x"}, {zeroValue: true}}
	values := []string{"", "old", "世界", "\xff\x00"}
	for _, old := range values {
		for _, replacement := range values {
			edges = append(edges,
				input{old: old, replacement: replacement, readRune: true},
				input{old: old, replacement: replacement, offset: math.MaxInt64})
		}
	}
	check.CompareSamples("godiff-strings-reader-reset", edges, sample, encode, actual, func(c input) string {
		if c.nilReceiver {
			return "panic"
		}
		return describe(state{cursor: 0, previous: -1}, c.replacement)
	})
}

func sample(rng *rand.Rand, i int) input {
	old := make([]byte, rng.Intn(65))
	replacement := make([]byte, rng.Intn(65))
	rng.Read(old)
	rng.Read(replacement)
	offset := int64(rng.Intn(len(old) + 2))
	if i%4 == 0 {
		offset = rng.Int63()
	}
	return input{old: string(old), replacement: string(replacement), offset: offset, readRune: i%2 == 1}
}

func prepare(c input) *strings.Reader {
	if c.zeroValue {
		return new(strings.Reader)
	}
	reader := strings.NewReader(c.old)
	if pos, err := reader.Seek(c.offset, io.SeekStart); err != nil || pos != c.offset {
		log.Fatalf("Seek(%d) = %d, %v", c.offset, pos, err)
	}
	if c.readRune {
		if _, _, err := reader.ReadRune(); err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
	return reader
}

func observe(reader *strings.Reader) state {
	// Read the private fields without invoking a method that changes them.
	fields := reflect.ValueOf(reader).Elem()
	return state{cursor: fields.FieldByName("i").Int(), previous: fields.FieldByName("prevRune").Int()}
}

func appendBytes(row []byte, text string) []byte {
	for i := range len(text) {
		row = strconv.AppendUint(append(row, ' '), uint64(text[i]), 10)
	}
	return row
}

func encode(c input) string {
	if c.nilReceiver {
		return string(appendBytes([]byte("0"), c.replacement))
	}
	s := observe(prepare(c))
	row := fmt.Appendf(nil, "1 %d %d %d", s.cursor, s.previous, len(c.old))
	return string(appendBytes(appendBytes(row, c.old), c.replacement))
}

func actual(c input) (out string) {
	var reader *strings.Reader
	if c.nilReceiver {
		defer func() {
			if p := recover(); p != nil {
				if _, ok := p.(runtime.Error); !ok {
					panic(p)
				}
				out = "panic"
			}
		}()
	} else {
		reader = prepare(c)
	}
	reader.Reset(c.replacement)
	s := observe(reader)
	raw, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	return describe(s, string(raw))
}

func describe(s state, text string) string {
	return string(appendBytes(fmt.Appendf(nil, "%d %d", s.cursor, s.previous), text))
}
