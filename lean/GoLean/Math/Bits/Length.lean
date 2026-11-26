import Init.Data.BitVec.Lemmas
import Init.Data.BitVec.Decidable

namespace GoLean.Math.Bits

/-- A length is minimal precisely when it fits in k bits exactly for k ≥ length.
    This includes zero, which fits in zero bits. -/
def IsBitLength (x n : Nat) : Prop := ∀ k, n ≤ k ↔ x < 2 ^ k

/-- Go's len8tab generator repeatedly shifts right until zero.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/make_tables.go#L85-L91 -/
private def byteLength : Nat → BitVec 8 → Nat
  | 0, _ => 0
  | steps + 1, x => if x == 0 then 0 else 1 + byteLength steps (x >>> 1)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L310 -/
def len8 (x : BitVec 8) : Nat := byteLength 8 x

theorem len8_spec (x : BitVec 8) : IsBitLength x.toNat (len8 x) := by
  have table : ∀ b : BitVec 8, len8 b = if b.toNat = 0 then 0 else b.toNat.log2 + 1 := by decide
  intro k
  rw [table]
  split
  · rename_i hz
    simp [hz, Nat.two_pow_pos]
  · exact Nat.log2_lt ‹_›

/-- Each wider implementation skips the low half when the high half is nonzero. -/
private def extendLength (half : Nat) (small : BitVec half → Nat)
    (x : BitVec (half + half)) : Nat :=
  if x.toNat < 2 ^ half then small (BitVec.ofNat half x.toNat)
  else half + small (BitVec.ofNat half (x.toNat / 2 ^ half))

private theorem extendLength_spec (small : BitVec w → Nat)
    (correct : ∀ b, IsBitLength b.toNat (small b)) (x : BitVec (w + w)) :
    IsBitLength x.toNat (extendLength w small x) := by
  intro k
  unfold extendLength
  split
  · rename_i h
    simpa [Nat.mod_eq_of_lt h] using correct (BitVec.ofNat w x.toNat) k
  · rename_i h
    have bound := x.isLt
    rw [Nat.pow_add] at bound
    have hi := (Nat.div_lt_iff_lt_mul (Nat.two_pow_pos w)).2 bound
    by_cases fits : w ≤ k
    · obtain ⟨d, rfl⟩ := Nat.le.dest fits
      rw [Nat.add_le_add_iff_left, correct, BitVec.toNat_ofNat, Nat.mod_eq_of_lt hi,
        Nat.div_lt_iff_lt_mul (Nat.two_pow_pos w), Nat.pow_add, Nat.mul_comm]
    · have := Nat.pow_le_pow_right (Nat.zero_lt_two) (Nat.le_of_lt (Nat.lt_of_not_ge fits))
      omega

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L315 -/
def len16 (x : BitVec 16) : Nat := extendLength 8 len8 x

theorem len16_spec (x : BitVec 16) : IsBitLength x.toNat (len16 x) :=
  extendLength_spec len8 len8_spec x

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L324 -/
def len32 (x : BitVec 32) : Nat := extendLength 16 len16 x

theorem len32_spec (x : BitVec 32) : IsBitLength x.toNat (len32 x) :=
  extendLength_spec len16 len16_spec x

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L337 -/
def len64 (x : BitVec 64) : Nat := extendLength 32 len32 x

theorem len64_spec (x : BitVec 64) : IsBitLength x.toNat (len64 x) :=
  extendLength_spec len32 len32_spec x

/-- Len on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L302 -/
def len := len64

theorem len_spec (x : BitVec 64) : IsBitLength x.toNat (len x) := len64_spec x

/-- Removing the leading zeros leaves the minimal unsigned representation. -/
def IsLeadingZeros (x : BitVec w) (n : Nat) : Prop :=
  n ≤ w ∧ IsBitLength x.toNat (w - n)

private theorem leading_spec (x : BitVec w) (n : Nat) (h : IsBitLength x.toNat n) :
    IsLeadingZeros x (w - n) := by
  have hn := (h w).2 x.isLt
  exact ⟨Nat.sub_le .., by rwa [show w - (w - n) = n by omega]⟩

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L28 -/
def leadingZeros8 (x : BitVec 8) : Nat := 8 - len8 x

theorem leadingZeros8_spec (x : BitVec 8) : IsLeadingZeros x (leadingZeros8 x) :=
  leading_spec x _ (len8_spec x)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L31 -/
def leadingZeros16 (x : BitVec 16) : Nat := 16 - len16 x

theorem leadingZeros16_spec (x : BitVec 16) : IsLeadingZeros x (leadingZeros16 x) :=
  leading_spec x _ (len16_spec x)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L34 -/
def leadingZeros32 (x : BitVec 32) : Nat := 32 - len32 x

theorem leadingZeros32_spec (x : BitVec 32) : IsLeadingZeros x (leadingZeros32 x) :=
  leading_spec x _ (len32_spec x)

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L37 -/
def leadingZeros64 (x : BitVec 64) : Nat := 64 - len64 x

theorem leadingZeros64_spec (x : BitVec 64) : IsLeadingZeros x (leadingZeros64 x) :=
  leading_spec x _ (len64_spec x)

/-- LeadingZeros on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L25 -/
def leadingZeros (x : BitVec 64) : Nat := 64 - len x

theorem leadingZeros_spec (x : BitVec 64) : IsLeadingZeros x (leadingZeros x) :=
  leading_spec x _ (len_spec x)

end GoLean.Math.Bits
