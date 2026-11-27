package slices_test

import (
	"slices"
	"testing"
	"unsafe"
)

func TestGoLeanOverlapsReview(t *testing.T) {
	words := make([][3]byte, 8)
	separate := make([][3]byte, 8)
	cases := []struct {
		a, b [][3]byte
		want bool
	}{
		{words[1:3], words[2:4], true},
		{words[2:4], words[1:3], true},
		{words[1:3], words[3:5], false},
		{words[3:5], words[1:3], false},
		{words[1:1], words[1:3], false},
		{words[1:3], separate[1:3], false},
	}
	for _, candidate := range []struct {
		name string
		run  func([][3]byte, [][3]byte) bool
		want int
	}{
		{"Go", slices.GoLeanOverlaps[[3]byte], 0},
		{"overlap uses capacity", capacityOverlap, 3},
		{"overlap includes the byte just beyond the interval", adjacentOverlap, 2},
	} {
		failures := 0
		for i, c := range cases {
			if got := candidate.run(c.a, c.b); got != c.want {
				failures++
				if candidate.want == 0 {
					t.Fatalf("Go fixture %d: got %t want %t", i, got, c.want)
				}
			}
		}
		if failures != candidate.want {
			t.Fatalf("%s disagreed on %d fixtures, want %d", candidate.name, failures, candidate.want)
		}
		t.Logf("%s disagreed on %d of %d fixtures", candidate.name, failures, len(cases))
	}
	var zero [4]struct{}
	if slices.GoLeanOverlaps(zero[:3], zero[1:]) {
		t.Fatal("zero-size elements have no overlapping storage bytes")
	}
	t.Log("nonempty shared zero-size ranges did not overlap")
}

func capacityOverlap(a, b [][3]byte) bool {
	return slices.GoLeanOverlaps(a[:cap(a)], b[:cap(b)])
}

func adjacentOverlap(a, b [][3]byte) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	x, y := uintptr(unsafe.Pointer(&a[0])), uintptr(unsafe.Pointer(&b[0]))
	return slices.GoLeanOverlaps(a, b) || x+uintptr(len(a))*3 == y || y+uintptr(len(b))*3 == x
}
