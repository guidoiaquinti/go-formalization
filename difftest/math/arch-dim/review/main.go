package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const (
	posInf    uint64 = 0x7ff0000000000000
	negInf    uint64 = 0xfff0000000000000
	negZero   uint64 = 0x8000000000000000
	sigA      uint64 = 0x7ff0000000000051
	sigB      uint64 = 0xfff0000000003752
	quietBit  uint64 = 0x0008000000000000
	quietA    uint64 = 0x7ff8000000000033
	quietB    uint64 = 0xfff8000000000044
	canonical uint64 = 0x7ff8000000000001
)

type fixture struct {
	a, b           uint64
	armMax, armMin uint64
	x86Max, x86Min uint64
}

func main() {
	cases := []fixture{
		{quietA, quietB, quietB, quietB, canonical, canonical},
		{quietB, quietA, quietA, quietA, canonical, canonical},
		{sigA, sigB, sigB | quietBit, sigB | quietBit, canonical, canonical},
		{sigB, sigA, sigA | quietBit, sigA | quietBit, canonical, canonical},
		{sigA, quietB, sigA | quietBit, sigA | quietBit, canonical, canonical},
		{quietA, sigB, sigB | quietBit, sigB | quietBit, canonical, canonical},
		{posInf, sigB, posInf, sigB | quietBit, posInf, canonical},
		{sigA, negInf, sigA | quietBit, negInf, canonical, negInf},
		{0, negZero, 0, negZero, 0, negZero},
		{negZero, 0, 0, negZero, 0, negZero},
		{negZero, negZero, negZero, negZero, negZero, negZero},
		{1, negZero + 1, 1, negZero + 1, 1, negZero + 1},
		{negZero, negZero + 1, negZero, negZero + 1, negZero, negZero + 1},
	}
	native := 2
	switch runtime.GOARCH {
	case "amd64":
		native = 0
	case "arm64":
		native = 1
	case "loong64", "s390x", "riscv64":
		fmt.Printf("Native %s outside the reviewed assembly profiles\n", runtime.GOARCH)
		native = -1
	}
	var model bytes.Buffer
	var wants []string
	f, err := os.CreateTemp("", "golean-arch-dim-review-")
	must(err)
	defer os.Remove(f.Name())
	nativeCount := 0

	var misses [3]int
	for _, c := range cases {
		rows := []string{fmt.Sprintf("%d,%d", c.x86Max, c.x86Min), fmt.Sprintf("%d,%d", c.armMax, c.armMin), "panic:not implemented"}
		for profile, want := range rows {
			fmt.Fprintln(&model, profile, c.a, c.b)
			wants = append(wants, want)
		}
		if native >= 0 {
			nativeCount++

			_, err = fmt.Fprintf(f, "%d %d %q\n", c.a, c.b, rows[native])
			must(err)
		}
		if rows[0] != rows[1] {
			misses[0]++
		}
		if c.a&^negZero > posInf && c.b&^negZero > posInf && c.armMax != (c.a|quietBit) {
			misses[1]++
		}
		if c.armMax == posInf && c.b&^negZero > posInf || c.armMin == negInf && c.a&^negZero > posInf {
			misses[2]++
		}
	}
	if misses != [3]int{8, 5, 2} {
		log.Fatalf("variant checks did not distinguish wanted rows: %v", misses)
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-math-arch-dim", len(wants), &model) {
		if got != wants[i] {
			log.Fatalf("row %d: Lean=%s want=%s", i, got, wants[i])
		}
	}
	if native >= 0 {
		must(os.Setenv("GOLEAN_MATH_ARCH_DIM_CASES", f.Name()))
		const data = "difftest/math/arch-dim/testdata/"
		must(check.TestInPackage("math", "TestGoLeanArchDim", data+"bridge_test.go", data+"arch_dim_test.go"))
	}
	fmt.Printf("archMax/Min: %d native fixture pairs, %d Lean profile pairs, three variants fail %v\n", nativeCount, len(wants), misses)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
