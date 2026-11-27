package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                    []int64
	A, B, Mode, FailAt, Mask int
	Want                     string
}

func wire(status string, values []int64, trace string) string {
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = fmt.Sprint(value)
	}
	return status + "/" + strings.Join(fields, ",") + "/" + trace
}

func main() {
	short := append([]int64{91, 9, 1}, slices.Repeat([]int64{9}, 47)...)
	short = append(short, 92)
	long := append([]int64{91, 9, 1}, slices.Repeat([]int64{9}, 48)...)
	long = append(long, 92)
	swapped := slices.Clone(long)
	swapped[1], swapped[2] = 1, 9

	prefix := append([]int64{100, 50, 10, 11, 0}, slices.Repeat([]int64{20}, 47)...)
	prefix = append(prefix, 92)
	full := slices.Clone(prefix)
	copy(full, []int64{0, 100, 50, 10, 11})
	during := slices.Clone(prefix)
	copy(during, []int64{100, 0, 50, 10, 11})
	six := make([]int64, 51)
	for i := range six {
		six[i] = int64(i)
	}
	for i := 0; i < 24; i += 4 {
		six[i], six[i+1] = six[i+1], six[i]
	}
	fiveFixed := slices.Clone(six)
	for i := 0; i < 20; i += 4 {
		fiveFixed[i], fiveFixed[i+1] = fiveFixed[i+1], fiveFixed[i]
	}
	cases := []fixture{
		{[]int64{91, 92}, 1, 1, 0, -1, 0, "ok:false/91,92/|ok:false/91,92/"},
		{[]int64{91, 7, 92}, 1, 2, 0, -1, 0, "ok:true/91,7,92/|ok:true/91,7,92/"},
		{short, 1, 50, 1, -1, 0, wire("ok:false", short, "") + "|" + wire("ok:false", short, "1,9")},
		{long, 1, 51, 1, 1, 0, wire("ok:true", swapped, "") + "|" + wire("panic:1", swapped, "1,9;9,9")},
		{prefix, 2, 52, 1, 4, 0, wire("ok:true", full, "") + "|" + wire("panic:4", during, "11,10;0,11;0,10;0,50;0,100")},
		{six, 0, 50, 0, -1, 0, wire("ok:false", fiveFixed, "") + "|" + wire("nil", six, "")},
	}
	early := cases[2]
	badShort := slices.Clone(early.Words)
	if early.B-early.A >= 49 {
		badShort[1], badShort[2] = badShort[2], badShort[1]
	}
	if slices.Equal(badShort, early.Words) {
		panic("early shift threshold mutation survived")
	}
	boundedLeft := slices.Clone(prefix)
	boundedLeft[3], boundedLeft[4] = boundedLeft[4], boundedLeft[3]
	boundedLeft[2], boundedLeft[3] = boundedLeft[3], boundedLeft[2]
	if slices.Equal(boundedLeft, full) || boundedLeft[0] != 100 {
		panic("range-bounded left-shift mutation survived")
	}
	if slices.Equal(long, swapped) {
		panic("initial swap rollback mutation survived")
	}
	fmt.Println("rejected shift at 49, left shift stopping at a and rollback after comparator panic")

	sixFixed := slices.Clone(fiveFixed)
	sixFixed[20], sixFixed[21] = sixFixed[21], sixFixed[20]
	if slices.Equal(sixFixed, fiveFixed) || sixFixed[20] != 20 {
		panic("six-repair mutation survived")
	}

	f, err := os.CreateTemp("", "review-partial-insertion-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.A, " ", c.B, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
		for _, word := range c.Words {
			fmt.Fprint(&input, " ", word)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(f.Close())
	for i, row := range check.LeanRows("godiff-slices-partial-insertion", len(cases), &input) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_PARTIAL_INSERTION_CASES", f.Name()))
	const data = "difftest/slices/partial-insertion/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanPartialInsertion", data+"bridge_test.go", data+"partial_insertion_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
