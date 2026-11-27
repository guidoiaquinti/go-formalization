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

func expected(c input) string {
	if c.Tag == 0 {
		return "nil-receiver"
	}
	unread := c.Stored[c.Offset:]
	size := len(unread)
	status := "eof"
	for i, b := range unread {
		if b == byte(c.Count) {
			size, status = i+1, "ok"
			break
		}
	}
	snapshot := fmt.Sprintf("%d:-1:%s:%s", c.Offset+size, wire(c.Stored), wire(c.Spare))
	return fmt.Sprintf("%s:%s:%s", status, wire(unread[:size]), snapshot)
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
	for i, got := range check.LeanRows("godiff-bytes-buffer-read-string", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.ReadString %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-read-string-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_READ_STRING_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{{Tag: 0}, {Tag: 1}}
	for delimiter := range 256 {
		raw := []byte{byte(delimiter), byte(delimiter + 1), byte(delimiter), byte(delimiter + 2)}
		for offset := 0; offset <= len(raw); offset++ {
			for spare := range 9 {
				edges = append(edges, input{Tag: 2, Offset: offset, Flag: offset % 5, Count: delimiter, Stored: raw, Spare: bytes.Repeat([]byte{byte(delimiter)}, spare)})
			}
		}
	}
	edges = append(edges,
		input{Tag: 2, Stored: []byte{0, 255, 128}, Spare: []byte{7}, Offset: 1, Flag: 2, Count: 0},
		input{Tag: 2, Stored: []byte{0, 255, 9, 255}, Spare: []byte{7}, Offset: 1, Flag: 2, Count: 255},
		input{Tag: 2, Stored: []byte("é\xff\x00"), Count: 0xc3},
	)
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferReadString", "difftest/bytes/buffer-read-string/testdata/bridge_test.go", "difftest/bytes/buffer-read-string/testdata/read_string_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.ReadString pinned Go overlay: all %d inputs agree.\n", len(all))
}
