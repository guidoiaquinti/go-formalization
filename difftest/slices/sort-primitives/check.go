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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Length  int
	Word    uint64
	Present bool
	Want    string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	file, err := os.CreateTemp("", "golean-slices-sort-primitives-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) error {
		var rows bytes.Buffer
		for _, c := range cases {
			present := 0
			if c.Present {
				present = 1
			}
			fmt.Fprintln(&rows, strconv.IntSize, c.Length, present, c.Word)
		}
		for i, got := range check.LeanRows("godiff-slices-sort-primitives", len(cases), &rows) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				return fmt.Errorf("%+v: Lean=%q", c, got)
			}
			if err := encoder.Encode(c); err != nil {
				return err
			}
		}
		fmt.Printf("sort primitives %s: %d Lean rows agree\n", label, len(cases))
		return nil
	}
	words := append(binary64.Boundaries(), 0x7f800000, 0x7f800001, 0xff800000, 0xffffffff, ^uint64(0))
	for bit := range 64 {
		words = append(words, uint64(1)<<bit)
	}
	lengths := []int{math.MinInt, -1, 0, 1, math.MaxInt}
	for bit := 1; bit < strconv.IntSize-1; bit++ {
		power := 1 << bit
		lengths = append(lengths, power-1, power, power+1)
	}
	var boundaries []fixture
	for _, word := range words {
		for _, length := range lengths {
			for _, present := range []bool{false, true} {
				boundaries = append(boundaries, fixture{Length: length, Word: word, Present: present})
			}
		}
	}
	if err := compare("boundaries", boundaries); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			word := rng.Uint64()
			if i%2 == 0 {
				word = word | (uint64(0x7ff0000000000000) | uint64(0x7f800000))
			}
			cases[i] = fixture{Length: int(rng.Uint64()), Word: word, Present: i%7 != 0}
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SLICES_SORT_PRIMITIVES_CASES", file.Name()); err != nil {
		return err
	}
	const data = "difftest/slices/sort-primitives/testdata/"
	return check.TestInPackage("slices", "TestGoLeanSortPrimitives", data+"bridge_test.go", data+"primitives_test.go")
}

func expected(c fixture) string {
	next := "panic"
	if c.Present {
		word := bitParityStage(c.Word, -13)
		word = bitParityStage(word, 7)
		word = bitParityStage(word, -17)
		next = fmt.Sprintf("%d,%d", word, word)
	}
	power := uint(1)
	for value := uint(c.Length); value != 0; value /= 2 {
		power *= 2
	}
	f32, f64 := math.Float32frombits(uint32(c.Word)), math.Float64frombits(c.Word)
	return fmt.Sprintf("%s/%d/%t/%t/false", next, power, math.IsNaN(float64(f32)), math.IsNaN(f64))
}

func bitParityStage(word uint64, offset int) uint64 {
	var result uint64
	for bit := range 64 {
		value := (word >> bit) & 1
		other := bit + offset
		if other >= 0 && other < 64 {
			value = (value + ((word >> other) & 1)) % 2
		}
		result |= value << bit
	}
	return result
}
