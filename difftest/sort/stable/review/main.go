package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Tag, N int
	Keys         []int
	Want         string
}

func main() {
	cases := []fixture{
		{0, 0, 3, []int{2, 1, 1}, "ok:0:1,2,0:L1,0;S1,0;L2,1;S2,1;L1,0"},
		{1, 0, 3, []int{2, 1, 1}, "ok:0:1,2,0:L1,0;S1,0;L2,1;S2,1;L1,0"},
		{2, 0, 2, []int{2, 1, 0}, "ok:1:1,0,2:L1,0;S1,0"},
		{2, 1, 2, []int{2, 1}, "nil-interface:0:0,1:"},
		{3, 0, 0, []int{2, 1, 1}, "ok:0:1,2,0:L1,0;L2,1;L1,0"},
		{3, 1, 0, []int{1}, "ok:0:0:"},
		{3, 1, 0, []int{2, 1}, "nil-less:0:0,1:"},
		{3, 3, 0, nil, "non-slice:0::"},
	}
	keys := make([]int, 21)
	ids := []string{"20"}
	var events []string
	for i := range 20 {
		keys[i] = i + 1
		ids = append(ids, strconv.Itoa(i))
		if i > 0 {
			events = append(events, fmt.Sprintf("L%d,%d", i, i-1))
		}
	}
	events = append(events, "L20,10", "L20,5", "L20,2", "L20,1", "L20,0")
	for i := 20; i > 0; i-- {
		events = append(events, fmt.Sprintf("S%d,%d", i, i-1))
	}
	for kind := range 2 {
		cases = append(cases, fixture{kind, 0, 21, keys,
			"ok:0:" + strings.Join(ids, ",") + ":" + strings.Join(events, ";")})
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d", c.Kind, c.Tag, c.N)
		for _, v := range c.Keys {
			fmt.Fprintf(&input, " %d", v)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-sort-stable", len(cases), &input)
	file, err := os.CreateTemp("", "golean-sort-stable-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_STABLE_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanStable", "difftest/sort/stable/testdata/stable_test.go"))
	if strings.HasPrefix(rows[0], "ok:0:2,1,0:") || strings.HasPrefix(rows[2], "ok:2:") ||
		strings.HasPrefix(rows[5], "nil-less:") || strings.Contains(rows[8], ";L20,19;") {
		log.Fatal("tie reordering, duplicate Len, eager nil Less or whole-range insertion mutant accepted")
	}
	fmt.Println("Ten independent observations matched pinned Go and Lean; four mutant behaviors rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
