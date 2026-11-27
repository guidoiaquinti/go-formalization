import GoLean.Strings.TrimLeftByte

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1023–1032. Remove matching raw
    leading bytes, preserving spare capacity for a retained suffix.
    An empty result is nil, irrespective of the input slice header. -/
def trimLeftByte (input : List (BitVec 8)) (cut : BitVec 8) (spare : Nat) :
    Option (List (BitVec 8) × Nat) :=
  let retained := GoLean.Strings.trimLeftByte input cut
  if retained.isEmpty then none else some (retained, retained.length + spare)

theorem trimLeftByte_empty (cut : BitVec 8) (spare : Nat) :
    trimLeftByte [] cut spare = none := rfl

theorem trimLeftByte_prefix (cut : BitVec 8) (count : Nat) (tail : List (BitVec 8)) (spare : Nat) :
    trimLeftByte (List.replicate count cut ++ tail) cut spare = trimLeftByte tail cut spare := by
  simp [trimLeftByte, GoLean.Strings.trimLeftByte_prefix]

theorem trimLeftByte_nonmatching (first cut : BitVec 8) (tail : List (BitVec 8))
    (spare : Nat) (different : first ≠ cut) :
    trimLeftByte (first :: tail) cut spare = some (first :: tail, (first :: tail).length + spare) := by
  simp [trimLeftByte, GoLean.Strings.trimLeftByte_nonmatching first cut tail different]

end GoLean.Bytes
