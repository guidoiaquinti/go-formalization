package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type preparation int

const (
	readAndReslice preparation = iota
	compact
	grow
	reset
)

type input struct {
	buffer   *bytes.Buffer
	offset   int
	capacity int
}

func prepare(mode preparation, left, added, spare int) input {
	capacity := 2*(left+added) + spare
	total := capacity
	if mode == readAndReslice {
		total -= added + spare
	}
	if mode == grow {
		added += capacity + 1
	}

	storage := make([]byte, total, capacity)
	base := unsafe.SliceData(storage)
	b := bytes.NewBuffer(storage)
	consumed := total - left
	b.Next(consumed)
	if mode == reset {
		b.Reset()
		left = 0
	}
	b.Write(make([]byte, added))
	length := left + added
	if b.Len() != length {
		panic(fmt.Sprintf("mode %d prepared %d unread bytes, want %d", mode, b.Len(), length))
	}
	if mode == compact && unsafe.SliceData(b.Bytes()) != base {
		panic("compaction did not restore the base of the same backing array")
	}
	if mode == grow && (unsafe.SliceData(b.Bytes()) == base || cap(b.Bytes()) <= capacity) {
		panic("growth did not allocate a larger backing array")
	}
	offset := 0
	if mode == readAndReslice {
		offset = consumed
	}
	if mode == grow {
		capacity = cap(b.Bytes())
	}
	readAgain := spare % (length + 1)
	b.Next(readAgain)
	return input{b, offset + readAgain, capacity}
}

func main() {
	edges := []input{{}, {buffer: new(bytes.Buffer)}}
	for mode := readAndReslice; mode <= reset; mode++ {
		for left := 1; left <= 32; left++ {
			for added := 1; added <= 32; added++ {
				for _, spare := range []int{0, 1, 7} {
					edges = append(edges, prepare(mode, left, added, spare))
				}
			}
		}
	}
	check.CompareSamples("godiff-bytes-buffer-cap", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) input {
	return prepare(preparation(i%4), 1+r.Intn(96), 1+r.Intn(96), r.Intn(32))
}

func encode(c input) string {
	if c.buffer == nil {
		return "0"
	}
	return fmt.Sprintf("1 %d", c.capacity)
}

func expected(c input) string {
	if c.buffer == nil {
		return "panic"
	}
	view := c.buffer.Bytes()
	return strconv.Itoa(cap(view) + c.offset)
}

func native(c input) (output string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if c.buffer != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			output = "panic"
		}
	}()
	return strconv.Itoa(c.buffer.Cap())
}
