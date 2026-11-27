package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	initial, consume, append, spare int
	reset, nilReceiver              bool
}

func geometry(c input) (length, offset int) {
	if c.reset {
		return c.append, 0
	}
	return c.initial + c.append, c.consume
}

func main() {
	var edges []input
	for initial := 0; initial <= 64; initial++ {
		for consume := 0; consume <= initial; consume++ {
			for _, count := range []int{0, 1, 7, 65} {
				for _, reset := range []bool{false, true} {
					edges = append(edges, input{initial, consume, count, (initial + consume) % 8, reset, false})
				}
			}
		}
	}
	edges = append(edges, input{nilReceiver: true})
	check.CompareSamples("godiff-bytes-buffer-len", edges, sample, encode, native, expected)
}

func sample(r *rand.Rand, i int) input {
	n := r.Intn(512)
	return input{n, r.Intn(n + 1), r.Intn(80), r.Intn(32), i%3 == 0, false}
}

func encode(c input) string {
	if c.nilReceiver {
		return "0"
	}
	length, offset := geometry(c)
	return fmt.Sprintf("1 %d %d", length, offset)
}

func expected(c input) string {
	if c.nilReceiver {
		return "panic"
	}
	length, offset := geometry(c)
	return strconv.Itoa(length - offset)
}

func native(c input) (output string) {
	var b *bytes.Buffer
	if !c.nilReceiver {
		if c.initial == 0 && c.append+c.spare == 0 {
			b = new(bytes.Buffer)
		} else {
			b = bytes.NewBuffer(make([]byte, c.initial, c.initial+c.append+c.spare))
		}
		b.Next(c.consume)
		if c.reset {
			b.Reset()
		}
		if c.append > 0 {
			b.Write(make([]byte, c.append))
		}
	}

	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !c.nilReceiver || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			output = "panic"
		}
	}()
	got := b.Len()
	if got != len(b.Bytes()) {
		panic(fmt.Sprintf("Buffer.Len()=%d, len(Buffer.Bytes())=%d", got, len(b.Bytes())))
	}
	return strconv.Itoa(got)
}
