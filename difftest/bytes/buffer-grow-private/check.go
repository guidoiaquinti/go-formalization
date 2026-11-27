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
	raw := append(bytes.Clone(c.Stored), c.Spare...)
	size, offset, flag := len(c.Stored), c.Offset, c.Flag
	reused := len(raw) > 0
	marker := "slice"
	if c.Tag == 1 {
		marker = "nil"
	}
	if offset > 0 && offset == size {
		size, offset, flag = 0, 0, 0
	}
	status := "bounds"
	switch {
	case c.Count < -size:
	case c.Count <= len(raw)-size:
		status = fmt.Sprintf("ok:%d", size)
		size += c.Count
	default:
		unread := raw[offset:size]
		switch {
		case c.Tag == 1 && c.Count <= 64:
			raw, marker, reused = make([]byte, 64), "slice", false
		case c.Count <= len(raw)/2-len(unread):
			copy(raw, unread)
		case len(raw) > math.MaxInt-len(raw)-c.Count:
			status = "too-large"
		default:
			requested := max(size+c.Count, 2*(len(raw)-offset))
			classes := [...]int{8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512}
			found := false
			for _, capacity := range classes {
				if capacity >= requested {
					raw, marker, reused = make([]byte, capacity), "slice", false
					copy(raw, unread)
					found = true
					break
				}
			}
			if !found {
				panic("growth request exceeds tested size classes")
			}
		}
		if status != "too-large" {
			status = fmt.Sprintf("ok:%d", len(unread))
			size, offset = len(unread)+c.Count, 0
		}
	}
	return fmt.Sprintf("%s:%d:%d:%d:%s:%t:%s:%s", status, flag, offset, len(raw), marker, reused, wire(raw[:size]), wire(raw[size:]))
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0, c.Count, math.MaxInt)
		} else {
			fmt.Fprintf(&modelInput, "%d %d %d %d %d %d", c.Tag, c.Count, math.MaxInt, len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-grow-private", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.grow %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-grow-private-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_GROW_PRIVATE_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{
		{Tag: 0, Count: -1}, {Tag: 0, Count: math.MinInt}, {Tag: 0, Count: 0}, {Tag: 0, Count: 1},
		{Tag: 1, Count: -1}, {Tag: 1, Count: 0},
		{Tag: 2, Stored: []byte{1, 2}, Spare: []byte{3}, Offset: 2, Flag: 1, Count: 0},
		{Tag: 2, Stored: []byte{1, 2, 3, 4, 5, 6}, Offset: 4, Flag: 2, Count: 1},
		{Tag: 2, Stored: []byte{1, 2}, Offset: 2, Flag: 1, Count: math.MaxInt},
	}
	for n := 1; n <= 256; n++ {
		edges = append(edges, input{Tag: 1, Count: n})
	}
	for length := range 65 {
		stored := make([]byte, length)
		for i := range stored {
			stored[i] = byte(i*91 + length)
		}
		for off := 0; off <= length; off++ {
			for spare := range 5 {
				retained := bytes.Repeat([]byte{231}, spare)
				for _, count := range []int{-length - 1, -length, -1, (off + spare) % 9} {
					edges = append(edges, input{Tag: 2, Offset: off, Flag: 3, Stored: stored, Spare: retained, Count: count})
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
			count := rng.Intn(385) - 129
			if i%7 == 0 && len(stored)+len(spare) > 0 {
				count = math.MaxInt - 2*(len(stored)+len(spare)) + 1
			}
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferGrowPrivate", "difftest/bytes/buffer-grow-private/testdata/bridge_test.go", "difftest/bytes/buffer-grow-private/testdata/grow_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.grow pinned Go overlay: all %d inputs agree.\n", len(all))
}
