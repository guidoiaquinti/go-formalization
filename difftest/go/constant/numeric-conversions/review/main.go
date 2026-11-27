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
	X    int64
	Big  string
	Tag  int
	Want string
}

func main() {
	cases := []fixture{
		{0, "0", 0, "0|0/1|0/512|0|0/1|0/512|0/1|nil:0"},
		{9007199254740993, "340282366920938463463374607431768211457", 1, "0|0/1|0/512|9007199254740993|9007199254740993/1|9007199254740993/512|340282366920938463463374607431768211457/1|typedNil:0"},
		{-1, "-340282366920938463463374607431768211457", 4, "0|0/1|0/512|-1|-1/1|-1/512|-340282366920938463463374607431768211457/1|4:0"},
		{math.MinInt64, "7", 3, "0|0/1|0/512|-9223372036854775808|-9223372036854775808/1|-9223372036854775808/512|7/1|3:0"},
	}
	f, err := os.CreateTemp("", "golean-constant-convert-review-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c.X, c.Big, c.Tag)
		must(enc.Encode(c))
	}
	must(f.Close())
	var rejected [8]int
	for i, got := range check.LeanRows("godiff-go-constant-numeric-conversions", len(cases), &input) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
		x, _ := strconv.ParseInt(c.Big, 10, 64)
		wrong := [8]string{"1", "0/0", "0/53", strconv.FormatInt(int64(int32(c.X)), 10),
			fmt.Sprintf("%d/1", uint64(c.X)), strconv.FormatFloat(float64(c.X), 'f', 0, 64) + "/512",
			fmt.Sprintf("%d/1", x), "nil:0"}
		for j, field := range strings.Split(got, "|") {
			if field != wrong[j] {
				rejected[j]++
			}
		}
	}
	if want := [8]int{4, 4, 4, 2, 2, 1, 2, 3}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_CONSTANT_CONVERSION_ROWS", f.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanNumericConversions",
		"difftest/go/constant/numeric-conversions/testdata/conversion_test.go"))
	fmt.Printf("constant conversions: %d pinned native and Lean cases, mutations %v\n", len(cases), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
