package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type scenario struct {
	index, length int
	decisions     [2]uint64
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []scenario
	for n := range 33 {
		for i := range 33 {
			for mask := range 16 {
				cases = append(cases, scenario{i, n, [2]uint64{uint64(mask), 0}})
			}
		}
	}
	for _, pair := range [][2]int{
		{0, math.MaxInt},
		{math.MaxInt/2 - 1, math.MaxInt},
		{math.MaxInt / 2, math.MaxInt},
		{math.MaxInt/2 + 1, math.MaxInt},
		{math.MaxInt, math.MaxInt},
		{math.MaxInt, 0},
	} {
		for _, mask := range [][2]uint64{{0, 0}, {math.MaxUint64, 0}, {math.MaxUint64, math.MaxUint64}} {
			cases = append(cases, scenario{pair[0], pair[1], mask})
		}
	}
	boundaries := len(cases)
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			n, index := r.Int(), r.Int()
			mask := [2]uint64{r.Uint64(), r.Uint64()}
			switch i % 4 {
			case 0:
				index = 0
				mask = [2]uint64{math.MaxUint64, math.MaxUint64}
			case 1:
				index = n / 4
			case 2:
				index = r.Intn(1024)
			}
			cases = append(cases, scenario{index, n, mask})
		}
	}

	var modelInput, fixture bytes.Buffer
	for _, c := range cases {
		row := fmt.Sprintf("%d %d %d %d", c.index, c.length, c.decisions[0], c.decisions[1])
		fmt.Fprintln(&modelInput, row)
		fmt.Fprintln(&fixture, row, c.expected())
	}
	for i, row := range check.LeanRows("godiff-container-heap-down", len(cases), &modelInput) {
		if want := cases[i].expected(); row != want {
			return fmt.Errorf("Lean down trace for %+v was %q, want %q", cases[i], row, want)
		}
	}

	dir, err := os.MkdirTemp("", "golean-heap-down-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_HEAP_DOWN_CASES", file); err != nil {
		return err
	}
	if err := check.TestInPackage("container/heap", "TestGoLeanDown",
		"difftest/container/heap/down/testdata/down_test.go"); err != nil {
		return err
	}
	fmt.Printf("container/heap.down: %d boundaries and %d random scenarios agree, seeds %v.\n",
		boundaries, len(cases)-boundaries, seeds)
	return nil
}

func (c scenario) expected() string {
	var trace bytes.Buffer
	calls := 0
	less := func(a, b int) bool {
		fmt.Fprintf(&trace, " L:%d:%d", a, b)
		bit := c.decisions[calls/64] >> (calls % 64) & 1
		calls++
		return bit != 0
	}
	index := c.index
	// Guard with the number of parents before computing any child index.
	for index < c.length/2 {
		child := 2*index + 1
		if child < c.length-1 && less(child+1, child) {
			child++
		}
		if !less(child, index) {
			break
		}
		fmt.Fprintf(&trace, " S:%d:%d", index, child)
		index = child
	}
	return strconv.FormatBool(index != c.index) + trace.String()
}
