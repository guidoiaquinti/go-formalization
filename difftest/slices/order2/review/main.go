package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	negativeZero := uint64(1) << 63
	nan := uint64(0x7ff8000000000007)
	cases := []struct {
		Words                          []uint64
		A, B, Kind, Mode, Reply, Count int
		HasCounter                     bool
		Want                           string
	}{
		{[]uint64{77, 9, 2, 88}, 1, 2, 0, 1, -7, math.MaxInt, true,
			fmt.Sprintf("ok:2,1/%d/|ok:2,1/%d/2,9", math.MinInt, math.MinInt)},
		{[]uint64{9, 9}, 0, 1, 0, 1, 0, 3, true, "ok:0,1/3/|ok:0,1/3/9,9"},
		{[]uint64{1, 9}, 0, 1, 0, 2, 66, 3, false, "ok:0,1/nil/|panic:66/nil/9,1"},
		{[]uint64{9, 1}, 0, 1, 0, 1, -1, 3, false, "nil/nil/|nil/nil/1,9"},
		{[]uint64{1, 9}, 0, 1, 0, 0, 0, 3, true, "ok:0,1/3/|nil/3/"},
		{[]uint64{0, negativeZero}, 0, 1, 1, 0, 0, 3, true, "ok:0,1/3/|nil/3/"},
		{[]uint64{0, nan}, 0, 1, 1, 0, 0, 3, true, "ok:1,0/4/|nil/3/"},
	}
	reversed := cases[0]
	x, y := reversed.Words[reversed.A], reversed.Words[reversed.B]
	reversedTrace := fmt.Sprintf("ok:2,1/%d/|ok:2,1/%d/%d,%d", reversed.Count+1, reversed.Count+1, x, y)

	tied := cases[1]
	nonStrict := tied.Want
	if tied.Words[tied.B] <= tied.Words[tied.A] && tied.Reply <= 0 {
		nonStrict = fmt.Sprintf("ok:1,0/%d/|ok:1,0/%d/9,9", tied.Count+1, tied.Count+1)
	}

	panicking := cases[2]
	counterFirst := panicking.Want
	if !panicking.HasCounter {
		counterFirst = "ok:0,1/nil/|nil/nil/"
	}

	unordered := cases[6]
	rawFloat := unordered.Want
	if !(math.Float64frombits(unordered.Words[unordered.B]) < math.Float64frombits(unordered.Words[unordered.A])) {
		rawFloat = fmt.Sprintf("ok:0,1/%d/|nil/%d/", unordered.Count, unordered.Count)
	}

	for _, mutant := range []struct {
		name    string
		fixture int
		output  string
	}{
		{"reversed comparator operands", 0, reversedTrace},
		{"non-strict tie comparison", 1, nonStrict},
		{"counter access before callback", 2, counterFirst},
		{"raw float comparison ignores NaN ordering", 6, rawFloat},
	} {
		if cases[mutant.fixture].Want == mutant.output {
			panic("undetected mutation: " + mutant.name)
		}
		fmt.Printf("%s rejected with %s\n", mutant.name, mutant.output)
	}
	file, err := os.CreateTemp("", "review-order2-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	var input bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.HasCounter {
			present = 1
		}
		fmt.Fprint(&input, strconv.IntSize, " ", c.Kind, " ", c.A, " ", c.B, " ", present, " ", c.Count, " ", c.Mode, " ", c.Reply)
		for _, word := range c.Words {
			fmt.Fprint(&input, " ", word)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(file.Close())
	for i, row := range check.LeanRows("godiff-slices-order2", len(cases), &input) {
		if strings.TrimSpace(row) != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q", i, row))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_ORDER2_CASES", file.Name()))
	const data = "difftest/slices/order2/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanOrder2", data+"bridge_test.go", data+"order2_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
