import GoLean.Math.Bits.Length

namespace GoLean.Math

/-- Go counts the high word first and only inspects the low word when the
    high word is zero. The sum fits in Go's int32 result.
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L65-L71
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L37-L38 -/
def lz (high low : BitVec 64) : Nat :=
  let count := Bits.leadingZeros64 high
  if count = 64 then count + Bits.leadingZeros64 low else count

/-- Removing the returned number of zeros leaves the minimum width able to
    represent the combined unsigned value. A zero value needs no bits. -/
theorem lz_spec (high low : BitVec 64) :
    lz high low ≤ 128 ∧
      Bits.IsBitLength (high.toNat * 2 ^ 64 + low.toNat) (128 - lz high low) := by
  have upper := Bits.leadingZeros64_spec high
  have lower := Bits.leadingZeros64_spec low
  have highBound := upper.1
  have zeroHigh : Bits.leadingZeros64 high = 64 ↔ high.toNat = 0 := by
    have zeroWidth := upper.2 0
    simp only [Nat.le_zero, Nat.pow_zero, Nat.lt_one_iff] at zeroWidth
    omega
  unfold lz
  dsimp only
  by_cases z : high.toNat = 0
  · rw [ite_eq_left (zeroHigh.mpr z), zeroHigh.mpr z]
    refine ⟨by have := lower.1; omega, ?_⟩
    intro k
    simpa [z, show 128 - (64 + Bits.leadingZeros64 low) =
      64 - Bits.leadingZeros64 low by omega] using lower.2 k
  · rw [ite_eq_right (mt zeroHigh.mp z)]
    refine ⟨by omega, ?_⟩
    intro k
    by_cases wide : 64 ≤ k
    · obtain ⟨width, rfl⟩ := Nat.le.dest wide
      have highSpec := upper.2 width
      have lowBound := low.isLt
      rw [Nat.pow_add, Nat.mul_comm (2 ^ 64)]
      omega
    · have powerBound := Nat.pow_le_pow_right Nat.zero_lt_two (show k ≤ 64 by omega)
      omega

end GoLean.Math
