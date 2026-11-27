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
	Raw             []byte
	Capacity, Count int
	NilInput        bool
}
type fixture struct {
	Input input
	Want  string
}

func expected(c input) string {
	requested := 2 * c.Capacity
	if c.Count <= math.MaxInt-len(c.Raw) {
		requested = max(requested, len(c.Raw)+c.Count)
	}
	if requested == 0 {
		return "nil"
	}
	classes := [...]int{8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512}
	for _, capacity := range classes {
		if capacity < requested {
			continue
		}
		payload := make([]string, len(c.Raw))
		for i, b := range c.Raw {
			payload[i] = strconv.Itoa(int(b))
		}
		spare := strings.TrimSuffix(strings.Repeat("0,", capacity-len(c.Raw)), ",")
		return fmt.Sprintf("%d:%s:%s", capacity, strings.Join(payload, ","), spare)
	}
	panic("raw requested capacity exceeds tested classes")
}

func batch(label string, cases []input, all *[]fixture) {
	var rows bytes.Buffer
	fixtures := make([]fixture, len(cases))
	for i, c := range cases {
		fmt.Fprintf(&rows, "%d %d %d", c.Capacity, c.Count, math.MaxInt)
		for _, b := range c.Raw {
			fmt.Fprintf(&rows, " %d", b)
		}
		fmt.Fprintln(&rows)
		fixtures[i] = fixture{c, expected(c)}
	}
	for i, got := range check.LeanRows("godiff-bytes-grow-slice", len(cases), &rows) {
		if got != fixtures[i].Want {
			log.Fatalf("%s %+v: got %q want %q", label, cases[i], got, fixtures[i].Want)
		}
	}
	*all = append(*all, fixtures...)
	fmt.Printf("growSlice %s: %d Lean inputs agree.\n", label, len(cases))
}

func main() {
	edges := []input{{NilInput: true}, {Count: -1}, {Count: 1}, {Count: 1, NilInput: true}, {Count: math.MinInt, NilInput: true}}
	for capacity := range 74 {
		for length := 0; length <= capacity; length++ {
			raw := make([]byte, length)
			for i := range raw {
				raw[i] = byte(i*17 + length)
			}
			for _, count := range []int{math.MinInt, -length, -1, 0, 1, 2*capacity - length, 2*capacity - length + 1, 255} {
				edges = append(edges, input{Raw: raw, Capacity: capacity, Count: count})
			}
			if length > 0 {
				edges = append(edges, input{Raw: raw, Capacity: capacity, Count: math.MaxInt})
			}
		}
	}
	var all []fixture
	batch("boundaries", edges, &all)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			capacity := rng.Intn(257)
			raw := make([]byte, rng.Intn(capacity+1))
			rng.Read(raw)
			cases[i] = input{Raw: raw, Capacity: capacity, Count: rng.Intn(513) - 256}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, &all)
	}
	f, err := os.CreateTemp("", "golean-grow-slice-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(all))
	must(f.Close())
	must(os.Setenv("GOLEAN_GROW_SLICE_CASES", f.Name()))
	const data = "difftest/bytes/grow-slice/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanGrowSlice", data+"bridge_test.go", data+"grow_slice_test.go"))
	fmt.Printf("growSlice pinned Go overlay: all %d inputs agree.\n", len(all))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
