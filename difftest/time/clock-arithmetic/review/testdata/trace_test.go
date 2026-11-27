package time_test

import (
	"math"
	"testing"
	"time"
)

func TestGoLeanArithmeticReview(t *testing.T) {
	a := time.FixedZone("same", 75)
	b := time.FixedZone("same", 75)
	if a == b {
		t.Fatal("fixture requires distinct location pointers")
	}
	base := time.GoLeanArithmeticTime(0, 0, a)
	for _, c := range []struct {
		wall    uint64
		seconds int64
		want    time.Duration
		exact   bool
	}{
		{854775807, 9223372036, math.MaxInt64, true},
		{854775808, 9223372036, math.MaxInt64, false},
		{145224192, -9223372037, math.MinInt64, true},
		{145224191, -9223372037, math.MinInt64, false},
		{1073741823, 0, math.MaxInt64, false},
	} {
		v := time.GoLeanArithmeticTime(c.wall, c.seconds, b)
		d := v.Sub(base)
		result := base.Add(d)
		if d != c.want || result.Equal(v) != c.exact || result.Location() != a {
			t.Fatalf("wall=%d seconds=%d: duration %d exact %t loc %p", c.wall, c.seconds, d, result.Equal(v), result.Location())
		}
	}
	marked := time.GoLeanArithmeticTime((1<<63)+(1<<30), math.MaxInt64, b)
	next := marked.Add(1)
	wall, ext, loc := time.GoLeanArithmeticFields(next)
	if wall != 1 || ext != 59453308801 || loc != b {
		t.Fatalf("reading overflow did not strip marker and retain wall instant/location: %d %d %p", wall, ext, loc)
	}
	if time.GoLeanArithmeticSubMono(1, math.MinInt64) != math.MaxInt64 ||
		time.GoLeanArithmeticSubMono(math.MinInt64, 1) != math.MinInt64 ||
		time.GoLeanArithmeticSubMono(math.MinInt64+1, 1) != math.MinInt64 {
		t.Fatal("monotonic subtraction endpoint transcript disagrees")
	}
}
