import Init.Data.Int.DivMod.Lemmas
import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Math

private def ldexpExponent (word : Nat) : Nat := word / 2 ^ 52 % 2048
private def ldexpFraction (word : Nat) : Nat := word % 2 ^ 52
private def ldexpSign (word : Nat) : Nat := word / 2 ^ 63 * 2 ^ 63

/-- Multiplication by 2^52 normalizes every nonzero binary64 subnormal exactly.
    The original leading bit moves to position 52, outside the payload field.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L59-L67 -/
private def normalizeLdexp (exponent payload : Nat) : Nat × Nat × Int :=
  if exponent = 0 then
    let k := payload.log2
    (k + 1, (payload <<< (52 - k)) % 2 ^ 52, -52)
  else (exponent, payload, 0)

/-- The final multiplication by 2^-53 in the denormal branch rounds the
    shifted 53-bit significand to nearest, with ties to even. -/
private def roundLdexp (significand shift : Nat) : Nat :=
  let q := significand >>> shift
  let divisor := 2 ^ shift
  let twiceRemainder := 2 * (significand % divisor)
  if twiceRemainder < divisor then q
  else if divisor < twiceRemainder then q + 1
  else q + q % 2

private def assembleLdexp (sign : Nat) (unbiased : Int) (payload : Nat) : Nat :=
  if unbiased < -1075 then sign
  else if 1023 < unbiased then sign + 2047 * 2 ^ 52
  else if unbiased < -1022 then
    let staged := unbiased + 53
    sign + roundLdexp (2 ^ 52 + payload) (53 - (staged + 1022)).toNat
  else
    sign + (unbiased + 1023).toNat * 2 ^ 52 + payload

/-- Source-shaped binary64 model of the Go fallback. `width` is the native
    Go int width. Exponent additions wrap, and the denormal branch stages the
    exponent before multiplying by 2^-53.
    https://github.com/golang/go/blob/go1.27.0/src/math/ldexp.go#L15-L51 -/
def ldexp (width word : Nat) (exp : Int) : Nat :=
  let exponent := ldexpExponent word
  let payload := ldexpFraction word
  if (exponent = 0 ∧ payload = 0) ∨ exponent = 2047 then word
  else
    let (normalExponent, normalPayload, adjustment) := normalizeLdexp exponent payload
    let adjusted := (exp + adjustment).bmod (2 ^ width)
    let unbiased := (adjusted + ((normalExponent : Int) - 1023)).bmod (2 ^ width)
    assembleLdexp (ldexpSign word) unbiased normalPayload

private theorem normalizedPayload_eq (f : Nat) (nonzero : f ≠ 0) (bound : f < 2 ^ 52) :
    2 ^ 52 + (f <<< (52 - f.log2)) % 2 ^ 52 = f <<< (52 - f.log2) := by
  let k := f.log2
  change 2 ^ 52 + (f <<< (52 - k)) % 2 ^ 52 = f <<< (52 - k)
  have kBound : k < 52 := (Nat.log2_lt nonzero).2 bound
  have power : 2 ^ k * 2 ^ (52 - k) = 2 ^ 52 := by
    rw [← Nat.pow_add, Nat.add_sub_of_le (Nat.le_of_lt kBound)]
  have low := Nat.mul_le_mul_right (2 ^ (52 - k)) (Nat.log2_self_le nonzero)
  have high := Nat.mul_lt_mul_of_pos_right (Nat.lt_log2_self (n := f))
    (Nat.two_pow_pos (52 - k))
  change 2 ^ k * _ ≤ _ at low
  change 2 ^ (k + 1) * _ > _ at high
  rw [power] at low
  rw [Nat.pow_succ, Nat.mul_assoc, Nat.mul_comm 2, ← Nat.mul_assoc, power] at high
  rw [Nat.shiftLeft_eq]
  omega

private theorem roundLdexp_spec (significand shift : Nat) :
    roundLdexp significand shift =
      let divisor := 2 ^ shift
      let q := significand / divisor
      let r := significand % divisor
      q + if divisor < 2 * r ∨ (divisor = 2 * r ∧ q % 2 = 1) then 1 else 0 := by
  simp only [roundLdexp, Nat.shiftRight_eq_div_pow]
  split
  · split <;> omega
  · split <;> (split <;> omega)

private def exactScaledWord (sign : Nat) (unbiased : Int) (coefficient : Nat) : Nat :=
  if unbiased < -1075 then sign
  else if 1023 < unbiased then sign + 2047 * 2 ^ 52
  else if unbiased < -1022 then
    let divisor := 2 ^ (-1022 - unbiased).toNat
    let q := coefficient / divisor
    let r := coefficient % divisor
    sign + q + if divisor < 2 * r ∨ (divisor = 2 * r ∧ q % 2 = 1) then 1 else 0
  else sign + (unbiased + 1023).toNat * 2 ^ 52 + (coefficient - 2 ^ 52)

private theorem assembleLdexp_eq (sign payload : Nat) (unbiased : Int) :
    assembleLdexp sign unbiased payload = exactScaledWord sign unbiased (2 ^ 52 + payload) := by
  by_cases underflow : unbiased < -1075
  · simp [assembleLdexp, exactScaledWord, underflow]
  by_cases overflow : 1023 < unbiased
  · simp [assembleLdexp, exactScaledWord, underflow, overflow]
  by_cases subnormal : unbiased < -1022
  · have shift : 53 - (unbiased + 53 + 1022) = -1022 - unbiased := by omega
    simp [assembleLdexp, exactScaledWord, underflow, overflow, subnormal, roundLdexp_spec,
      shift, Nat.add_assoc]
  · simp [assembleLdexp, exactScaledWord, underflow, overflow, subnormal]

/-- Exact returned word. Separate the leading exponent and integer
    coefficient of the original value, scale using one wrapped addition,
    and round a subnormal result to the nearest even integer coefficient. -/
theorem ldexp_spec (width word : Nat) (exp : Int) :
    ldexp width word exp =
      let magnitude := word % 2 ^ 63
      let sign := word - magnitude
      let exponent : Nat := magnitude / 2 ^ 52
      let payload : Nat := magnitude % 2 ^ 52
      if magnitude = 0 ∨ exponent = 2047 then word
      else
        let highest : Int := if exponent = 0 then (payload.log2 : Int) - 1074
          else Int.ofNat exponent - 1023
        let coefficient := if exponent = 0 then payload <<< (52 - payload.log2) else 2 ^ 52 + payload
        exactScaledWord sign ((exp + highest).bmod (2 ^ width)) coefficient := by
  have es : ldexpExponent word = word % 2 ^ 63 / 2 ^ 52 := by unfold ldexpExponent; omega
  have fs : ldexpFraction word = word % 2 ^ 63 % 2 ^ 52 := by unfold ldexpFraction; omega
  have ss : ldexpSign word = word - word % 2 ^ 63 := by unfold ldexpSign; omega
  have zs : (word % 2 ^ 63 / 2 ^ 52 = 0 ∧ word % 2 ^ 63 % 2 ^ 52 = 0) ↔
      word % 2 ^ 63 = 0 := by omega
  dsimp only
  by_cases special : word % 2 ^ 63 = 0 ∨ word % 2 ^ 63 / 2 ^ 52 = 2047
  · simp only [ldexp, es, fs, zs, special, ite_true]
  simp only [ldexp, es, fs, ss, zs, special, ite_false]
  by_cases eZero : word % 2 ^ 63 / 2 ^ 52 = 0
  · have fNonzero : word % 2 ^ 63 % 2 ^ 52 ≠ 0 := by omega
    have fBound : word % 2 ^ 63 % 2 ^ 52 < 2 ^ 52 := Nat.mod_lt _ (by decide)
    have shifted := normalizedPayload_eq _ fNonzero fBound
    simp only [normalizeLdexp, eZero, ite_true, Int.bmod_add_bmod]
    have adjustment (k : Nat) :
        exp + (-52 : Int) + (((k + 1 : Nat) : Int) - 1023) = exp + ((k : Int) - 1074) := by
      omega
    rw [adjustment, assembleLdexp_eq, shifted]
  · simp only [normalizeLdexp, eZero, ite_false, Int.add_zero, Int.bmod_add_bmod]
    exact assembleLdexp_eq _ _ _

/-- At the native integer extrema the wrapped exponent sends signed halves
    to signed infinity and signed twos to signed zero on both Go word widths. -/
theorem ldexp_extremeExponent (width : Nat) (native : width = 32 ∨ width = 64) :
    let minExponent : Int := -(2 ^ (width - 1))
    let maxExponent : Int := 2 ^ (width - 1) - 1
    ldexp width 0x3fe0000000000000 minExponent = 0x7ff0000000000000 ∧
    ldexp width 0xbfe0000000000000 minExponent = 0xfff0000000000000 ∧
    ldexp width 0x4000000000000000 maxExponent = 0 ∧
    ldexp width 0xc000000000000000 maxExponent = 0x8000000000000000 := by
  rcases native with rfl | rfl <;> decide

end GoLean.Math
