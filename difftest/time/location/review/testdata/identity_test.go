package time_test

import (
	"testing"
	"time"
)

func TestGoLeanLocationIdentityControl(t *testing.T) {
	first := time.FixedZone("UTC", 0)
	second := time.FixedZone("UTC", 0)
	if first == second || first == time.UTC || second == time.UTC {
		t.Fatal("fixture requires equal-name locations with different pointers")
	}
	const rawWall = uint64(9223372044370968581)
	const rawExt = ^uint64(0)
	v := time.GoLeanLocatedTime(rawWall, rawExt, first)
	assert := func(label string, actual time.Time, wall, ext uint64, loc *time.Location) {
		t.Helper()
		w, e, l := time.GoLeanLocatedParts(actual)
		if w != wall || e != ext || l != loc {
			t.Fatalf("%s: got wall %d ext %d loc %p, expected %d %d %p", label, w, e, l, wall, ext, loc)
		}
	}
	const seconds = 59453308807
	assert("strip", time.GoLeanStripMono(v), 5, seconds, first)
	assert("UTC", v.UTC(), 5, seconds, nil)
	assert("Local", v.Local(), 5, seconds, time.Local)
	assert("second equal-name location", v.In(second), 5, seconds, second)
	assert("setLoc nil", time.GoLeanSetLoc(v, nil), 5, seconds, nil)
	assert("original", v, rawWall, rawExt, first)
	if v.Location() != first || v.In(second).Location() != second || v.UTC().Location() != time.UTC {
		t.Fatal("Location did not preserve raw identity or fall back to UTC")
	}
	func() {
		defer func() {
			if p := recover(); p != "time: missing Location in call to Time.In" {
				t.Fatalf("nil In panic: %v", p)
			}
		}()
		_ = v.In(nil)
	}()
}
