import GoLean.Math.Bits.Length

namespace GoLean.Math

/-- The FMA helper extracts the raw sign, exponent and fraction. It adds
    the implicit leading bit when the exponent field is nonzero and
    normalizes the fraction using a leading-zero count otherwise.
    https://github.com/golang/go/blob/go1.27.0/src/math/fma.go#L73-L91
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L11-L16 -/
def fmaSplit (word : BitVec 64) : Nat × Int × Nat :=
  let bits := word.toNat
  let sign := bits >>> 63
  let exponent := (bits >>> 52) &&& 2047
  let payload := bits &&& (2 ^ 52 - 1)
  if exponent = 0 then
    let shift := Bits.leadingZeros64 (BitVec.ofNat 64 payload) - 11
    (sign, 1 - (shift : Int), payload <<< shift)
  else (sign, exponent, payload ||| 2 ^ 52)

private theorem splitLength (word : BitVec 64) (nonzero : word.toNat ≠ 0) :
    Bits.len64 word = word.toNat.log2 + 1 := by
  have upper := (Bits.len64_spec word (word.toNat.log2 + 1)).mpr Nat.lt_log2_self
  have lower := (Nat.log2_lt nonzero).mpr
    ((Bits.len64_spec word (Bits.len64 word)).mp (Nat.le_refl _))
  omega

/-- The helper keeps the raw sign. A nonzero exponent field stays intact
    and a nonzero subnormal is shifted so that its top bit is at position
    52. A signed zero has zero significand and exponent -52. -/
theorem fmaSplit_spec (word : BitVec 64) :
    fmaSplit word =
      let bits := word.toNat
      let sign := bits / 2 ^ 63
      let magnitude := bits % 2 ^ 63
      let e : Nat := magnitude / 2 ^ 52
      let f := magnitude % 2 ^ 52
      if e ≠ 0 then (sign, (e : Int), 2 ^ 52 + f)
      else if f = 0 then (sign, -52, 0)
      else
        let k : Nat := f.log2
        (sign, (k : Int) - 51, f * 2 ^ (52 - k)) := by
  unfold fmaSplit
  dsimp only
  rw [Nat.shiftRight_eq_div_pow, Nat.shiftRight_eq_div_pow,
    show (2047 : Nat) = 2 ^ 11 - 1 by decide, Nat.and_two_pow_sub_one_eq_mod,
    Nat.and_two_pow_sub_one_eq_mod]
  have exponents : word.toNat / 2 ^ 52 % 2 ^ 11 = word.toNat % 2 ^ 63 / 2 ^ 52 := by omega
  have fractions : word.toNat % 2 ^ 52 = word.toNat % 2 ^ 63 % 2 ^ 52 := by omega
  rw [exponents, fractions]
  let f := word.toNat % 2 ^ 63 % 2 ^ 52
  let e := word.toNat % 2 ^ 63 / 2 ^ 52
  change (if e = 0 then
      let shift := Bits.leadingZeros64 (BitVec.ofNat 64 f) - 11
      (word.toNat / 2 ^ 63, 1 - (shift : Int), f <<< shift)
    else (word.toNat / 2 ^ 63, (e : Int), f ||| 2 ^ 52)) =
      if e ≠ 0 then (word.toNat / 2 ^ 63, (e : Int), 2 ^ 52 + f)
      else if f = 0 then (word.toNat / 2 ^ 63, -52, 0)
      else (word.toNat / 2 ^ 63, (f.log2 : Int) - 51, f * 2 ^ (52 - f.log2))
  by_cases ez : e = 0
  · rw [ite_eq_left ez, ite_eq_right (not_not_intro ez)]
    by_cases fz : f = 0
    · rw [ite_eq_left fz]
      have zeroLength : Bits.len64 (0#64) = 0 := by decide
      change (word.toNat / 2 ^ 63, 1 - ((Bits.leadingZeros64 (BitVec.ofNat 64 f) - 11 : Nat) : Int),
        f <<< (Bits.leadingZeros64 (BitVec.ofNat 64 f) - 11)) = _
      simp [fz, Bits.leadingZeros64, zeroLength]
    · rw [ite_eq_right fz]
      let payload := BitVec.ofNat 64 f
      have fBound : f < 2 ^ 52 := Nat.mod_lt _ (by decide)
      have payloadValue : payload.toNat = f := by
        exact Nat.mod_eq_of_lt (by omega)
      have length := splitLength payload (by simpa only [payloadValue] using fz)
      have topBound : f.log2 < 52 := (Nat.log2_lt fz).mpr fBound
      have shift : Bits.leadingZeros64 payload - 11 = 52 - f.log2 := by
        unfold Bits.leadingZeros64
        rw [length, payloadValue]
        omega
      change (word.toNat / 2 ^ 63, 1 - ((Bits.leadingZeros64 payload - 11 : Nat) : Int),
        f <<< (Bits.leadingZeros64 payload - 11)) =
        (word.toNat / 2 ^ 63, (f.log2 : Int) - 51, f * 2 ^ (52 - f.log2))
      rw [shift, Nat.shiftLeft_eq]
      have exponentEq : 1 - ((52 - f.log2 : Nat) : Int) = (f.log2 : Int) - 51 := by omega
      rw [exponentEq]
  · rw [ite_eq_right ez, ite_eq_left ez]
    have fBound : f < 2 ^ 52 := Nat.mod_lt _ (by decide)
    rw [Nat.or_two_pow_eq_add_of_lt fBound, Nat.add_comm f]

end GoLean.Math
