package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := []string{"", "é", "\xff\x80\x00", "é\xff\x00"}
	for n := 1; n <= 512; n++ {
		raw := make([]byte, n)
		for i := range raw {
			raw[i] = byte(i*103 + n)
		}
		edges = append(edges, string(raw), strings.Repeat("\xff", n))
	}
	check.CompareSamples("godiff-bytes-new-buffer-string", edges, sample, encode, native, expected)
}
func sample(rng *rand.Rand, i int) string {
	raw := make([]byte, rng.Intn(513))
	rng.Read(raw)
	if i%3 == 0 {
		for j := 0; j < len(raw); j += 5 {
			raw[j] = 0
		}
	}
	return string(raw)
}
func encode(s string) string {
	row := []byte("0")
	for i := range len(s) {
		row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
	}
	return string(row)
}

//go:noinline
func constructed(s string) *bytes.Buffer { return bytes.NewBufferString(s) }

//go:noinline
func firstByte(s string) byte { return s[0] }

func native(s string) string {
	b := constructed(s)
	view := b.Bytes()
	if b.UnreadByte() == nil || b.UnreadRune() == nil {
		panic("NewBufferString accepted unread before reading")
	}
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	result := fmt.Sprintf("%s:%d:%d:%d:0:%s:%s", marker, b.Len(), b.Cap(), b.Cap()-cap(view), wire(view), wire(view[len(view):cap(view)]))
	if len(view) > 0 {
		original := firstByte(s)
		view[0] ^= 0xff
		if firstByte(s) != original {
			panic("NewBufferString destination aliases input")
		}
	}
	return result
}
func expected(s string) string {
	if len(s) == 0 {
		return "slice:0:0:0:0::"
	}
	classes := [...]int{8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512}
	for _, capacity := range classes {
		if capacity >= len(s) {
			return fmt.Sprintf("slice:%d:%d:0:0:%s:%s", len(s), capacity, wire([]byte(s)), wire(make([]byte, capacity-len(s))))
		}
	}
	panic("input exceeds tested capacity classes")
}
func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}
