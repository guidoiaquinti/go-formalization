package binary_test

import (
	"encoding/binary"
	"fmt"
	"slices"
	"testing"
)

func TestGoLeanEnsureReview(t *testing.T) {
	var rejected [3]int
	for _, c := range []struct {
		count  int
		want   []byte
		reused bool
	}{
		{0, []byte{11, 22}, true},
		{1, []byte{11, 22, 33}, true},
		{2, []byte{11, 22, 33, 44}, true},
		{3, []byte{11, 22, 33, 44, 0}, false},
		{8, []byte{11, 22, 33, 44, 0, 0, 0, 0, 0, 0}, false},
	} {
		backing := []byte{91, 11, 22, 33, 44, 92}
		out, pos := binary.GoLeanEnsure(backing[1:3:5], c.count)
		if !slices.Equal(out, c.want) || !slices.Equal(pos, c.want[2:]) || (&out[0] == &backing[1]) != c.reused || len(pos) != c.count || len(out) != 2+c.count {
			t.Fatalf("%+v: out %v, suffix %v, original %v", c, out, pos, backing)
		}
		if !slices.Equal(backing, []byte{91, 11, 22, 33, 44, 92}) {
			t.Fatal("ensure mutated original")
		}
		zeroed := slices.Clone(c.want)
		clear(zeroed[2:])
		if !slices.Equal(out, zeroed) {
			rejected[0]++
		}
		if !c.reused && !slices.Equal(out, zeroed) {
			rejected[1]++
		}
		if len(pos) > 0 {
			pos[0] = 117
			if out[2] != 117 || (backing[3] == 117) != c.reused || backing[0] != 91 || backing[5] != 92 {
				t.Fatal("suffix does not write through to output with the proper original alias")
			}
			out[len(out)-1] = 118
			if pos[len(pos)-1] != 118 {
				t.Fatal("output does not write through to suffix")
			}
			rejected[2]++
		}
	}
	nilOut, nilPos := binary.GoLeanEnsure(nil, 0)
	emptyOut, emptyPos := binary.GoLeanEnsure(make([]byte, 0), 0)
	if nilOut != nil || nilPos != nil || emptyOut == nil || emptyPos == nil {
		t.Fatal("zero-length slice nilness")
	}
	func() {
		backing := []byte{7, 9, 8}
		defer func() {
			if got := fmt.Sprint(recover()); got != "cannot be negative" || !slices.Equal(backing, []byte{7, 9, 8}) {
				t.Fatalf("negative count: panic %s, backing %v", got, backing)
			}
		}()
		binary.GoLeanEnsure(backing[:2], -1)
	}()
	if rejected != [3]int{4, 2, 4} {
		t.Fatalf("unexpected mutant rejections %v", rejected)
	}
	t.Logf("mutant rejections %v", rejected)
}
