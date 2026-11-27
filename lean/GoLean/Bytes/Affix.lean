import GoLean.Strings.CutPrefix
import GoLean.Strings.CutSuffix

namespace GoLean.Bytes

/-- A nil slice is absent. A non-nil slice contains its visible bytes and capacity. -/
abbrev Slice := Option (List (BitVec 8) × Nat)

def contents (src : Slice) : List (BitVec 8) := (src.getD ([], 0)).1

/-- The length guard and byte equality are the same as for strings.
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L19-L22
    https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L583-L591 -/
def hasPrefix (src : Slice) (pre : List (BitVec 8)) : Bool :=
  Strings.hasPrefix (contents src) pre

def hasSuffix (src : Slice) (suffix : List (BitVec 8)) : Bool :=
  Strings.hasSuffix (contents src) suffix

/-- Removing leading bytes also subtracts their count from capacity. -/
private def dropStart (src : Slice) (n : Nat) : Slice :=
  src.map fun (bytes, cap) => (bytes.drop n, cap - n)

/-- A two-index slice with start zero preserves capacity. -/
private def dropEnd (src : Slice) (n : Nat) : Slice :=
  src.map fun (bytes, cap) => (bytes.take (bytes.length - n), cap)

/-- https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L1401-L1409 -/
def cutPrefix (src : Slice) (pre : List (BitVec 8)) : Slice × Bool :=
  if hasPrefix src pre then (dropStart src pre.length, true) else (src, false)

/-- https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L1411-L1422 -/
def cutSuffix (src : Slice) (suffix : List (BitVec 8)) : Slice × Bool :=
  if hasSuffix src suffix then (dropEnd src suffix.length, true) else (src, false)

/-- https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L864-L871 -/
def trimPrefix (src : Slice) (pre : List (BitVec 8)) : Slice :=
  if hasPrefix src pre then dropStart src pre.length else src

/-- https://github.com/golang/go/blob/go1.27.0/src/bytes/bytes.go#L873-L880 -/
def trimSuffix (src : Slice) (suffix : List (BitVec 8)) : Slice :=
  if hasSuffix src suffix then dropEnd src suffix.length else src

theorem hasPrefix_spec (src : Slice) (pre : List (BitVec 8)) :
    hasPrefix src pre = true ↔ ∃ remainder, contents src = pre ++ remainder :=
  Strings.hasPrefix_spec (contents src) pre

theorem cutPrefix_spec (src : Slice) (pre : List (BitVec 8)) :
    (cutPrefix src pre).2 = hasPrefix src pre ∧
    ((cutPrefix src pre).2 = true →
      (cutPrefix src pre).1 = src.map (fun (bytes, cap) => (bytes.drop pre.length, cap - pre.length))) ∧
    ((cutPrefix src pre).2 = false → (cutPrefix src pre).1 = src) := by
  by_cases h : hasPrefix src pre = true <;> simp [cutPrefix, dropStart, h]

theorem trimPrefix_spec (src : Slice) (pre : List (BitVec 8)) :
    trimPrefix src pre = (cutPrefix src pre).1 := by
  by_cases h : hasPrefix src pre = true <;> simp [trimPrefix, cutPrefix, h]

theorem hasSuffix_spec (src : Slice) (suffix : List (BitVec 8)) :
    hasSuffix src suffix = true ↔ ∃ remainder, contents src = remainder ++ suffix :=
  Strings.hasSuffix_spec (contents src) suffix

theorem cutSuffix_spec (src : Slice) (suffix : List (BitVec 8)) :
    (cutSuffix src suffix).2 = hasSuffix src suffix ∧
    ((cutSuffix src suffix).2 = true →
      (cutSuffix src suffix).1 = src.map (fun (bytes, cap) => (bytes.take (bytes.length - suffix.length), cap))) ∧
    ((cutSuffix src suffix).2 = false → (cutSuffix src suffix).1 = src) := by
  by_cases h : hasSuffix src suffix = true <;> simp [cutSuffix, dropEnd, h]

theorem trimSuffix_spec (src : Slice) (suffix : List (BitVec 8)) :
    trimSuffix src suffix = (cutSuffix src suffix).1 := by
  by_cases h : hasSuffix src suffix = true <;> simp [trimSuffix, cutSuffix, h]

end GoLean.Bytes
