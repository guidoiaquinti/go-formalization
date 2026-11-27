package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	initial                string
	reserve, request, mode int
	touched                bool
}

func main() {
	var edges []input
	for _, initial := range []string{"", "abc", "é🚀", "\x00\xff"} {
		for _, reserve := range []int{0, 1, 8, 64, 1024, 32768} {
			for mode := range 3 {
				for _, touched := range []bool{false, true} {
					c := input{initial: initial, reserve: reserve, mode: mode, touched: touched}
					b := prepare(c)
					free := b.Cap() - b.Len()
					for _, n := range []int{math.MinInt, -1, 0, 1, free - 1, free, free + 1, 4097, 32768, 32769} {
						c.request = n
						edges = append(edges, c)
					}
				}
			}
		}
	}
	check.CompareSamples("godiff-strings-builder-grow", edges, sample, encode, actual, expected)
}

func sample(r *rand.Rand, i int) input {
	data := make([]byte, r.Intn(129))
	r.Read(data)
	c := input{initial: string(data), request: r.Intn(4097) - 512, mode: r.Intn(3), touched: i%5 == 0}
	if i%3 == 0 {
		c.reserve = r.Intn(513)
	}
	if i%4 == 0 {
		b := prepare(c)
		c.request = b.Cap() - b.Len() + i%3 - 1
	}
	return c
}

func prepare(c input) *strings.Builder {
	b := new(strings.Builder)
	if c.touched || c.reserve != 0 {
		b.Grow(c.reserve)
	}
	if c.initial != "" {
		b.WriteString(c.initial)
	}
	return b
}

func appendBytes(row []byte, value string) []byte {
	for i := range len(value) {
		row = strconv.AppendUint(append(row, ' '), uint64(value[i]), 10)
	}
	return row
}

func identity(c input) bool { return c.touched || c.reserve != 0 || c.initial != "" }

// On an empty byte slice, Grow uses a single no-scan allocation with the
// same roundupsize input as MakeNoZero, without invoking Builder.Grow.
func roundedCapacity(request int) int { return cap(slices.Grow([]byte(nil), request)) }

func encode(c input) string {
	b := prepare(c)
	allocationRequest, rounded := 0, 0
	if c.mode != 2 && !(c.mode == 1 && identity(c)) && c.request > b.Cap()-b.Len() {
		allocationRequest = 2*b.Cap() + c.request
		rounded = roundedCapacity(allocationRequest)
	}
	capacity := b.Cap()
	if capacity == 0 {
		capacity = -1
	}
	id := 0
	if identity(c) {
		id = 1
	}
	row := fmt.Appendf(nil, "%d %d %d %d %d %d", c.mode, id, capacity, c.request, allocationRequest, rounded)
	return string(appendBytes(row, c.initial))
}

func describe(status string, identified bool, capacity int, contents string) string {
	return string(appendBytes(fmt.Appendf(nil, "%s %t %d", status, identified, capacity), contents))
}

func actual(c input) (row string) {
	origin := prepare(c)
	b := origin
	var copied strings.Builder
	switch c.mode {
	case 1:
		copied = *origin
		b = &copied
	case 2:
		b = nil
	}
	snapshot := origin.String()
	priorCap := origin.Cap()
	observed := func(status string) string {
		addr := reflect.ValueOf(b).Elem().FieldByName("addr")
		return describe(status, !addr.IsNil(), b.Cap(), b.String())
	}
	defer func() {
		if snapshot != c.initial {
			log.Fatal("Grow changed a previously returned string")
		}
		if c.mode == 1 && (origin.String() != c.initial || origin.Cap() != priorCap) {
			log.Fatal("Grow on a copied builder changed the original")
		}
		if p := recover(); p != nil {
			_, runtimePanic := p.(runtime.Error)
			switch {
			case c.mode == 2 && runtimePanic:
				row = "nil"
			case c.mode == 1 && p == "strings: illegal use of non-zero Builder copied by value":
				row = observed("copied")
			case p == "strings.Builder.Grow: negative count":
				row = observed("negative")
			default:
				panic(p)
			}
		}
	}()
	b.Grow(c.request)
	pointer := reflect.ValueOf(b)
	if pointer.Elem().FieldByName("addr").Pointer() != pointer.Pointer() {
		log.Fatal("successful Grow did not record receiver identity")
	}
	return observed("ok")
}

func expected(c input) string {
	if c.mode == 2 {
		return "nil"
	}
	b := prepare(c)
	capacity := b.Cap()
	if c.mode == 1 && identity(c) {
		return describe("copied", true, capacity, c.initial)
	}
	if c.request < 0 {
		return describe("negative", true, capacity, c.initial)
	}
	if capacity-len(c.initial) < c.request {
		capacity = roundedCapacity(2*capacity + c.request)
	}
	return describe("ok", true, capacity, c.initial)
}
