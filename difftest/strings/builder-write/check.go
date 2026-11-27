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
	initial, appended string
	reserve, mode     int
	nilSlice, touched bool
}

func main() {
	var edges []input
	for _, initial := range []string{"", "abc", "猫", "\x00\xff"} {
		for _, appended := range []string{"", "a", "🚀", "\xff\x00"} {
			for _, reserve := range []int{0, 64} {
				for mode := range 3 {
					edges = append(edges, input{initial: initial, appended: appended, reserve: reserve, mode: mode})
					if appended == "" {
						for _, touched := range []bool{false, true} {
							edges = append(edges, input{initial: initial, reserve: reserve, mode: mode, nilSlice: true, touched: touched})
						}
						edges = append(edges, input{initial: initial, reserve: reserve, mode: mode, touched: true})
					}
				}
			}
		}
	}
	for b := range 256 {
		edges = append(edges, input{appended: string([]byte{byte(b)})})
	}
	check.CompareSamples("godiff-strings-builder-write", edges, sample, encode, actual, expected)
}

func sample(r *rand.Rand, i int) input {
	data := make([]byte, r.Intn(129))
	r.Read(data)
	split := r.Intn(len(data) + 1)
	reserve := 0
	if i%3 == 0 {
		reserve = r.Intn(513)
	}
	return input{initial: string(data[:split]), appended: string(data[split:]), reserve: reserve,
		mode: r.Intn(3), nilSlice: split == len(data) && i%2 == 0, touched: i%5 == 0}
}

func prepare(c input) *strings.Builder {
	b := new(strings.Builder)
	if c.touched {
		b.Grow(0)
	}
	if c.reserve != 0 {
		b.Grow(c.reserve)
	}
	if c.initial != "" {
		b.WriteString(c.initial)
	}
	return b
}

func identified(b *strings.Builder) bool {
	return !reflect.ValueOf(b).Elem().FieldByName("addr").IsNil()
}

func appendBytes(row []byte, value string) []byte {
	for i := range len(value) {
		row = strconv.AppendUint(append(row, ' '), uint64(value[i]), 10)
	}
	return row
}

func encode(c input) string {
	b := prepare(c)
	capacity := b.Cap()
	if capacity == 0 {
		capacity = -1
	}
	identity := 0
	if c.touched || c.reserve != 0 || c.initial != "" {
		identity = 1
	}
	nilSlice := 0
	if c.nilSlice {
		nilSlice = 1
	}
	row := fmt.Appendf(nil, "%d %d %d %d %d", c.mode, identity, capacity, len(c.initial), nilSlice)
	return string(appendBytes(appendBytes(row, c.initial), c.appended))
}

func describe(status string, identity bool, contents string) string {
	return string(appendBytes(fmt.Appendf(nil, "%s %t", status, identity), contents))
}

func actual(c input) (row string) {
	origin := prepare(c)
	b := origin
	var receiverCopy strings.Builder
	switch c.mode {
	case 1:
		receiverCopy = *origin
		b = &receiverCopy
	case 2:
		b = nil
	}
	snapshot := origin.String()
	defer func() {
		if snapshot != c.initial {
			log.Fatal("write changed a previously returned string")
		}
		if c.mode == 1 && origin.String() != c.initial {
			log.Fatal("write to copied receiver changed the original")
		}
		if p := recover(); p != nil {
			_, runtimePanic := p.(runtime.Error)
			switch {
			case c.mode == 2 && runtimePanic:
				row = "nil"
			case c.mode == 1 && p == "strings: illegal use of non-zero Builder copied by value":
				row = describe("copied", identified(b), b.String())
			default:
				panic(p)
			}
		}
	}()
	var p []byte
	if !c.nilSlice {
		p = []byte(c.appended)
	}
	n, err := b.Write(p)
	if err != nil {
		log.Fatalf("Write returned %v", err)
	}
	pointer := reflect.ValueOf(b)
	if pointer.Elem().FieldByName("addr").Pointer() != pointer.Pointer() {
		log.Fatal("successful Write did not record the receiver identity")
	}
	for i := range p {
		p[i] ^= 0xff
	}
	return describe(fmt.Sprintf("ok %d", n), true, b.String())
}

func expected(c input) string {
	if c.mode == 2 {
		return "nil"
	}
	if c.mode == 1 && (c.touched || c.reserve != 0 || c.initial != "") {
		return describe("copied", true, c.initial)
	}
	return describe(fmt.Sprintf("ok %d", len(c.appended)), true, c.initial+c.appended)
}
