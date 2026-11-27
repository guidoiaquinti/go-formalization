import GoLean.Math.Trunc

namespace GoLean.Math

/-- The finite negative branch requires a zero fractional part from Modf.
    Its exact finite subtraction is zero precisely when Trunc keeps the
    word. Negative infinity also reaches Modf, whose NaN fraction fails
    the test. Negative zero and NaNs fail the original comparison.
    https://github.com/golang/go/blob/go1.27.0/src/math/gamma.go#L216-L222
    https://github.com/golang/go/blob/go1.27.0/src/math/modf.go#L14-L18 -/
def isNegInt (word : BitVec 64) : Bool :=
  let bits := word.toNat
  let magnitude := bits % 2 ^ 63
  if bits ≥ 2 ^ 63 ∧ magnitude ≠ 0 ∧ magnitude < 2047 * 2 ^ 52 then
    trunc word == bits
  else false

/-- A strictly negative integer has sign one, a finite exponent from
    1023 onward, and no bits below the encoded unit position. -/
theorem isNegInt_spec (word : BitVec 64) :
    isNegInt word = true ↔
      let bits := word.toNat
      let magnitude := bits % 2 ^ 63
      let e := magnitude / 2 ^ 52
      bits / 2 ^ 63 = 1 ∧ 1023 ≤ e ∧ e < 2047 ∧
        magnitude % 2 ^ (1075 - e) = 0 := by
  unfold isNegInt
  dsimp only
  let bits := word.toNat
  let magnitude := bits % 2 ^ 63
  let e := magnitude / 2 ^ 52
  change (if 2 ^ 63 ≤ bits ∧ magnitude ≠ 0 ∧ magnitude < 2047 * 2 ^ 52
      then trunc word == bits else false) = true ↔
      bits / 2 ^ 63 = 1 ∧ 1023 ≤ e ∧ e < 2047 ∧ magnitude % 2 ^ (1075 - e) = 0
  have bound : bits < 2 ^ 64 := word.isLt
  by_cases negativeFinite : 2 ^ 63 ≤ bits ∧ magnitude ≠ 0 ∧ magnitude < 2047 * 2 ^ 52
  · rw [ite_eq_left negativeFinite, beq_iff_eq]
    have sign : bits / 2 ^ 63 = 1 := by omega
    have finite : e < 2047 := by omega
    have exponent : bits / 2 ^ 52 % 2048 = e := by omega
    simp only [sign, finite, true_and]
    rw [trunc_spec]
    dsimp only
    rw [exponent]
    by_cases small : e < 1023
    · rw [ite_eq_left small]
      change bits - magnitude = bits ↔ _
      omega
    · rw [ite_eq_right small]
      simp only [show 1023 ≤ e by omega, true_and]
      by_cases fractional : e < 1075
      · rw [ite_eq_left fractional]
        let unit := 2 ^ (1075 - e)
        change bits / unit * unit = bits ↔ magnitude % unit = 0
        have remainder : magnitude % unit = bits % unit :=
          Nat.mod_mod_of_dvd bits (Nat.pow_dvd_pow 2 (by omega : 1075 - e ≤ 63))
        rw [remainder]
        have split := Nat.mod_add_div' bits unit
        omega
      · rw [ite_eq_right fractional]
        simp [bits, show 1075 - e = 0 by omega, Nat.mod_one]
  · rw [ite_eq_right negativeFinite]
    simp only [Bool.false_eq_true, false_iff]
    intro h
    apply negativeFinite
    have hs := h.1
    have lower := h.2.1
    have upper := h.2.2.1
    exact ⟨by omega, by omega, by omega⟩

end GoLean.Math
