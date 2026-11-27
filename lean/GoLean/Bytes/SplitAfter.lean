import GoLean.Bytes.Split
import GoLean.Bytes.SplitAfterN

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:435–437. An unlimited split retains
    each matched separator and follows the shared capacity rules. -/
def splitAfter (input separator : List Nat) (spare : Nat) : Option (List (List Nat × Nat)) :=
  splitAfterN input separator spare (-1)

theorem splitAfter_eq_splitAfterN (input separator : List Nat) (spare : Nat) :
    splitAfter input separator spare = splitAfterN input separator spare (-1) := rfl

theorem splitAfter_noSeparator (input : List Nat) (spare : Nat) :
    splitAfter input [] spare = split input [] spare :=
  splitAfterN_noSeparator input spare (-1)

theorem splitAfter_nonnil (input separator : List Nat) (spare : Nat) :
    (splitAfter input separator spare).isSome = true := by
  have h := GoLean.Strings.splitAfter_nonnil input separator
  unfold splitAfter splitAfterN withSplitCaps
  rw [GoLean.Strings.splitAfterN_unbounded]
  cases p : GoLean.Strings.splitAfter input separator <;> simp [p] at *

end GoLean.Bytes
