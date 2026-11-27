package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	const sign = uint64(1) << 63
	const inf = uint64(0x7ff0000000000000)
	bits := math.Float64bits
	type fixture struct {
		word uint64
		want [3]uint64
	}
	same := func(word, want uint64) fixture {
		return fixture{word, [3]uint64{want, want, want}}
	}
	cases := []fixture{
		same(0, bits(1)), same(sign, bits(1)),
		same(inf, inf), same(sign|inf, 0),
		same(inf+1, inf+1), same(sign|inf|0x8000000000005, sign|inf|0x8000000000005),
		same(bits(-0x1p-53), bits(1)-1), same(bits(0x1p-53), bits(1)),
		same(bits(-744), 2),
		{bits(-745), [3]uint64{1, 1, 0}},
		{bits(709.5), [3]uint64{inf, inf, 0x7fe81e9b4b52d0c9}},
		{0x4026e5532e4f6b00, [3]uint64{0x40f6e0862114ffce, 0x40f6e0862114ffcf, 0x40f6e0862114ffd0}},
	}

	profile, err := nativeProfile()
	if err != nil {
		log.Fatal(err)
	}
	var inputs bytes.Buffer
	var unscaledFailures, exchangedFailures int
	for _, c := range cases {
		got := bits(math.Exp(math.Float64frombits(c.word)))
		if want := c.want[profile]; got != want {
			log.Fatalf("Exp(%#x)=%#x want %#x", c.word, got, want)
		}
		if got > 0 && got < uint64(1)<<52 {
			// Flushing the subnormal result changes the expected word.
			unscaledFailures++
		}
		other := 2
		if profile == 2 {
			other = 0
		}
		if c.want[other] != got {
			exchangedFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-exp", 3*len(cases), &inputs) {
		if want := strconv.FormatUint(cases[i/3].want[i%3], 10); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	expectedFlushFailures := 2
	if profile == 2 {
		expectedFlushFailures = 1
	}
	if unscaledFailures != expectedFlushFailures || exchangedFailures != 3 {
		log.Fatalf("flush subnormal failed %d, swapped profile failed %d", unscaledFailures, exchangedFailures)
	}
	fmt.Printf("Exp passed %d fixtures; flushing subnormal outputs failed %d; swapping profiles failed %d\n",
		len(cases), unscaledFailures, exchangedFailures)
}

func nativeProfile() (int, error) {
	switch runtime.GOARCH {
	case "arm64":
		return 2, nil
	case "amd64":
		dir, err := os.MkdirTemp("", "golean-exp-review-")
		if err != nil {
			return 0, err
		}
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "profile")
		if err := os.Setenv("GOLEAN_MATH_EXP_PROFILE", path); err != nil {
			return 0, err
		}
		const fixtures = "difftest/math/exp/testdata/"
		if err := check.TestInPackage("math", "TestGoLeanExpProfile", fixtures+"profile_amd64_test.go", fixtures+"observe_profile_amd64_test.go"); err != nil {
			return 0, err
		}
		value, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		fused, err := strconv.ParseBool(string(value))
		if fused {
			return 1, err
		}
		return 0, err
	default:
		return 0, fmt.Errorf("unmodeled Exp architecture %s", runtime.GOARCH)
	}
}
