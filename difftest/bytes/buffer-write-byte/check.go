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
	Byte              byte
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
	raw := append(bytes.Clone(c.Stored), c.Spare...)
	size, offset, reused := len(c.Stored), c.Offset, true
	switch {
	case len(c.Spare) != 0:
	case c.Tag == 1:
		raw, size, offset, reused = make([]byte, 64), 0, 0, false
	case size == c.Offset && c.Offset != 0:
		size, offset = 0, 0
	default:
		unread := c.Stored[c.Offset:]
		size, offset = len(unread), 0
		if size+1 <= len(raw)/2 {
			copy(raw, unread)
		} else {
			requested := max(len(raw)+1, 2*len(unread))
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
	}
	raw[size] = c.Byte
	return fmt.Sprintf("ok:0:%d:%d:slice:%t:%s:%s", offset, len(raw), reused, wire(raw[:size+1]), wire(raw[size+1:]))
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0, c.Byte)
		} else {
			fmt.Fprintf(&modelInput, "%d %d %d %d %d", c.Tag, c.Byte, len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-write-byte", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.WriteByte %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-write-byte-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_WRITE_BYTE_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{{Tag: 0}}
	for value := range 256 {
		edges = append(edges, input{Tag: 1, Flag: value - 128, Byte: byte(value)})
	}
	for length := range 65 {
		stored := make([]byte, length)
		for i := range stored {
			stored[i] = byte(i*91 + length)
		}
		for off := 0; off <= length; off++ {
			for spare := range 8 {
				edges = append(edges, input{Tag: 2, Offset: off, Flag: (off+length)%8 - 1, Stored: stored, Spare: bytes.Repeat([]byte{231}, spare), Byte: byte(off + len(stored) + spare)})
			}
		}
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			stored, spare := make([]byte, rng.Intn(257)), make([]byte, rng.Intn(6))
			rng.Read(stored)
			rng.Read(spare)
			cases[i] = input{Tag: 2, Offset: rng.Intn(len(stored) + 1), Flag: rng.Intn(256) - 128, Stored: stored, Spare: spare, Byte: byte(rng.Intn(256))}
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferWriteByte", "difftest/bytes/buffer-write-byte/testdata/bridge_test.go", "difftest/bytes/buffer-write-byte/testdata/write_byte_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.WriteByte pinned Go overlay: all %d inputs agree.\n", len(all))
}
