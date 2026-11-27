import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- Go clears the binary64 sign bit using AND NOT and reinterprets the word.
    The model operates directly on the input and output representations.
    https://github.com/golang/go/blob/go1.27.0/src/math/abs.go#L7-L15
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L15
    https://github.com/golang/go/blob/go1.27.0/src/math/unsafe.go#L32-L41 -/
def abs (bits : BitVec 64) : BitVec 64 := bits &&& ~~~(0x8000000000000000 : BitVec 64)

/-- Taking the remainder keeps the exponent and fraction and clears the sign. -/
theorem abs_spec (bits : BitVec 64) :
    (abs bits).toNat = bits.toNat % (2 ^ 63) := by
  rw [abs, BitVec.toNat_and]
  change bits.toNat &&& (2 ^ 63 - 1) = _
  exact Nat.and_two_pow_sub_one_eq_mod _ _

end GoLean.Math
