package main

import (
	"bytes"
	"fmt"
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
	check.CompareSamples("godiff-bytes-buffer-reset", edges, sample, encode, native, expected)
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

func native(c input) (output string) {
	buffer, storage := prepare(c)
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if c.receiverTag != 0 || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			output = "panic"
		}
	}()
	buffer.Reset()
	view := buffer.Bytes()
	backing := bytes.Clone(view[:cap(view)])
	length, capacity := len(view), cap(view)
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	unreadByteRejected := buffer.UnreadByte() != nil
	runeBuffer, _ := prepare(c)
	runeBuffer.Reset()
	unreadRuneRejected := runeBuffer.UnreadRune() != nil
	if capacity > 0 {
		buffer.WriteByte(0xa5)
	}
	return describe(marker, length, capacity, backing, unreadByteRejected, unreadRuneRejected, storage)
}

func expected(c input) string {
	if c.receiverTag == 0 {
		return "panic"
	}
	marker := "slice"
	if c.receiverTag == 1 {
		marker = "nil"
	}
	backing := append(bytes.Clone(c.stored), c.spare...)
	reused := bytes.Clone(backing)
	if len(reused) > 0 {
		reused[0] = 0xa5
	}
	return describe(marker, 0, len(backing), backing, true, true, reused)
}

func describe(marker string, length, capacity int, backing []byte, unreadByteRejected, unreadRuneRejected bool, reused []byte) string {
	return fmt.Sprintf("%s:%d:%d:%s:%t:%t:%s", marker, length, capacity, wireBytes(backing),
		unreadByteRejected, unreadRuneRejected, wireBytes(reused))
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
