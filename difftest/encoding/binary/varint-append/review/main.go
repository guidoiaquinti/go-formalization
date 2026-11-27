package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os/exec"
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
		{0, 0, 0, nil, "false:0:", "false:0:"},
		{63, 0, 0, []byte{99, 88}, "true:63:63,88", "true:126:126,88"},
		{64, 1, 1, []byte{7, 9, 55}, "true:9,64:7,9,64", "false:9,128,1:7,9,128"},
		{128, 1, 1, []byte{7, 9, 55}, "false:9,128,1:7,9,128", "false:9,128,2:7,9,128"},
		{16384, 1, 1, []byte{7, 9, 55, 66}, "false:9,128,128,1:7,9,128,128", "false:9,128,128,2:7,9,128,128"},
		{^uint64(0), 1, 1, []byte{7, 9, 55}, "false:9,255,255,255,255,255,255,255,255,255,1:7,9,255", "true:9,1:7,9,1"},
		{1 << 63, 0, 0, []byte{22, 33}, "false:128,128,128,128,128,128,128,128,128,1:128,128", "false:255,255,255,255,255,255,255,255,255,1:255,255"},
	}
	var rows bytes.Buffer
	reserveFailures, signedReserveFailures, signedFailures := 0, 0, 0
	for i, c := range cases {
		unsigned := apply(c, func(dst []byte) []byte { return binary.AppendUvarint(dst, c.word) })
		signed := apply(c, func(dst []byte) []byte { return binary.AppendVarint(dst, int64(c.word)) })
		if unsigned != c.unsigned || signed != c.signed {
			log.Fatalf("native %d: unsigned=%s signed=%s", i, unsigned, signed)
		}
		payload := binary.AppendUvarint(nil, c.word)
		if apply(c, func(dst []byte) []byte { return append(dst, payload...) }) != c.unsigned {
			reserveFailures++
		}
		signedPayload := binary.AppendVarint(nil, int64(c.word))
		if apply(c, func(dst []byte) []byte { return append(dst, signedPayload...) }) != c.signed {
			signedReserveFailures++
		}
		if apply(c, func(dst []byte) []byte { return binary.AppendUvarint(dst, c.word<<1) }) != c.signed {
			signedFailures++
		}
		fmt.Fprint(&rows, c.word, " ", c.offset, " ", c.length)
		for _, b := range c.backing {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
	}
	if reserveFailures != 4 || signedReserveFailures != 4 || signedFailures != 2 {
		log.Fatalf("mutation counts unsigned reserve=%d signed reserve=%d sign fold=%d", reserveFailures, signedReserveFailures, signedFailures)
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-binary-varint-append")
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
	fmt.Printf("Append varint reviewer control: %d exact native/Lean paired fixtures; unsigned bulk-reserve, signed bulk-reserve and unsigned-fold alternatives fail %d/%d/%d.\n", len(cases), reserveFailures, signedReserveFailures, signedFailures)
}

func apply(c fixture, operation func([]byte) []byte) string {
	backing := slices.Clone(c.backing)
	output := operation(backing[c.offset : c.offset+c.length : len(backing)])
	reuse := len(backing) > c.offset && &output[0] == &backing[c.offset]
	return fmt.Sprintf("%t:%s:%s", reuse, wire(output), wire(backing))
}

func wire(raw []byte) string {
	items := make([]string, len(raw))
	for i, b := range raw {
		items[i] = strconv.Itoa(int(b))
	}
	return strings.Join(items, ",")
}
