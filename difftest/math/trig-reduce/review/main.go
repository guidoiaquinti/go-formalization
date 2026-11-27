package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	word, octant uint64
	residual     float64
}

func main() {
	// The expected pairs were computed independently with a decimal Machin
	// expansion of pi, keeping 410 digits before subtracting the multiple.
	fixtures := []fixture{
		{0x0000000000000001, 0, 4.9406564584124654e-324},
		{0x0010000000000000, 0, 2.2250738585072014e-308},
		{0x3fd0000000000000, 0, 0.25},
		{0x3fe921fb54442d17, 0, 0.78539816339744817},
		{0x3fe921fb54442d18, 0, 0.78539816339744828},
		{0x3fe921fb54442d19, 2, -0.78539816339744828},
		{0x3ff921fb54442d18, 2, -6.123233995736766e-17},
		{0x41bfffffffffffff, 4, -0.33266991106280741},
		{0x41c0000000000000, 4, -0.33266985145816264},
		{0x41c0000000000001, 4, -0.33266973224887308},
		{0x41cfffffffffffff, 0, -0.66533982212561482},
		{0x435fffffffffffff, 0, 0.055799108050987152},
		{0x4360000000000000, 6, -0.6565898723337027},
		{0x4360000000000001, 0, -0.51057150630818582},
		{0x4760000000000000, 6, -0.19372038052768475},
		{0x47f0000000000000, 2, -0.22466624209610969},
		{0x4e67c0ab20123456, 6, 0.47779663171047076},
		{0x55d3ba8225f404ee, 6, -0.32612048618130263},
		{0x5ff0000000000000, 2, -0.13737665396280732},
		{0x73efffffffffffff, 4, -0.069553071215476978},
		{0x7fe0000000000000, 4, -0.59816590407208436},
		{0x7fefffffffffffff, 4, -0.0049619751507872728},
	}
	var inputs strings.Builder
	for _, c := range fixtures {
		fmt.Fprintln(&inputs, c.word)
	}
	results := check.LeanRows("godiff-math-trig-reduce", len(fixtures), strings.NewReader(inputs.String()))
	var cases bytes.Buffer
	var collapsedOctant, roundedPeriod int
	for i, c := range fixtures {
		var octant, residualWord uint64
		if _, err := fmt.Sscan(results[i], &octant, &residualWord); err != nil {
			log.Fatal(err)
		}
		actual := math.Float64frombits(residualWord)
		if octant != c.octant || math.IsNaN(actual) || math.Abs(actual-c.residual) > 0x1p-50 ||
			c.word < 0x3fe921fb54442d18 && residualWord != c.word {
			log.Fatalf("Lean input %#016x reduced to (%d, %.17g), want (%d, %.17g)", c.word, octant, actual, c.octant, c.residual)
		}
		fmt.Fprintln(&cases, c.word, octant, residualWord, c.octant, math.Float64bits(c.residual))
		if c.octant != 0 {
			collapsedOctant++
		}
		if math.Abs(math.Remainder(math.Float64frombits(c.word), math.Pi/2)-c.residual) > 0x1p-50 {
			roundedPeriod++
		}
	}
	dir, err := os.MkdirTemp("", "golean-trig-reduce-review-")
	must(err)
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	must(os.WriteFile(path, cases.Bytes(), 0600))
	must(os.Setenv("GOLEAN_MATH_TRIG_REDUCE_CASES", path))
	const overlay = "difftest/math/trig-reduce/testdata/"
	must(check.TestInPackage("math", "TestGoLeanTrigReduce", overlay+"bridge_test.go", overlay+"trig_reduce_test.go"))
	if collapsedOctant == 0 || roundedPeriod == 0 {
		log.Fatalf("controls did not distinguish zero octant (%d) or rounded period (%d)", collapsedOctant, roundedPeriod)
	}
	fmt.Printf("trigReduce passed %d native and Lean fixtures; zeroed octant failed %d, rounded period failed %d\n",
		len(fixtures), collapsedOctant, roundedPeriod)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
