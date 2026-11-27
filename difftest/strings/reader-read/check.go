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
	text         string
	dst          []byte
	offset       int64
	primeRune    bool
	nilReceiver  bool
	zeroReceiver bool
}

func main() {
	edges := []input{{nilReceiver: true}, {zeroReceiver: true}}
	for _, text := range []string{"", "abc", "猫", "\xff\x00"} {
		for _, offset := range []int64{0, 1, int64(len(text)), int64(len(text)) + 1, math.MaxInt64} {
			for _, dst := range [][]byte{nil, {}, {0xfa}, {1, 2, 3, 4, 5}} {
				edges = append(edges, input{text: text, dst: dst, offset: offset},
					input{text: text, dst: dst, offset: offset, primeRune: true})
			}
		}
	}
	check.CompareSamples("godiff-strings-reader-read", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	src := make([]byte, rng.Intn(65))
	dst := make([]byte, rng.Intn(65))
	rng.Read(src)
	rng.Read(dst)
	offset := int64(rng.Intn(len(src) + 2))
	if i%4 == 0 {
		offset = rng.Int63()
	}
	return input{text: string(src), dst: dst, offset: offset, primeRune: i%2 == 1}
}

func start(c input) (cursor, previous int64) {
	cursor, previous = c.offset, -1
	if c.zeroReceiver {
		previous = 0
	}
	if c.primeRune && cursor < int64(len(c.text)) {
		_, width := utf8.DecodeRuneInString(c.text[cursor:])
		previous, cursor = cursor, cursor+int64(width)
	}
	return
}

func appendBytes(row, data []byte) []byte {
	for _, b := range data {
		row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
	}
	return row
}

func encode(c input) string {
	tag := 1
	if c.nilReceiver {
		tag = 0
	}
	cursor, previous := start(c)
	row := fmt.Appendf(nil, "%d %d %d %d", tag, cursor, previous, len(c.dst))
	return string(appendBytes(appendBytes(row, c.dst), []byte(c.text)))
}

func describe(n int, eof bool, cursor, previous int64, dst []byte) string {
	return string(appendBytes(fmt.Appendf(nil, "%d %t %d %d", n, eof, cursor, previous), dst))
}

func actual(c input) (result string) {
	var r *strings.Reader
	switch {
	case c.nilReceiver:
		defer func() {
			if p := recover(); p != nil {
				if _, ok := p.(runtime.Error); !ok {
					panic(p)
				}
				result = "panic"
			}
		}()
	case c.zeroReceiver:
		r = new(strings.Reader)
	default:
		r = strings.NewReader(c.text)
		if position, err := r.Seek(c.offset, io.SeekStart); err != nil || position != c.offset {
			log.Fatalf("Seek(%d) = %d, %v", c.offset, position, err)
		}
		if c.primeRune {
			if _, _, err := r.ReadRune(); err != nil && err != io.EOF {
				log.Fatal(err)
			}
		}
	}
	buffer := make([]byte, len(c.dst)+2)
	const head, tail = 0xfa, 0xfb
	buffer[0], buffer[len(buffer)-1] = head, tail
	var dst []byte
	if c.dst != nil {
		dst = buffer[1 : len(buffer)-1]
		copy(dst, c.dst)
	}
	n, err := r.Read(dst)
	if err != nil && err != io.EOF {
		log.Fatalf("unexpected Read error: %v", err)
	}
	if buffer[0] != head || buffer[len(buffer)-1] != tail {
		log.Fatal("Read wrote outside the destination")
	}
	fields := reflect.ValueOf(r).Elem()
	if text := fields.FieldByName("s").String(); text != c.text {
		log.Fatalf("Read changed backing string from %q to %q", c.text, text)
	}
	return describe(n, err == io.EOF, fields.FieldByName("i").Int(), fields.FieldByName("prevRune").Int(), dst)
}

func expected(c input) string {
	if c.nilReceiver {
		return "panic"
	}
	cursor, previous := start(c)
	if cursor >= int64(len(c.text)) {
		return describe(0, true, cursor, previous, c.dst)
	}
	end := min(len(c.text), int(cursor)+len(c.dst))
	taken := c.text[cursor:end]
	dst := append([]byte(taken), c.dst[len(taken):]...)
	return describe(len(taken), false, int64(end), -1, dst)
}
