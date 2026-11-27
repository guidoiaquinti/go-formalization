import Init.Data.List.TakeDrop

namespace GoLean.Strings

/-- One call to IndexByte followed by the two source substring assignments. -/
private def nextLine (s : List (BitVec 8)) : List (BitVec 8) × List (BitVec 8) :=
  let before := s.takeWhile (· != 10)
  match s.dropWhile (· != 10) with
  | [] => (before, [])
  | newline :: rest => (before ++ [newline], rest)

/--
Go 1.27.0 strings/iter.go:18–33. A callback budget counts the call that
returns false. The second result is the iterator's captured string after
that invocation. A zero budget represents skipping the invocation.
-/
def lines (s : List (BitVec 8)) : Nat → List (List (BitVec 8)) × List (BitVec 8)
  | 0 => ([], s)
  | budget + 1 =>
    if s.isEmpty then ([], []) else
      let (line, remaining) := nextLine s
      let (more, rest) := lines remaining budget
      (line :: more, rest)

private theorem nextLine_reconstruct (s : List (BitVec 8)) :
    (nextLine s).1 ++ (nextLine s).2 = s := by
  have split := List.takeWhile_append_dropWhile (p := (· != (10 : BitVec 8))) (l := s)
  unfold nextLine
  split <;> rename_i h <;> rw [h] at split <;> simpa [List.append_assoc] using split

/-- A callback stop preserves the exact bytes that have yet to be yielded. -/
theorem lines_reconstruct (s : List (BitVec 8)) (budget : Nat) :
    (lines s budget).1.flatten ++ (lines s budget).2 = s := by
  induction budget generalizing s with
  | zero => simp [lines]
  | succ budget ih =>
    cases s with
    | nil => simp [lines]
    | cons b rest =>
      simp only [lines, List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte,
        List.flatten_cons, List.append_assoc]
      rw [ih]
      exact nextLine_reconstruct (b :: rest)

end GoLean.Strings
