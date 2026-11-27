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
	stored, spare, payload []byte
	offset                 int
	receiverTag            int
}

func main() {
	edges := []input{{receiverTag: 0}, {receiverTag: 1}, {receiverTag: 1, payload: []byte{0xff}}, {receiverTag: 2}}
	for length := 0; length <= 16; length++ {
		for offset := 0; offset <= length; offset++ {
			for spare := 0; spare <= 16; spare++ {
				for _, count := range []int{0, max(0, spare-1), spare, spare + 1} {
					edges = append(edges, input{fill(length, 17), fill(spare, 95), fill(count, 219), offset, 2})
				}
			}
		}
	}
	check.CompareSamples("godiff-bytes-buffer-available-buffer", edges, sample, encode, native, expected)
}

func fill(length int, start byte) []byte {
	out := make([]byte, length)
	for i := range out {
		out[i] = start + byte(i)
	}
	return out
}

func sample(r *rand.Rand, i int) input {
	c := input{stored: make([]byte, r.Intn(64)), spare: make([]byte, r.Intn(32)), payload: make([]byte, r.Intn(48)), receiverTag: 2}
	r.Read(c.stored)
	r.Read(c.spare)
	r.Read(c.payload)
	c.offset = r.Intn(len(c.stored) + 1)
	if i%17 == 0 {
		c.stored, c.spare = nil, nil
		c.offset = 0
		c.receiverTag = 1
	}
	return c
}

func encode(c input) string {
	if c.receiverTag == 0 {
		return "0"
	}
	row := fmt.Appendf(nil, "%d %d %d %d", c.receiverTag, len(c.stored), len(c.spare), c.offset)
	for _, block := range [][]byte{c.stored, c.spare, c.payload} {
		for _, b := range block {
			row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
		}
	}
	return string(row)
}

func native(c input) (output string) {
	var buffer *bytes.Buffer
	switch c.receiverTag {
	case 1:
		buffer = new(bytes.Buffer)
	case 2:
		storage := make([]byte, len(c.stored)+len(c.spare))
		copy(storage, c.stored)
		copy(storage[len(c.stored):], c.spare)
		buffer = bytes.NewBuffer(storage[:len(c.stored)])
		buffer.Next(c.offset)
	}
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if c.receiverTag != 0 || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			output = "panic"
		}
	}()
	view := buffer.AvailableBuffer()
	before := bytes.Clone(view[:cap(view)])
	length, capacity := len(view), cap(view)
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	data := append(view, c.payload...)
	viewAfter := buffer.AvailableBuffer()
	tailAfter := bytes.Clone(viewAfter[:cap(viewAfter)])
	unread := bytes.Clone(buffer.Bytes())
	n, err := buffer.Write(data)
	if err != nil {
		panic(err)
	}
	return describe(marker, length, capacity, before, tailAfter, unread, n, buffer.Bytes())
}

func expected(c input) string {
	if c.receiverTag == 0 {
		return "panic"
	}
	marker := "slice"
	if c.receiverTag == 1 {
		marker = "nil"
	}
	tailAfter := bytes.Clone(c.spare)
	if len(c.payload) <= len(c.spare) {
		copy(tailAfter, c.payload)
	}
	unread := c.stored[c.offset:]
	result := append(bytes.Clone(unread), c.payload...)
	return describe(marker, 0, len(c.spare), c.spare, tailAfter, unread, len(c.payload), result)
}

func describe(marker string, length, capacity int, before, after, unread []byte, n int, result []byte) string {
	return fmt.Sprintf("%s:%d:%d:%s:%s:%s:%d:%s", marker, length, capacity,
		wireBytes(before), wireBytes(after), wireBytes(unread), n, wireBytes(result))
}

func wireBytes(input []byte) string {
	var row []byte
	for i, b := range input {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(b), 10)
	}
	return string(row)
}
