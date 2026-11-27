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
	Failed            bool
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
	call, status, marker := "-", "ok", "slice"
	if c.Tag == 1 {
		marker = "nil"
	}
	n, offset := 0, c.Offset
	unread := c.Stored[c.Offset:]
	if len(unread) > 0 {
		call = wire(unread) + "/" + wire(c.Spare)
		switch {
		case c.Count > len(unread):
			status = "invalid"
		case c.Failed:
			n, status = c.Count, "writer"
		case c.Count < len(unread):
			n, status = c.Count, "short"
		default:
			n = c.Count
		}
		offset += n
	}
	if status == "ok" {
		c.Spare = append(bytes.Clone(c.Stored), c.Spare...)
		c.Stored = nil
		offset = 0
	}
	return fmt.Sprintf("%s:%d:%s:%s:%d:0:%s:%s", status, n, call, marker, offset, wire(c.Stored), wire(c.Spare))
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0)
		} else {
			failed := 0
			if c.Failed {
				failed = 1
			}
			fmt.Fprintf(&modelInput, "%d %d %d %d %d %d", c.Tag, c.Count, failed, len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-write-to", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.WriteTo %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-write-to-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_WRITE_TO_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{{Tag: 0}, {Tag: 1}, {Tag: 1, Count: 99, Failed: true}}
	for length := range 33 {
		stored := make([]byte, length)
		for i := range stored {
			stored[i] = byte(length + i*83)
		}
		for offset := 0; offset <= length; offset++ {
			for spare := range 5 {
				for _, count := range []int{0, 1, max(0, length-offset-1), length - offset, length - offset + 1} {
					for _, failed := range []bool{false, true} {
						edges = append(edges, input{Tag: 2, Stored: stored, Spare: bytes.Repeat([]byte{255}, spare), Offset: offset, Flag: offset%6 - 1, Count: count, Failed: failed})
					}
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
			cases[i] = input{Tag: 2, Offset: rng.Intn(len(stored) + 1), Flag: rng.Intn(256) - 128, Stored: stored, Spare: spare, Count: count, Failed: i%2 == 0}
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferWriteTo", "difftest/bytes/buffer-write-to/testdata/bridge_test.go", "difftest/bytes/buffer-write-to/testdata/write_to_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.WriteTo pinned Go overlay: all %d inputs agree.\n", len(all))
}
