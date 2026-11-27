import GoLean.Strings.IndexFunc
import GoLean.Strings.LastIndexFunc

namespace GoLean.Strings

/-- Go 1.27.0 strings/strings.go:868–874. Drop the bytes before the first
    false callback result, retaining the calls that located it. -/
def trimLeftFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (List (BitVec 8) × List Nat) := do
  let (offset, calls) ← indexFuncWithTruth input predicate false
  return (if offset < 0 then [] else input.drop offset.toNat, calls)

/-- Go 1.27.0 strings/strings.go:878–887. Keep the final non-trimmed rune
    by decoding its width from the original bytes after the reverse search. -/
def trimRightFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (List (BitVec 8) × List Nat) := do
  let (offset, calls) ← lastIndexFuncWithTruth input predicate false
  if offset < 0 then return ([], calls)
  let width := (GoLean.Unicode.Utf8.decodeRuneInString
    ((input.drop offset.toNat).map BitVec.toNat)).2
  return (input.take (offset.toNat + width), calls)

/-- Go 1.27.0 strings/strings.go:891–893. The right scan sees only the suffix
    left by the first scan. None records a nil-callback panic. -/
def trimFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (List (BitVec 8) × List Nat) := do
  let (remaining, leftCalls) ← trimLeftFunc input predicate
  let (result, rightCalls) ← trimRightFunc remaining predicate
  return (result, leftCalls ++ rightCalls)

theorem trimFunc_empty (predicate : Option (Nat → Bool)) :
    trimFunc [] predicate = some ([], []) := rfl

theorem trimFunc_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    trimFunc (head :: tail) none = none := rfl

theorem trimFunc_fullLeft (input : List (BitVec 8)) (predicate : Option (Nat → Bool))
    (calls : List Nat) (trimmed : indexFuncWithTruth input predicate false = some (-1, calls)) :
    trimFunc input predicate = some ([], calls) := by
  simp [trimFunc, trimLeftFunc, trimmed, trimRightFunc, lastIndexFuncWithTruth]

theorem trimLeftFunc_empty (predicate : Option (Nat → Bool)) :
    trimLeftFunc [] predicate = some ([], []) := rfl

theorem trimLeftFunc_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    trimLeftFunc (head :: tail) none = none := rfl

theorem trimLeftFunc_at (input : List (BitVec 8)) (predicate : Option (Nat → Bool))
    (offset : Nat) (calls : List Nat)
    (found : indexFuncWithTruth input predicate false = some (Int.ofNat offset, calls)) :
    trimLeftFunc input predicate = some (input.drop offset, calls) := by
  have nonnegative : ¬(offset : Int) < 0 := by omega
  simp [trimLeftFunc, found, nonnegative]

theorem trimLeftFunc_all (input : List (BitVec 8)) (predicate : Option (Nat → Bool))
    (calls : List Nat) (missing : indexFuncWithTruth input predicate false = some (-1, calls)) :
    trimLeftFunc input predicate = some ([], calls) := by
  simp [trimLeftFunc, missing]

theorem trimRightFunc_empty (predicate : Option (Nat → Bool)) :
    trimRightFunc [] predicate = some ([], []) := rfl

theorem trimRightFunc_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    trimRightFunc (head :: tail) none = none := rfl

theorem trimRightFunc_at (input : List (BitVec 8)) (predicate : Option (Nat → Bool))
    (offset : Nat) (calls : List Nat)
    (found : lastIndexFuncWithTruth input predicate false = some (Int.ofNat offset, calls)) :
    trimRightFunc input predicate =
      some (input.take (offset + (GoLean.Unicode.Utf8.decodeRuneInString
        ((input.drop offset).map BitVec.toNat)).2), calls) := by
  have nonnegative : ¬(offset : Int) < 0 := by omega
  simp [trimRightFunc, found, nonnegative]

theorem trimRightFunc_all (input : List (BitVec 8)) (predicate : Option (Nat → Bool))
    (calls : List Nat) (missing : lastIndexFuncWithTruth input predicate false = some (-1, calls)) :
    trimRightFunc input predicate = some ([], calls) := by
  simp [trimRightFunc, missing]

end GoLean.Strings
