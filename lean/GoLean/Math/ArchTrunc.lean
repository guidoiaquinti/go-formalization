import GoLean.Math.Trunc

namespace GoLean.Math

/-- Go 1.27.0 src/math/floor_arm64.s:21–26 and src/math/floor_amd64.s:57–76.
    FRINTZD rounds toward zero. The amd64 magnitude guard returns the word
    unchanged, otherwise integer conversion discards fractional bits and
    the original sign is restored even when the integer part is zero. -/
def archTrunc (amd64 : Bool) (word : BitVec 64) : Nat :=
  if amd64 && (0x4330000000000000 ≤ word.toNat % 2 ^ 63) then word.toNat
  else trunc word

theorem archTrunc_spec (amd64 : Bool) (word : BitVec 64) :
    archTrunc amd64 word =
      let bits := word.toNat
      let biased := bits / 2 ^ 52 % 2048
      if biased < 1023 then bits - bits % 2 ^ 63
      else if biased < 1075 then
        let quantum := 2 ^ (1075 - biased)
        bits / quantum * quantum
      else bits := by
  unfold archTrunc
  dsimp only
  split
  · rename_i guard
    have large : 0x4330000000000000 ≤ word.toNat % 2 ^ 63 := by
      simpa using (Bool.and_eq_true_iff.mp guard).2
    have exp : word.toNat / 2 ^ 52 % 2048 = word.toNat % 2 ^ 63 / 2 ^ 52 := by omega
    have integral : ¬ word.toNat % 2 ^ 63 / 2 ^ 52 < 1075 := by omega
    have notSmall : ¬ word.toNat % 2 ^ 63 / 2 ^ 52 < 1023 := by omega
    simp [exp, integral, notSmall]
  · exact trunc_spec word

end GoLean.Math
