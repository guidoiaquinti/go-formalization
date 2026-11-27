package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	nilReceiver bool
	write       bool
	grow        int
	content     string
}

func main() {
	edges := []input{{nilReceiver: true}, {}, {write: true}, {grow: 17}, {write: true, grow: 17, content: "abc"}}
	for b := range 256 {
		edges = append(edges, input{write: true, content: string([]byte{byte(b)})})
	}
	check.CompareSamples("godiff-strings-builder-reset", edges, sample, modelInput, actual, func(c input) string {
		if c.nilReceiver {
			return "panic"
		}
		return "0 0"
	})
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	switch i % 3 {
	case 0:
		return input{grow: 1 + rng.Intn(512)}
	case 1:
		return input{write: true, content: string(raw)}
	default:
		return input{write: true, grow: rng.Intn(512), content: string(raw)}
	}
}

func prepare(c input, b *strings.Builder) {
	if c.grow != 0 {
		b.Grow(c.grow)
	}
	if c.write {
		b.WriteString(c.content)
	}
}

func modelInput(c input) string {
	if c.nilReceiver {
		return "0"
	}
	var b strings.Builder
	prepare(c, &b)
	capacity := b.Cap()
	if capacity == 0 {
		capacity = -1
	}
	identity := 0
	if c.grow != 0 || c.write {
		identity = 1
	}
	row := fmt.Appendf(nil, "1 %d %d", identity, capacity)
	for i := range len(c.content) {
		row = strconv.AppendUint(append(row, ' '), uint64(c.content[i]), 10)
	}
	return string(row)
}

func actual(c input) (out string) {
	if c.nilReceiver {
		defer func() {
			if p := recover(); p != nil {
				if _, ok := p.(runtime.Error); !ok {
					panic(p)
				}
				out = "panic"
			}
		}()
		var b *strings.Builder
		b.Reset()
		return "returned"
	}
	var b strings.Builder
	prepare(c, &b)
	saved := b.String()
	b.Reset()
	row := fmt.Appendf(nil, "%d %d", b.Cap(), b.Len())
	s := b.String()
	for i := range len(s) {
		row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
	}
	b.WriteString("reused")
	if b.String() != "reused" || saved != c.content {
		log.Fatalf("after Reset and reuse: builder=%q, saved=%q, original=%q", b.String(), saved, c.content)
	}
	return string(row)
}
