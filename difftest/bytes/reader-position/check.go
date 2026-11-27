package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Previous, Whence, Script, Count int
	Cursor, Offset                       int64
	Raw, Spare                           []byte
}

type fixture struct {
	Input input
	Want  string
}

func expectedSeek(c input) string {
	if c.Tag == 0 {
		return "nilReader"
	}
	var base int64
	switch c.Whence {
	case 0:
	case 1:
		base = c.Cursor
	case 2:
		base = int64(len(c.Raw))
	default:
		return fmt.Sprintf("whence:0:%d:-1", c.Cursor)
	}
	absolute := new(big.Int).Add(big.NewInt(base), big.NewInt(c.Offset))
	if absolute.Sign() < 0 || !absolute.IsInt64() {
		return fmt.Sprintf("negative:0:%d:-1", c.Cursor)
	}
	position := absolute.Int64()
	return fmt.Sprintf("ok:%d:%d:-1", position, position)
}

func expectedWrite(c input) string {
	if c.Tag == 0 {
		return "nilReader"
	}
	wire := func(status string, count, cursor int64, trace string) string {
		return fmt.Sprintf("%s:%d:%d:-1:%s", status, count, cursor, trace)
	}
	if c.Cursor >= int64(len(c.Raw)) {
		return wire("ok", 0, c.Cursor, "none")
	}
	if c.Script == 3 {
		return wire("nilWriter", 0, c.Cursor, "none")
	}
	unread := c.Raw[c.Cursor:]
	fields := make([]string, len(unread))
	for i, value := range unread {
		fields[i] = strconv.Itoa(int(value))
	}
	trace := fmt.Sprintf("%d:%s", len(unread)+len(c.Spare), strings.Join(fields, ","))
	switch {
	case c.Script == 2:
		return wire("writerPanic", 0, c.Cursor, trace)
	case c.Count > len(unread):
		return wire("invalidCount", 0, c.Cursor, trace)
	}
	status := "ok"
	if c.Script == 1 {
		status = "writerError"
	} else if c.Count < len(unread) {
		status = "shortWrite"
	}
	return wire(status, int64(c.Count), c.Cursor+int64(c.Count), trace)
}

func batch(label string, cases []input, all *[]fixture) {
	var protocol bytes.Buffer
	rows := make([]fixture, len(cases))
	for i, c := range cases {
		fmt.Fprintf(&protocol, "%d %d %d %d %d %d %d %d", c.Tag, len(c.Raw), c.Cursor, c.Previous,
			c.Offset, c.Whence, c.Script, c.Count)
		for _, section := range [][]byte{c.Raw, c.Spare} {
			for _, b := range section {
				fmt.Fprintf(&protocol, " %d", b)
			}
		}
		fmt.Fprintln(&protocol)
		rows[i] = fixture{c, expectedSeek(c) + "|" + expectedWrite(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-reader-position", len(cases), &protocol) {
		if got != rows[i].Want {
			log.Fatalf("%s row %d (%+v): got %q want %q", label, i, cases[i], got, rows[i].Want)
		}
	}
	*all = append(*all, rows...)
	fmt.Printf("Reader position %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	edges := []input{{}, {Whence: -1, Script: 3}, {Tag: 1, Previous: 5}}
	for _, raw := range [][]byte{{}, {255}, {0x61, 0xe2, 0x82, 0xac, 0xff}} {
		for _, cursor := range []int64{0, 1, int64(len(raw)), math.MaxInt64} {
			for _, offset := range []int64{math.MinInt64, -1, 0, 1, math.MaxInt64} {
				for _, whence := range []int{-1, 0, 1, 2, 3} {
					for _, count := range []int{math.MinInt, -1, 0, 1, len(raw), len(raw) + 1, math.MaxInt} {
						for script := range 4 {
							edges = append(edges, input{2, 9, whence, script, count, cursor, offset, raw, []byte{0x82, 0xac}})
						}
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
			if i%23 == 0 {
				continue
			}
			raw, spare := make([]byte, rng.Intn(65)), make([]byte, rng.Intn(10))
			rng.Read(raw)
			rng.Read(spare)
			tag := 2
			if i%19 == 0 {
				tag, raw, spare = 1, nil, nil
			}
			cursor, offset := int64(rng.Intn(len(raw)+3)), int64(rng.Intn(131)-65)
			if i%7 == 0 {
				cursor = rng.Int63()
			}
			if i%3 == 0 {
				offset = int64(rng.Uint64())
			}
			count := rng.Intn(len(raw)+5) - 2
			cases[i] = input{tag, rng.Intn(23) - 3, rng.Intn(5) - 1, rng.Intn(4), count, cursor, offset, raw, spare}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	file, err := os.CreateTemp("", "golean-reader-position-")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(all))
	must(file.Close())
	must(os.Setenv("GOLEAN_READER_POSITION_CASES", file.Name()))
	const data = "difftest/bytes/reader-position/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReaderPosition", data+"bridge_test.go", data+"reader_position_test.go"))
	fmt.Printf("Reader position pinned Go overlay: all %d inputs agree.\n", len(all))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
