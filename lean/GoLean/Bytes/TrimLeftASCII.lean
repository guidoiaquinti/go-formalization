import GoLean.Strings.TrimLeftASCII

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1034–1046. Scan raw table entries from
    the left. The outer None denotes a nil-table panic. An empty suffix
    returns a nil slice. A retained suffix preserves spare capacity. -/
def trimLeftASCII (input : List (BitVec 8)) (table : Option (BitVec 8 → Bool)) (spare : Nat) :
    Option (Option (List (BitVec 8) × Nat)) :=
  (GoLean.Strings.trimLeftASCII input table).map fun retained =>
    if retained.isEmpty then none else some (retained, retained.length + spare)

theorem trimLeftASCII_empty (table : Option (BitVec 8 → Bool)) (spare : Nat) :
    trimLeftASCII [] table spare = some none := rfl

theorem trimLeftASCII_nil (head : BitVec 8) (tail : List (BitVec 8)) (spare : Nat) :
    trimLeftASCII (head :: tail) none spare = none := rfl

theorem trimLeftASCII_skip (head : BitVec 8) (tail : List (BitVec 8)) (spare : Nat)
    (member : BitVec 8 → Bool) (selected : member head = true) :
    trimLeftASCII (head :: tail) (some member) spare =
      trimLeftASCII tail (some member) spare := by
  simp [trimLeftASCII, GoLean.Strings.trimLeftASCII_skip head tail member selected]

theorem trimLeftASCII_retain (head : BitVec 8) (tail : List (BitVec 8)) (spare : Nat)
    (member : BitVec 8 → Bool) (missing : member head = false) :
    trimLeftASCII (head :: tail) (some member) spare =
      some (some (head :: tail, (head :: tail).length + spare)) := by
  simp [trimLeftASCII, GoLean.Strings.trimLeftASCII_retain head tail member missing]

end GoLean.Bytes
