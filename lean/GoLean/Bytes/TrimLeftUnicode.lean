import GoLean.Strings.Trim

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1048–1060, with membership from
    lines 976–983. Drop each matching rune's raw byte width. Returning a
    nil slice after full removal also discards the remaining capacity. -/
def trimLeftUnicode (input cutset : List (BitVec 8)) (spare : Nat) :
    Option (List (BitVec 8) × Nat) :=
  let retained := GoLean.Strings.trimLeftUnicode input cutset
  if retained.isEmpty then none else some (retained, retained.length + spare)

theorem trimLeftUnicode_empty (cutset : List (BitVec 8)) (spare : Nat) :
    trimLeftUnicode [] cutset spare = none := rfl

theorem trimLeftUnicode_stop (head : BitVec 8) (tail cutset : List (BitVec 8)) (spare : Nat)
    (missing : GoLean.Strings.containsRune (cutset.map BitVec.toNat)
      (Int.ofNat (GoLean.Unicode.Utf8.decodeRuneInString ((head :: tail).map BitVec.toNat)).1) = false) :
    trimLeftUnicode (head :: tail) cutset spare =
      some (head :: tail, (head :: tail).length + spare) := by
  simp [trimLeftUnicode, GoLean.Strings.trimLeftUnicode_stop head tail cutset missing]

theorem trimLeftUnicode_capacity (input cutset : List (BitVec 8)) (spare : Nat)
    (retained : List (BitVec 8)) (capacity : Nat)
    (result : trimLeftUnicode input cutset spare = some (retained, capacity)) :
    capacity = retained.length + spare := by
  simp only [trimLeftUnicode] at result
  split at result
  · contradiction
  · cases result
    rfl

end GoLean.Bytes
