package constant

import (
	"math"
	"math/big"
	"testing"
)

func TestGoLeanSignReview(t *testing.T) {
	minusZero := floatVal{new(big.Float).SetFloat64(math.Copysign(0, -1))}
	if !minusZero.val.Signbit() || minusZero.Kind() != Float || Sign(minusZero) != 0 {
		t.Fatal("negative floating zero should have Float kind and sign zero")
	}
	if got := Sign(complexVal{ratVal{new(big.Rat).SetFrac64(-1, 9)}, int64Val(7)}); got != -1 {
		t.Fatalf("mixed component signs must combine to -1, got %d", got)
	}
	assertPanic := func(value Value, want string) {
		t.Helper()
		defer func() {
			if got := recover(); got != want {
				t.Fatalf("got panic %v, expected %q", got, want)
			}
		}()
		_ = Sign(value)
	}
	assertPanic(complexVal{MakeBool(true), MakeString("x")}, "true not numeric")
	assertPanic(complexVal{int64Val(1), MakeString("x")}, `"x" not numeric`)
}
