import GoLean.Math.RoundBits

namespace GoLean.Math

open RoundBits

/-- The generic implementation adds a shifted half unit before clearing
    fractional bits. The exponent guards keep the unsigned subtraction exact.
    https://github.com/golang/go/blob/go1.27.0/src/math/floor.go#L90-L119
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L17 -/
def roundWord (word : BitVec 64) : Nat :=
  let bits := word.toNat
  let e := (bits >>> 52) &&& 2047
  if e < 1023 then
    let sign := bits &&& 2 ^ 63
    if e == 1022 then sign ||| (1023 * 2 ^ 52) else sign
  else if e < 1075 then
    let shift := e - 1023
    (bits + (2 ^ 51 >>> shift)) &&& (2 ^ 64 - 1 - ((2 ^ 52 - 1) >>> shift))
  else bits

/-- For magnitudes from one up to 2^52, the result advances one integer
    quantum exactly when the fractional residue reaches half a quantum.
    Below one the result is signed zero or signed one. Integral, infinite,
    and NaN encodings at larger exponents pass through unchanged. -/
theorem roundWord_spec (word : BitVec 64) :
    roundWord word =
      let bits := word.toNat
      let e := bits / 2 ^ 52 % 2048
      if e < 1023 then
        bits / 2 ^ 63 * 2 ^ 63 + if e = 1022 then 1023 * 2 ^ 52 else 0
      else if e < 1075 then
        let quantum := 2 ^ (1075 - e)
        let half := 2 ^ (1074 - e)
        (bits / quantum + if half ≤ bits % quantum then 1 else 0) * quantum
      else bits := by
  unfold roundWord
  dsimp only
  rw [Nat.shiftRight_eq_div_pow, show (2047 : Nat) = 2 ^ 11 - 1 by decide,
    Nat.and_two_pow_sub_one_eq_mod]
  let bits := word.toNat
  let e := bits / 2 ^ 52 % 2048
  change (if e < 1023 then
      if e == 1022 then (bits &&& 2 ^ 63) ||| (1023 * 2 ^ 52) else bits &&& 2 ^ 63
    else if e < 1075 then
      (bits + (2 ^ 51 >>> (e - 1023))) &&&
        (2 ^ 64 - 1 - ((2 ^ 52 - 1) >>> (e - 1023)))
    else bits) = _
  have bound : bits < 2 ^ 64 := word.isLt
  by_cases small : e < 1023
  · rw [ite_eq_left small, ite_eq_left small]
    rw [show 2 ^ 63 = (2 : Nat) ^ 64 - 2 ^ 63 by decide,
      clearLow bits 63 bound (by decide)]
    rw [show (2 : Nat) ^ 64 - 2 ^ 63 = 2 ^ 63 by decide]
    change (if e == 1022 then (bits / 2 ^ 63 * 2 ^ 63) ||| (1023 * 2 ^ 52)
      else bits / 2 ^ 63 * 2 ^ 63) =
      bits / 2 ^ 63 * 2 ^ 63 + (if e = 1022 then 1023 * 2 ^ 52 else 0)
    have sign : bits / 2 ^ 63 = 0 ∨ bits / 2 ^ 63 = 1 := by omega
    rcases sign with zero | one
    · by_cases eq : e = 1022 <;> simp [zero, eq]
    · by_cases eq : e = 1022 <;> simp [one, eq]
  · rw [ite_eq_right small, ite_eq_right small]
    by_cases fractional : e < 1075
    · rw [ite_eq_left fractional, ite_eq_left fractional]
      have exponent : 51 - (e - 1023) = 1074 - e := by omega
      have exponent' : 52 - (e - 1023) = 1075 - e := by omega
      rw [shiftedPow 51 (e - 1023) (by omega), shiftedMask 52 (e - 1023) (by omega),
        exponent, exponent']
      let half := 2 ^ (1074 - e)
      let quantum := 2 ^ (1075 - e)
      have hhalf : 0 < half := Nat.two_pow_pos _
      have twice : quantum = half * 2 := by
        dsimp [quantum, half]
        rw [show 1075 - e = (1074 - e) + 1 by omega, Nat.pow_succ]
      have boundedQuantum : quantum ≤ 2 ^ 64 := Nat.pow_le_pow_right (by decide) (by omega)
      have halfBound : half ≤ 2 ^ 51 := Nat.pow_le_pow_right (by decide) (by omega)
      have roundedBound : bits + half < 2 ^ 64 := by
        have gap : bits / 2 ^ 52 % 2048 < 1075 := fractional
        omega
      change (bits + half) &&& (2 ^ 64 - 1 - (quantum - 1)) = _
      rw [show 2 ^ 64 - 1 - (quantum - 1) = 2 ^ 64 - quantum by omega,
        clearLow _ (1075 - e) roundedBound (by omega)]
      change ((bits + half) / quantum) * quantum =
        (bits / quantum + if half ≤ bits % quantum then 1 else 0) * quantum
      have hquantum : 0 < quantum := by omega
      have hl : half < quantum := by omega
      rw [Nat.add_div hquantum, Nat.div_eq_of_lt hl, Nat.mod_eq_of_lt hl, Nat.add_zero]
      simp only [show (quantum ≤ bits % quantum + half) ↔ half ≤ bits % quantum by omega]
    · rw [ite_eq_right fractional, ite_eq_right fractional]

end GoLean.Math
