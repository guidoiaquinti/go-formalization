package main

import (
	"bytes"
	"fmt"
	"io"
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
}

func main() {
	edges := []input{{receiverTag: 0}, {receiverTag: 1}, {receiverTag: 2}}
	for length := 0; length <= 16; length++ {
		for offset := 0; offset <= length; offset++ {
			for spare := 0; spare <= 16; spare++ {
				edges = append(edges, input{fill(length, 17), fill(spare, 95), offset, 2, false})
			}
		}
	}
	for _, token := range []string{"a", "é", "�", "🚀", "\xff"} {
		_, width := utf8.DecodeRuneInString(token)
		edges = append(edges, input{[]byte(token + "X"), []byte{9, 8}, width, 2, true})
	}
	for value := range 256 {
		edges = append(edges, input{[]byte{'q', byte(value)}, []byte{9, 8}, 1, 2, false})
		edges = append(edges, input{[]byte{byte(value)}, []byte{9, 8}, 1, 2, true})
	}
	for raw := range 1 << 16 {
		edges = append(edges, input{[]byte{byte(raw >> 8), byte(raw)}, []byte{7}, 0, 2, false})
	}
	for _, token := range []string{"é", "�", "🚀", "\xe0\xa0\x80", "\xed\xa0\x80", "\xf4\x8f\xbf\xbf", "\xf4\x90\x80\x80", "\xe2\x82", "\xf0\x9f\x9a", "\xc0\xaf", "\xff"} {
		edges = append(edges, input{[]byte("a" + token + "z"), []byte{8, 9}, 1, 2, false})
	}
	check.CompareSamples("godiff-bytes-buffer-read-rune", edges, sample, encode, native, expected)
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
	return c
}

func encode(c input) string {
	if c.receiverTag == 0 {
		return "0"
	}
	lastRead := 0
	if c.runeRead {
		lastRead = c.offset
	} else if c.offset > 0 {
		lastRead = -1
	}
	row := fmt.Appendf(nil, "%d %d %d %d", c.receiverTag, len(c.stored), c.offset, lastRead)
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

func readRune(buffer *bytes.Buffer) (r rune, width int, status string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if buffer != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			status = "nil-receiver"
		}
	}()
	r, width, err := buffer.ReadRune()
	switch err {
	case nil:
		return r, width, "ok"
	case io.EOF:
		return r, width, "eof"
	default:
		panic(err)
	}
}

func native(c input) string {
	buffer, _ := prepare(c)
	value, width, status := readRune(buffer)
	if buffer == nil {
		return status
	}
	view := buffer.Bytes()
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	offset, capacity := buffer.Cap()-cap(view), cap(view)
	unread, spare := bytes.Clone(view), bytes.Clone(view[len(view):cap(view)])
	byteBack := "rejected"
	if buffer.UnreadByte() == nil {
		read, err := buffer.ReadByte()
		if err != nil {
			panic(err)
		}
		byteBack = strconv.Itoa(int(read))
	}
	runeBuffer, _ := prepare(c)
	readRune(runeBuffer)
	runeBack := "rejected"
	if runeBuffer.UnreadRune() == nil {
		start := runeBuffer.Cap() - cap(runeBuffer.Bytes())
		retry, retryWidth, retryStatus := readRune(runeBuffer)
		if retryStatus != "ok" {
			panic("UnreadRune retry did not succeed")
		}
		runeBack = fmt.Sprintf("%d,%d,%d", start, retry, retryWidth)
	}
	return describe(status, value, width, marker, offset, capacity, unread, spare, byteBack, runeBack)
}

func firstRuneByRange(input []byte) (value rune, width int) {
	width = len(input)
	for offset, r := range string(input) {
		if offset != 0 {
			width = offset
			break
		}
		value = r
	}
	return value, width
}

func expected(c input) string {
	if c.receiverTag == 0 {
		return "nil-receiver"
	}
	marker := "slice"
	if c.receiverTag == 1 {
		marker = "nil"
	}
	if c.offset == len(c.stored) {
		backing := append(bytes.Clone(c.stored), c.spare...)
		return describe("eof", 0, 0, marker, 0, len(backing), nil, backing, "rejected", "rejected")
	}
	value, width := firstRuneByRange(c.stored[c.offset:])
	offset := c.offset + width
	byteBack := strconv.Itoa(int(c.stored[offset-1]))
	runeBack := fmt.Sprintf("%d,%d,%d", c.offset, value, width)
	return describe("ok", value, width, marker, offset, len(c.stored)+len(c.spare)-offset,
		c.stored[offset:], c.spare, byteBack, runeBack)
}

func describe(status string, value rune, width int, marker string, offset, capacity int, unread, spare []byte, byteBack, runeBack string) string {
	return fmt.Sprintf("%s:%d:%d:%s:%d:%d:%s:%s:%s:%s", status, value, width, marker, offset, capacity,
		wireBytes(unread), wireBytes(spare), byteBack, runeBack)
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
