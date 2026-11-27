import GoLean.Bytes.IndexFuncPrivate
import GoLean.Bytes.LastIndexFuncPrivate
import GoLean.Unicode.Utf8.DecodeRune

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:835–862. The left scan returns a
    nil slice when all bytes are trimmed. Otherwise only the left
    offset reduces capacity. The outer None denotes a callback panic. -/
def trimFunc (input : List (BitVec 8)) (spare : Nat) (predicate : Option (Nat → Bool)) :
    Option (Option (List (BitVec 8) × Nat) × List Nat) := do
  let (left, leftCalls) ← indexFuncPrivate input predicate false
  if left < 0 then return (none, leftCalls)
  let remaining := input.drop left.toNat
  let (right, rightCalls) ← lastIndexFuncPrivate remaining predicate false
  let endpoint := if right < 0 then 0 else
    let width := (GoLean.Unicode.Utf8.decodeRune ((remaining.drop right.toNat).map BitVec.toNat)).2
    right.toNat + width
  return (some (remaining.take endpoint, remaining.length + spare), leftCalls ++ rightCalls)

theorem trimFunc_empty (spare : Nat) (predicate : Option (Nat → Bool)) :
    trimFunc [] spare predicate = some (none, []) := rfl

theorem trimFunc_nil (lead : BitVec 8) (rest : List (BitVec 8)) (spare : Nat) :
    trimFunc (lead :: rest) spare none = none := rfl

theorem trimFunc_fullLeft (input : List (BitVec 8)) (spare : Nat)
    (predicate : Option (Nat → Bool)) (calls : List Nat)
    (trimmed : indexFuncPrivate input predicate false = some (-1, calls)) :
    trimFunc input spare predicate = some (none, calls) := by
  simp [trimFunc, trimmed]

end GoLean.Bytes
