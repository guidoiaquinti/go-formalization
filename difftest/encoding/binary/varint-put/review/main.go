package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type fixture struct {
	word             uint64
	offset, length   int
	backing          []byte
	unsigned, signed string
}

func main() {
	cases := []fixture{
		{0, 0, 0, nil, "bounds:", "bounds:"},
		{63, 1, 1, []byte{7, 9, 55}, "1:7,63,55", "1:7,126,55"},
		{64, 1, 1, []byte{7, 9, 55}, "1:7,64,55", "bounds:7,128,55"},
		{128, 1, 1, []byte{7, 9, 55}, "bounds:7,128,55", "bounds:7,128,55"},
		{16384, 1, 2, []byte{7, 9, 55, 66}, "bounds:7,128,128,66", "bounds:7,128,128,66"},
		{^uint64(0), 1, 1, []byte{7, 9, 55}, "bounds:7,255,55", "1:7,1,55"},
		{1 << 63, 0, 2, []byte{22, 33, 44}, "bounds:128,128,44", "bounds:255,255,44"},
		{^uint64(63), 0, 1, []byte{89, 90}, "bounds:192,90", "1:127,90"},
		{^uint64(64), 0, 1, []byte{89, 90}, "bounds:191,90", "bounds:129,90"},
	}
	var rows bytes.Buffer
	unsignedPrechecks, signedPrechecks, signedFold, capacityFailures := 0, 0, 0, 0
	for i, c := range cases {
		unsigned := apply(c, func(dst []byte) int { return binary.PutUvarint(dst, c.word) })
		signed := apply(c, func(dst []byte) int { return binary.PutVarint(dst, int64(c.word)) })
		if unsigned != c.unsigned || signed != c.signed {
			log.Fatalf("native %d: unsigned=%s want %s signed=%s want %s", i, unsigned, c.unsigned, signed, c.signed)
		}
		for index, payload := range [][]byte{binary.AppendUvarint(nil, c.word), binary.AppendVarint(nil, int64(c.word))} {
			got := apply(c, func(dst []byte) int {
				if len(dst) < len(payload) {
					return -1
				}
				return copy(dst, payload)
			})
			if index == 0 && got != c.unsigned {
				unsignedPrechecks++
			}
			if index == 1 && got != c.signed {
				signedPrechecks++
			}
		}
		if apply(c, func(dst []byte) int { return binary.PutUvarint(dst, c.word<<1) }) != c.signed {
			signedFold++
		}
		if apply(c, func(dst []byte) int { return binary.PutUvarint(dst[:cap(dst)], c.word) }) != c.unsigned {
			capacityFailures++
		}
		fmt.Fprint(&rows, c.word, " ", c.offset, " ", c.length)
		for _, b := range c.backing {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
	}
	if unsignedPrechecks != 6 || signedPrechecks != 5 || signedFold != 4 || capacityFailures != 6 {
		log.Fatalf("mutation counts precheck unsigned/signed=%d/%d signed=%d capacity=%d", unsignedPrechecks, signedPrechecks, signedFold, capacityFailures)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-binary-varint-put")
	cmd.Stdin = &rows
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	results := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(results) != len(cases) {
		log.Fatalf("got %d Lean rows", len(results))
	}
	for i, got := range results {
		if want := cases[i].unsigned + "|" + cases[i].signed; got != want {
			log.Fatalf("Lean fixture %d got %q want %q", i, got, want)
		}
	}
	fmt.Printf("Put varint reviewer control: %d exact native/Lean paired fixtures, precheck mutations fail %d/%d, signed-fold fails %d, capacity-bound fails %d.\n", len(cases), unsignedPrechecks, signedPrechecks, signedFold, capacityFailures)
}

func apply(c fixture, operation func([]byte) int) (result string) {
	backing := slices.Clone(c.backing)
	count := -1
	defer func() {
		if fault := recover(); fault != nil {
			err, ok := fault.(runtime.Error)
			if !ok || !strings.HasPrefix(err.Error(), "runtime error: index out of range") {
				panic(fault)
			}
		}
		status := "bounds"
		if count >= 0 {
			status = strconv.Itoa(count)
		}
		result = status + ":" + wire(backing)
	}()
	count = operation(backing[c.offset : c.offset+c.length : len(backing)])
	return
}

func wire(raw []byte) string {
	items := make([]string, len(raw))
	for i, b := range raw {
		items[i] = strconv.Itoa(int(b))
	}
	return strings.Join(items, ",")
}
