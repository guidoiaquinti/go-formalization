package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Offset, Flag int
	Stored, Spare     []byte
	Count             int
}

type fixture struct {
	Input input
	Want  string
}

func wire(raw []byte) string {
	values := make([]string, len(raw))
	for i, b := range raw {
		values[i] = strconv.Itoa(int(b))
	}
	return strings.Join(values, ",")
}

func destination(count int) []byte {
	if count == 0 {
		return nil
	}
	dst := make([]byte, count)
	for i := range dst {
		dst[i] = byte(167 + i*31)
	}
	return dst
}

func expected(c input) string {
	if c.Tag == 0 {
		return "panic"
	}
	dst := destination(c.Count)
	original := append(bytes.Clone(c.Stored), c.Spare...)
	stored, spare, offset := c.Stored, c.Spare, c.Offset
	unread := c.Stored[c.Offset:]
	status, marker, flag := "ok", "slice", 0
	if c.Tag == 1 {
		marker = "nil"
	}
	n := copy(dst, unread)
	if len(unread) == 0 {
		stored, spare, offset = nil, original, 0
		if len(dst) > 0 {
			status = "eof"
		}
	} else if n > 0 {
		flag = -1
		offset += n
	}
	return fmt.Sprintf("%s:%d:%s:%d:%d:%s:%s:%s", status, n, marker, offset, flag, wire(dst), wire(stored), wire(spare))
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0)
		} else {
			fmt.Fprintf(&modelInput, "%d %d %d %d %d", c.Tag, len(c.Stored), len(c.Spare), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare, destination(c.Count)} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-read", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.Read %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-read-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_READ_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{{Tag: 0}, {Tag: 0, Count: 1}, {Tag: 1}, {Tag: 1, Count: 1}}
	for length := range 33 {
		stored := make([]byte, length)
		for i := range stored {
			stored[i] = byte(i*29 + length)
		}
		for offset := 0; offset <= length; offset++ {
			for spare := range 5 {
				for _, count := range []int{0, 1, max(0, length-offset-1), length - offset, length - offset + 1, 64} {
					edges = append(edges, input{Tag: 2, Stored: stored, Spare: bytes.Repeat([]byte{255}, spare), Offset: offset, Flag: offset%6 - 1, Count: count})
				}
			}
		}
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			stored, spare := make([]byte, rng.Intn(251)), make([]byte, rng.Intn(6))
			rng.Read(stored)
			rng.Read(spare)
			count := rng.Intn(256)
			cases[i] = input{Tag: 2, Offset: rng.Intn(len(stored) + 1), Flag: rng.Intn(256) - 128, Stored: stored, Spare: spare, Count: count}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(all); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	if err := check.TestInPackage("bytes", "TestGoLeanBufferRead", "difftest/bytes/buffer-read/testdata/bridge_test.go", "difftest/bytes/buffer-read/testdata/read_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.Read pinned Go overlay: all %d inputs agree.\n", len(all))
}
