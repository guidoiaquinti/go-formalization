import GoLean.Math.Log

namespace GoLean.Math

/-- Go 1.27.0 src/math/log10.go:22–37. Use Frexp, special-case
    the half significand, and otherwise reconstruct from Log(fraction),
    reciprocal ln(2), and the integer exponent. The arm64 compiler
    contracts only the final multiplication and addition. -/
def log2Word (arm64 : Bool) (word : BitVec 64) : Option Nat :=
  let (fraction, exponent) := frexp word.toNat
  if fraction = 0x3fe0000000000000 then
    some (roundSigned ((exponent - 1) * 2 ^ 1074) 1)
  else
    (logWord arm64 (BitVec.ofNat 64 fraction)).map fun logarithm =>
      if 0x7ff0000000000000 ≤ logarithm % 2 ^ 63 then logarithm
      else multiplyAddWords arm64 logarithm 0x3ff71547652b82fe
        (roundSigned (exponent * 2 ^ 1074) 1)

theorem log2Word_halfMantissa (arm64 : Bool) (word : BitVec 64) (exp : Int)
    (half : frexp word.toNat = (0x3fe0000000000000, exp)) :
    log2Word arm64 word = some (roundSigned ((exp - 1) * 2 ^ 1074) 1) := by
  simp only [log2Word, half, ite_true]

theorem log2Word_negative (arm64 : Bool) (word : BitVec 64)
    (fraction : 2 ^ 63 ≤ (frexp word.toNat).1)
    (nonzero : (frexp word.toNat).1 % 2 ^ 63 ≠ 0)
    (fits : (frexp word.toNat).1 < 2 ^ 64) :
    log2Word arm64 word = none := by
  unfold log2Word
  dsimp only
  rw [ite_eq_right (by omega), logWord_negative arm64 _ (by simpa [Nat.mod_eq_of_lt fits] using fraction)
    (by simpa [Nat.mod_eq_of_lt fits] using nonzero)]
  rfl

theorem log2Word_edges (arm64 : Bool) :
    log2Word arm64 (0 : BitVec 64) = some 0xfff0000000000000 ∧
    log2Word arm64 (BitVec.ofNat 64 0x8000000000000000) = some 0xfff0000000000000 ∧
    log2Word arm64 (BitVec.ofNat 64 0x7ff0000000000000) = some 0x7ff0000000000000 := by
  cases arm64 <;> decide

end GoLean.Math
