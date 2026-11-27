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
		same(0, 0), same(sign, sign), same(1, 1), same(sign|1, sign|1),
		same(bits(0x1p-55), bits(0x1p-55)), same(bits(-0x1p-55), bits(-0x1p-55)),
		same(inf, inf), same(sign|inf, inf+1), same(inf+1, inf+1),
		same(bits(-1), sign|inf), same(bits(-1)+1, inf+1),
		same(bits(-1)-1, 0xc0425e4f7b2737fa),
		same(bits(0x1p-54), bits(0x1p-54)),
		same(bits(-0x1p-54), bits(-0x1p-54)),
		same(bits(1), 0x3fe62e42fefa39ef), same(bits(3), 0x3ff62e42fefa39ef),
		{0x40244723a8990e04, [3]uint64{0x40034898a9861628, 0x40034898a9861628, 0x40034898a9861629}},
		{0x40355ede9103df39, [3]uint64{0x4008dcaa6335b691, 0x4008dcaa6335b690, 0x4008dcaa6335b691}},
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
	var directSubtractFailures, exchangedFailures int
	for _, c := range cases {
		x := math.Float64frombits(c.word)
		got := c.want[profile]
		fmt.Fprintln(&fixtureRows, c.word, observe(got), 0, true)
		if math.Abs(x) < 0x1p-54 && bits(math.Log(1+x)) != got {
			directSubtractFailures++
		}
		if observe(c.want[(profile+1)%3]) != observe(got) {
			exchangedFailures++
		}
		for mode := range c.want {
			fmt.Fprintln(&inputs, mode, c.word)
		}
	}
	for i, row := range check.LeanRows("godiff-math-log1p", 3*len(cases), &inputs) {
		if want := observe(cases[i/3].want[i%3]); row != want {
			log.Fatalf("case %d, profile %d: Lean=%s want=%s", i/3, i%3, row, want)
		}
	}
	if directSubtractFailures != 5 || exchangedFailures != []int{1, 2, 1}[profile] {
		log.Fatalf("direct subtraction failed %d, exchanged profile failed %d", directSubtractFailures, exchangedFailures)
	}
	dir, err := os.MkdirTemp("", "golean-private-log1p-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixtureRows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_MATH_LOG1P_CASES", path); err != nil {
		log.Fatal(err)
	}
	const fixtures = "difftest/math/log1p/testdata/"
	if err := check.TestInPackage("math", "TestGoLeanLog1p", fixtures+"bridge_test.go", fixtures+"log1p_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Private log1p passed %d fixtures; Log(1+x) failed %d tiny probes; exchanging profiles failed %d\n",
		len(cases), directSubtractFailures, exchangedFailures)
}
