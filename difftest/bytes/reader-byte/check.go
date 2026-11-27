package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Offset, Previous int
	Cursor                int64
	Raw, Spare            []byte
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

func expected(c input) string {
	if c.Tag == 0 {
		return "panic|panic"
	}
	marker := "slice"
	if c.Tag == 1 {
		marker = "nil"
	}
	describe := func(cursor int64, previous int) string {
		return fmt.Sprintf("%s:%d:%d:%s:%s", marker, cursor, previous, wire(c.Raw), wire(c.Spare))
	}
	read := "0:EOF:" + describe(c.Cursor, -1)
	if c.Cursor < int64(len(c.Raw)) {
		read = fmt.Sprintf("%d:ok:%s", c.Raw[c.Cursor], describe(c.Cursor+1, -1))
	}
	unread := "bytes.Reader.UnreadByte: at beginning of slice:" + describe(0, c.Previous)
	if c.Cursor > 0 {
		unread = "ok:" + describe(c.Cursor-1, -1)
	}
	return read + "|" + unread
}

func batch(label string, cases []input, all *[]fixture) {
	var protocol bytes.Buffer
	rows := make([]fixture, len(cases))
	for i, c := range cases {
		fmt.Fprintf(&protocol, "%d %d %d %d", c.Tag, len(c.Raw), c.Cursor, c.Previous)
		for _, section := range [][]byte{c.Raw, c.Spare} {
			for _, value := range section {
				fmt.Fprintf(&protocol, " %d", value)
			}
		}
		fmt.Fprintln(&protocol)
		rows[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-reader-byte", len(cases), &protocol) {
		if got != rows[i].Want {
			log.Fatalf("%s row %d: got %q want %q", label, i, got, rows[i].Want)
		}
	}
	*all = append(*all, rows...)
	fmt.Printf("Reader byte %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	edges := []input{{}}
	for _, cursor := range []int64{0, 1, math.MaxInt64} {
		for _, previous := range []int{-1, 0, 7} {
			edges = append(edges, input{Tag: 1, Cursor: cursor, Previous: previous})
			edges = append(edges, input{Tag: 2, Cursor: cursor, Previous: previous})
		}
	}
	for value := range 256 {
		for _, cursor := range []int64{0, 1, 2, 3, math.MaxInt64} {
			for _, previous := range []int{-1, 0, 7} {
				edges = append(edges, input{2, 3, previous, cursor, []byte{byte(value), byte(255 - value)}, []byte{0x7a, 0xe9}})
			}
		}
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			switch {
			case i%23 == 0:
				continue
			case i%19 == 0:
				cases[i] = input{Tag: 1, Cursor: rng.Int63(), Previous: rng.Intn(16) - 1}
			default:
				raw, spare := make([]byte, rng.Intn(65)), make([]byte, rng.Intn(17))
				rng.Read(raw)
				rng.Read(spare)
				cursor := int64(rng.Intn(len(raw) + 4))
				if i%7 == 0 {
					cursor = rng.Int63()
				}
				cases[i] = input{2, rng.Intn(17), rng.Intn(16) - 1, cursor, raw, spare}
			}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	file, err := os.CreateTemp("", "golean-reader-byte-")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(all))
	must(file.Close())
	must(os.Setenv("GOLEAN_READER_BYTE_CASES", file.Name()))
	const data = "difftest/bytes/reader-byte/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReaderByte", data+"bridge_test.go", data+"reader_byte_test.go"))
	fmt.Printf("Reader byte pinned Go overlay: all %d inputs agree.\n", len(all))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
