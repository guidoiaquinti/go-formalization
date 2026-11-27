package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	input, want float64
	exact       bool
}

func agrees(got float64, test fixture) bool {
	if math.IsNaN(test.want) {
		return math.IsNaN(got)
	}
	if test.exact {
		return math.Float64bits(got) == math.Float64bits(test.want)
	}
	return math.Abs(got-test.want) <= 0x1p-51
}

func main() {
	negativeZero := math.Copysign(0, -1)
	// Finite golden values were obtained independently with platform libm.
	tests := []fixture{
		{1.0 / 3, 0x1.4978fa3269ee1p-2, false}, {3, 0x1.3fc176b7a8560p+0, false},
		{-1.0 / 7, -0x1.229aec47638dcp-3, false}, {-7, -0x1.6dcc57bb565fdp+0, false},
		{1.0 / 16, 0x1.ff55bb72cfdeap-5, false}, {16, 0x1.82250768ac529p+0, false},
		{0x1.126145e9ecd58p-2, 0x1.0c152382d7367p-2, false},
		{0x1.a827999fcef34p-2, 0x1.921fb54442d1ap-2, false},
		{1, math.Pi / 4, false},
		{1.5, 0x1.f730bd281f69bp-1, false},
		{2.5, 0x1.30b6d796a4da8p+0, false},
		{100, 0x1.8f905eb2def22p+0, false},
		{0x1p-26, 0x1.fffffffffffffp-27, false},
		{0, 0, true}, {negativeZero, negativeZero, true},
		{math.Inf(1), math.Pi / 2, true}, {math.Inf(-1), -math.Pi / 2, true},
		{math.NaN(), math.NaN(), true},
	}
	const profiles = 2
	var input strings.Builder
	var cases bytes.Buffer
	var lostSign, lostShift int
	for _, test := range tests {
		for profile := range profiles {
			fmt.Fprintf(&input, "%d %d\n", profile, math.Float64bits(test.input))
		}
		got := math.Atan(test.input)
		if !agrees(got, test) {
			log.Fatalf("Go Atan(%g) = %.17g, want %.17g", test.input, got, test.want)
		}
		observed := "NaN"
		if !math.IsNaN(got) {
			observed = strconv.FormatUint(math.Float64bits(got), 10)
		}
		fmt.Fprintln(&cases, math.Float64bits(test.input), observed, math.Float64bits(test.want), test.exact)
		if !agrees(math.Abs(got), test) {
			lostSign++
		}
		x := math.Abs(test.input)
		if x > 2.414213562373095 && !agrees(-math.Copysign(math.Atan(1/x), test.input), test) {
			lostShift++
		}
	}
	rows := check.LeanRows("godiff-math-atan", profiles*len(tests), strings.NewReader(input.String()))
	values := make([]float64, len(rows))
	for i, row := range rows {
		got := math.NaN()
		if row != "NaN" {
			word, err := strconv.ParseUint(row, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			got = math.Float64frombits(word)
		}
		values[i] = got
		test := tests[i/profiles]
		if !agrees(got, test) {
			log.Fatalf("Lean profile %d Atan(%g) = %.17g, want %.17g", i%profiles, test.input, got, test.want)
		}
	}
	for i := 0; i < 6; i += 2 {
		for profile := range profiles {
			sum := values[profiles*i+profile] + values[profiles*(i+1)+profile]
			if math.Abs(sum-math.Copysign(math.Pi/2, tests[i].input)) > 0x1p-51 {
				log.Fatalf("Lean profile %d failed reciprocal identity at %g", profile, tests[i].input)
			}
		}
	}
	dir, err := os.MkdirTemp("", "golean-private-atan-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_ATAN_CASES", path))
	const overlay = "difftest/math/atan/testdata/"
	must(check.TestInPackage("math", "TestGoLeanAtan", overlay+"bridge_test.go", overlay+"atan_test.go"))
	if lostSign == 0 || lostShift == 0 {
		log.Fatal("controls did not distinguish both faults")
	}
	fmt.Printf("Private atan passed %d native and %d Lean profile fixtures; losing signs failed %d, omitting the large-argument offset failed %d\n",
		len(tests), len(rows), lostSign, lostShift)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
