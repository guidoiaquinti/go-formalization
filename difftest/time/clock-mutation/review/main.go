package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type raw struct{ wall, ext uint64 }

type fixture struct {
	Wall, Ext, Argument uint64
	Location            int
	Want                string
}

type example struct {
	state    fixture
	set, add raw
}

func main() {
	const flag = uint64(1 << 63)
	const minusOne = ^uint64(0)
	const first = uint64(59453308800)
	const packedSecond = uint64(1 << 30)
	const spare = uint64(1 << 60)
	cases := []example{
		{fixture{17, flag + 1, minusOne, 2, ""}, raw{17, flag + 1}, raw{17, flag}},
		{fixture{17, flag, minusOne, 1, ""}, raw{17, flag}, raw{17, flag + 1}},
		{fixture{spare | 17, first + 1, 7, 2, ""}, raw{flag | spare | packedSecond | 17, 7}, raw{spare | 17, first + 8}},
		{fixture{flag | packedSecond | 17, 123, minusOne, 1, ""}, raw{flag | packedSecond | 17, minusOne}, raw{flag | 17, 123}},
		{fixture{flag | 17, 123, minusOne, 0, ""}, raw{flag | 17, minusOne}, raw{17, first - 1}},
		{fixture{17, first - 1, 2, 2, ""}, raw{17, first - 1}, raw{17, first + 1}},
		{fixture{spare | 17, flag - 1, 1, 1, ""}, raw{spare | 17, flag - 1}, raw{spare | 17, flag - 1}},
	}
	f, err := os.CreateTemp("", "golean-clock-mut-control-")
	must(err)
	defer os.Remove(f.Name())
	enc := gob.NewEncoder(f)
	render := func(a, b raw, loc int) string {
		return fmt.Sprintf("%d,%d,%d|%d,%d,%d", a.wall, a.ext, loc, b.wall, b.ext, loc)
	}
	var input bytes.Buffer
	for _, c := range cases {
		v := c.state
		fmt.Fprintln(&input, v.Wall, v.Ext, v.Location, v.Argument)
		v.Want = render(c.set, c.add, v.Location)
		must(enc.Encode(v))
	}
	must(f.Close())
	var count [5]int
	for i, got := range check.LeanRows("godiff-time-clock-mutation", len(cases), &input) {
		c := cases[i]
		s, a, loc := c.set, c.add, c.state.Location
		if want := render(s, a, loc); got != want {
			log.Fatalf("row %d: got %q want %q", i, got, want)
		}
		replaceWall, overflowMin, clampExact, keepFlag := s, a, a, a
		if c.state.Wall&flag == 0 && s.wall&flag != 0 {
			replaceWall.wall &^= spare
		}
		if c.state.Wall&flag == 0 && c.state.Ext == flag && c.state.Argument == minusOne {
			overflowMin.ext = flag
		}
		if a.ext == flag {
			clampExact.ext = flag + 1
		}
		if c.state.Wall&flag != 0 && a.wall&flag == 0 {
			keepFlag.wall |= flag
		}
		for j, wrong := range []string{
			render(replaceWall, a, loc), render(s, overflowMin, loc), render(s, clampExact, loc),
			render(s, keepFlag, loc), render(s, a, 0),
		} {
			if wrong != got {
				count[j]++
			}
		}
	}
	if want := [5]int{1, 1, 1, 1, 6}; count != want {
		log.Fatalf("mutant counts %v want %v", count, want)
	}
	must(os.Setenv("GOLEAN_TIME_CLOCK_MUTATION_ROWS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockMutation",
		"difftest/time/clock-mutation/testdata/bridge_test.go",
		"difftest/time/clock-mutation/testdata/mutation_test.go"))
	fmt.Printf("clock mutations: %d literal pinned Go and Lean cases, mutation counts %v\n", len(cases), count)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
