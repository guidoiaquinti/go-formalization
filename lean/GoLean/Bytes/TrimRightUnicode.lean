import GoLean.Strings.Trim

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1095–1107, with membership at
    lines 976–983. Scan backward, retaining an exact prefix with the
    same capacity and nilness even when everything was removed. -/
def trimRightUnicode (input : Option (List (BitVec 8) × Nat)) (cutset : List (BitVec 8)) :
    Option (List (BitVec 8) × Nat) :=
  input.map fun (bytes, capacity) => (GoLean.Strings.trimRightUnicode bytes cutset, capacity)

theorem trimRightUnicode_empty (capacity : Option Nat) (cutset : List (BitVec 8)) :
    trimRightUnicode (capacity.map fun cap => ([], cap)) cutset =
      capacity.map fun cap => ([], cap) := by
  cases capacity <;> rfl

theorem trimRightUnicode_capacity (input : Option (List (BitVec 8) × Nat))
    (cutset : List (BitVec 8)) :
    (trimRightUnicode input cutset).map Prod.snd = input.map Prod.snd := by
  cases input with
  | none => rfl
  | some slice => cases slice; rfl

theorem trimRightUnicode_length (input cutset : List (BitVec 8)) (capacity : Nat) :
    ((trimRightUnicode (some (input, capacity)) cutset).map (fun slice => slice.1.length)).getD 0
      ≤ input.length :=
  GoLean.Strings.trimRightUnicode_length input cutset

end GoLean.Bytes
