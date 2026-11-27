import GoLean.Bytes.SplitN
import GoLean.Strings.Split

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:429. The wrapper keeps no separator
    bytes and uses the unbounded limit of the same scanner as SplitN. -/
def split (input separator : List Nat) (spare : Nat) : Option (List (List Nat × Nat)) :=
  splitN input separator spare (-1)

theorem split_eq_splitN (input separator : List Nat) (spare : Nat) :
    split input separator spare = splitN input separator spare (-1) := rfl

theorem split_nonnil (input separator : List Nat) (spare : Nat) :
    (split input separator spare).isSome = true := by
  have h := GoLean.Strings.split_nonnil input separator
  rw [GoLean.Strings.split_eq_splitN] at h
  unfold split splitN withSplitCaps
  cases p : GoLean.Strings.splitN input separator (-1) <;> simp [p] at *

end GoLean.Bytes
