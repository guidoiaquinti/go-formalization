package color

import (
	"runtime"
	"testing"
)

func TestGoLeanSqDiffControl(t *testing.T) {
	cases := []struct{ x, y, want uint32 }{
		{65535, 65535, 0},
		{1, 0, 0},
		{0, 1, 0},
		{2, 0, 1},
		{0, 2, 1},
		{3, 0, 2},
		{0, 3, 2},
		{40005, 40000, 6},
		{40000, 40005, 6},
		{65535, 0, 1073709056},
		{0, 65535, 1073709056},
	}
	var real, earlyDivision, roundUp int
	for _, c := range cases {
		if sqDiff(c.x, c.y) != c.want {
			real++
		}
		d := max(c.x, c.y) - min(c.x, c.y)
		half := d / 2
		if half*half != c.want {
			earlyDivision++
		}
		if (d*d+3)/4 != c.want {
			roundUp++
		}
	}
	if real != 0 || earlyDivision != 6 || roundUp != 8 {
		t.Fatalf("sqDiff failures: Go=%d early-division=%d round-up=%d", real, earlyDivision, roundUp)
	}
	t.Logf("%s %s/%s: sqDiff Go=%d/%d early-division=%d/%d round-up=%d/%d",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, real, len(cases),
		earlyDivision, len(cases), roundUp, len(cases))
}
