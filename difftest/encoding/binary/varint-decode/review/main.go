package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strings"
)

type fixture struct {
	input []byte
	u     uint64
	s     int64
	count int
}

func main() {
	continued := bytes.Repeat([]byte{128}, 10)
	maxWord := append(bytes.Repeat([]byte{255}, 9), 1)
	maxSigned := bytes.Clone(maxWord)
	maxSigned[0] = 254
	cases := []fixture{
		{nil, 0, 0, 0},
		{[]byte{0, 255, 255}, 0, 0, 1},
		{[]byte{0x81, 0, 0xff}, 1, -1, 2},
		{continued, 0, 0, 0},
		{append(bytes.Clone(continued), 0), 0, 0, -11},
		{append(bytes.Clone(continued), 255, 0), 0, 0, -11},
		{append(bytes.Clone(continued[:9]), 2, 128), 0, 0, -10},
		{append(bytes.Clone(continued[:9]), 1, 128), 1 << 63, 1 << 62, 10},
		{maxWord, math.MaxUint64, math.MinInt64, 10},
		{maxSigned, math.MaxUint64 - 1, math.MaxInt64, 10},
		{continued[:9], 0, 0, 0},
	}
	var input bytes.Buffer
	tenthFailures, countFailures, signFailures := 0, 0, 0
	for i, c := range cases {
		u, n := binary.Uvarint(c.input)
		s, m := binary.Varint(c.input)
		if u != c.u || s != c.s || n != c.count || m != c.count {
			log.Fatalf("fixture %d: native %d:%d|%d:%d", i, u, n, s, m)
		}
		if prematureCount(c.input) != c.count {
			tenthFailures++
		}
		if c.count > 0 && len(c.input) != c.count {
			countFailures++
		}
		if c.count > 0 && int64(c.u>>1) != c.s {
			signFailures++
		}
		fmt.Fprint(&input, len(c.input))
		for _, b := range c.input {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	if tenthFailures != 3 || countFailures != 3 || signFailures != 2 {
		log.Fatalf("mutation counts tenth=%d length=%d sign=%d", tenthFailures, countFailures, signFailures)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-binary-varint-decode")
	cmd.Stdin = &input
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	rows := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(rows) != len(cases) {
		log.Fatalf("got %d Lean rows", len(rows))
	}
	for i, got := range rows {
		c := cases[i]
		if want := fmt.Sprintf("%d:%d|%d:%d", c.u, c.count, c.s, c.count); got != want {
			log.Fatalf("Lean fixture %d got %q want %q", i, got, want)
		}
	}
	fmt.Printf("Varint decoder reviewer control: %d exact native and Lean fixtures; early tenth-byte rejection, full-input count and missing sign fail %d/%d/%d.\n", len(cases), tenthFailures, countFailures, signFailures)
}

func prematureCount(input []byte) int {
	if len(input) >= 10 {
		allContinued := true
		for _, b := range input[:10] {
			allContinued = allContinued && b >= 128
		}
		if allContinued {
			return -10
		}
	}
	_, count := binary.Uvarint(input)
	return count
}
