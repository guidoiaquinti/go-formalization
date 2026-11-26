import Init.Data.BitVec.Decidable

namespace GoLean.Math.Bits

/-- Go generates each entry of `pop8tab` by clearing the lowest set bit with
    `x &= x - 1` until the byte is zero. A byte has at most eight set bits.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/make_tables.go#L69-L75 -/
private def countSetBits : Nat → BitVec 8 → Nat
  | 0, _ => 0
  | remaining + 1, x =>
    if x == 0 then 0 else 1 + countSetBits remaining (x &&& (x - 1))

/-- Go returns `int(pop8tab[x])` for every uint8, with no preconditions.
    Nat represents the result exactly, from 0 for 0x00 to 8 for 0xff.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L125 -/
def onesCount8 (x : BitVec 8) : Nat := countSetBits 8 x

/-- The return value counts exactly the set bit positions from 0 through 7. -/
theorem onesCount8_spec (x : BitVec 8) :
    onesCount8 x = (List.range 8).countP x.getLsbD := by
  revert x
  decide

end GoLean.Math.Bits
