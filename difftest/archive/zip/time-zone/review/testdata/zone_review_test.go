package zip

import (
	"testing"
	"time"
)

func TestGoLeanZipZoneReview(t *testing.T) {
	before := timeZone(14*time.Hour - 1)
	after := timeZone(14*time.Hour + 1)
	if before != after || before != time.FixedZone("", 50400) ||
		before == timeZone(13*time.Hour) {
		t.Fatal("rounded whole-hour aliases do not match the native cache")
	}
	zero := timeZone(99 * time.Hour)
	if zero != timeZone(-99*time.Hour) || zero != timeZone(0) ||
		zero != time.FixedZone("", 0) || zero == time.UTC {
		t.Fatal("out-of-range and rounded zero offsets did not share the unnamed zero cache")
	}
	quarter := timeZone(-450 * time.Second)
	if quarter == timeZone(-15*time.Minute) || quarter == time.FixedZone("", -900) {
		t.Fatal("a fractional-hour location unexpectedly reused a pointer")
	}
}
