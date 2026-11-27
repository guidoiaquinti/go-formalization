package main

import (
	"io"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	contents    string
	cursor      int64
	nilReceiver bool
	zero        bool
}

func main() {
	edges := []input{{nilReceiver: true}, {zero: true}}
	for _, contents := range []string{"", "\x00", "\xff\xc0\x80", "\u00e9\U0001f680", "abcdef"} {
		for _, cursor := range []int64{0, 1, int64(len(contents)), int64(len(contents)) + 1, 1 << 32, math.MaxInt64} {
			edges = append(edges, input{contents: contents, cursor: cursor})
		}
	}
	check.CompareSamples("godiff-strings-reader-len", edges, sample, modelInput, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	cursor := int64(rng.Intn(len(raw) + 1))
	if i%2 == 0 {
		cursor = rng.Int63()
	}
	return input{contents: string(raw), cursor: cursor}
}

func modelInput(c input) string {
	if c.nilReceiver {
		return "0"
	}
	row := strconv.AppendInt([]byte("1 "), c.cursor, 10)
	for i := range len(c.contents) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.contents[i]), 10)
	}
	return string(row)
}

func actual(c input) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); c.nilReceiver && ok {
				result = "panic"
			} else {
				panic(p)
			}
		}
	}()
	var reader *strings.Reader
	switch {
	case c.zero:
		reader = new(strings.Reader)
	case !c.nilReceiver:
		reader = strings.NewReader(c.contents)
		if position, err := reader.Seek(c.cursor, io.SeekStart); err != nil || position != c.cursor {
			log.Fatalf("Seek(%d) returned position %d, error %v", c.cursor, position, err)
		}
	}
	return strconv.Itoa(reader.Len())
}

func expected(c input) string {
	if c.nilReceiver {
		return "panic"
	}
	return strconv.Itoa(len(c.contents[min(c.cursor, int64(len(c.contents))):]))
}
