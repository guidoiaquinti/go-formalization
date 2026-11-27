import Init.Data.BitVec.Lemmas

namespace GoLean.Math

private def nanEncoding (x : BitVec 64) : Bool :=
  x.toNat / 2^52 % 2^11 == 2047 && x.toNat % 2^52 != 0

/-- Binary64 numeric equality on raw encodings. NaNs are unordered, signed zeros
    are equal, and all other encodings have unique numeric values. -/
private def binary64Equal (x y : BitVec 64) : Bool :=
  !(nanEncoding x || nanEncoding y) &&
    (x == y || (x.toNat % 2^63 == 0 && y.toNat % 2^63 == 0))

/-- Go 1.27.0 implements IsNaN by comparing the operand with itself for floating
    inequality. This model uses binary64 equality on encodings for that comparison.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L33-L40 -/
def isNaN (x : BitVec 64) : Bool :=
  !binary64Equal x x

/-- The two encoding ranges contain all and only NaNs, for either sign. -/
theorem isNaN_spec (x : BitVec 64) :
    isNaN x = true ↔
      (0x7ff0000000000000 < x.toNat ∧ x.toNat < 0x8000000000000000) ∨
        0xfff0000000000000 < x.toNat := by
  have bound : x.toNat < 2^64 := x.isLt
  simp [isNaN, binary64Equal, nanEncoding, Bool.and_eq_true]
  omega

end GoLean.Math
