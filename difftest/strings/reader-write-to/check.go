package main

import (
	"errors"
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
	text                  string
	cursor                int64
	count, script         int
	stringWriter, nilRead bool
	primeRune             bool
}

var writerError = errors.New("writer error")
var panicMarker = new(int)

func main() {
	edges := []input{
		{nilRead: true},
		{text: "猫a", cursor: 3, count: 1, stringWriter: true, primeRune: true},
		{text: "猫", cursor: 3, script: 2, primeRune: true},
		{text: "ab", cursor: 1, script: 2, primeRune: true},
	}
	for _, text := range []string{"", "abc", "猫", "\xff\x00"} {
		for _, cursor := range []int64{0, 1, int64(len(text)), math.MaxInt64} {
			for _, count := range []int{math.MinInt, -1, 0, 1, len(text), len(text) + 1, math.MaxInt} {
				for script := range 4 {
					for _, method := range []bool{false, true} {
						edges = append(edges, input{text, cursor, count, script, method, false, false})
					}
				}
			}
		}
	}
	check.CompareSamples("godiff-strings-reader-write-to", edges, sample, encode, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
	text := make([]byte, rng.Intn(65))
	rng.Read(text)
	cursor := int64(rng.Intn(len(text) + 2))
	return input{string(text), cursor, rng.Intn(len(text)+5) - 2, rng.Intn(4), rng.Intn(2) != 0, false, false}
}

func appendBytes(row []byte, text string) []byte {
	for i := range len(text) {
		row = strconv.AppendUint(append(row, ' '), uint64(text[i]), 10)
	}
	return row
}

func encode(c input) string {
	live, method := 1, 0
	if c.nilRead {
		live = 0
	}
	if c.stringWriter {
		method = 1
	}
	previous := -1
	if c.primeRune {
		previous = 0
	}
	row := fmt.Appendf(nil, "%d %d %d %d %d %d", live, c.cursor, previous, method, c.script, c.count)
	return string(appendBytes(row, c.text))
}

type writerProbe struct {
	r       *strings.Reader
	input   input
	trace   string
	payload string
}

func (p *writerProbe) write(method, text string) (int, error) {
	if p.trace != "none" {
		log.Fatal("callback invoked more than once")
	}
	fields := reflect.ValueOf(p.r).Elem()
	if fields.FieldByName("i").Int() != p.input.cursor || fields.FieldByName("prevRune").Int() != -1 {
		log.Fatal("callback observed the wrong cursor or rune marker")
	}
	p.trace, p.payload = method, text
	switch p.input.script {
	case 1:
		return p.input.count, writerError
	case 2:
		panic(panicMarker)
	default:
		return p.input.count, nil
	}
}

func (p *writerProbe) Write(data []byte) (int, error) {
	return p.write("bytes", string(data))
}

type stringProbe struct{ *writerProbe }

func (p stringProbe) WriteString(s string) (int, error) {
	return p.write("string", s)
}

func describe(status string, count, cursor, previous int64, trace, payload string) string {
	row := fmt.Appendf(nil, "%s %d %d %d %s", status, count, cursor, previous, trace)
	return string(appendBytes(row, payload))
}

func actual(c input) (row string) {
	var r *strings.Reader
	if !c.nilRead {
		r = strings.NewReader(c.text)
		if c.primeRune {
			if _, size, err := r.ReadRune(); int64(size) != c.cursor || err != nil {
				log.Fatalf("priming rune returned %d, %v", size, err)
			}
		} else if pos, err := r.Seek(c.cursor, io.SeekStart); pos != c.cursor || err != nil {
			log.Fatalf("preparing cursor %d returned %d, %v", c.cursor, pos, err)
		}
	}
	p := &writerProbe{r: r, input: c, trace: "none"}
	var w io.Writer = p
	if c.script == 3 {
		w = nil
	} else if c.stringWriter {
		w = stringProbe{p}
	}
	status := "ok"
	var count int64
	defer func() {
		if recovered := recover(); recovered != nil {
			_, runtimePanic := recovered.(runtime.Error)
			switch {
			case c.nilRead && runtimePanic:
				row = "nilReader"
				return
			case c.script == 3 && runtimePanic:
				status = "nilWriter"
			case recovered == panicMarker:
				status = "writerPanic"
			case recovered == "strings.Reader.WriteTo: invalid WriteString count":
				status = "invalidCount"
			default:
				panic(recovered)
			}
		}
		fields := reflect.ValueOf(r).Elem()
		if fields.FieldByName("s").String() != c.text {
			log.Fatal("WriteTo changed backing bytes")
		}
		row = describe(status, count, fields.FieldByName("i").Int(), fields.FieldByName("prevRune").Int(), p.trace, p.payload)
	}()
	var err error
	count, err = r.WriteTo(w)
	switch err {
	case nil:
	case io.ErrShortWrite:
		status = "shortWrite"
	case writerError:
		status = "writerError"
	default:
		log.Fatalf("unexpected writer error: %v", err)
	}
	return
}

func expected(c input) string {
	if c.nilRead {
		return "nilReader"
	}
	if c.cursor >= int64(len(c.text)) {
		return describe("ok", 0, c.cursor, -1, "none", "")
	}
	if c.script == 3 {
		return describe("nilWriter", 0, c.cursor, -1, "none", "")
	}
	trace := "bytes"
	if c.stringWriter {
		trace = "string"
	}
	suffix := c.text[c.cursor:]
	if c.script == 2 {
		return describe("writerPanic", 0, c.cursor, -1, trace, suffix)
	}
	if c.count > len(suffix) {
		return describe("invalidCount", 0, c.cursor, -1, trace, suffix)
	}
	status := "ok"
	if c.script == 1 {
		status = "writerError"
	} else if c.count < len(suffix) {
		status = "shortWrite"
	}
	return describe(status, int64(c.count), c.cursor+int64(c.count), -1, trace, suffix)
}
