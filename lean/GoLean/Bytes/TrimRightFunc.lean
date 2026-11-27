import GoLean.Bytes.LastIndexFuncPrivate
import GoLean.Unicode.Utf8.DecodeRune

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:847–856. The final two-index slice
    preserves original capacity and nilness. The outer None denotes panic. -/
def trimRightFunc (input : Option (List (BitVec 8) × Nat))
    (predicate : Option (Nat → Bool)) :
    Option (Option (List (BitVec 8) × Nat) × List Nat) := do
  let (bytes, capacity) := input.getD ([], 0)
  let (offset, calls) ← lastIndexFuncPrivate bytes predicate false
  let endpoint := if offset < 0 then 0 else
    let width := (GoLean.Unicode.Utf8.decodeRune ((bytes.drop offset.toNat).map BitVec.toNat)).2
    offset.toNat + width
  return (input.map (fun _ => (bytes.take endpoint, capacity)), calls)

theorem trimRightFunc_empty (input : Option Nat) (predicate : Option (Nat → Bool)) :
    trimRightFunc (input.map fun cap => ([], cap)) predicate =
      some (input.map fun cap => ([], cap), []) := by
  cases input <;> rfl

theorem trimRightFunc_nil (lead : BitVec 8) (rest : List (BitVec 8)) (capacity : Nat) :
    trimRightFunc (some (lead :: rest, capacity)) none = none := rfl

theorem trimRightFunc_all (bytes : List (BitVec 8)) (capacity : Nat)
    (predicate : Option (Nat → Bool)) (calls : List Nat)
    (missing : lastIndexFuncPrivate bytes predicate false = some (-1, calls)) :
    trimRightFunc (some (bytes, capacity)) predicate = some (some ([], capacity), calls) := by
  simp [trimRightFunc, missing]

theorem trimRightFunc_at (bytes : List (BitVec 8)) (capacity offset : Nat)
    (predicate : Option (Nat → Bool)) (calls : List Nat)
    (found : lastIndexFuncPrivate bytes predicate false = some (Int.ofNat offset, calls)) :
    trimRightFunc (some (bytes, capacity)) predicate =
      some (some (bytes.take (offset + (GoLean.Unicode.Utf8.decodeRune
        ((bytes.drop offset).map BitVec.toNat)).2), capacity), calls) := by
  have nonnegative : ¬ (offset : Int) < 0 := by omega
  simp [trimRightFunc, found, nonnegative]

end GoLean.Bytes
