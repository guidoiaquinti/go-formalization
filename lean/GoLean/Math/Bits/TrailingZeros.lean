import Init.Data.BitVec.Bitblast
import Init.Data.BitVec.Decidable

namespace GoLean.Math.Bits

/-- Zero returns the full width. Otherwise the reported bit is set and all
    less significant bits are clear. -/
def IsTrailingZeros (x : BitVec w) (n : Nat) : Prop :=
  if x = 0 then n = w else x.getLsbD n = true ∧ ∀ j < n, x.getLsbD j = false

/-- In two's complement the bits above the first one flip when x is negated. -/
private theorem isolatedBit (x : BitVec w) (hx : x ≠ 0) :
    x &&& (-x) = BitVec.twoPow w x.ctz.toNat := by
  let k := x.ctz.toNat
  have setBit : x.getLsbD k = true := BitVec.getLsbD_true_ctz_of_ne_zero hx
  have kw : k < w := by
    apply Nat.lt_of_not_ge
    intro h
    rw [BitVec.getLsbD_of_ge x k (by omega)] at setBit
    contradiction
  apply BitVec.eq_of_getLsbD_eq
  intro i hi
  rw [BitVec.getLsbD_and, BitVec.getLsbD_neg, BitVec.getLsbD_twoPow]
  change _ = (decide (k < w) && decide (k = i))
  by_cases lower : i < k
  · simp [BitVec.getLsbD_false_of_lt_ctz lower, show k ≠ i by omega]
  · by_cases eq : k = i
    · subst i
      have none : ¬ ∃ j < k, x.getLsbD j = true := by
        rintro ⟨j, hj, h⟩
        rw [BitVec.getLsbD_false_of_lt_ctz hj] at h
        contradiction
      simp [setBit, none, kw]
    · have someBit : ∃ j < i, x.getLsbD j = true := ⟨k, by omega, setBit⟩
      simp [hi, someBit, eq]

private def decode32 (bit : BitVec 32) : Nat :=
  let table : Array Nat := #[
    0, 1, 28, 2, 29, 14, 24, 3, 30, 22, 20, 15, 25, 17, 4, 8,
    31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9]
  table[((bit * 0x077cb531) >>> 27).toNat]!

private def decode64 (bit : BitVec 64) : Nat :=
  let table : Array Nat := #[
    0, 1, 56, 2, 57, 49, 28, 3, 61, 58, 42, 50, 38, 29, 17, 4,
    62, 47, 59, 36, 45, 43, 51, 22, 53, 39, 33, 30, 24, 18, 12, 5,
    63, 55, 48, 27, 60, 41, 37, 16, 46, 35, 44, 21, 52, 32, 23, 11,
    54, 26, 40, 15, 34, 20, 31, 10, 25, 14, 19, 9, 13, 8, 7, 6]
  table[((bit * 0x03f79d71b4ca8b09) >>> 58).toNat]!

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L72 -/
def trailingZeros16 (x : BitVec 16) : Nat :=
  if x = 0 then 16 else decode32 ((x &&& (-x)).setWidth 32)

theorem trailingZeros16_spec (x : BitVec 16) : IsTrailingZeros x (trailingZeros16 x) := by
  by_cases hx : x = 0
  · simp [IsTrailingZeros, trailingZeros16, hx]
  · have setBit := BitVec.getLsbD_true_ctz_of_ne_zero hx
    have bound : x.ctz.toNat < 16 := by
      apply Nat.lt_of_not_ge
      intro h
      rw [BitVec.getLsbD_of_ge x _ (by omega)] at setBit
      contradiction
    have table : ∀ j : Fin 16, decode32 ((BitVec.twoPow 16 j.val).setWidth 32) = j.val := by decide
    simp only [IsTrailingZeros, trailingZeros16, hx, ↓reduceIte, isolatedBit x hx]
    rw [table ⟨_, bound⟩]
    exact ⟨setBit, fun _ => BitVec.getLsbD_false_of_lt_ctz⟩

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L81 -/
def trailingZeros32 (x : BitVec 32) : Nat :=
  if x = 0 then 32 else decode32 (x &&& (-x))

theorem trailingZeros32_spec (x : BitVec 32) : IsTrailingZeros x (trailingZeros32 x) := by
  by_cases hx : x = 0
  · simp [IsTrailingZeros, trailingZeros32, hx]
  · have setBit := BitVec.getLsbD_true_ctz_of_ne_zero hx
    have bound : x.ctz.toNat < 32 := by
      apply Nat.lt_of_not_ge
      intro h
      rw [BitVec.getLsbD_of_ge x _ (by omega)] at setBit
      contradiction
    have table : ∀ j : Fin 32, decode32 (BitVec.twoPow 32 j.val) = j.val := by decide
    simp only [IsTrailingZeros, trailingZeros32, hx, ↓reduceIte, isolatedBit x hx]
    rw [table ⟨_, bound⟩]
    exact ⟨setBit, fun _ => BitVec.getLsbD_false_of_lt_ctz⟩

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L90 -/
def trailingZeros64 (x : BitVec 64) : Nat :=
  if x = 0 then 64 else decode64 (x &&& (-x))

theorem trailingZeros64_spec (x : BitVec 64) : IsTrailingZeros x (trailingZeros64 x) := by
  by_cases hx : x = 0
  · simp [IsTrailingZeros, trailingZeros64, hx]
  · have setBit := BitVec.getLsbD_true_ctz_of_ne_zero hx
    have bound : x.ctz.toNat < 64 := by
      apply Nat.lt_of_not_ge
      intro h
      rw [BitVec.getLsbD_of_ge x _ (by omega)] at setBit
      contradiction
    have table : ∀ j : Fin 64, decode64 (BitVec.twoPow 64 j.val) = j.val := by decide
    simp only [IsTrailingZeros, trailingZeros64, hx, ↓reduceIte, isolatedBit x hx]
    rw [table ⟨_, bound⟩]
    exact ⟨setBit, fun _ => BitVec.getLsbD_false_of_lt_ctz⟩

/-- TrailingZeros on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L59 -/
def trailingZeros := trailingZeros64

theorem trailingZeros_spec (x : BitVec 64) : IsTrailingZeros x (trailingZeros x) :=
  trailingZeros64_spec x

end GoLean.Math.Bits
