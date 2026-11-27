import GoLean.Math.Sin

namespace GoLean.Math

/-- Go 1.27.0 src/math/trig_reduce.go:31–78. For positive finite inputs,
    return the input below rounded π/4 and otherwise apply Payne-Hanek
    reduction. Pinned trigonometric callers use the range starting at 2^29. -/
def trigReduceWord (word : BitVec 64) : Nat × Nat :=
  if word.toNat < 0x3fe921fb54442d18 then (0, word.toNat)
  else reduceLargeSin word.toNat

theorem trigReduceWord_octant (word : BitVec 64) :
    (trigReduceWord word).1 < 8 := by
  unfold trigReduceWord
  split
  · change 0 < 8
    decide
  · simp only [reduceLargeSin]
    exact Nat.mod_lt _ (by decide)

theorem trigReduceWord_even (word : BitVec 64) :
    (trigReduceWord word).1 % 2 = 0 := by
  unfold trigReduceWord
  split
  · rfl
  · simp only [reduceLargeSin]
    omega

theorem trigReduceWord_small (word : BitVec 64) (small : word.toNat < 0x3fe921fb54442d18) :
    trigReduceWord word = (0, word.toNat) := by
  simp [trigReduceWord, small]

end GoLean.Math
