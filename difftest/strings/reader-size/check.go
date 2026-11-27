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
	text        string
	cursor      int64
	nilReceiver bool
	zeroValue   bool
}

func main() {
	edges := []input{{nilReceiver: true}, {zeroValue: true}}
	for _, text := range []string{"", "a", "こんにちは", "\x00\xff"} {
		for _, cursor := range []int64{0, 1, int64(len(text)), int64(len(text)) + 1, math.MaxInt64} {
			edges = append(edges, input{text: text, cursor: cursor})
		}
	}
	check.CompareSamples("godiff-strings-reader-size", edges, func(rng *rand.Rand, i int) input {
		raw := make([]byte, rng.Intn(65))
		rng.Read(raw)
		cursor := int64(rng.Intn(len(raw) + 2))
		if i%2 == 0 {
			cursor = rng.Int63()
		}
		return input{text: string(raw), cursor: cursor}
	}, encode, actual, func(c input) string {
		if c.nilReceiver {
			return "panic"
		}
		return strconv.Itoa(len([]byte(c.text)))
	})
}

func encode(c input) string {
	if c.nilReceiver {
		return "0"
	}
	row := strconv.AppendInt([]byte("1 "), c.cursor, 10)
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func actual(c input) (result string) {
	var reader *strings.Reader
	if c.nilReceiver {
		defer func() {
			if p := recover(); p != nil {
				if _, ok := p.(runtime.Error); !ok {
					panic(p)
				}
				result = "panic"
			}
		}()
	} else {
		reader = new(strings.Reader)
		if !c.zeroValue {
			reader.Reset(c.text)
			if position, err := reader.Seek(c.cursor, io.SeekStart); err != nil || position != c.cursor {
				log.Fatalf("Seek(%d) returned position %d, error %v", c.cursor, position, err)
			}
		}
	}
	return strconv.FormatInt(reader.Size(), 10)
}
