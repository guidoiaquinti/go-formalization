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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	NilReceiver, NilStorage bool
	Length, Spare, Offset   int
}

type fixture struct {
	Input input
	Want  string
}

func batch(label string, cases []input, all *[]fixture) {
	var model bytes.Buffer
	rows := make([]fixture, len(cases))
	for i, c := range cases {
		want := "panic"
		if c.NilReceiver {
			fmt.Fprintln(&model, 0)
		} else {
			fmt.Fprintf(&model, "1 %d %d\n", c.Length, c.Offset)
			want = strconv.FormatBool(c.Offset > c.Length || c.Offset == c.Length)
		}
		rows[i] = fixture{c, want}
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-empty", len(cases), &model) {
		if got != rows[i].Want {
			log.Fatalf("%s %+v: got %s want %s", label, cases[i], got, rows[i].Want)
		}
	}
	*all = append(*all, rows...)
	fmt.Printf("Buffer.empty %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	edges := []input{{NilReceiver: true}}
	for _, offset := range []int{math.MinInt, -1, 0, 1, math.MaxInt} {
		edges = append(edges, input{NilStorage: true, Offset: offset})
	}
	for length := range 65 {
		for _, offset := range []int{math.MinInt, -1, 0, length - 1, length, length + 1, math.MaxInt} {
			for spare := range 9 {
				edges = append(edges, input{Length: length, Spare: spare, Offset: offset})
			}
		}
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			length := rng.Intn(513)
			offset := rng.Intn(1025) - 256
			if i%4 == 0 {
				offset = length + i%3 - 1
			}
			cases[i] = input{Length: length, Spare: rng.Intn(128), Offset: offset}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	f, err := os.CreateTemp("", "golean-buffer-empty-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(all))
	must(f.Close())
	must(os.Setenv("GOLEAN_BUFFER_EMPTY_CASES", f.Name()))
	const data = "difftest/bytes/buffer-empty/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferEmpty", data+"bridge_test.go", data+"buffer_empty_test.go"))
	fmt.Printf("Buffer.empty pinned Go overlay: all %d inputs agree.\n", len(all))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
