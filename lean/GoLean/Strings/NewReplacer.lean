import Init.Data.BitVec.Basic

namespace GoLean.Strings

structure NewReplacerState where
  oldNew : Option (List (List (BitVec 8)))
  onceIsZero : Bool
  implementationIsNil : Bool

/-- Go 1.27.0 src/strings/replace.go:32–37. A returned Replacer
    owns a copy of the argument slice and has yet to build its implementation.
    None represents the odd-argument panic. -/
def newReplacer (oldNew : List (List (BitVec 8))) : Option NewReplacerState :=
  if oldNew.length % 2 == 1 then none
  else some ⟨if oldNew.isEmpty then none else some oldNew, true, true⟩

theorem newReplacer_odd (oldNew : List (List (BitVec 8))) (odd : oldNew.length % 2 = 1) :
    newReplacer oldNew = none := by simp [newReplacer, odd]

theorem newReplacer_empty : newReplacer [] = some ⟨none, true, true⟩ := rfl

theorem newReplacer_even (first second : List (BitVec 8)) (rest : List (List (BitVec 8)))
    (even : rest.length % 2 = 0) :
    newReplacer (first :: second :: rest) = some ⟨some (first :: second :: rest), true, true⟩ := by
  have parity : (rest.length + 1 + 1) % 2 = 0 := by omega
  simp [newReplacer, parity]

end GoLean.Strings
