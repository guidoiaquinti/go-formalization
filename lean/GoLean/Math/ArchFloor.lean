import GoLean.Math.FloorCeil

namespace GoLean.Math

/-- Observable floor rounding for the arm64 FRINTMD instruction and
    amd64's magnitude guard, integer conversion and conditional decrement.
    NaN words represent the NaN class.
    https://github.com/golang/go/blob/go1.27.0/src/math/floor_arm64.s#L8-L13
    https://github.com/golang/go/blob/go1.27.0/src/math/floor_amd64.s#L10-L29 -/
def archFloor (amd64 : Bool) (word : BitVec 64) : Nat :=
  let magnitude := word.toNat % 2 ^ 63
  if amd64 && (magnitude = 0 || 0x4330000000000000 ≤ magnitude) then word.toNat
  else floorWord word

theorem archFloor_spec (amd64 : Bool) (word : BitVec 64) :
    archFloor amd64 word =
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
  have result := floorWord_spec word
  unfold archFloor
  dsimp only
  split
  · rename_i guard
    have magnitude : word.toNat % 2 ^ 63 = 0 ∨
        0x4330000000000000 ≤ word.toNat % 2 ^ 63 := by
      simpa using (Bool.and_eq_true_iff.mp guard).2
    have exp : word.toNat / 2 ^ 52 % 2048 = word.toNat % 2 ^ 63 / 2 ^ 52 := by omega
    rcases magnitude with zero | large
    · simp [exp, zero]
    · have integral : ¬ word.toNat % 2 ^ 63 / 2 ^ 52 < 1075 := by omega
      have notSmall : ¬ word.toNat % 2 ^ 63 / 2 ^ 52 < 1023 := by omega
      simp [exp, integral, notSmall]
  · exact result

end GoLean.Math
