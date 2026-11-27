import Init.Data.BitVec.Lemmas

namespace GoLean.Math

/-- The Go fallback first compares the magnitude against one, then extracts
    and unbiases the exponent. It clears fraction bits with an AND NOT mask.
    The magnitude guard ensures the unsigned exponent subtraction is exact.
    https://github.com/golang/go/blob/go1.27.0/src/math/floor.go#L55-L80
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L16 -/
def trunc (word : BitVec 64) : Nat :=
  let bits := word.toNat
  if bits % 2 ^ 63 < 1023 * 2 ^ 52 then
    bits / 2 ^ 63 * 2 ^ 63
  else
    let exponent := ((bits >>> 52) &&& 2047) - 1023
    if exponent < 52 then
      bits &&& (2 ^ 64 - 2 ^ (52 - exponent))
    else bits

private theorem clearLow (bits k : Nat) (hbits : bits < 2 ^ 64) (hk : k ≤ 64) :
    bits &&& (2 ^ 64 - 2 ^ k) = bits / 2 ^ k * 2 ^ k := by
  have power : 2 ^ (64 - k) * 2 ^ k = (2 : Nat) ^ 64 := by
    rw [← Nat.pow_add, Nat.sub_add_cancel hk]
  have mask : 2 ^ 64 - 2 ^ k = (2 ^ (64 - k) - 1) * 2 ^ k := by
    rw [Nat.sub_mul, Nat.one_mul, power]
  have highBound : bits / 2 ^ k < 2 ^ (64 - k) := by
    rw [Nat.div_lt_iff_lt_mul (Nat.two_pow_pos k), power]
    exact hbits
  have upper : (bits &&& (2 ^ 64 - 2 ^ k)) / 2 ^ k = bits / 2 ^ k := by
    rw [Nat.and_div_two_pow, mask, Nat.mul_div_cancel _ (Nat.two_pow_pos k),
      Nat.and_two_pow_sub_one_of_lt_two_pow highBound]
  have lower : (bits &&& (2 ^ 64 - 2 ^ k)) % 2 ^ k = 0 := by
    rw [Nat.and_mod_two_pow, mask, Nat.mul_mod_left, Nat.and_zero]
  have combine := Nat.mod_add_div' (bits &&& (2 ^ 64 - 2 ^ k)) (2 ^ k)
  rw [upper, lower, Nat.zero_add] at combine
  exact combine.symm

/-- Words of magnitude below one become signed zero. In the fractional
    exponent range, division removes the bits below the unit position. The
    generic implementation preserves all other words, including NaN encodings. -/
theorem trunc_spec (word : BitVec 64) :
    trunc word =
      let bits := word.toNat
      let biased := bits / 2 ^ 52 % 2048
      if biased < 1023 then bits - bits % 2 ^ 63
      else if biased < 1075 then
        let quantum := 2 ^ (1075 - biased)
        bits / quantum * quantum
      else bits := by
  unfold trunc
  dsimp only
  rw [Nat.shiftRight_eq_div_pow, show (2047 : Nat) = 2 ^ 11 - 1 by decide,
    Nat.and_two_pow_sub_one_eq_mod]
  have branch : (word.toNat % 2 ^ 63 < 1023 * 2 ^ 52) ↔
      word.toNat / 2 ^ 52 % 2 ^ 11 < 1023 := by omega
  by_cases small : word.toNat / 2 ^ 52 % 2 ^ 11 < 1023
  · rw [ite_eq_left (branch.mpr small), ite_eq_left small]
    omega
  · rw [ite_eq_right (mt branch.mp small)]
    have lower : 1023 ≤ word.toNat / 2 ^ 52 % 2 ^ 11 := by omega
    by_cases fractional : word.toNat / 2 ^ 52 % 2 ^ 11 < 1075
    · rw [ite_eq_left (show word.toNat / 2 ^ 52 % 2 ^ 11 - 1023 < 52 by omega),
        ite_eq_right small, ite_eq_left fractional, clearLow _ _ word.isLt (by omega)]
      rw [show 52 - (word.toNat / 2 ^ 52 % 2 ^ 11 - 1023) =
        1075 - word.toNat / 2 ^ 52 % 2 ^ 11 by omega]
    · rw [ite_eq_right (show ¬ word.toNat / 2 ^ 52 % 2 ^ 11 - 1023 < 52 by omega),
        ite_eq_right small, ite_eq_right fractional]

end GoLean.Math
