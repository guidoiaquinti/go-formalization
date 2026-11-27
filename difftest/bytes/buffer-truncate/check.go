package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	stored, spare []byte
	offset        int
	receiverTag   int
	runeRead      bool
	count         int
}

func main() {
	var edges []input
	for _, count := range []int{math.MinInt, -1, 0, 1, math.MaxInt} {
		edges = append(edges, input{receiverTag: 0, count: count}, input{receiverTag: 1, count: count}, input{receiverTag: 2, count: count})
	}
	for length := 0; length <= 16; length++ {
		for offset := 0; offset <= length; offset++ {
			for spare := 0; spare <= 8; spare++ {
				unread := length - offset
				for _, count := range []int{math.MinInt, -1, 0, 1, unread / 2, unread, unread + 1, math.MaxInt} {
					edges = append(edges, input{fill(length, 17), fill(spare, 95), offset, 2, false, count})
				}
			}
		}
	}
	for _, token := range []string{"a", "é", "�", "🚀", "\xff"} {
		_, width := utf8.DecodeRuneInString(token)
		for _, count := range []int{-1, 0, 1, 2, 3} {
			edges = append(edges, input{[]byte(token + "XY"), []byte{9, 8}, width, 2, true, count})
		}
	}
	check.CompareSamples("godiff-bytes-buffer-truncate", edges, sample, encode, native, expected)
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
	if len(c.stored) > 0 && i%3 == 0 {
		c.runeRead = true
		_, c.offset = utf8.DecodeRune(c.stored)
	}
	if i%17 == 0 {
		c.stored, c.spare = nil, nil
		c.offset = 0
		c.receiverTag = 1
		c.runeRead = false
	}
	unread := len(c.stored) - c.offset
	c.count = r.Intn(unread+5) - 2
	return c
}

func encode(c input) string {
	if c.receiverTag == 0 {
		return fmt.Sprintf("0 %d", c.count)
	}
	lastRead := 0
	if c.runeRead {
		lastRead = c.offset
	} else if c.offset > 0 {
		lastRead = -1
	}
	row := fmt.Appendf(nil, "%d %d %d %d %d", c.receiverTag, len(c.stored), c.offset, lastRead, c.count)
	for _, block := range [][]byte{c.stored, c.spare} {
		for _, b := range block {
			row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
		}
	}
	return string(row)
}

func prepare(c input) (*bytes.Buffer, []byte) {
	switch c.receiverTag {
	case 0:
		return nil, nil
	case 1:
		return new(bytes.Buffer), nil
	}
	storage := make([]byte, len(c.stored)+len(c.spare))
	copy(storage, c.stored)
	copy(storage[len(c.stored):], c.spare)
	buffer := bytes.NewBuffer(storage[:len(c.stored)])
	if c.runeRead {
		_, width, err := buffer.ReadRune()
		if err != nil || width != c.offset {
			panic("unexpected ReadRune setup")
		}
	} else {
		buffer.Next(c.offset)
	}
	return buffer, storage
}

func apply(buffer *bytes.Buffer, count int) (status string) {
	defer func() {
		if value := recover(); value != nil {
			if value == "bytes.Buffer: truncation out of range" && buffer != nil {
				status = "bounds"
				return
			}
			err, ok := value.(runtime.Error)
			if buffer != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			status = "nil-receiver"
		}
	}()
	buffer.Truncate(count)
	return "ok"
}

func native(c input) string {
	buffer, _ := prepare(c)
	status := apply(buffer, c.count)
	if buffer == nil {
		return status
	}
	view := buffer.Bytes()
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	offset := buffer.Cap() - cap(view)
	capacity := cap(view)
	unread := bytes.Clone(view)
	spare := bytes.Clone(view[len(view):cap(view)])
	unreadByteRejected := buffer.UnreadByte() != nil
	runeBuffer, _ := prepare(c)
	apply(runeBuffer, c.count)
	unreadRuneRejected := runeBuffer.UnreadRune() != nil
	return describe(status, marker, offset, capacity, unread, spare, unreadByteRejected, unreadRuneRejected)
}

func expected(c input) string {
	if c.receiverTag == 0 {
		return "nil-receiver"
	}
	marker := "slice"
	if c.receiverTag == 1 {
		marker = "nil"
	}
	status := "ok"
	offset, end := c.offset, len(c.stored)
	switch {
	case c.count == 0:
		offset, end = 0, 0
	case c.count < 0 || c.count > len(c.stored)-offset:
		status = "bounds"
	default:
		end = offset + c.count
	}
	spare := append(bytes.Clone(c.stored[end:]), c.spare...)
	return describe(status, marker, offset, len(c.stored)+len(c.spare)-offset, c.stored[offset:end], spare, true, true)
}

func describe(status, marker string, offset, capacity int, unread, spare []byte, unreadByteRejected, unreadRuneRejected bool) string {
	return fmt.Sprintf("%s:%s:%d:%d:%s:%s:%t:%t", status, marker, offset, capacity,
		wireBytes(unread), wireBytes(spare), unreadByteRejected, unreadRuneRejected)
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
