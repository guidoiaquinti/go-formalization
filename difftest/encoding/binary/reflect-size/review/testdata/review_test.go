package binary_test

import (
	"encoding/binary"
	"reflect"
	"runtime"
	"testing"
)

type reviewPacked struct {
	A uint8
	B uint64
	_ uint16
}
type reviewUnsupported struct{ Address *uint16 }

func TestGoLeanReflectSizeReview(t *testing.T) {
	live := new(uint64)
	values := []struct {
		value              any
		byType, data, size int
	}{
		{reviewPacked{}, 11, 11, 11},
		{[0]*uint16{}, -1, -1, -1},
		{[]*uint16{}, -1, -1, -1},
		{[0]uint16{}, 0, 0, 0},
		{&live, -1, -1, -1},
		{&[]uint16{1, 2}, -1, -1, 4},
		{(*uint32)(nil), -1, -1, -1},
		{[]reviewPacked{{}, {}, {}}, -1, 33, 33},
		{reviewUnsupported{}, -1, -1, -1},
	}
	var rejects [4]int
	for i, c := range values {
		typ := reflect.TypeOf(c.value)
		value := reflect.ValueOf(c.value)
		byType := binary.GoLeanReflectTypeSize(typ)
		cold, warm := binary.GoLeanReflectValueSize(value)
		size := binary.Size(c.value)
		if byType != c.byType || cold != c.data || warm != cold || size != c.size {
			t.Fatalf("case %d: type=%d, cold/warm=%d/%d, Size=%d", i, byType, cold, warm, size)
		}
		switch i {
		case 0:
			if int(typ.Size()) != size {
				rejects[0]++
			}
		case 1, 2:
			if size != 0 {
				rejects[1]++
			}
		case 4:
			if size != 8 {
				rejects[2]++
			}
		case 5:
			if size != -1 {
				rejects[3]++
			}
		}
	}
	cold, warm := binary.GoLeanReflectValueSize(reflect.Value{})
	if cold != -1 || warm != -1 || binary.Size(nil) != -1 {
		t.Fatal("invalid reflection value")
	}
	func() {
		defer func() {
			value := recover()
			if err, ok := value.(runtime.Error); !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				t.Fatalf("sizeof(nil) panic = %v", value)
			}
		}()
		binary.GoLeanReflectTypeSize(nil)
	}()
	if rejects != [4]int{1, 2, 1, 1} {
		t.Fatalf("unexpected mutant rejections %v", rejects)
	}
	t.Logf("mutant rejections %v", rejects)
}
