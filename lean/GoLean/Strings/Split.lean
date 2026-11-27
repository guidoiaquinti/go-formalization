import GoLean.Strings.SplitN

namespace GoLean.Strings

/-- Go 1.27.0 strings/strings.go:342. Retain no separator bytes and use
    a negative limit for the complete split. -/
def split (input separator : List Nat) : Option (List (List Nat)) :=
  genSplit input separator 0 (-1)

theorem split_eq_splitN (input separator : List Nat) :
    split input separator = splitN input separator (-1) := rfl

theorem split_nonnil (input separator : List Nat) :
    (split input separator).isSome = true := by
  simp only [split, genSplit]
  simp
  split <;> rfl

theorem split_empty : split [] [] = some [] := by
  simp [split, genSplit, explode, countRunes]
  rfl

end GoLean.Strings
