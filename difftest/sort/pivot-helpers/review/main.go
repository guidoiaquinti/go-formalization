package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, A, B, C, Initial int
	Script                 []bool
	Want                   string
}

func main() {
	if strconv.IntSize != 64 {
		return
	}
	fixtures := []fixture{
		{0, 2, 5, 0, math.MaxInt, []bool{true}, "5:2:-9223372036854775808:5,2"},
		{0, 5, 2, 0, -4, []bool{false}, "5:2:-4:2,5"},
		{1, 3, 5, 7, 7, []bool{true, false, true}, "5:-1:9:5,3;7,3;3,5"},
		{1, 3, 5, 7, 0, []bool{false, false, false}, "5:-1:0:5,3;7,5;5,3"},
		{2, 5, 0, 0, 0, []bool{true, true, false}, "6:-1:2:5,4;6,4;6,5"},
		{2, 1, 0, 0, 0, []bool{false, false, false}, "1:-1:0:1,0;2,1;1,0"},
		{3, 0, 7, 0, 0, nil, "2:1:0:"},
		{3, 0, 8, 0, 0, nil, "4:1:0:4,2;6,4;4,2"},
		{3, 0, 49, 0, 0, nil, "24:1:0:24,12;36,24;24,12"},
		{3, 0, 50, 0, 0, nil, "24:1:0:12,11;13,12;12,11;24,23;25,24;24,23;36,35;37,36;36,35;24,12;36,24;24,12"},
		{3, 0, 50, 0, 0, []bool{true, true, true, true, true, true, true, true, true, true, true, true}, "24:2:0:12,11;13,11;13,12;24,23;25,23;25,24;36,35;37,35;37,36;24,12;36,12;36,24"},
	}
	var input bytes.Buffer
	for _, f := range fixtures {
		fmt.Fprint(&input, f.Kind, " ", f.A, " ", f.B, " ", f.C, " ", f.Initial)
		for _, less := range f.Script {
			fmt.Fprint(&input, " ", map[bool]int{true: 1}[less])
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-sort-pivot-helpers", len(fixtures), &input)
	f, err := os.CreateTemp("", "golean-sort-pivot-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	failures := [3]int{}
	for i, fixture := range fixtures {
		if rows[i] != fixture.Want {
			log.Fatalf("fixture %d: Lean got %q want %q", i, rows[i], fixture.Want)
		}
		must(enc.Encode(fixture))
		if fixture.Kind < 3 {
			got := strings.SplitN(rows[i], ":", 4)
			got[2] = strconv.Itoa(fixture.Initial)
			if strings.Join(got, ":") != fixture.Want {
				failures[0]++
			}
		}
		if fixture.Kind == 3 && fixture.B-fixture.A == 50 && rows[i] != rows[8] {
			failures[1]++
		}
		if fixture.Kind == 3 && strings.Contains(rows[i], ":2:0:") {
			failures[2]++
		}
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_SORT_PIVOT_CASES", f.Name()))
	const bridge = "difftest/sort/pivot-helpers/testdata/"
	must(check.TestInPackage("sort", "TestGoLeanPivotHelpers", bridge+"bridge_test.go", bridge+"pivot_helpers_test.go"))
	if failures != [3]int{3, 2, 1} {
		log.Fatalf("unexpected mutant failure counts: %v", failures)
	}
	fmt.Printf("11 fixed traces agree with Lean and both Go variants. Mutant failures: unchanged count %d, omitted ninther %d, dropped decreasing hint %d.\n", failures[0], failures[1], failures[2])
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
