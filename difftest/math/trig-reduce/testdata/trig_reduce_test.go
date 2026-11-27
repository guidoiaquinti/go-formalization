package math_test

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

func TestGoLeanTrigReduce(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_MATH_TRIG_REDUCE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var count int
	var maxError float64
	for {
		var input, modelOctant, modelAngle, refOctant, refAngle uint64
		n, err := fmt.Fscan(reader, &input, &modelOctant, &modelAngle, &refOctant, &refAngle)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		x := math.Float64frombits(input)
		if !(x > 0) || math.IsInf(x, 0) {
			t.Fatalf("input %#016x is outside the positive finite domain", input)
		}
		octant, angle := math.GoLeanTrigReduce(x)
		if octant != modelOctant || math.Float64bits(angle) != modelAngle {
			t.Fatalf("input %#016x: Go=(%d, %#016x), Lean=(%d, %#016x)",
				input, octant, math.Float64bits(angle), modelOctant, modelAngle)
		}
		if x < math.Pi/4 && (octant != 0 || math.Float64bits(angle) != input) {
			t.Fatalf("input %#016x: small guard returned (%d, %#016x)", input, octant, math.Float64bits(angle))
		}
		if octant >= 8 || octant%2 != 0 || octant != refOctant {
			t.Fatalf("input %#016x: octant=%d want %d in {0,2,4,6}", input, octant, refOctant)
		}
		want := math.Float64frombits(refAngle)
		delta := math.Abs(angle - want)
		if math.IsNaN(delta) || delta > 0x1p-50 {
			t.Fatalf("input %#016x: angle=%.17g reference=%.17g absolute difference=%g", input, angle, want, delta)
		}
		maxError = max(maxError, delta)
		count++
	}
	if count == 0 {
		t.Fatal("no private trigReduce comparison rows")
	}
	t.Logf("private trigReduce matched %d Lean rows and reference octants; maximum absolute remainder error %g", count, maxError)
}
