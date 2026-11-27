package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

func expected(c input) string {
	if c.Tag == 0 {
		return "nil-receiver"
	}
	snapshot := fmt.Sprintf("%d:%d:%s:%s", c.Offset, c.Flag, wire(c.Stored), wire(c.Spare))
	if c.Count < 0 {
		return "bounds:" + snapshot
	}
	unread := c.Stored[c.Offset:]
	size := min(len(unread), c.Count)
	status := "ok"
	if size < c.Count {
		status = "eof"
	}
	if c.Tag == 1 {
		return status + ":nil:" + snapshot
	}
	spare := append(bytes.Clone(unread[size:]), c.Spare...)
	return fmt.Sprintf("%s:slice:%d:%s:%s:%s", status, len(unread)+len(c.Spare), wire(unread[:size]), wire(spare), snapshot)
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0, c.Count)
		} else {
			fmt.Fprintf(&modelInput, "%d %d %d %d %d", c.Tag, c.Count, len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-peek", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.Peek %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-peek-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_PEEK_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{{Tag: 0, Count: math.MinInt}, {Tag: 0}, {Tag: 0, Count: math.MaxInt}}
	for _, count := range []int{math.MinInt, -1, 0, 1, math.MaxInt} {
		edges = append(edges, input{Tag: 1, Count: count})
	}
	for length := range 33 {
		stored := make([]byte, length)
		for i := range stored {
			stored[i] = byte(i*29 + length)
		}
		for offset := 0; offset <= length; offset++ {
			for spare := range 5 {
				for _, count := range []int{math.MinInt, -1, 0, 1, length - offset, length - offset + 1, math.MaxInt} {
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
			count := rng.Intn(513) - 256
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferPeek", "difftest/bytes/buffer-peek/testdata/bridge_test.go", "difftest/bytes/buffer-peek/testdata/peek_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.Peek pinned Go overlay: all %d inputs agree.\n", len(all))
}
