import GoLean.Math.RoundBits

namespace GoLean.Math

open RoundBits

/-- The generic implementation biases the bit pattern by a half unit minus
    one bit, adding back the bit when the integer part is odd. Below one
    it rounds magnitudes strictly above one half to signed one.
    The explicit upper guard collapses Go's shifts of 52 bits or more,
    which make both the bias and the fraction-clearing mask zero.
    https://github.com/golang/go/blob/go1.27.0/src/math/floor.go#L128-L158
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L17 -/
def roundToEvenWord (word : BitVec 64) : Nat :=
  let bits := word.toNat
  let e := (bits >>> 52) &&& 2047
  if e ≥ 1023 then
    let shift := e - 1023
    if shift < 52 then
      let odd := (bits >>> (52 - shift)) &&& 1
      (bits + ((2 ^ 51 - 1 + odd) >>> shift)) &&&
        (2 ^ 64 - 1 - ((2 ^ 52 - 1) >>> shift))
    else bits
  else
    let sign := bits &&& 2 ^ 63
    if e == 1022 && (bits &&& (2 ^ 52 - 1)) != 0 then sign ||| (1023 * 2 ^ 52)
    else sign

/-- Adding one quantum rounds upward exactly above the midpoint or at a
    midpoint whose truncated result is odd. At exponent 1023 the odd
    biased exponent contributes the implicit leading one. At larger
    fractional exponents the exponent field and the implicit leading one
    both contribute even multiples of the quantum. -/
theorem roundToEvenWord_spec (word : BitVec 64) :
    roundToEvenWord word =
      let bits := word.toNat
      let e := bits / 2 ^ 52 % 2048
      if e < 1023 then
        bits / 2 ^ 63 * 2 ^ 63 +
          if e = 1022 ∧ bits % 2 ^ 52 ≠ 0 then 1023 * 2 ^ 52 else 0
      else if e < 1075 then
        let quantum := 2 ^ (1075 - e)
        let half := 2 ^ (1074 - e)
        let whole := bits / quantum
        let fraction := bits % quantum
        (whole + if half < fraction ∨ (fraction = half ∧ whole % 2 = 1)
          then 1 else 0) * quantum
      else bits := by
  unfold roundToEvenWord
  dsimp only
  rw [Nat.shiftRight_eq_div_pow, show (2047 : Nat) = 2 ^ 11 - 1 by decide,
    Nat.and_two_pow_sub_one_eq_mod]
  let bits := word.toNat
  let e := bits / 2 ^ 52 % 2048
  change (if e ≥ 1023 then
      if e - 1023 < 52 then
        let odd := (bits >>> (52 - (e - 1023))) &&& 1
        (bits + ((2 ^ 51 - 1 + odd) >>> (e - 1023))) &&&
          (2 ^ 64 - 1 - ((2 ^ 52 - 1) >>> (e - 1023)))
      else bits
    else if e == 1022 && (bits &&& (2 ^ 52 - 1)) != 0 then
      (bits &&& 2 ^ 63) ||| (1023 * 2 ^ 52)
    else bits &&& 2 ^ 63) =
      if e < 1023 then
        bits / 2 ^ 63 * 2 ^ 63 +
          (if e = 1022 ∧ bits % 2 ^ 52 ≠ 0 then 1023 * 2 ^ 52 else 0)
      else if e < 1075 then
        let quantum := 2 ^ (1075 - e)
        let half := 2 ^ (1074 - e)
        let whole := bits / quantum
        let fraction := bits % quantum
        (whole + if half < fraction ∨ (fraction = half ∧ whole % 2 = 1)
          then 1 else 0) * quantum
      else bits
  have bound : bits < 2 ^ 64 := word.isLt
  by_cases small : e < 1023
  · rw [ite_eq_right (show ¬ e ≥ 1023 by omega), ite_eq_left small,
      Nat.and_two_pow_sub_one_eq_mod,
      show (2 : Nat) ^ 63 = 2 ^ 64 - 2 ^ 63 by decide,
      clearLow bits 63 bound (by decide),
      show (2 : Nat) ^ 64 - 2 ^ 63 = 2 ^ 63 by decide]
    have sign : bits / 2 ^ 63 = 0 ∨ bits / 2 ^ 63 = 1 := by omega
    rcases sign with zero | one
    · by_cases eq : e = 1022 <;> by_cases residue : bits % 2 ^ 52 = 0 <;>
        simp [zero, eq, residue]
    · by_cases eq : e = 1022 <;> by_cases residue : bits % 2 ^ 52 = 0 <;>
        simp [one, eq, residue]
  · rw [ite_eq_left (show e ≥ 1023 by omega), ite_eq_right small]
    by_cases fractional : e < 1075
    · rw [ite_eq_left (show e - 1023 < 52 by omega), ite_eq_left fractional]
      have exponent : 51 - (e - 1023) = 1074 - e := by omega
      have exponent' : 52 - (e - 1023) = 1075 - e := by omega
      rw [shiftedMask 52 (e - 1023) (by omega), exponent', Nat.shiftRight_eq_div_pow,
        show (1 : Nat) = 2 ^ 1 - 1 by decide, Nat.and_two_pow_sub_one_eq_mod]
      let half := 2 ^ (1074 - e)
      let quantum := 2 ^ (1075 - e)
      let odd := bits / quantum % 2
      let bias := ((2 ^ 51 - 1 + odd) >>> (e - 1023))
      change (bits + bias) &&& (2 ^ 64 - 1 - (quantum - 1)) = _
      have oddCases : odd = 0 ∨ odd = 1 := by omega
      have biasForm : bias = if odd = 1 then half else half - 1 := by
        rcases oddCases with even | oddOne
        · simp only [bias, even, Nat.add_zero, Nat.zero_ne_one, ite_false]
          rw [shiftedMask 51 (e - 1023) (by omega), exponent]
        · simp only [bias, oddOne, ite_true]
          rw [show 2 ^ 51 - 1 + 1 = (2 : Nat) ^ 51 by decide,
            shiftedPow 51 (e - 1023) (by omega), exponent]
      have hhalf : 0 < half := Nat.two_pow_pos _
      have twice : quantum = half * 2 := by
        dsimp [quantum, half]
        rw [show 1075 - e = (1074 - e) + 1 by omega, Nat.pow_succ]
      have biasBound : bias ≤ half := by split at biasForm <;> omega
      have halfBound : half ≤ 2 ^ 51 := Nat.pow_le_pow_right (by decide) (by omega)
      have boundedQuantum : quantum ≤ 2 ^ 64 := Nat.pow_le_pow_right (by decide) (by omega)
      have roundedBound : bits + bias < 2 ^ 64 := by
        have gap : bits / 2 ^ 52 % 2048 < 1075 := fractional
        omega
      rw [show 2 ^ 64 - 1 - (quantum - 1) = 2 ^ 64 - quantum by omega,
        clearLow _ (1075 - e) roundedBound (by omega)]
      change ((bits + bias) / quantum) * quantum =
        (bits / quantum + if half < bits % quantum ∨ (bits % quantum = half ∧ odd = 1)
          then 1 else 0) * quantum
      have hquantum : 0 < quantum := by omega
      have biasBelow : bias < quantum := by omega
      rw [Nat.add_div hquantum, Nat.div_eq_of_lt biasBelow, Nat.mod_eq_of_lt biasBelow,
        Nat.add_zero]
      have condition : quantum ≤ bits % quantum + bias ↔
          half < bits % quantum ∨ (bits % quantum = half ∧ odd = 1) := by
        by_cases h : odd = 1
        · simp only [h, ite_true] at biasForm
          simp only [h, and_true]
          omega
        · simp only [h, ite_false] at biasForm
          simp only [h, and_false, or_false]
          omega
      simp only [condition]
    · rw [ite_eq_right (show ¬ e - 1023 < 52 by omega), ite_eq_right fractional]

end GoLean.Math
