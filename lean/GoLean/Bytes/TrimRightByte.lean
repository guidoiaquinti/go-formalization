import Init.Data.BitVec.Lemmas

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:1078–1083. Remove matching raw bytes
    from the end. A two-index prefix slice retains nilness and capacity. -/
def trimRightByte (input : Option (List (BitVec 8) × Nat)) (cut : BitVec 8) :
    Option (List (BitVec 8) × Nat) :=
  input.map fun (bytes, capacity) => ((bytes.reverse.dropWhile (· == cut)).reverse, capacity)

theorem trimRightByte_capacity (input : Option (List (BitVec 8) × Nat)) (cut : BitVec 8) :
    (trimRightByte input cut).map Prod.snd = input.map Prod.snd := by
  cases input with
  | none => rfl
  | some slice => cases slice; rfl

theorem trimRightByte_empty (input : Option Nat) (cut : BitVec 8) :
    trimRightByte (input.map fun capacity => ([], capacity)) cut =
      input.map fun capacity => ([], capacity) := by
  cases input <;> rfl

theorem trimRightByte_skipLast (input : List (BitVec 8)) (cut : BitVec 8) (capacity : Nat) :
    trimRightByte (some (input ++ [cut], capacity)) cut =
      trimRightByte (some (input, capacity)) cut := by
  simp [trimRightByte]

theorem trimRightByte_keepLast (input : List (BitVec 8)) (last cut : BitVec 8)
    (capacity : Nat) (different : last ≠ cut) :
    trimRightByte (some (input ++ [last], capacity)) cut = some (input ++ [last], capacity) := by
  have noMatch : (last == cut) = false := by simpa using different
  simp [trimRightByte, noMatch]

end GoLean.Bytes
