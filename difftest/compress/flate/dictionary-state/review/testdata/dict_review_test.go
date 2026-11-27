package flate

import (
	"runtime"
	"slices"
	"testing"
)

func TestGoLeanDictReview(t *testing.T) {
	backing := []byte{41, 42, 43, 44, 45, 46}
	d := dictDecoder{hist: backing[:4:6], rdPos: 1, wrPos: 2}
	write := d.writeSlice()
	if len(write) != 2 || cap(write) != 4 {
		t.Fatalf("writeSlice length %d and capacity %d", len(write), cap(write))
	}
	write[0] = 101
	write[:cap(write)][3] = 102
	if d.hist[2] != 101 || backing[5] != 102 {
		t.Fatal("writeSlice failed write-through for visible or spare cells")
	}
	backing[3] = 103
	if write[1] != 103 {
		t.Fatal("backing change did not reach writeSlice")
	}
	d.writeMark(2)
	func() {
		defer func() {
			if _, ok := recover().(runtime.Error); !ok {
				t.Fatal("full visible window did not raise a runtime error despite spare capacity")
			}
		}()
		d.writeByte(255)
	}()
	if d.rdPos != 1 || d.wrPos != 4 || d.full || len(d.hist) != 4 || cap(d.hist) != 6 ||
		&d.hist[0] != &backing[0] || !slices.Equal(backing, []byte{41, 42, 101, 103, 45, 102}) {
		t.Fatal("failed write changed state or backing")
	}
	flush := d.readFlush()
	if !slices.Equal(flush, []byte{42, 101, 103}) || cap(flush) != 5 ||
		d.rdPos != 0 || d.wrPos != 0 || !d.full || d.histSize() != 4 {
		t.Fatal("flush failed to return unread bytes and wrap both cursors")
	}
	flush[0] = 202
	if backing[1] != 202 {
		t.Fatal("flush copied unread bytes")
	}
	backing[2] = 203
	if flush[1] != 203 {
		t.Fatal("flush failed to observe backing change")
	}
}
