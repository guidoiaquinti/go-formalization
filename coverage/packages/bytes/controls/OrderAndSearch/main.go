package main

import (
	"bytes"
	"cmp"
	"fmt"
	"log"
)

func main() {
	backing := []byte{0x80, 0x00, 0xff, 0xff}
	compareCases := []struct {
		a, b []byte
		want int
	}{
		{nil, []byte{}, 0},
		{[]byte{}, nil, 0},
		{[]byte{0x80}, []byte{0x7f}, 1},
		{[]byte{0x00}, []byte{0xff}, -1},
		{[]byte{0xff}, []byte{0x00, 0xff}, 1},
		{[]byte{0xff, 0x00}, []byte{0xfe}, 1},
		{backing[:1], backing[:2], -1},
		{backing[:3], backing[:3:3], 0},
		{backing[2:3], backing[3:], 0},
	}
	var nilSensitive, lengthEquality, signedBytes, lengthFirst int
	for i, c := range compareCases {
		if eq, order := bytes.Equal(c.a, c.b), bytes.Compare(c.a, c.b); eq != (c.want == 0) || order != c.want {
			log.Fatalf("comparison %d: Equal=%t Compare=%d, want order %d", i, eq, order, c.want)
		}
		if (bytes.Equal(c.a, c.b) && (c.a == nil) == (c.b == nil)) != (c.want == 0) {
			nilSensitive++
		}
		if (len(c.a) == len(c.b)) != (c.want == 0) {
			lengthEquality++
		}
		signed := cmp.Compare(len(c.a), len(c.b))
		for j := range min(len(c.a), len(c.b)) {
			if c.a[j] != c.b[j] {
				signed = cmp.Compare(int8(c.a[j]), int8(c.b[j]))
				break
			}
		}
		if signed != c.want {
			signedBytes++
		}
		order := cmp.Compare(len(c.a), len(c.b))
		if order == 0 {
			order = bytes.Compare(c.a, c.b)
		}
		if order != c.want {
			lengthFirst++
		}
	}

	searchCases := []struct {
		data        []byte
		key         byte
		first, last int
	}{
		{nil, 0, -1, -1},
		{[]byte{}, 0xff, -1, -1},
		{backing[:1], 0xff, -1, -1},
		{backing[1:3], 0xff, 1, 1},
		{backing, 0xff, 2, 3},
		{[]byte{0, 0xff, 0xff, 0}, 0, 0, 3},
		{[]byte{1, 2, 1, 1, 2}, 1, 0, 3},
		{[]byte{0xab}, 0xab, 0, 0},
	}
	var swapMatches, zeroIfMissing, wrongReverseOffset int
	for i, c := range searchCases {
		first, last := bytes.IndexByte(c.data, c.key), bytes.LastIndexByte(c.data, c.key)
		if first != c.first || last != c.last {
			log.Fatalf("search %d: first=%d last=%d, want %d and %d", i, first, last, c.first, c.last)
		}
		if first != last {
			swapMatches++
		}
		if first == -1 {
			zeroIfMissing++
		} else if len(c.data)-1-last != c.last {
			wrongReverseOffset++
		}
	}
	if nilSensitive != 2 || lengthEquality != 2 || signedBytes != 3 || lengthFirst != 1 ||
		swapMatches != 3 || zeroIfMissing != 3 || wrongReverseOffset != 4 {
		log.Fatalf("unexpected mutation failures: nil=%d equality=%d signed=%d length-first=%d swap=%d missing=%d reverse=%d",
			nilSensitive, lengthEquality, signedBytes, lengthFirst, swapMatches, zeroIfMissing, wrongReverseOffset)
	}
	fmt.Printf("bytes: %d comparison and %d search fixtures pass, mutant failures nil=%d equality=%d signed=%d length-first=%d swap=%d missing=%d reverse=%d\n",
		len(compareCases), len(searchCases), nilSensitive, lengthEquality, signedBytes, lengthFirst, swapMatches, zeroIfMissing, wrongReverseOffset)
}
