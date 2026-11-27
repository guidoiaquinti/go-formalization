import Init.Data.BitVec.Basic

namespace GoLean.Strings

private def matchingLength : List (BitVec 8) → List (BitVec 8) → Nat
  | x :: xs, y :: ys => if x = y then 1 + matchingLength xs ys else 0
  | _, _ => 0

/-- Go 1.27.0 src/strings/search.go:91–98. Compare bytes from the right
    until either input is exhausted or the first mismatch is encountered. -/
def longestCommonSuffix (a b : List (BitVec 8)) : Nat :=
  matchingLength a.reverse b.reverse

private theorem matchingLength_spec (a b : List (BitVec 8)) (n : Nat) :
    n ≤ matchingLength a b ↔ n ≤ a.length ∧ n ≤ b.length ∧ a.take n = b.take n := by
  induction n generalizing a b with
  | zero => simp
  | succ n ih =>
    cases a with
    | nil => simp [matchingLength]
    | cons x xs =>
      cases b with
      | nil => simp [matchingLength]
      | cons y ys =>
        by_cases same : x = y
        · subst y
          simpa [matchingLength, Nat.add_comm 1, Nat.succ_le_succ_iff] using ih xs ys
        · simp [matchingLength, same]

theorem longestCommonSuffix_spec (a b : List (BitVec 8)) (n : Nat) :
    n ≤ longestCommonSuffix a b ↔
      n ≤ a.length ∧ n ≤ b.length ∧ a.reverse.take n = b.reverse.take n := by
  simpa [longestCommonSuffix] using matchingLength_spec a.reverse b.reverse n

end GoLean.Strings
