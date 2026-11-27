import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Math

/-- Fields of the input binary64 word. Only the high sign bit contributes to
    sign; the middle eleven bits encode the biased exponent. -/
private def signField (word : Nat) := word / 2 ^ 63 * 2 ^ 63
private def exponentField (word : Nat) := word / 2 ^ 52 % 2048
private def fractionField (word : Nat) := word % 2 ^ 52

/-- The pure-Go helper multiplies a nonzero subnormal by 2^52 exactly. For a
    fraction with highest set bit k, the resulting normal exponent field is
    k+1 and its fraction field contains the shifted low bits.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L61-L67 -/
private def normalizeFrexp (exponent fraction : Nat) : Nat × Nat × Int :=
  if exponent = 0 then
    let k := fraction.log2
    (k + 1, (fraction <<< (52 - k)) % 2 ^ 52, -52)
  else (exponent, fraction, 0)

/-- Models the native Go fallback's result word and signed exponent. The
    word arithmetic extracts and assembles disjoint binary64 fields.
    https://github.com/golang/go/blob/go1.27.0/src/math/frexp.go#L17-L39 -/
def frexp (word : Nat) : Nat × Int :=
  let exponent := exponentField word
  let fraction := fractionField word
  if (exponent = 0 ∧ fraction = 0) ∨ exponent = 2047 then (word, (0 : Int))
  else
    let (e, f, adjustment) := normalizeFrexp exponent fraction
    (signField word + 1022 * 2 ^ 52 + f, adjustment + (e : Int) - 1023 + 1)

theorem normalizedPayload (f : Nat) (nonzero : f ≠ 0) (bound : f < 2 ^ 52) :
    (f <<< (52 - f.log2)) % 2 ^ 52 = (f - 2 ^ f.log2) <<< (52 - f.log2) := by
  let k := f.log2
  change (f <<< (52 - k)) % 2 ^ 52 = (f - 2 ^ k) <<< (52 - k)
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
  rw [Nat.shiftLeft_eq, Nat.shiftLeft_eq, Nat.sub_mul, power]
  omega

/-- Exact returned representation and exponent, including zero, infinities
    and NaN payloads. In a subnormal, removing the highest set bit leaves
    precisely the bits that will follow the new implicit leading one. -/
theorem frexp_spec (word : Nat) :
    frexp word =
      let magnitude := word % 2 ^ 63
      let sign := word - magnitude
      let exponent := magnitude / 2 ^ 52
      let fraction := magnitude % 2 ^ 52
      if magnitude = 0 ∨ exponent = 2047 then (word, (0 : Int))
      else if exponent = 0 then
        let k := fraction.log2
        (sign + 1022 * 2 ^ 52 + ((fraction - 2 ^ k) <<< (52 - k)), (k : Int) - 1073)
      else (sign + 1022 * 2 ^ 52 + fraction, (exponent : Int) - 1022) := by
  have es : exponentField word = word % 2 ^ 63 / 2 ^ 52 := by
    unfold exponentField
    omega
  have fs : fractionField word = word % 2 ^ 63 % 2 ^ 52 := by
    unfold fractionField
    omega
  have ss : signField word = word - word % 2 ^ 63 := by
    unfold signField
    omega
  have zs : (word % 2 ^ 63 / 2 ^ 52 = 0 ∧ word % 2 ^ 63 % 2 ^ 52 = 0) ↔
      word % 2 ^ 63 = 0 := by omega
  unfold frexp
  dsimp only
  rw [es, fs, ss]
  simp only [zs]
  split
  next _ => rfl
  next nonSpecial =>
    unfold normalizeFrexp
    split
    next eZero =>
      have fNonzero : word % 2 ^ 63 % 2 ^ 52 ≠ 0 := by omega
      dsimp only
      rw [normalizedPayload _ fNonzero (Nat.mod_lt _ (by decide))]
      congr 1
      omega
    next _ =>
      dsimp only
      congr 1
      omega

end GoLean.Math
