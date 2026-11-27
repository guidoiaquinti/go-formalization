package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	code, variant int
	flag          bool
}

func oracle(c fixture) string {
	names := []string{"Unknown", "Bool", "String", "Int", "Float", "Complex"}
	name := fmt.Sprintf("Kind(%d)", c.code)
	if c.code >= 0 && c.code < len(names) {
		name = names[c.code]
	}
	boolean := "false"
	if c.flag {
		boolean = "true"
	}
	operand := []string{"false", boolean, "panic:<nil> not a Bool", "panic:\"x\" not a Bool",
		"panic:" + strconv.Itoa(c.code) + " not a Bool"}
	return fmt.Sprintf("unknownVal|boolVal:%s|0|unknown|unknown|1|%s|%s|%s|%s",
		boolean, boolean, boolean, operand[c.variant], name)
}

func main() {
	file, err := os.CreateTemp("", "golean-constant-scalars-")
	must(err)
	defer os.Remove(file.Name())
	defer file.Close()
	out := bufio.NewWriter(file)
	compare := func(label string, fixtures []fixture) {
		var input bytes.Buffer
		for _, c := range fixtures {
			flag := 0
			if c.flag {
				flag = 1
			}
			fmt.Fprintln(&input, c.code, flag, c.variant)
		}
		for i, row := range check.LeanRows("godiff-go-constant-scalars", len(fixtures), &input) {
			c := fixtures[i]
			want := oracle(c)
			if row != want {
				log.Fatalf("%s %+v: Lean %q want %q", label, c, row, want)
			}
			_, err := fmt.Fprintf(out, "%d %t %d\t%s\n", c.code, c.flag, c.variant, want)
			must(err)
		}
		fmt.Printf("%s: %d observations\n", label, len(fixtures))
	}
	var fixed []fixture
	for _, code := range []int{math.MinInt, math.MinInt + 1, -10, -1, 0, 1, 2, 3, 4, 5, 6, 7, 10, math.MaxInt - 1, math.MaxInt} {
		for variant := range 5 {
			for _, flag := range []bool{false, true} {
				fixed = append(fixed, fixture{code, variant, flag})
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{int(r.Uint64()), r.Intn(5), r.Intn(2) == 1}
		}
		compare(fmt.Sprintf("seed=%d", seed), cases)
	}
	must(out.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_CONSTANT_SCALARS", file.Name()))
	must(check.TestInPackage("go/constant", "TestGoLeanScalars", "difftest/go/constant/scalars/testdata/scalars_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
