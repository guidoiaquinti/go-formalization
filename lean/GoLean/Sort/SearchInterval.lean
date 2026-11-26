namespace GoLean.Sort

/-- Go 1.27.0 src/sort/search.go:58-72. For nonnegative Go int bounds,
the unsigned sum fits in a word and its right shift is division by two. -/
def searchInterval (test : Nat → Bool) (lo hi : Nat) : Nat :=
  if lo < hi then
    let mid := (lo + hi) / 2
    if test mid then searchInterval test lo mid
    else searchInterval test (mid + 1) hi
  else lo
termination_by hi - lo
decreasing_by all_goals omega

theorem searchInterval_partition (test : Nat → Bool) (lo hi cut : Nat)
    (left : lo ≤ cut) (right : cut ≤ hi)
    (partition : ∀ i, lo ≤ i → i < hi → (test i = true ↔ cut ≤ i)) :
    searchInterval test lo hi = cut := by
  rw [searchInterval]
  split
  · rename_i openInterval
    let mid := (lo + hi) / 2
    have midLeft : lo ≤ mid := by dsimp [mid]; omega
    have midRight : mid < hi := by dsimp [mid]; omega
    have atMid := partition mid midLeft midRight
    change (if test mid then searchInterval test lo mid
      else searchInterval test (mid + 1) hi) = cut
    split
    · rename_i yes
      exact searchInterval_partition test lo mid cut left (atMid.mp yes)
        (fun i h₁ h₂ => partition i h₁ (by omega))
    · rename_i no
      have above : mid < cut := by
        exact Nat.lt_of_not_ge (fun h => no (atMid.mpr h))
      exact searchInterval_partition test (mid + 1) hi cut (by omega) right
        (fun i h₁ h₂ => partition i (by omega) h₂)
  · omega
termination_by hi - lo
decreasing_by all_goals dsimp [mid] at * <;> omega

end GoLean.Sort
