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
	type fixture struct {
		input uint64
		want  [3]uint64
	}
	same := func(input, want uint64) fixture {
		return fixture{input, [3]uint64{want, want, want}}
	}
	cases := []fixture{
		same(0, 0), same(sign, sign),
		same(inf, inf+1), same(inf|sign, inf+1), same(inf+7, inf+1),
		same(1, 1), same(sign|1, sign|1),
		same(0x3fe921fb54442d18, 0x3fe6a09e667f3bcc),
		same(0x3ff921fb54442d18, 0x3ff0000000000000),
		same(0x400921fb54442d18, 0x3ca1a62633145c00),
		same(0x41bfffffffffffff, 0x3fd4e67c4aa3dd2f),
		same(0x41c0000000000000, 0x3fd4e67c0e2622df),
		same(0x41c0000000000001, 0x3fd4e67b952aadfe),
		same(0x47f0000000000000, 0x3fef321f756b2844),
		same(0x5ff0000000000000, 0x3fefb2d203b97e5e),
		same(0x7fe0000000000000, 0x3fe205248cbdb760),
		same(0x7fefffffffffffff, 0x3f7452fc98b34eb0),
		{0x41c6a429749b4154, [3]uint64{0xbfd8f8662ca17b64, 0xbfd8f8662ca17b65, 0xbfd8f8662ca17b65}},
		{0x41b5ffd3dd1671f2, [3]uint64{0xbfe18072017f5704, 0xbfe18072017f5703, 0xbfe18072017f5703}},
	}
	observe := func(word uint64) string {
		if math.IsNaN(math.Float64frombits(word)) {
			return "NaN"
		}
		return strconv.FormatUint(word, 10)
	}
	var inputs, fixtureRows bytes.Buffer
	profile := nativeProfile()
	var foldedFailures, reductionFailures int
	for _, c := range cases {
		x := math.Float64frombits(c.input)
		want := c.want[profile]
		got := observe(want)
		exact := math.IsNaN(x) || math.IsInf(x, 0) || x == 0
		fmt.Fprintln(&fixtureRows, c.input, got, want, exact)
		if got != observe(math.Float64bits(math.Sin(math.Abs(x)))) {
			foldedFailures++
		}
		if x >= 0x1p29 && !math.IsInf(x, 0) && math.Abs(math.Float64frombits(want)-math.Sin(math.Mod(x, 2*math.Pi))) > 0x1p-30 {
			reductionFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.input)
		}
	}
	for i, row := range check.LeanRows("godiff-math-sin", 3*len(cases), &inputs) {
		if want := observe(cases[i/3].want[i%3]); row != want {
			log.Fatalf("fixture %d profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if foldedFailures != 2 || reductionFailures != 7 {
		log.Fatalf("lost-sign mutant failed %d; rounded-period mutant failed %d", foldedFailures, reductionFailures)
	}
	dir, err := os.MkdirTemp("", "golean-private-sin-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtureRows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_SIN_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/sin/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanSin", fixtures+"bridge_test.go", fixtures+"sin_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Private sin passed %d fixtures; lost sign fails %d; rounded-period reduction fails %d\n",
		len(cases), foldedFailures, reductionFailures)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine Sin build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unreviewed architecture %s", runtime.GOARCH)
		return 0
	}
}
