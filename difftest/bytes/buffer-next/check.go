package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
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
	check.CompareSamples("godiff-bytes-buffer-next", edges, sample, encode, native, expected)
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

func apply(buffer *bytes.Buffer, count int) (view []byte, status string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok {
				panic(value)
			}
			switch {
			case buffer == nil && err.Error() == "runtime error: invalid memory address or nil pointer dereference":
				status = "nil-receiver"
			case buffer != nil && strings.HasPrefix(err.Error(), "runtime error: slice bounds out of range"):
				status = "bounds"
			default:
				panic(value)
			}
		}
	}()
	return buffer.Next(count), "ok"
}

func native(c input) string {
	buffer, storage := prepare(c)
	view, status := apply(buffer, c.count)
	if buffer == nil {
		return status
	}
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	visible := bytes.Clone(view)
	capacity := cap(view)
	for i := range view {
		view[i] ^= 1
	}
	unread := bytes.Clone(buffer.Bytes())
	backing := bytes.Clone(storage)
	offset := buffer.Cap() - cap(buffer.Bytes())
	unreturned := "rejected"
	if buffer.UnreadByte() == nil {
		b, err := buffer.ReadByte()
		if err != nil {
			panic(err)
		}
		unreturned = strconv.Itoa(int(b))
	}
	runeBuffer, _ := prepare(c)
	apply(runeBuffer, c.count)
	unreadRuneRejected := runeBuffer.UnreadRune() != nil
	return describe(status, marker, capacity, visible, offset, unread, backing, unreturned, unreadRuneRejected)
}

func expected(c input) string {
	if c.receiverTag == 0 {
		return "nil-receiver"
	}
	status, marker := "ok", "slice"
	capacity, size := 0, 0
	if c.count < 0 {
		status, marker = "bounds", "nil"
	} else {
		size = min(c.count, len(c.stored)-c.offset)
		capacity = len(c.stored) + len(c.spare) - c.offset
		if c.receiverTag == 1 {
			marker = "nil"
		}
	}
	offset := c.offset + size
	visible := c.stored[c.offset:offset]
	backing := append(bytes.Clone(c.stored), c.spare...)
	for i := c.offset; i < offset; i++ {
		backing[i] ^= 1
	}
	unreturned := "rejected"
	if size > 0 {
		unreturned = strconv.Itoa(int(backing[offset-1]))
	}
	return describe(status, marker, capacity, visible, offset, c.stored[offset:], backing, unreturned, true)
}

func describe(status, marker string, capacity int, visible []byte, offset int, unread, backing []byte, unreturned string, unreadRuneRejected bool) string {
	return fmt.Sprintf("%s:%s:%d:%s:%d:%s:%s:%s:%t", status, marker, capacity, wireBytes(visible),
		offset, wireBytes(unread), wireBytes(backing), unreturned, unreadRuneRejected)
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
