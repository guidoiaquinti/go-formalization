package math_test

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
	"unsafe"
)

func TestGoLeanNormalize(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_MATH_NORMALIZE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := bufio.NewReader(f)
	for {
		var word, want uint64
		var adjustment, exponent int
		n, err := fmt.Fscan(input, &word, &want, &adjustment, &exponent)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		var x float64
		binary.NativeEndian.PutUint64((*[8]byte)(unsafe.Pointer(&x))[:], word)
		y, e := math.GoLeanNormalize(x)
		got := binary.NativeEndian.Uint64((*[8]byte)(unsafe.Pointer(&y))[:])
		gotExponent := math.GoLeanIlogb(x)
		if got != want || e != adjustment || gotExponent != exponent {
			t.Fatalf("input %#016x: Go=(%#016x %d %d) want=(%#016x %d %d)",
				word, got, e, gotExponent, want, adjustment, exponent)
		}
	}
}
