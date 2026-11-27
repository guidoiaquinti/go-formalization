package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Previous        int
	Cursor, Offset       int64
	Sharing, Size        int
	Raw, Spare, Separate []byte
}

type fixture struct {
	Input input
	Want  string
}

func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, value := range raw {
		fields[i] = strconv.Itoa(int(value))
	}
	return strings.Join(fields, ",")
}

func observe(c input, positioned bool) string {
	source := c.Cursor
	if positioned {
		source = c.Offset
	}
	storage := append(slices.Clone(c.Raw), c.Spare...)
	destination := slices.Clone(c.Separate)
	if c.Sharing >= 0 {
		destination = storage[c.Sharing : c.Sharing+c.Size]
	}
	status, n := "ok", 0
	switch {
	case positioned && source < 0:
		status = "bytes.Reader.ReadAt: negative offset"
	case c.Tag == 0:
		return "panic"
	case source >= int64(len(c.Raw)):
		status = "EOF"
	default:
		n = min(len(destination), len(c.Raw)-int(source))
		for i := range n {
			destination[i] = c.Raw[int(source)+i]
		}
		if positioned && n < len(destination) {
			status = "EOF"
		}
		if !positioned {
			c.Cursor += int64(n)
			c.Previous = -1
		}
	}
	state := "receiver-nil"
	if c.Tag != 0 {
		tag := "slice"
		if c.Tag == 1 {
			tag = "nil"
		}
		state = fmt.Sprintf("%s:%d:%d:%s:%s", tag, c.Cursor, c.Previous,
			wire(storage[:len(c.Raw)]), wire(storage[len(c.Raw):]))
	}
	return fmt.Sprintf("%d:%s:%s:%s", n, status, wire(destination), state)
}

func batch(label string, cases []input, all *[]fixture) {
	var protocol bytes.Buffer
	rows := make([]fixture, len(cases))
	for i, c := range cases {
		fmt.Fprintf(&protocol, "%d %d %d %d %d %d %d %d", c.Tag, len(c.Raw), len(c.Spare),
			c.Cursor, c.Previous, c.Offset, c.Sharing, c.Size)
		for _, section := range [][]byte{c.Raw, c.Spare, c.Separate} {
			for _, b := range section {
				fmt.Fprintf(&protocol, " %d", b)
			}
		}
		fmt.Fprintln(&protocol)
		rows[i] = fixture{c, observe(c, false) + "|" + observe(c, true)}
	}
	for i, got := range check.LeanRows("godiff-bytes-reader-read", len(cases), &protocol) {
		if got != rows[i].Want {
			log.Fatalf("%s row %d (%+v): got %q want %q", label, i, cases[i], got, rows[i].Want)
		}
	}
	*all = append(*all, rows...)
	fmt.Printf("Reader copy %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	var edges []input
	for _, offset := range []int64{math.MinInt64, -1, 0, 1, math.MaxInt64} {
		for _, tag := range []int{0, 1, 2} {
			for _, destination := range [][]byte{nil, {}, {0xa1, 0xb2}} {
				edges = append(edges, input{Tag: tag, Previous: 7, Offset: offset, Sharing: -1, Separate: destination})
			}
		}
	}
	for length := range 8 {
		raw, spare := make([]byte, length), []byte{0xa1, 0xb2}
		for i := range raw {
			raw[i] = byte(0xf0 - 29*i)
		}
		for _, cursor := range []int64{0, 1, int64(length), int64(length + 1), math.MaxInt64} {
			for _, offset := range []int64{math.MinInt64, -1, 0, 1, int64(length) - 1, int64(length), math.MaxInt64} {
				for start := range length + len(spare) + 1 {
					for size := range length + len(spare) - start + 1 {
						edges = append(edges, input{2, 7, cursor, offset, start, size, raw, spare, nil})
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
			raw, spare := make([]byte, rng.Intn(48)), make([]byte, rng.Intn(12))
			rng.Read(raw)
			rng.Read(spare)
			c := input{Tag: 2, Previous: rng.Intn(18) - 2, Cursor: int64(rng.Intn(len(raw) + 4)),
				Offset: int64(rng.Intn(len(raw)+7) - 3), Raw: raw, Spare: spare, Sharing: -1}
			switch {
			case i%23 == 0:
				c.Tag, c.Raw, c.Spare = 0, nil, nil
			case i%19 == 0:
				c.Tag, c.Raw, c.Spare = 1, nil, nil
			case i%7 == 0:
				c.Cursor = rng.Int63()
			}
			if c.Tag == 2 && i%2 == 0 {
				capacity := len(raw) + len(spare)
				c.Sharing = rng.Intn(capacity + 1)
				c.Size = rng.Intn(capacity - c.Sharing + 1)
			} else {
				c.Separate = make([]byte, rng.Intn(54))
				rng.Read(c.Separate)
			}
			cases[i] = c
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	file, err := os.CreateTemp("", "golean-reader-copy-")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(all))
	must(file.Close())
	must(os.Setenv("GOLEAN_READER_COPY_CASES", file.Name()))
	const data = "difftest/bytes/reader-read/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReaderCopy", data+"bridge_test.go", data+"reader_read_test.go"))
	fmt.Printf("Reader copy pinned Go overlay: all %d inputs agree.\n", len(all))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
