package main

import (
	"fmt"
	"log"
	"math/big"
	"math/bits"
	"runtime"
)

type impl func(uint64, uint64, uint64) uint64

func check(fn impl, w uint, hi, lo, y uint64) (ok bool, detail string) {
	defer func() {
		if p := recover(); p != nil {
			e, runtimePanic := p.(runtime.Error)
			ok = y == 0 && runtimePanic && e.Error() == "runtime error: integer divide by zero"
			detail = fmt.Sprint(p)
		}
	}()
	r := fn(hi, lo, y)
	if y == 0 || r >= y {
		return false, fmt.Sprintf("r=%d", r)
	}
	n := new(big.Int).SetUint64(hi)
	n.Lsh(n, w).Add(n, new(big.Int).SetUint64(lo))
	// An integral, nonnegative q exists exactly when n-r is nonnegative
	// and has no remainder on division by y.
	delta := new(big.Int).Sub(n, new(big.Int).SetUint64(r))
	rem := new(big.Int)
	quotient, _ := new(big.Int).QuoRem(delta, new(big.Int).SetUint64(y), rem)
	return quotient.Sign() >= 0 && rem.Sign() == 0, fmt.Sprintf("r=%d", r)
}
func main() {
	for _, w := range []uint{32, 64} {
		max := uint64(1)<<w - 1
		top := uint64(1) << (w - 1)
		mid := uint64(1) << (w / 2)
		words := []uint64{0, 1, 2, 3, mid - 1, mid, mid + 1, top - 1, top, top + 1, max - 1, max}
		real := func(h, l, y uint64) uint64 { return uint64(bits.Rem32(uint32(h), uint32(l), uint32(y))) }
		rawDiv := func(h, l, y uint64) uint64 { _, r := bits.Div32(uint32(h), uint32(l), uint32(y)); return uint64(r) }
		if w == 64 {
			real = bits.Rem64
			rawDiv = func(h, l, y uint64) uint64 { _, r := bits.Div64(h, l, y); return r }
		}
		if w == bits.UintSize {
			real = func(h, l, y uint64) uint64 { return uint64(bits.Rem(uint(h), uint(l), uint(y))) }
		}
		fns := []impl{real, rawDiv, func(h, l, y uint64) uint64 { return l % y }, func(h, l, y uint64) uint64 {
			if y == 0 {
				return 0
			}
			return real(h, l, y)
		}}
		want := [...]struct{ panic, value int }{{0, 0}, {0, 792}, {0, 756}, {144, 0}}
		for i, f := range fns {
			var num, zeros, ordinary, badPanic, badValue int
			var witness string
			for _, h := range words {
				for _, l := range words {
					for _, y := range words {
						ok, detail := check(f, w, h, l, y)
						num++
						if y == 0 {
							zeros++
							if !ok {
								badPanic++
							}
						} else {
							ordinary++
							if !ok {
								badValue++
							}
						}
						if !ok && witness == "" {
							witness = fmt.Sprintf("(%d,%d,%d) %s", h, l, y, detail)
						}
					}
				}
			}
			fmt.Printf("%s uint=%d width=%d fn=%d total=%d zero=%d positive=%d zeroFail=%d positiveFail=%d first=%s\n", runtime.Version(), bits.UintSize, w, i, num, zeros, ordinary, badPanic, badValue, witness)
			if badPanic != want[i].panic || badValue != want[i].value {
				log.Fatalf("width %d control %d: expected %d panic and %d value failures", w, i, want[i].panic, want[i].value)
			}
		}
	}
}
