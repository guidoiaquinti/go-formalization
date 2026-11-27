import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- The FMA helper returns the one-bit indicator of an all-zero input.
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L9-L16 -/
def fmaZero (word : BitVec 64) : BitVec 64 :=
  if word = 0 then 1 else 0

/-- The branchless formula in the Go source comments folds the low bit
    into the upper 63 bits. Subtracting one sets the top bit exactly when
    that folded word was zero. -/
theorem fmaZero_spec (word : BitVec 64) :
    fmaZero word = (((word >>> (1 : Nat)) ||| (word &&& 1)) - 1) >>> (63 : Nat) := by
  let folded := word.toNat / 2 ||| word.toNat % 2
  have wordBound := word.isLt
  have foldedBound : folded < 2 ^ 63 := by
    apply Nat.or_lt_two_pow <;> omega
  have foldedZero : folded = 0 ↔ word.toNat = 0 := by
    unfold folded
    rw [Nat.or_eq_zero_iff]
    omega
  have foldedWord : ((word >>> (1 : Nat)) ||| (word &&& 1)).toNat = folded := by
    change (word.toNat >>> 1) ||| (word.toNat &&& 1) = folded
    rw [Nat.shiftRight_eq_div_pow, show (1 : Nat) = 2 ^ 1 - 1 by decide,
      Nat.and_two_pow_sub_one_eq_mod]
  apply BitVec.eq_of_toNat_eq
  rw [BitVec.toNat_ushiftRight, BitVec.toNat_sub, foldedWord, Nat.shiftRight_eq_div_pow]
  change (fmaZero word).toNat = (2 ^ 64 - 1 + folded) % 2 ^ 64 / 2 ^ 63
  unfold fmaZero
  by_cases isZero : word = 0
  · rw [ite_eq_left isZero, foldedZero.mpr (by simp [isZero])]
    decide
  · rw [ite_eq_right isZero]
    have foldedNonzero : folded ≠ 0 := by
      intro z
      apply isZero
      exact BitVec.eq_of_toNat_eq (foldedZero.mp z)
    change 0 = (2 ^ 64 - 1 + folded) % 2 ^ 64 / 2 ^ 63
    omega

end GoLean.Math
