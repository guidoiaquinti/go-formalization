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
	if size == 0 {
		return status + ":nil:" + snapshot
	}
	classes := [...]int{8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512}
	for _, capacity := range classes {
		if capacity >= size {
			return fmt.Sprintf("%s:slice:%d:%s:%s:%s", status, capacity, wire(unread[:size]), wire(make([]byte, capacity-size)), snapshot)
		}
	}
	panic("result larger than tested capacity classes")
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
	for i, got := range check.LeanRows("godiff-bytes-buffer-read-bytes", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.ReadBytes %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-read-bytes-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_READ_BYTES_CASES", path); err != nil {
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
	for size := 1; size <= 512; size++ {
		edges = append(edges, input{Tag: 2, Stored: bytes.Repeat([]byte{0x80}, size), Spare: []byte{0xfe}, Count: 0xff})
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferReadBytes", "difftest/bytes/buffer-read-bytes/testdata/bridge_test.go", "difftest/bytes/buffer-read-bytes/testdata/read_bytes_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.ReadBytes pinned Go overlay: all %d inputs agree.\n", len(all))
}
