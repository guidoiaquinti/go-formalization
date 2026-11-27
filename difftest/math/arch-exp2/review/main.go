package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	type fixture struct{ x, result float64 }
	cases := []fixture{
		{math.Inf(-1), 0}, {math.Inf(1), math.Inf(1)}, {math.Copysign(0, -1), 1},
		{-1075, 0}, {-1074.5, 0}, {math.Nextafter(-1074, math.Inf(-1)), 0},
		{math.Float64frombits(0x7ff8000000000042), math.Float64frombits(0x7ff8000000000042)},
	}
	for _, exponent := range []int{-1074, -1073, -1022, -1021, -10, 0, 10, 1022, 1023, 1024} {
		cases = append(cases, fixture{float64(exponent), math.Ldexp(1, exponent)})
	}
	var model bytes.Buffer
	for _, p := range []int{0, 2} {
		for _, c := range cases {
			fmt.Fprintf(&model, "%d %d\n", p, math.Float64bits(c.x))
		}
	}
	rows := check.LeanRows("godiff-math-arch-exp2", 2*len(cases), &model)
	var fixtureRows bytes.Buffer
	for i, c := range cases {
		want := strconv.FormatUint(math.Float64bits(c.result), 10)
		if rows[i] != "panic:not implemented" || rows[i+len(cases)] != want {
			log.Fatalf("archExp2 %#016x: stub %s, arm %s, want %s", math.Float64bits(c.x), rows[i], rows[i+len(cases)], want)
		}
		native := rows[i]
		if runtime.GOARCH == "arm64" {
			native = want
		}
		fmt.Fprintf(&fixtureRows, "%d %q %d true\n", math.Float64bits(c.x), native, math.Float64bits(c.result))
	}
	if runtime.GOARCH == "loong64" || runtime.GOARCH == "riscv64" {
		fmt.Printf("%d profile rows passed; native port outside scope.\n", len(rows))
		return
	}
	out, err := os.CreateTemp("", "golean-exp2-review-*")
	must(err)
	defer os.Remove(out.Name())
	_, err = out.Write(fixtureRows.Bytes())
	must(err)
	must(out.Close())
	must(os.Setenv("GOLEAN_MATH_ARCH_EXP2_CASES", out.Name()))
	const data = "difftest/math/arch-exp2/testdata/"
	must(check.TestInPackage("math", "TestGoLeanArchExp2", data+"bridge_test.go", data+"arch_exp2_test.go"))
	fmt.Printf("archExp2: %d native, %d model observations passed.\n", len(cases), len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
