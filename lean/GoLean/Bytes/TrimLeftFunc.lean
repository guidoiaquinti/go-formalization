import GoLean.Bytes.IndexFuncPrivate

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:837–843. The first false predicate
    keeps the original suffix and its remaining capacity. No match
    yields nil. The outer None represents a callback panic. -/
def trimLeftFunc (input : List (BitVec 8)) (spare : Nat) (predicate : Option (Nat → Bool)) :
    Option (Option (List (BitVec 8) × Nat) × List Nat) := do
  let (offset, calls) ← indexFuncPrivate input predicate false
  if offset < 0 then return (none, calls)
  let suffix := input.drop offset.toNat
  return (some (suffix, suffix.length + spare), calls)

theorem trimLeftFunc_empty (spare : Nat) (predicate : Option (Nat → Bool)) :
    trimLeftFunc [] spare predicate = some (none, []) := rfl

theorem trimLeftFunc_nil (lead : BitVec 8) (rest : List (BitVec 8)) (spare : Nat) :
    trimLeftFunc (lead :: rest) spare none = none := rfl

theorem trimLeftFunc_at (input : List (BitVec 8)) (spare offset : Nat)
    (predicate : Option (Nat → Bool)) (calls : List Nat)
    (found : indexFuncPrivate input predicate false = some (Int.ofNat offset, calls)) :
    trimLeftFunc input spare predicate =
      some (some (input.drop offset, (input.drop offset).length + spare), calls) := by
  have nonnegative : ¬ (offset : Int) < 0 := by omega
  simp [trimLeftFunc, found, nonnegative]

theorem trimLeftFunc_all (input : List (BitVec 8)) (spare : Nat)
    (predicate : Option (Nat → Bool)) (calls : List Nat)
    (missing : indexFuncPrivate input predicate false = some (-1, calls)) :
    trimLeftFunc input spare predicate = some (none, calls) := by
  simp [trimLeftFunc, missing]

end GoLean.Bytes
