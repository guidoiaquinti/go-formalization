package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	x, y uint16
}

func expected(c input) uint64 {
	d := int64(c.x) - int64(c.y)
	return uint64(d*d) / 4
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cases []input
	edges := []uint16{0, 1, 2, 3, 255, 256, 32767, 32768, 65534, 65535}
	for _, x := range edges {
		for _, y := range edges {
			cases = append(cases, input{x, y})
		}
	}
	for d := range 65536 {
		cases = append(cases, input{0, uint16(d)}, input{uint16(d), 0})
	}
	seeds := []int64{17, 2718, 65537}
	for _, seed := range seeds {
		r := rand.New(rand.NewSource(seed))
		for range 20000 {
			cases = append(cases, input{uint16(r.Uint32()), uint16(r.Uint32())})
		}
	}

	fixture, err := os.CreateTemp("", "golean-sq-diff-")
	if err != nil {
		return err
	}
	defer os.Remove(fixture.Name())
	defer fixture.Close()
	var leanInput bytes.Buffer
	fixtureRows := bufio.NewWriter(fixture)
	for _, c := range cases {
		fmt.Fprintln(&leanInput, c.x, c.y)
		fmt.Fprintln(fixtureRows, c.x, c.y, expected(c))
	}
	if err := fixtureRows.Flush(); err != nil {
		return err
	}
	if err := fixture.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_SQ_DIFF_CASES", fixture.Name()); err != nil {
		return err
	}
	if err := check.TestInPackage("image/color", "TestGoLeanSqDiff", "difftest/image/color/sq-diff/testdata/sq_diff_test.go"); err != nil {
		return err
	}
	for i, row := range check.LeanRows("godiff-image-color-sq-diff", len(cases), &leanInput) {
		if want := strconv.FormatUint(expected(cases[i]), 10); row != want {
			return fmt.Errorf("sqDiff(%d, %d): Lean=%q, want %q", cases[i].x, cases[i].y, row, want)
		}
	}
	fmt.Printf("image/color.sqDiff: %d boundary pairs, 131072 difference sweeps, and 60000 sampled pairs agree; seeds %v\n", len(edges)*len(edges), seeds)
	return nil
}
