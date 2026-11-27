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
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
func batch(label string, cases []input, all *[]fixture) {
	var model bytes.Buffer
	rows := make([]fixture, len(cases))
	for i, c := range cases {
		fmt.Fprintf(&model, "%d %d %d %d", c.Tag, len(c.Raw), c.Cursor, c.Previous)
		for _, raw := range [][]byte{c.Raw, c.Spare} {
			for _, v := range raw {
				fmt.Fprintf(&model, " %d", v)
			}
		}
		fmt.Fprintln(&model)
		marker := "slice"
		if c.Tag == 0 {
			marker = "nil"
		}
		initial := fmt.Sprintf("%s:0:-1:%s:%s", marker, wire(c.Raw), wire(c.Spare))
		rows[i] = fixture{c, fmt.Sprintf("%s|%s|panic|%d|panic", initial, initial, len(c.Raw))}
	}
	for i, got := range check.LeanRows("godiff-bytes-reader-setup", len(cases), &model) {
		if got != rows[i].Want {
			log.Fatalf("%s row %d: got %q want %q", label, i, got, rows[i].Want)
		}
	}
	*all = append(*all, rows...)
	fmt.Printf("Reader setup %s: %d Lean inputs agree.\n", label, len(cases))
}
func main() {
	edges := []input{{Tag: 0}, {Tag: 0, Cursor: math.MaxInt64, Previous: 42}}
	for length := range 33 {
		b := make([]byte, length)
		for i := range b {
			b[i] = byte(i*149 + length)
		}
		for offset := range 8 {
			for _, cursor := range []int64{0, int64(length), int64(length) + 1, math.MaxInt64} {
				for _, previous := range []int{-1, 0, 7} {
					edges = append(edges, input{1, offset, previous, cursor, b, bytes.Repeat([]byte{255}, offset)})
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
			if i%19 == 0 {
				continue
			}
			raw, spare := make([]byte, rng.Intn(129)), make([]byte, rng.Intn(33))
			rng.Read(raw)
			rng.Read(spare)
			cases[i] = input{1, rng.Intn(32), rng.Intn(16) - 1, rng.Int63(), raw, spare}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	f, err := os.CreateTemp("", "golean-reader-setup-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(all))
	must(f.Close())
	must(os.Setenv("GOLEAN_READER_SETUP_CASES", f.Name()))
	const data = "difftest/bytes/reader-setup/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReaderSetup", data+"bridge_test.go", data+"reader_setup_test.go"))
	fmt.Printf("Reader setup pinned Go overlay: all %d inputs agree.\n", len(all))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
