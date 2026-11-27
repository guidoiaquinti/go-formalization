package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	stored, spare []byte
	offset        int
	receiverTag   int
}

func main() {
	edges := []input{{receiverTag: 0}, {receiverTag: 1}, {receiverTag: 2}}
	for length := 0; length <= 16; length++ {
		for offset := 0; offset <= length; offset++ {
			for spare := 0; spare <= 16; spare++ {
				edges = append(edges, input{fill(length, 17), fill(spare, 95), offset, 2})
			}
		}
	}
	check.CompareSamples("godiff-bytes-buffer-string", edges, sample, encode, native, expected)
}

func fill(length int, start byte) []byte {
	out := make([]byte, length)
	for i := range out {
		out[i] = start + byte(i)
	}
	return out
}

func sample(r *rand.Rand, i int) input {
	c := input{stored: make([]byte, r.Intn(64)), spare: make([]byte, r.Intn(32)), receiverTag: 2}
	r.Read(c.stored)
	r.Read(c.spare)
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
	row := fmt.Appendf(nil, "%d %d %d", c.receiverTag, len(c.stored), c.offset)
	for _, block := range [][]byte{c.stored, c.spare} {
		for _, b := range block {
			row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
		}
	}
	return string(row)
}

func native(c input) string {
	var buffer *bytes.Buffer
	switch c.receiverTag {
	case 0:
		return wireBytes([]byte(buffer.String()))
	case 1:
		buffer = new(bytes.Buffer)
	case 2:
		storage := make([]byte, len(c.stored)+len(c.spare))
		copy(storage, c.stored)
		copy(storage[len(c.stored):], c.spare)
		buffer = bytes.NewBuffer(storage[:len(c.stored)])
		buffer.Next(c.offset)
	}

	original := buffer.String()
	view := buffer.Bytes()
	for i := range view {
		view[i] ^= 1
	}
	changed := buffer.String()
	overwrite := bytes.Repeat([]byte{0xa5}, len(c.stored)+len(c.spare))
	buffer.Reset()
	buffer.Write(overwrite)
	return describe(original, changed, buffer.String())
}

func expected(c input) string {
	if c.receiverTag == 0 {
		return wireBytes([]byte("<nil>"))
	}
	unread := c.stored[c.offset:]
	mutated := make([]byte, len(unread))
	for i, value := range unread {
		mutated[i] = value ^ 1
	}
	overwrite := bytes.Repeat([]byte{0xa5}, len(c.stored)+len(c.spare))
	return describe(string(unread), string(mutated), string(overwrite))
}

func describe(original, changed, current string) string {
	return fmt.Sprintf("%s:%s:%s", wireBytes([]byte(original)), wireBytes([]byte(changed)),
		wireBytes([]byte(current)))
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
