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
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text        string
	offset      int64
	readRune    bool
	nilReceiver bool
	zeroValue   bool
}

func main() {
	edges := []input{{nilReceiver: true}, {zeroValue: true}}
	for _, text := range []string{"", "x", "xZ", "猫Z", "\xff\x00"} {
		for _, offset := range []int64{0, 1, int64(len(text)), int64(len(text)) + 1, math.MaxInt64} {
			edges = append(edges, input{text: text, offset: offset}, input{text: text, offset: offset, readRune: true})
		}
	}
	check.CompareSamples("godiff-strings-reader-unread-byte", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	offset := int64(rng.Intn(len(raw) + 2))
	if i%4 == 0 {
		offset = rng.Int63()
	}
	return input{text: string(raw), offset: offset, readRune: i%2 == 1}
}

func encode(c input) string {
	if c.nilReceiver {
		return "0"
	}
	cursor, previous := start(c)
	row := fmt.Appendf(nil, "1 %d %d", cursor, previous)
	for i := range len(c.text) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
	}
	return string(row)
}

func actual(c input) (row string) {
	var reader *strings.Reader
	switch {
	case c.nilReceiver:
		defer func() {
			if p := recover(); p != nil {
				if _, ok := p.(runtime.Error); !ok {
					panic(p)
				}
				row = "panic"
			}
		}()
	case c.zeroValue:
		reader = new(strings.Reader)
	default:
		reader = strings.NewReader(c.text)
		if position, err := reader.Seek(c.offset, io.SeekStart); err != nil || position != c.offset {
			log.Fatalf("Seek(%d) = %d, %v", c.offset, position, err)
		}
		if c.readRune {
			if _, _, err := reader.ReadRune(); err != nil && err != io.EOF {
				log.Fatal(err)
			}
		}
	}
	err := reader.UnreadByte()
	fields := reflect.ValueOf(reader).Elem()
	if text := fields.FieldByName("s").String(); text != c.text {
		log.Fatalf("UnreadByte changed backing string from %q to %q", c.text, text)
	}
	status := "ok"
	if err != nil {
		status = err.Error()
	}
	return fmt.Sprintf("%d %d %s", fields.FieldByName("i").Int(), fields.FieldByName("prevRune").Int(), status)
}

func start(c input) (cursor, previous int64) {
	cursor, previous = c.offset, -1
	if c.zeroValue {
		previous = 0
	}
	if c.readRune && cursor < int64(len(c.text)) {
		_, width := utf8.DecodeRuneInString(c.text[cursor:])
		previous, cursor = cursor, cursor+int64(width)
	}
	return
}

func expected(c input) string {
	if c.nilReceiver {
		return "panic"
	}
	cursor, previous := start(c)
	if cursor == 0 {
		return fmt.Sprintf("0 %d strings.Reader.UnreadByte: at beginning of string", previous)
	}
	return fmt.Sprintf("%d -1 ok", cursor-1)
}
