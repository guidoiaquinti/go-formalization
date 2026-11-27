package binary_test

import (
	"encoding/binary"
	"testing"
)

type reviewDefinedByte byte
type reviewDefinedBytes []byte
type reviewDefinedIntPointer *int32
type reviewDefinedInts []int32
type reviewByteAlias = uint8

func TestGoLeanIntSizeReview(t *testing.T) {
	var rejected [4]int
	for _, c := range []struct {
		input                 any
		want                  int
		nilMutant, kindMutant int
	}{
		{(*int32)(nil), 4, -1, 4},
		{(*float64)(nil), 8, -1, 8},
		{reviewDefinedIntPointer(nil), 0, 0, 4},
		{reviewDefinedInts{11, 22}, 0, 0, 8},
		{[]reviewDefinedByte{1, 2}, 0, 0, 2},
		{reviewDefinedBytes{11, 22}, 0, 0, 2},
		{[]int32{1, 2}, 8, 8, 8},
	} {
		size, raw := binary.GoLeanIntDataSize(c.input)
		if size != c.want || raw != nil {
			t.Fatalf("%T size=%d, raw=%v", c.input, size, raw)
		}
		if size != c.nilMutant {
			rejected[0]++
		}
		if size != c.kindMutant {
			rejected[1]++
		}
	}
	backing := []byte{91, 11, 22, 33, 92}
	for _, length := range []int{0, 2} {
		original := []reviewByteAlias(backing[1 : 1+length : 4])
		size, raw := binary.GoLeanIntDataSize(original)
		if size != length || len(raw) != length || cap(raw) != 3 || &raw[:cap(raw)][0] != &backing[1] {
			t.Fatal("alias/capacity changed")
		}
		if cap(raw) != len(raw) {
			rejected[2]++
		}
		raw[:cap(raw)][2] = 117
		if backing[3] != 117 || backing[0] != 91 || backing[4] != 92 {
			t.Fatal("write through or view boundaries")
		}
		rejected[3]++
	}
	size, raw := binary.GoLeanIntDataSize([]byte(nil))
	if size != 0 || raw != nil {
		t.Fatal("nil slice changed")
	}
	if rejected != [4]int{2, 4, 2, 2} {
		t.Fatalf("incorrect mutant rejection counts: %v", rejected)
	}
	t.Logf("mutant rejections: %v", rejected)
}
