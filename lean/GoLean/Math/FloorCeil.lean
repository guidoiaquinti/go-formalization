import GoLean.Math.Trunc

namespace GoLean.Math

private def biased (word : BitVec 64) : Nat :=
  word.toNat / 2 ^ 52 % 2048

private def unitWord (word : BitVec 64) : Nat :=
  if biased word < 1023 then 1023 * 2 ^ 52
  else 2 ^ (1075 - biased word)

/-- The generic floor splits the magnitude into integer and fractional parts.
    Clearing fractional bits models the integer part. On a negative input,
    adding one to that magnitude adds one encoded unit. The unit for a zero
    integer part is the encoding of one. A carry at a power-of-two threshold
    advances the exponent field. Unchanged words include zeros and specials.
    https://github.com/golang/go/blob/go1.27.0/src/math/floor.go#L7-L33 -/
def floorWord (word : BitVec 64) : Nat :=
  let whole := trunc word
  if word.msb && whole != word.toNat then whole + unitWord word else whole

/-- The generic ceil negates the argument and the result of Floor.
    https://github.com/golang/go/blob/go1.27.0/src/math/floor.go#L36-L52 -/
def ceilWord (word : BitVec 64) : Nat :=
  floorWord (word ^^^ 0x8000000000000000) ^^^ (2 ^ 63)

theorem floorWord_spec (word : BitVec 64) :
    floorWord word =
      let bits := word.toNat
      let e := bits / 2 ^ 52 % 2048
      if e < 1023 then
        let signedZero := bits - bits % 2 ^ 63
        if word.msb && bits % 2 ^ 63 != 0 then signedZero + 1023 * 2 ^ 52
        else signedZero
      else if e < 1075 then
        let unit := 2 ^ (1075 - e)
        (bits / unit + (if word.msb && bits % unit != 0 then 1 else 0)) * unit
      else bits := by
  unfold floorWord unitWord biased
  rw [trunc_spec]
  dsimp only
  split
  · rename_i small
    have eq : (word.toNat - word.toNat % 2 ^ 63 = word.toNat) ↔
        word.toNat % 2 ^ 63 = 0 := by omega
    simp [bne, eq]
  · rename_i large
    split
    · rename_i fractional
      let unit := 2 ^ (1075 - word.toNat / 2 ^ 52 % 2048)
      have eq : (word.toNat / unit * unit = word.toNat) ↔ word.toNat % unit = 0 := by
        have h := Nat.mod_add_div' word.toNat unit
        omega
      change (if word.msb && (word.toNat / unit * unit != word.toNat)
        then word.toNat / unit * unit + unit else word.toNat / unit * unit) =
        (word.toNat / unit + (if word.msb && word.toNat % unit != 0 then 1 else 0)) * unit
      by_cases z : word.toNat % unit = 0
      · simp [z, eq.mpr z]
      · cases word.msb <;> simp [z, mt eq.mp z, Nat.add_mul]

    · rename_i integral
      simp

theorem ceilWord_spec (word : BitVec 64) :
    ceilWord word =
      let flipped := word ^^^ (0x8000000000000000 : BitVec 64)
      let bits := flipped.toNat
      let e := bits / 2 ^ 52 % 2048
      let rounded :=
        if e < 1023 then
          let signedZero := bits - bits % 2 ^ 63
          if flipped.msb && bits % 2 ^ 63 != 0 then signedZero + 1023 * 2 ^ 52
          else signedZero
        else if e < 1075 then
          let unit := 2 ^ (1075 - e)
          (bits / unit + (if flipped.msb && bits % unit != 0 then 1 else 0)) * unit
        else bits
      rounded ^^^ (2 ^ 63) := by
  rw [ceilWord, floorWord_spec]

end GoLean.Math
