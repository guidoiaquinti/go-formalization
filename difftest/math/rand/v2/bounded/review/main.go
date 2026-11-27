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
	Bound uint64
	Words []uint64
	Want  string
}

func ok(value uint64, pulls int) string { return fmt.Sprintf("ok:%d:%d", value, pulls) }
func invalid(method string) string      { return "invalid argument to " + method + ":0:0" }
func same(value string) []string {
	row := make([]string, 10)
	for i := range row {
		row[i] = value
	}
	return row
}

func main() {
	// 3 * inverse = 2 * 2^64 + 1, so this draw sits exactly at threshold.
	const inverse = uint64(0xaaaaaaaaaaaaaaab)
	zero := []string{ok(uint64(^uint(0)), 1), ok(math.MaxUint32, 1), invalid("Int64N"), invalid("Uint64N"),
		invalid("Int32N"), invalid("Uint32N"), invalid("IntN"), invalid("UintN"), invalid("N"), invalid("N")}
	narrow := same(ok(1, 1))
	narrow[8] = invalid("N")
	truncated := same(ok(14, 1))
	truncated[8], truncated[9] = invalid("N"), invalid("N")
	max := []string{ok(math.MaxUint64-1, 2), ok(math.MaxUint32-1, 2), invalid("Int64N"), ok(math.MaxUint64-1, 2),
		invalid("Int32N"), ok(math.MaxUint32-1, 2), invalid("IntN"), ok(uint64(^uint(0))-1, 2), invalid("N"), ok(math.MaxUint16-1, 2)}
	cases := []fixture{
		{3, []uint64{0, 0, inverse}, strings.Join(same(ok(2, 3)), "|")},
		{3, []uint64{0, 0}, strings.Join(same("callback:0:3"), "|")},
		{8, []uint64{13, 0}, strings.Join(same(ok(5, 1)), "|")},
		{8, nil, strings.Join(same("callback:0:1"), "|")},
		{0, []uint64{math.MaxUint64}, strings.Join(zero, "|")},
		{128, []uint64{129}, strings.Join(narrow, "|")},
		{65536, []uint64{65550}, strings.Join(truncated, "|")},
		{math.MaxUint64, []uint64{0, math.MaxUint64}, strings.Join(max, "|")},
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d", strconv.IntSize, c.Bound)
		for _, v := range c.Words {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	var rejects [4]int
	for i, got := range check.LeanRows("godiff-rand-v2-bounded", len(cases), &wire) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].Want)
		}
		for j, value := range strings.Split(got, "|") {
			if i == 0 && value != ok(0, 1) {
				rejects[0]++
			}
			if i == 0 && value != "callback:0:4" {
				rejects[1]++
			}
			if i == 1 && value != "callback:0:2" {
				rejects[2]++
			}
			if i == 4 && j >= 2 && value != ok(0, 1) {
				rejects[3]++
			}
		}
	}
	if rejects != [4]int{10, 10, 10, 8} {
		log.Fatalf("unexpected control counts %v", rejects)
	}
	f, err := os.CreateTemp("", "golean-rand-bounded-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_RAND_BOUNDED_CASES", f.Name()))
	const dir = "difftest/math/rand/v2/bounded/testdata/"
	must(check.TestInPackage("math/rand/v2", "TestGoLeanRandBounded", dir+"bridge_test.go", dir+"bounded_test.go"))
	fmt.Printf("Bounded rand: %d independent native and Lean fixtures, mutant rejections %v\n", len(cases), rejects)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
