import GoLean.Math.Frexp

namespace GoLean.Math

private def exponent (word : Nat) := word / 2 ^ 52 % 2048
private def fraction (word : Nat) := word % 2 ^ 52
private def sign (word : Nat) := word / 2 ^ 63 * 2 ^ 63

/-- Go normalizes finite nonzero subnormals by multiplying by 2^52.
    This multiplication is exact. A highest fraction bit at position k
    becomes the implicit leading bit for biased exponent k+1.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L59-L67 -/
def normalize (word : Nat) : Nat × Int :=
  if exponent word = 0 then
    let f := fraction word
    (sign word + (f.log2 + 1) * 2 ^ 52 + ((f <<< (52 - f.log2)) % 2 ^ 52), -52)
  else (word, 0)

/-- Go extracts the normalized biased exponent and adds the adjustment.
    https://github.com/golang/go/blob/go1.27.0/src/math/logb.go#L47-L52 -/
def ilogbFinite (word : Nat) : Int :=
  let (normalized, adjustment) := normalize word
  (exponent normalized : Int) - 1023 + adjustment

theorem normalize_spec (word : Nat) (nonzero : word % 2 ^ 63 ≠ 0) :
    normalize word =
      let magnitude := word % 2 ^ 63
      if magnitude < 2 ^ 52 then
        let k := magnitude.log2
        (word - magnitude + (k + 1) * 2 ^ 52 +
          ((magnitude - 2 ^ k) <<< (52 - k)), -52)
      else (word, 0) := by
  have e : exponent word = word % 2 ^ 63 / 2 ^ 52 := by unfold exponent; omega
  have f : fraction word = word % 2 ^ 63 % 2 ^ 52 := by unfold fraction; omega
  have s : sign word = word - word % 2 ^ 63 := by unfold sign; omega
  simp only [normalize, e, f, s]
  by_cases subnormal : word % 2 ^ 63 < 2 ^ 52
  · rw [ite_eq_left (by omega), ite_eq_left subnormal, Nat.mod_eq_of_lt subnormal]
    rw [normalizedPayload _ nonzero subnormal]
  · rw [ite_eq_right (by omega), ite_eq_right subnormal]

theorem ilogbFinite_spec (word : Nat) (nonzero : word % 2 ^ 63 ≠ 0) :
    ilogbFinite word =
      let magnitude := word % 2 ^ 63
      if magnitude < 2 ^ 52 then (magnitude.log2 : Int) - 1074
      else (magnitude / 2 ^ 52 : Nat) - (1023 : Int) := by
  have e : exponent word = word % 2 ^ 63 / 2 ^ 52 := by unfold exponent; omega
  have f : fraction word = word % 2 ^ 63 % 2 ^ 52 := by unfold fraction; omega
  by_cases subnormal : word % 2 ^ 63 < 2 ^ 52
  · have eb : exponent word = 0 := by omega
    have fb : fraction word = word % 2 ^ 63 := f.trans (Nat.mod_eq_of_lt subnormal)
    have logBound : (word % 2 ^ 63).log2 < 52 := (Nat.log2_lt nonzero).2 subnormal
    simp only [ilogbFinite, normalize, eb, ite_true, fb, subnormal]
    let magnitude := word % 2 ^ 63
    change magnitude.log2 < 52 at logBound
    let low := (magnitude <<< (52 - magnitude.log2)) % 2 ^ 52
    have lowBound : low < 2 ^ 52 := Nat.mod_lt _ (by decide)
    have normalizedExponent :
        exponent (sign word + (magnitude.log2 + 1) * 2 ^ 52 + low) =
          magnitude.log2 + 1 := by
      unfold exponent sign
      change (word / 2 ^ 63 * 2 ^ 63 + (magnitude.log2 + 1) * 2 ^ 52 + low) / 2 ^ 52 % 2048 = _
      omega
    change (exponent (sign word + (magnitude.log2 + 1) * 2 ^ 52 + low) : Int)
      - 1023 + -52 = (magnitude.log2 : Int) - 1074
    rw [normalizedExponent]
    omega
  · simp [ilogbFinite, normalize, e, subnormal]

end GoLean.Math
