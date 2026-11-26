import Init.Data.BitVec.Lemmas

namespace GoLean.Math.Bits

/-- Go distinguishes these two runtime.Error outcomes. -/
inductive DivisionPanic where
  | divideByZero
  | quotientOverflow
  deriving DecidableEq, Repr

/-- A double-width numerator with hi as its most significant word. -/
def dividend (hi lo : BitVec w) : Nat := hi.toNat * 2 ^ w + lo.toNat

/-- Model the observable outcome of the widened Div32 and the long division
    in Div64. The zero-divisor error takes precedence, including at hi = y = 0. -/
private def divideWord (hi lo y : BitVec w) : Except DivisionPanic (BitVec w × BitVec w) :=
  if y.toNat = 0 then .error .divideByZero
  else if y.toNat ≤ hi.toNat then .error .quotientOverflow
  else
    let n := dividend hi lo
    .ok (BitVec.ofNat w (n / y.toNat), BitVec.ofNat w (n % y.toNat))

/-- The two panic classes and the unique Euclidean quotient and remainder. -/
def IsDivision (hi lo y : BitVec w) (result : Except DivisionPanic (BitVec w × BitVec w)) : Prop :=
  match result with
  | .error .divideByZero => y.toNat = 0
  | .error .quotientOverflow => 0 < y.toNat ∧ y.toNat ≤ hi.toNat
  | .ok (q, r) => hi.toNat < y.toNat ∧ dividend hi lo = y.toNat * q.toNat + r.toNat ∧ r < y

private theorem divideWord_spec (hi lo y : BitVec w) : IsDivision hi lo y (divideWord hi lo y) := by
  unfold divideWord
  split
  · exact ‹_›
  · split
    · exact ⟨by omega, ‹_›⟩
    · have pos : 0 < y.toNat := by omega
      have top := Nat.mul_le_mul_right (2 ^ w) (show hi.toNat + 1 ≤ y.toNat by omega)
      rw [Nat.add_one_mul] at top
      have low := lo.isLt
      have bound : dividend hi lo < 2 ^ w * y.toNat := by
        rw [Nat.mul_comm]
        unfold dividend
        omega
      have hq := (Nat.div_lt_iff_lt_mul pos).2 bound
      have hr := Nat.mod_lt (dividend hi lo) pos
      have hrw := Nat.lt_trans hr y.isLt
      simp only [IsDivision, BitVec.toNat_ofNat, BitVec.lt_def]
      rw [Nat.mod_eq_of_lt hq, Nat.mod_eq_of_lt hrw]
      exact ⟨by omega, (Nat.div_add_mod ..).symm, hr⟩

/-- Remainder accepts an overflowing quotient. Only a zero divisor panics. -/
private def remainderWord (hi lo y : BitVec w) : Except DivisionPanic (BitVec w) :=
  if y.toNat = 0 then .error .divideByZero
  else .ok (BitVec.ofNat w (dividend hi lo % y.toNat))

def IsRemainder (hi lo y : BitVec w) (result : Except DivisionPanic (BitVec w)) : Prop :=
  match result with
  | .error .divideByZero => y.toNat = 0
  | .error .quotientOverflow => False
  | .ok r => r < y ∧ ∃ q, dividend hi lo = y.toNat * q + r.toNat

private theorem remainderWord_spec (hi lo y : BitVec w) :
    IsRemainder hi lo y (remainderWord hi lo y) := by
  unfold remainderWord
  split
  · exact ‹_›
  · have hr := Nat.mod_lt (dividend hi lo) (show 0 < y.toNat by omega)
    have hrw := Nat.lt_trans hr y.isLt
    simp only [IsRemainder, BitVec.lt_def, BitVec.toNat_ofNat]
    rw [Nat.mod_eq_of_lt hrw]
    exact ⟨hr, _, (Nat.div_add_mod ..).symm⟩

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L505 -/
def div32 (hi lo y : BitVec 32) := divideWord hi lo y

theorem div32_spec (hi lo y : BitVec 32) : IsDivision hi lo y (div32 hi lo y) :=
  divideWord_spec hi lo y

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L518 -/
def div64 (hi lo y : BitVec 64) := divideWord hi lo y

theorem div64_spec (hi lo y : BitVec 64) : IsDivision hi lo y (div64 hi lo y) :=
  divideWord_spec hi lo y

/-- Div on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L492 -/
def div (hi lo y : BitVec 64) := div64 hi lo y

theorem div_spec (hi lo y : BitVec 64) : IsDivision hi lo y (div hi lo y) :=
  divideWord_spec hi lo y

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L583 -/
def rem32 (hi lo y : BitVec 32) := remainderWord hi lo y

theorem rem32_spec (hi lo y : BitVec 32) : IsRemainder hi lo y (rem32 hi lo y) :=
  remainderWord_spec hi lo y

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L590 -/
def rem64 (hi lo y : BitVec 64) := remainderWord hi lo y

theorem rem64_spec (hi lo y : BitVec 64) : IsRemainder hi lo y (rem64 hi lo y) :=
  remainderWord_spec hi lo y

/-- Rem on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L573 -/
def rem (hi lo y : BitVec 64) := rem64 hi lo y

theorem rem_spec (hi lo y : BitVec 64) : IsRemainder hi lo y (rem hi lo y) :=
  remainderWord_spec hi lo y

end GoLean.Math.Bits
