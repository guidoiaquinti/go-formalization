import GoLean.Strings.Split

namespace GoLean.Strings

/-- Go 1.27.0 strings/strings.go:354–356. Retain the full separator and
    use a negative limit for the complete split. -/
def splitAfter (input separator : List Nat) : Option (List (List Nat)) :=
  genSplit input separator separator.length (-1)

theorem splitAfter_nonnil (input separator : List Nat) :
    (splitAfter input separator).isSome = true := by
  simp only [splitAfter, genSplit]
  simp
  split <;> rfl

theorem splitAfter_noSeparator (input : List Nat) :
    splitAfter input [] = split input [] := rfl

theorem splitAfter_empty : splitAfter [] [] = some [] := by
  rw [splitAfter_noSeparator, split_empty]

end GoLean.Strings
