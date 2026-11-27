namespace GoLean.Math

private def signWeight : Nat := 0x8000000000000000
private def infinity : Nat := 0x7ff0000000000000
private def nanWord : Nat := 0x7ff8000000000001

/-- Signed rank of a non-NaN word in binary64 order. Both zeros have rank
    zero. Within either sign, consecutive magnitudes are one rank apart. -/
private def nextRank (word : Nat) : Int :=
  if word < signWeight then (word : Int) else (signWeight : Int) - word

/-- Go tests NaNs first, then numerical equality, then zero. The remaining
    branches increment or decrement the integer encoding.
    https://github.com/golang/go/blob/go1.27.0/src/math/nextafter.go#L31-L53
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L7-L15
    https://github.com/golang/go/blob/go1.27.0/src/math/bits.go#L30-L40 -/
def nextafter (x y : Nat) : Nat :=
  let a := nextRank x
  let b := nextRank y
  if infinity < a.natAbs ∨ infinity < b.natAbs then nanWord
  else if a = b then x
  else if a = 0 then if y < signWeight then 1 else signWeight + 1
  else if (a < b ↔ 0 < a) then x + 1
  else x - 1

private def wordAtRank (rank : Int) (negativeZero : Bool) : Nat :=
  if rank < 0 ∨ (rank = 0 ∧ negativeZero = true) then
    signWeight + rank.natAbs
  else rank.natAbs

/-- For unequal non-NaN inputs, the output is the adjacent rank toward
    the target. A step reaching rank zero keeps the input's sign. -/
theorem nextafter_spec (x y : Nat) :
    nextafter x y =
      if infinity < (nextRank x).natAbs ∨ infinity < (nextRank y).natAbs then nanWord
      else if nextRank x = nextRank y then x
      else wordAtRank (nextRank x + (if nextRank x < nextRank y then 1 else -1))
        (decide (signWeight ≤ x)) := by
  simp only [nextafter]
  split
  · rfl
  · split
    · rfl
    · rename_i hDifferent
      by_cases hx : x < signWeight <;> by_cases hy : y < signWeight
      all_goals simp only [nextRank, hx, hy, ite_true, ite_false, wordAtRank] at hDifferent ⊢
      all_goals unfold signWeight at *
      all_goals split <;> split <;> try split <;> try split <;> try split
      all_goals simp_all [infinity] <;> omega

theorem nextafter_wordBound (x y : Nat) (hx : x < 2^64) :
    nextafter x y < 2^64 := by
  simp only [nextafter]
  split
  · decide
  · rename_i h
    have ax : (nextRank x).natAbs ≤ infinity := by omega
    have bounded : x ≤ signWeight + infinity := by
      unfold nextRank at ax
      split at ax <;> omega
    unfold signWeight infinity at bounded
    split
    · exact hx
    · split
      · split <;> decide
      · split <;> omega

end GoLean.Math
