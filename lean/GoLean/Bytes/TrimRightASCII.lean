import GoLean.Strings.TrimRightASCII

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1085–1093. Each trailing byte selects
    one of 256 table entries. Outer None denotes a nil-table panic.
    Shortening the slice preserves nilness and capacity. -/
def trimRightASCII (input : Option (List (BitVec 8) × Nat))
    (table : Option (BitVec 8 → Bool)) : Option (Option (List (BitVec 8) × Nat)) :=
  match input with
  | none => some none
  | some (bytes, capacity) =>
    (GoLean.Strings.trimRightASCII bytes table).map fun retained => some (retained, capacity)

theorem trimRightASCII_empty (capacity : Option Nat) (table : Option (BitVec 8 → Bool)) :
    trimRightASCII (capacity.map fun cap => ([], cap)) table =
      some (capacity.map fun cap => ([], cap)) := by
  cases capacity <;> rfl

theorem trimRightASCII_nil (head : BitVec 8) (tail : List (BitVec 8)) (capacity : Nat) :
    trimRightASCII (some (head :: tail, capacity)) none = none := rfl

theorem trimRightASCII_skip (leading : List (BitVec 8)) (last : BitVec 8) (capacity : Nat)
    (member : BitVec 8 → Bool) (selected : member last = true) :
    trimRightASCII (some (leading ++ [last], capacity)) (some member) =
      trimRightASCII (some (leading, capacity)) (some member) := by
  simp [trimRightASCII, GoLean.Strings.trimRightASCII_skip leading last member selected]

theorem trimRightASCII_retain (leading : List (BitVec 8)) (last : BitVec 8) (capacity : Nat)
    (member : BitVec 8 → Bool) (missing : member last = false) :
    trimRightASCII (some (leading ++ [last], capacity)) (some member) =
      some (some (leading ++ [last], capacity)) := by
  simp [trimRightASCII, GoLean.Strings.trimRightASCII_retain leading last member missing]

end GoLean.Bytes
