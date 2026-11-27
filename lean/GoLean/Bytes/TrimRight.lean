import GoLean.Strings.Trim

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1065–1108. The shared backward scanner
    removes trailing cutset members. The two-index prefix slice preserves
    the original capacity and nilness, including when no bytes remain. -/
def trimRight (input : Option (List (BitVec 8) × Nat)) (cutset : List (BitVec 8)) :
    Option (List (BitVec 8) × Nat) :=
  input.map fun (bytes, capacity) => (GoLean.Strings.trimRight bytes cutset, capacity)

theorem trimRight_nil (cutset : List (BitVec 8)) : trimRight none cutset = none := rfl

theorem trimRight_emptyCutset (input : Option (List (BitVec 8) × Nat)) :
    trimRight input [] = input := by
  cases input with
  | none => rfl
  | some slice => cases slice; simp [trimRight, GoLean.Strings.trimRight_emptyCutset]

theorem trimRight_capacity (input : Option (List (BitVec 8) × Nat)) (cutset : List (BitVec 8)) :
    (trimRight input cutset).map Prod.snd = input.map Prod.snd := by
  cases input with
  | none => rfl
  | some slice => cases slice; rfl

theorem trimRight_singleASCII (input : List (BitVec 8)) (byte : BitVec 8) (capacity : Nat)
    (ascii : byte.toNat < 128) :
    trimRight (some (input, capacity)) [byte] =
      some ((input.reverse.dropWhile (· == byte)).reverse, capacity) := by
  simp [trimRight, GoLean.Strings.trimRight_singleASCII input byte ascii]

end GoLean.Bytes
