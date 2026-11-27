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
		same(0, sign|inf), same(sign, sign|inf),
		same(inf, inf), same(sign|inf, inf+1),
		same(inf+1, inf+1), same(sign|inf|2, inf+1),
		same(sign|1, inf+1), same(bits(-1), inf+1),
		same(bits(1), 0), same(bits(0.5), 0xbfe62e42fefa39ef),
		same(bits(2), 0x3fe62e42fefa39ef),
		same(0x0010000000000000, 0xc086232bdd7abcd2),
		same(1, 0xc0874385446d71c3),
		same(2, 0xc0873df9b3adb335),
		{0x3ff2c83fc3f1dc35, [3]uint64{0x3fc4856f5eba8920, 0x3fc4856f5eba8920, 0x3fc4856f5eba891f}},
		{0x4016da8e15370af1, [3]uint64{0x3ffbe2973a74cc5a, 0x3ffbe2973a74cc59, 0x3ffbe2973a74cc59}},
	}
	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
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
	var forcedLeadingBitFailures, exchangedFailures int
	for _, c := range cases {
		got := observe(c.want[profile])
		fmt.Fprintln(&fixtureRows, c.word, got, 0, true)
		if c.word == 1 || c.word == 2 {
			// Both probes share this result under the forced-leading-bit model.
			if got != observe(0xc08628b76e3a7b61) {
				forcedLeadingBitFailures++
			}
		}
		if got != observe(c.want[(profile+1)%3]) {
			exchangedFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode+2, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-log", 3*len(cases), &inputs) {
		if want := observe(cases[i/3].want[i%3]); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if forcedLeadingBitFailures != 2 || exchangedFailures != []int{1, 1, 2}[profile] {
		log.Fatalf("forced leading bit failed %d, exchanged profile failed %d", forcedLeadingBitFailures, exchangedFailures)
	}
	dir, err := os.MkdirTemp("", "golean-private-log-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtureRows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_LOG_CASES", path); err != nil {
		log.Fatal(err)
	}
	const data = "difftest/math/log/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanLog", data+"bridge_test.go", data+"log_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Private log passed %d fixtures; forced leading bit failed %d; exchanging profiles failed %d\n",
		len(cases), forcedLeadingBitFailures, exchangedFailures)
}
