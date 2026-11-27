package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
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
		same(bits(-745), 1),
		same(bits(709.5), 0x7fe81e9b4b52d0c9),
		{0xc024a51a9d053b7c, [3]uint64{0x3f013ddc8a82e66a, 0x3f013ddc8a82e66a, 0x3f013ddc8a82e669}},
	}

	profile := 0
	if runtime.GOARCH == "arm64" {
		profile = 2
	} else {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("no build metadata")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				profile = 1
			}
		}
	}
	var inputs, fixtureRows bytes.Buffer
	var unscaledFailures, exchangedFailures int
	for _, c := range cases {
		got := c.want[profile]
		fmt.Fprintln(&fixtureRows, c.word, got, got, true)
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
			fmt.Fprintln(&inputs, mode+3, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-exp", 3*len(cases), &inputs) {
		if want := strconv.FormatUint(cases[i/3].want[i%3], 10); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if unscaledFailures != 2 || exchangedFailures != 1 {
		log.Fatalf("flush subnormal failed %d, swapped profile failed %d", unscaledFailures, exchangedFailures)
	}
	dir, err := os.MkdirTemp("", "golean-private-exp-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtureRows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_EXP_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/exp/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanExp", fixtures+"bridge_test.go", fixtures+"exp_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Private exp passed %d fixtures; flushing subnormal outputs failed %d; swapping profiles failed %d\n",
		len(cases), unscaledFailures, exchangedFailures)
}
