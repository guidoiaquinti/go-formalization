import GoLean.Math.Trunc

namespace GoLean.Math

/-- The Go helper rejects magnitudes at or beyond 2^53, then requires
    an exact integer and tests the low bit of its integer conversion.
    Modf gets the integer part from Trunc, so its finite fractional part
    is zero exactly when truncation keeps the original word. Magnitudes
    below one are never odd integers, and NaNs fail the zero-fraction test.
    https://github.com/golang/go/blob/go1.27.0/src/math/pow.go#L7-L19
    https://github.com/golang/go/blob/go1.27.0/src/math/modf.go#L14-L18 -/
def isOddInt (word : BitVec 64) : Bool :=
  let bits := word.toNat
  let magnitude := bits % 2 ^ 63
  if magnitude ≥ 1076 * 2 ^ 52 ∨ magnitude < 1023 * 2 ^ 52 then false
  else
    let e := (magnitude >>> 52) &&& 2047
    let significand := (bits &&& (2 ^ 52 - 1)) ||| 2 ^ 52
    trunc word == bits && (((significand >>> (1075 - e)) &&& 1) == 1)

private theorem significand_remainder (bits : Nat) (h : k ≤ 52) :
    (2 ^ 52 + bits % 2 ^ 52) % 2 ^ k = bits % 2 ^ k := by
  have divisible : 2 ^ k ∣ 2 ^ 52 := Nat.pow_dvd_pow 2 h
  rw [Nat.add_comm, Nat.add_mod, Nat.mod_eq_zero_of_dvd divisible, Nat.add_zero,
    Nat.mod_mod, Nat.mod_mod_of_dvd _ divisible]

/-- Decode the unsigned significand by dividing by its binary place value.
    A true result has zero fractional residue and an odd quotient. -/
theorem isOddInt_spec (word : BitVec 64) :
    isOddInt word = true ↔
      let magnitude := word.toNat % 2 ^ 63
      let e := magnitude / 2 ^ 52
      let coefficient := 2 ^ 52 + magnitude % 2 ^ 52
      let unit := 2 ^ (1075 - e)
      1023 ≤ e ∧ e ≤ 1075 ∧ coefficient % unit = 0 ∧ coefficient / unit % 2 = 1 := by
  unfold isOddInt
  dsimp only
  let bits := word.toNat
  let magnitude := bits % 2 ^ 63
  let e := magnitude / 2 ^ 52
  let f := magnitude % 2 ^ 52
  have field : (magnitude >>> 52) &&& 2047 = e := by
    rw [Nat.shiftRight_eq_div_pow, show (2047 : Nat) = 2 ^ 11 - 1 by decide,
      Nat.and_two_pow_sub_one_eq_mod]
    exact Nat.mod_eq_of_lt (by omega)
  have fraction : bits % 2 ^ 52 = f := by omega
  have fracBound : f < 2 ^ 52 := Nat.mod_lt _ (by decide)
  rw [field, Nat.and_two_pow_sub_one_eq_mod, fraction, Nat.or_two_pow_eq_add_of_lt fracBound,
    Nat.add_comm f, Nat.shiftRight_eq_div_pow, show (1 : Nat) = 2 ^ 1 - 1 by decide,
    Nat.and_two_pow_sub_one_eq_mod]
  change (if 1076 * 2 ^ 52 ≤ magnitude ∨ magnitude < 1023 * 2 ^ 52 then false
      else trunc word == bits && (((2 ^ 52 + f) / 2 ^ (1075 - e) % 2) == 1)) = true ↔
      1023 ≤ e ∧ e ≤ 1075 ∧ (2 ^ 52 + f) % 2 ^ (1075 - e) = 0 ∧
        (2 ^ 52 + f) / 2 ^ (1075 - e) % 2 = 1
  by_cases middle : 1023 ≤ e ∧ e ≤ 1075
  · rw [ite_eq_right (show ¬ (1076 * 2 ^ 52 ≤ magnitude ∨ magnitude < 1023 * 2 ^ 52) by omega)]
    simp only [Bool.and_eq_true, beq_iff_eq, middle.1, middle.2, true_and]
    have remainder := significand_remainder bits (k := 1075 - e) (by omega)
    rw [fraction] at remainder
    rw [remainder]
    have exponent : bits / 2 ^ 52 % 2048 = e := by omega
    rw [trunc_spec]
    dsimp only
    rw [exponent, ite_eq_right (by omega : ¬ e < 1023)]
    by_cases fractional : e < 1075
    · rw [ite_eq_left fractional]
      change (bits / 2 ^ (1075 - e) * 2 ^ (1075 - e) = bits ∧ _) ↔ _
      have split := Nat.mod_add_div' bits (2 ^ (1075 - e))
      have unchanged : bits / 2 ^ (1075 - e) * 2 ^ (1075 - e) = bits ↔
          bits % 2 ^ (1075 - e) = 0 := by omega
      simp only [unchanged]
    · have equal : e = 1075 := by omega
      simp [equal, bits, Nat.mod_one]
  · rw [ite_eq_left (show 1076 * 2 ^ 52 ≤ magnitude ∨ magnitude < 1023 * 2 ^ 52 by omega)]
    simp only [Bool.false_eq_true, false_iff]
    exact fun yes => middle ⟨yes.1, yes.2.1⟩

end GoLean.Math
