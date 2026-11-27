package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	stored, spare []byte
	leading       int
	nilSlice      bool
}

func main() {
	edges := []input{{nilSlice: true}}
	for length := range 33 {
		for spare := range 17 {
			for offset := range 8 {
				edges = append(edges, input{stored: pattern(length, offset+17), spare: pattern(spare, offset+31), leading: offset})
			}
		}
	}
	check.CompareSamples("godiff-bytes-new-buffer", edges, sample, encode, native, expected)
}
func pattern(n, start int) []byte {
	raw := make([]byte, n)
	for i := range raw {
		raw[i] = byte(start + i*127)
	}
	return raw
}
func sample(rng *rand.Rand, i int) input {
	if i%23 == 0 {
		return input{nilSlice: true}
	}
	stored := make([]byte, rng.Intn(513))
	spare := make([]byte, rng.Intn(257))
	rng.Read(stored)
	rng.Read(spare)
	return input{stored, spare, rng.Intn(32), false}
}
func encode(c input) string {
	if c.nilSlice {
		return "0 0"
	}
	row := fmt.Appendf(nil, "1 %d", len(c.stored))
	for _, part := range [][]byte{c.stored, c.spare} {
		for _, b := range part {
			row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
		}
	}
	return string(row)
}
func native(c input) string {
	var slice []byte
	var first *byte
	if !c.nilSlice {
		capacity := len(c.stored) + len(c.spare)
		backing := make([]byte, c.leading+capacity+5)
		for i := range backing {
			backing[i] = 0xee
		}
		copy(backing[c.leading:], c.stored)
		copy(backing[c.leading+len(c.stored):], c.spare)
		slice = backing[c.leading : c.leading+len(c.stored) : c.leading+capacity]
		if capacity > 0 {
			first = &slice[:capacity][0]
		}
	}
	b := bytes.NewBuffer(slice)
	view := b.Bytes()
	if b.Cap() > 0 && first != &view[:cap(view)][0] {
		panic("NewBuffer did not adopt backing")
	}
	if b.UnreadByte() == nil || b.UnreadRune() == nil {
		panic("NewBuffer accepted unread before reading")
	}
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	return fmt.Sprintf("%s:%d:%d:%d:0:%s:%s", marker, b.Len(), b.Cap(), b.Cap()-cap(view), wire(view), wire(view[len(view):cap(view)]))
}
func expected(c input) string {
	marker := "slice"
	if c.nilSlice {
		marker = "nil"
	}
	return fmt.Sprintf("%s:%d:%d:0:0:%s:%s", marker, len(c.stored), len(c.stored)+len(c.spare), wire(c.stored), wire(c.spare))
}
func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}
