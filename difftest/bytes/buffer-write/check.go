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
	Tag, Offset, Flag    int
	Stored, Spare        []byte
	Data                 []byte
	AliasStart, AliasLen int
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
		return "panic"
	}
	old := append(bytes.Clone(c.Stored), c.Spare...)
	raw := old
	payload := c.Data
	if c.AliasStart >= 0 {
		payload = old[c.AliasStart : c.AliasStart+c.AliasLen]
	}
	reservation := len(payload)
	size, offset, reused := len(c.Stored), c.Offset, true
	if len(c.Spare) < reservation {
		if size == offset && offset != 0 {
			size, offset = 0, 0
		}
		if len(raw)-size < reservation {
			unread := raw[offset:size]
			switch {
			case c.Tag == 1 && reservation <= 64:
				raw, reused = make([]byte, 64), false
			case len(unread)+reservation <= len(raw)/2:
				copy(raw, unread)
			default:
				requested := max(size+reservation, 2*(len(raw)-offset))
				classes := [...]int{8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512}
				for _, capacity := range classes {
					if capacity >= requested {
						raw, reused = make([]byte, capacity), false
						copy(raw, unread)
						break
					}
				}
				if reused {
					panic("growth request exceeds tested size classes")
				}
			}
			size, offset = len(unread), 0
		}
	}
	count := copy(raw[size:], payload)
	length := size + count
	marker := "slice"
	if c.Tag == 1 && count == 0 {
		marker = "nil"
	}
	return fmt.Sprintf("ok:%d:0:%d:%d:%s:%t:%s:%s", count, offset, len(raw), marker, reused && len(raw) > 0, wire(raw[:length]), wire(raw[length:]))
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0)
		} else {
			fmt.Fprintf(&modelInput, "%d %d %d %d %d %d %d", c.Tag, len(c.Stored), len(c.Spare), c.Offset, c.Flag, c.AliasStart, c.AliasLen)
			for _, block := range [][]byte{c.Stored, c.Spare, c.Data} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-write", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.Write %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-write-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_WRITE_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{
		{Tag: 0, AliasStart: -1}, {Tag: 0, Data: []byte{1}, AliasStart: -1},
		{Tag: 1, AliasStart: -1}, {Tag: 1, Data: []byte{}, AliasStart: -1},
		{Tag: 2, Stored: []byte{1, 2}, Flag: 1, AliasStart: 0, AliasLen: 2},
		{Tag: 2, Stored: []byte{1, 2, 3, 4, 5, 6, 7, 8}, Offset: 6, Flag: 1, AliasStart: 0, AliasLen: 2},
		{Tag: 2, Stored: []byte{1, 2}, Offset: 2, Flag: 2, AliasStart: -1},
		{Tag: 2, Stored: []byte{1, 2}, Spare: []byte{3, 4}, Flag: -1, AliasStart: 1, AliasLen: 2},
	}
	for n := 1; n <= 256; n++ {
		edges = append(edges, input{Tag: 1, Data: bytes.Repeat([]byte{255}, n), AliasStart: -1})
	}
	for length := range 65 {
		stored := make([]byte, length)
		for i := range stored {
			stored[i] = byte(i*91 + length)
		}
		for off := 0; off <= length; off++ {
			for spare := range 5 {
				retained := bytes.Repeat([]byte{231}, spare)
				edges = append(edges, input{Tag: 2, Offset: off, Flag: 3, Stored: stored, Spare: retained, Data: []byte{1, 2, 255}, AliasStart: -1})
				capacity := length + spare
				edges = append(edges, input{Tag: 2, Offset: off, Flag: -1, Stored: stored, Spare: retained, AliasStart: 0, AliasLen: min(capacity, 1+(off+spare)%7)})
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
			payload := make([]byte, rng.Intn(256))
			rng.Read(payload)
			c := input{Tag: 2, Offset: rng.Intn(len(stored) + 1), Flag: rng.Intn(256) - 128, Stored: stored, Spare: spare, Data: payload, AliasStart: -1}
			if i%3 == 0 {
				capacity := len(stored) + len(spare)
				c.AliasStart = rng.Intn(capacity + 1)
				c.AliasLen = rng.Intn(capacity - c.AliasStart + 1)
				c.Data = nil
			}
			cases[i] = c
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferWrite", "difftest/bytes/buffer-write/testdata/bridge_test.go", "difftest/bytes/buffer-write/testdata/write_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.Write pinned Go overlay: all %d inputs agree.\n", len(all))
}
