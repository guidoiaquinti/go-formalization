package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/bessel"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	profile := nativeProfile()
	edges := make([]float64, 0, 66)
	for _, word := range binary64.Boundaries() {
		edges = append(edges, binary64.Store(word))
	}
	for _, center := range []float64{2, 2.8571, 4.5454, 8, 20, 0x1p512, math.MaxFloat64} {
		for _, x := range []float64{math.Nextafter(center, 0), center, math.Nextafter(center, math.Inf(1))} {
			edges = append(edges, x, -x)
		}
	}
	dir, err := os.MkdirTemp("", "golean-math-pone-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fixture := bufio.NewWriter(file)
	compare := func(label string, cases []float64) error {
		var input bytes.Buffer
		for _, x := range cases {
			fmt.Fprintln(&input, profile, binary64.Load(x))
		}
		for i, row := range check.LeanRows("godiff-math-pone", len(cases), &input) {
			x := cases[i]
			if _, err := fmt.Fprintln(fixture, binary64.Load(x), row, reference(x)); err != nil {
				return err
			}
		}
		fmt.Printf("pone %s: %d Lean rows.\n", label, len(cases))
		return nil
	}
	if err := compare("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]float64, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		if err := compare(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_MATH_PONE_CASES", path); err != nil {
		return err
	}
	const data = "difftest/math/pone/testdata/"
	return check.TestInPackage("math", "TestGoLeanPone", data+"bridge_test.go", data+"pone_test.go")
}

func sample(rng *rand.Rand, i int) float64 {
	switch i % 6 {
	case 0:
		return binary64.Store(rng.Uint64())
	case 1:
		return 20*rng.Float64() - 2
	case 2:
		centers := []float64{2, 2.8571, 4.5454, 8, 20, 0x1p512}
		bits := binary64.Load(centers[rng.Intn(len(centers))])
		return binary64.Store(bits + uint64(rng.Intn(129)-64))
	case 3:
		return math.Ldexp(1+rng.Float64(), 1+rng.Intn(1023))
	case 4:
		return 1000 * rng.Float64()
	default:
		return 2 + 8*rng.Float64()
	}
}

func reference(x float64) string {
	if math.IsNaN(x) || x < 2 {
		return "panic"
	}
	if math.IsInf(x, 1) {
		return strconv.FormatUint(binary64.Load(1), 10)
	}
	return strconv.FormatUint(binary64.Load(bessel.POne(x)), 10)
}

func nativeProfile() int {
	switch runtime.GOARCH {
	case "arm64":
		return 2
	case "amd64":
		info, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("cannot determine pone build profile")
		}
		for _, setting := range info.Settings {
			if setting.Key == "GOAMD64" && (setting.Value == "v3" || setting.Value == "v4") {
				return 1
			}
		}
		return 0
	default:
		log.Fatalf("unmodeled pone architecture %s", runtime.GOARCH)
		return 0
	}
}
