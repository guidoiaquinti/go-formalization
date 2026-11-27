package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const word uint64 = 0x123456789abcdef0

type variant struct {
	name   string
	append func([]byte, uint64) []byte
	layout []byte
}

func render(data []byte) string {
	parts := make([]string, len(data))
	for i, value := range data {
		parts[i] = fmt.Sprint(value)
	}
	return strings.Join(parts, ",")
}

func main() {
	exe := "lean/.lake/build/bin/godiff-binary-endian-append"
	if len(os.Args) == 2 {
		exe = os.Args[1]
	}
	variants := []variant{
		{"big16", func(p []byte, n uint64) []byte { return binary.BigEndian.AppendUint16(p, uint16(n)) }, []byte{0xde, 0xf0}},
		{"big32", func(p []byte, n uint64) []byte { return binary.BigEndian.AppendUint32(p, uint32(n)) }, []byte{0x9a, 0xbc, 0xde, 0xf0}},
		{"big64", binary.BigEndian.AppendUint64, []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}},
		{"little16", func(p []byte, n uint64) []byte { return binary.LittleEndian.AppendUint16(p, uint16(n)) }, []byte{0xf0, 0xde}},
		{"little32", func(p []byte, n uint64) []byte { return binary.LittleEndian.AppendUint32(p, uint32(n)) }, []byte{0xf0, 0xde, 0xbc, 0x9a}},
		{"little64", binary.LittleEndian.AppendUint64, []byte{0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12}},
	}
	for index, v := range variants {
		width := len(v.layout)
		for _, extra := range []int{-1, 0, 1} {
			spare := width + extra
			before := append([]byte{0xa1, 0xb2, 0xc3}, bytes.Repeat([]byte{0x59}, spare)...)
			storage := slices.Clone(before)
			got := v.append(storage[1:3:len(storage)], word)
			want := append([]byte{0xb2, 0xc3}, v.layout...)
			reuse := extra >= 0
			old := slices.Clone(before)
			if reuse {
				copy(old[3:], v.layout)
			}
			if !bytes.Equal(got, want) || !bytes.Equal(storage, old) || (&got[0] == &storage[1]) != reuse {
				panic(fmt.Sprintf("%s spare=%d native got %x, backing %x", v.name, spare, got, storage))
			}
			row := fmt.Sprintf("%d 1 2", word)
			for _, value := range before {
				row += fmt.Sprintf(" %d", value)
			}
			cmd := exec.Command(exe)
			cmd.Stdin = strings.NewReader(row + "\n")
			raw, err := cmd.CombinedOutput()
			if err != nil {
				panic(fmt.Sprintf("%s: %v: %s", v.name, err, raw))
			}
			fields := strings.Split(strings.TrimSpace(string(raw)), "|")
			wantField := fmt.Sprintf("%t:%s:%s", reuse, render(want), render(old))
			if len(fields) != len(variants) || fields[index] != wantField {
				panic(fmt.Sprintf("%s spare=%d Lean got %q want %q", v.name, spare, raw, wantField))
			}
			// Subsequent writes through the returned slice distinguish reuse from growth.
			got[0] = 0xdd
			if (storage[1] == 0xdd) != reuse {
				panic(v.name + " alias write")
			}
		}
		for _, empty := range [][]byte{nil, make([]byte, 0)} {
			if got := v.append(empty, word); !bytes.Equal(got, v.layout) {
				panic(v.name + " empty destination")
			}
		}
	}
	fmt.Println("18 literal layout/alias Go and Lean fixtures and 12 native nil/empty fixtures passed")
}
