package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	prior, suffix string
	reserve       int
	reset         bool
	nilReceiver   bool
}

func main() {
	edges := []input{{nilReceiver: true}, {}, {reserve: 64},
		{prior: "old", reset: true}, {prior: "old", suffix: "猫", reset: true},
		{prior: "猫", suffix: "\xff\x00", reserve: 129}}
	for b := range 256 {
		edges = append(edges, input{suffix: string([]byte{byte(b)})})
	}
	check.CompareSamples("godiff-strings-builder-cap", edges, sample, encode, actual, expected)
}

func sample(r *rand.Rand, i int) input {
	data := make([]byte, r.Intn(65))
	r.Read(data)
	split := r.Intn(len(data) + 1)
	return input{prior: string(data[:split]), suffix: string(data[split:]),
		reserve: r.Intn(513), reset: i%2 == 0}
}

func prepare(b *strings.Builder, c input) bool {
	identity := false
	if c.prior != "" {
		b.WriteString(c.prior)
		identity = true
	}
	if c.reset {
		b.Reset()
		identity = false
	}
	if c.reserve != 0 {
		b.Grow(c.reserve)
		identity = true
	}
	if c.suffix != "" {
		b.WriteString(c.suffix)
		identity = true
	}
	return identity
}

func encode(c input) string {
	if c.nilReceiver {
		return "0"
	}
	var b strings.Builder
	identity := prepare(&b, c)
	identityBit := 0
	if identity {
		identityBit = 1
	}
	capacity := backingCapacity(&b)
	if capacity == 0 {
		capacity = -1
	}
	row := fmt.Appendf(nil, "1 %d %d", identityBit, capacity)
	text := b.String()
	for i := range len(text) {
		row = strconv.AppendUint(append(row, ' '), uint64(text[i]), 10)
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
		return strconv.Itoa(b.Cap())
	}
	var b strings.Builder
	prepare(&b, c)
	before, length := b.String(), b.Len()
	capacity := b.Cap()
	if b.String() != before || b.Len() != length || capacity < length {
		log.Fatalf("Cap returned %d and altered the builder or fell below its length %d", capacity, length)
	}
	return strconv.Itoa(capacity)
}

// Observe the pinned slice header without invoking Builder.Cap.
func backingCapacity(b *strings.Builder) int {
	return reflect.ValueOf(b).Elem().FieldByName("buf").Cap()
}

func expected(c input) string {
	if c.nilReceiver {
		return "panic"
	}
	var b strings.Builder
	prepare(&b, c)
	return strconv.Itoa(backingCapacity(&b))
}
