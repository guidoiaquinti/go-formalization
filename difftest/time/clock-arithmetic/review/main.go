package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixed := []struct{ input, want string }{
		{"0 2 1 0 1 2 -1", "999999999 1 1 1000000000 1"},
		{"9223372037928517632 -2 1 9223372045444710400 10 2 1", "9223372037928517633 -1 1 -12 -12"},
		{"9223372037928517632 9223372036854775807 2 0 59453308801 1 1", "1 59453308801 2 0 9223371977401467006"},
		{"9223372036854775808 0 1 999999999 59453308799 2 -1", "999999999 59453308799 1 1 -59453308799"},
		{"18446744073635809791 42 2 9223372036854775808 0 1 1", "0 68043243392 2 42 42"},
		{"854775807 9223372036 1 0 0 2 0", "854775807 9223372036 1 9223372036854775807 9223372036"},
		{"854775808 9223372036 2 0 0 1 0", "854775808 9223372036 2 9223372036854775807 9223372036"},
		{"145224192 -9223372037 2 0 0 1 0", "145224192 -9223372037 2 -9223372036854775808 -9223372037"},
		{"145224191 -9223372037 2 0 0 1 0", "145224191 -9223372037 2 -9223372036854775808 -9223372037"},
		{"9223372036854775808 9223372036854775807 0 9223372036854775808 -1 1 0", "9223372036854775808 9223372036854775807 0 9223372036854775807 9223372036854775807"},
		{"9223372036854775808 -9223372036854775808 1 9223372036854775808 1 2 0", "9223372036854775808 -9223372036854775808 1 -9223372036854775808 -9223372036854775808"},
		{"1073741823 0 2 0 -1 1 999999999", "1073741822 1 2 9223372036854775807 1"},
	}
	f, err := os.CreateTemp("", "golean-time-arithmetic-review")
	must(err)
	defer os.Remove(f.Name())
	var input bytes.Buffer
	for _, c := range fixed {
		fmt.Fprintln(&input, c.input)
		_, err := fmt.Fprintln(f, c.input, c.want)
		must(err)
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-time-clock-arithmetic", len(fixed), &input) {
		if got != fixed[i].want {
			log.Fatalf("fixture %d: Lean %q, expected %q", i, got, fixed[i].want)
		}
	}
	must(os.Setenv("GOLEAN_CLOCK_ARITHMETIC", f.Name()))
	const bridge = "difftest/time/clock-arithmetic/testdata/bridge_test.go"
	must(check.TestInPackage("time", "TestGoLeanClockArithmetic", bridge, "difftest/time/clock-arithmetic/testdata/arithmetic_test.go"))
	must(check.TestInPackage("time", "TestGoLeanArithmeticReview", bridge, "difftest/time/clock-arithmetic/review/testdata/trace_test.go"))
	fmt.Printf("%d literal Go/Lean arithmetic controls and native boundary traces passed\n", len(fixed))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
