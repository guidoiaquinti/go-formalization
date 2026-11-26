import Init.Data.BitVec.Decidable

namespace GoLean.Math.Bits

/-- Compute an entry of Go's `ntz8tab` using the table generator's low-bit test
    and unsigned right shift. The remaining width bounds the scan on zero.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/make_tables.go#L61-L67 -/
private def countLowZeros : Nat → BitVec 8 → Nat
  | 0, _ => 0
  | remaining + 1, x =>
    if x &&& 1 == 0 then 1 + countLowZeros remaining (x >>> 1) else 0

/-- Go 1.27.0 returns `int(ntz8tab[x])`. Every uint8 is valid. The Go int
    result fits in 0 through 8, so Nat represents its value exactly.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L67 -/
def trailingZeros8 (x : BitVec 8) : Nat := countLowZeros 8 x

theorem trailingZeros8_zero : trailingZeros8 0 = 8 := by decide

/-- For a nonzero byte, removing exactly the reported power of two leaves an
    odd number. Equivalently, the lowest n+1 bits have value 2^n. -/
theorem trailingZeros8_spec (x : BitVec 8) (hx : x ≠ 0) :
    x.toNat % 2 ^ (trailingZeros8 x + 1) = 2 ^ trailingZeros8 x := by
  revert x
  decide

end GoLean.Math.Bits
