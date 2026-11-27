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
	Rune              rune
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
	encoded := []byte(string(c.Rune))
	reservation := 4
	if uint32(c.Rune) < 128 {
		reservation = 1
	}
	raw := append(bytes.Clone(c.Stored), c.Spare...)
	size, offset, reused := len(c.Stored), c.Offset, true
	if len(c.Spare) < reservation {
		if size == offset && offset != 0 {
			size, offset = 0, 0
		}
		if len(raw)-size < reservation {
			unread := raw[offset:size]
			switch {
			case c.Tag == 1:
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
	copy(raw[size:], encoded)
	length := size + len(encoded)
	return fmt.Sprintf("ok:%d:0:%d:%d:slice:%t:%s:%s", len(encoded), offset, len(raw), reused, wire(raw[:length]), wire(raw[length:]))
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0, c.Rune)
		} else {
			fmt.Fprintf(&modelInput, "%d %d %d %d %d", c.Tag, c.Rune, len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-write-rune", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.WriteRune %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-write-rune-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_WRITE_RUNE_CASES", path); err != nil {
		log.Fatal(err)
	}
	runes := []rune{math.MinInt32, -1, 0, 0x7f, 0x80, 0x7ff, 0x800, 0xd7ff, 0xd800, 0xdfff, 0xe000, 0xfffd, 0xffff, 0x10000, 0x10ffff, 0x110000, math.MaxInt32}
	edges := []input{{Tag: 0, Rune: 'x'}, {Tag: 0, Rune: -1}}
	for _, value := range runes {
		edges = append(edges, input{Tag: 1, Flag: 4, Rune: value})
	}
	for length := range 65 {
		stored := make([]byte, length)
		for i := range stored {
			stored[i] = byte(i*91 + length)
		}
		for off := 0; off <= length; off++ {
			for spare := range 8 {
				value := runes[(off+len(stored)+spare)%len(runes)]
				edges = append(edges, input{Tag: 2, Offset: off, Flag: (off+length)%8 - 1, Stored: stored, Spare: bytes.Repeat([]byte{231}, spare), Rune: value})
			}
		}
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			stored, spare := make([]byte, rng.Intn(254)), make([]byte, rng.Intn(6))
			rng.Read(stored)
			rng.Read(spare)
			value := rune(rng.Uint32())
			switch i % 4 {
			case 0:
				value = rune(rng.Intn(128))
			case 1:
				value = runes[rng.Intn(len(runes))]
			case 2:
				value = rune(rng.Intn(0x110000))
			}
			cases[i] = input{Tag: 2, Offset: rng.Intn(len(stored) + 1), Flag: rng.Intn(256) - 128, Stored: stored, Spare: spare, Rune: value}
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferWriteRune", "difftest/bytes/buffer-write-rune/testdata/bridge_test.go", "difftest/bytes/buffer-write-rune/testdata/write_rune_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.WriteRune pinned Go overlay: all %d inputs agree.\n", len(all))
}
