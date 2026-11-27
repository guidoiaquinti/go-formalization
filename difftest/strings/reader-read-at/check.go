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
	text        string
	dst         []byte
	offset      int64
	nilReceiver bool
}

func main() {
	var edges []input
	for _, text := range []string{"", "abc", "猫", "\xff\x00"} {
		for _, offset := range []int64{math.MinInt64, -1, 0, 1, int64(len(text)), int64(len(text)) + 1, math.MaxInt64} {
			for _, dst := range [][]byte{nil, {}, {0xfa}, {1, 2, 3, 4, 5}} {
				edges = append(edges, input{text: text, dst: dst, offset: offset})
			}
		}
	}
	for _, offset := range []int64{-1, 0, 5} {
		edges = append(edges, input{nilReceiver: true, offset: offset, dst: []byte{42}})
	}
	check.CompareSamples("godiff-strings-reader-read-at", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, i int) input {
	src := make([]byte, rng.Intn(65))
	dst := make([]byte, rng.Intn(65))
	rng.Read(src)
	rng.Read(dst)
	offset := int64(rng.Intn(len(src)+3) - 1)
	if i%4 == 0 {
		offset = int64(rng.Uint64())
	}
	return input{text: string(src), dst: dst, offset: offset}
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
	row := fmt.Appendf(nil, "%d %d %d", tag, c.offset, len(c.dst))
	return string(appendBytes(appendBytes(row, c.dst), []byte(c.text)))
}

func describe(status string, n int, dst []byte) string {
	return string(appendBytes(fmt.Appendf(nil, "%s %d", status, n), dst))
}

type state struct {
	text             string
	cursor, previous int64
}

func observe(r *strings.Reader) state {
	fields := reflect.ValueOf(r).Elem()
	return state{fields.FieldByName("s").String(), fields.FieldByName("i").Int(), fields.FieldByName("prevRune").Int()}
}

func actual(c input) (result string) {
	var r *strings.Reader
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
		r = strings.NewReader(c.text)
		if _, _, err := r.ReadRune(); err != nil && err != io.EOF {
			log.Fatal(err)
		}
		before := observe(r)
		defer func() {
			if after := observe(r); after != before {
				log.Fatalf("ReadAt changed the reader from %+v to %+v", before, after)
			}
		}()
	}
	buffer := make([]byte, len(c.dst)+2)
	head, tail := byte(0xfa), byte(0xfb)
	buffer[0], buffer[len(buffer)-1] = head, tail
	var dst []byte
	if c.dst != nil {
		dst = buffer[1 : len(buffer)-1]
		copy(dst, c.dst)
	}
	n, err := r.ReadAt(dst, c.offset)
	if buffer[0] != head || buffer[len(buffer)-1] != tail {
		log.Fatal("ReadAt wrote outside the destination")
	}
	status := "ok"
	switch {
	case err == io.EOF:
		status = "eof"
	case err != nil:
		if err.Error() != "strings.Reader.ReadAt: negative offset" {
			log.Fatalf("unexpected ReadAt error: %v", err)
		}
		status = "negative"
	}
	return describe(status, n, dst)
}

func expected(c input) string {
	if c.offset < 0 {
		return describe("negative", 0, c.dst)
	}
	if c.nilReceiver {
		return "panic"
	}
	if c.offset >= int64(len(c.text)) {
		return describe("eof", 0, c.dst)
	}
	end := min(len(c.text), int(c.offset)+len(c.dst))
	taken := c.text[c.offset:end]
	dst := append([]byte(taken), c.dst[len(taken):]...)
	if len(taken) < len(c.dst) {
		return describe("eof", len(taken), dst)
	}
	return describe("ok", len(taken), dst)
}
