package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	code, operand int
	flag          bool
	want          string
}

func main() {
	cases := []fixture{
		{0, 0, false, "unknownVal|boolVal:false|0|unknown|unknown|1|false|false|false|Unknown"},
		{1, 1, true, "unknownVal|boolVal:true|0|unknown|unknown|1|true|true|true|Bool"},
		{2, 2, true, "unknownVal|boolVal:true|0|unknown|unknown|1|true|true|panic:<nil> not a Bool|String"},
		{3, 3, false, "unknownVal|boolVal:false|0|unknown|unknown|1|false|false|panic:\"x\" not a Bool|Int"},
		{4, 4, true, "unknownVal|boolVal:true|0|unknown|unknown|1|true|true|panic:4 not a Bool|Float"},
		{5, 1, false, "unknownVal|boolVal:false|0|unknown|unknown|1|false|false|false|Complex"},
		{6, 0, true, "unknownVal|boolVal:true|0|unknown|unknown|1|true|true|false|Kind(6)"},
		{-1, 4, false, "unknownVal|boolVal:false|0|unknown|unknown|1|false|false|panic:-1 not a Bool|Kind(-1)"},
		{math.MinInt, 0, false, "unknownVal|boolVal:false|0|unknown|unknown|1|false|false|false|Kind(" + strconv.Itoa(math.MinInt) + ")"},
	}
	mutants := [...]string{"nil", "boolVal:false", "1", "Unknown", "", "0", "false", "TRUE", "false", "Unknown"}
	var rejected [len(mutants)]int
	native, err := os.CreateTemp("", "golean-constant-scalar-control-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(native.Name())
	var input bytes.Buffer
	for _, c := range cases {
		flag := 0
		if c.flag {
			flag = 1
		}
		fmt.Fprintln(&input, c.code, flag, c.operand)
		if _, err := fmt.Fprintf(native, "%d %t %d\t%s\n", c.code, c.flag, c.operand, c.want); err != nil {
			log.Fatal(err)
		}
	}
	if err := native.Close(); err != nil {
		log.Fatal(err)
	}
	for i, got := range check.LeanRows("godiff-go-constant-scalars", len(cases), &input) {
		if got != cases[i].want {
			log.Fatalf("Lean row %d got %q want %q", i, got, cases[i].want)
		}
		fields := strings.Split(got, "|")
		for j, mutant := range mutants {
			if fields[j] != mutant {
				rejected[j]++
			}
		}
	}
	wantRejected := [len(mutants)]int{9, 4, 9, 9, 9, 9, 4, 9, 5, 8}
	if rejected != wantRejected {
		log.Fatalf("got mutation counts %v, want %v", rejected, wantRejected)
	}
	if err := os.Setenv("GOLEAN_CONSTANT_SCALARS", native.Name()); err != nil {
		log.Fatal(err)
	}
	if err := check.TestInPackage("go/constant", "TestGoLeanScalars", "difftest/go/constant/scalars/testdata/scalars_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("scalar review: %d pinned native and Lean rows, mutant rejection counts %v\n", len(cases), rejected)
}
