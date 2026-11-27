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
	status, marker := "ok", "slice"
	if c.Tag == 1 {
		marker = "nil"
	}
	if c.Flag == 0 {
		status = "sentinel-error"
	} else {
		if c.Offset > 0 {
			c.Offset--
		}
		c.Flag = 0
	}
	unread := c.Stored[c.Offset:]
	return fmt.Sprintf("%s:%d:%d:%d:%s:%s:%s:%s", status, c.Flag, c.Offset, len(unread)+len(c.Spare), marker, wire(unread), wire(c.Spare), wire(c.Stored))
}

func batch(label string, cases []input, all *[]fixture) {
	var modelInput bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		if c.Tag == 0 {
			fmt.Fprintln(&modelInput, 0)
		} else {
			fmt.Fprintf(&modelInput, "%d %d %d %d", c.Tag, len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&modelInput, " %d", b)
				}
			}
			fmt.Fprintln(&modelInput)
		}
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-unread-byte", len(cases), &modelInput) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: Lean=%q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("Buffer.UnreadByte %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	dir, err := os.MkdirTemp("", "golean-buffer-unread-byte-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BUFFER_UNREAD_BYTE_CASES", path); err != nil {
		log.Fatal(err)
	}
	edges := []input{{Tag: 0}}
	stored := make([]byte, 128)
	for i := range stored {
		stored[i] = byte(i * 17)
	}
	for flag := -128; flag < 128; flag++ {
		edges = append(edges, input{Tag: 1, Flag: flag})
		for off := 0; off <= len(stored); off++ {
			edges = append(edges, input{Tag: 2, Offset: off, Flag: flag, Stored: stored, Spare: []byte{255, 8}})
		}
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			stored, spare := make([]byte, rng.Intn(150)), make([]byte, rng.Intn(9))
			rng.Read(stored)
			rng.Read(spare)
			cases[i] = input{Tag: 2, Offset: rng.Intn(len(stored) + 1), Flag: rng.Intn(256) - 128, Stored: stored, Spare: spare}
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
	if err := check.TestInPackage("bytes", "TestGoLeanBufferUnreadByte", "difftest/bytes/buffer-unread-byte/testdata/bridge_test.go", "difftest/bytes/buffer-unread-byte/testdata/unread_byte_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Buffer.UnreadByte pinned Go overlay: all %d inputs agree.\n", len(all))
}
