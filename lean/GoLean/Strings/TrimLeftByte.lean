import Init.Data.BitVec.Lemmas

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:1010–1015. The helper compares raw
    bytes and accepts every byte value, including values above ASCII. -/
def trimLeftByte (input : List (BitVec 8)) (byte : BitVec 8) : List (BitVec 8) :=
  input.dropWhile (· == byte)

theorem trimLeftByte_prefix (byte : BitVec 8) (count : Nat) (tail : List (BitVec 8)) :
    trimLeftByte (List.replicate count byte ++ tail) byte = trimLeftByte tail byte := by
  induction count with
  | zero => rfl
  | succ count ih => simp [List.replicate_succ, trimLeftByte]

theorem trimLeftByte_nonmatching (first byte : BitVec 8) (tail : List (BitVec 8))
    (different : first ≠ byte) :
    trimLeftByte (first :: tail) byte = first :: tail := by
  have noMatch : (first == byte) = false := by simpa using different
  simp [trimLeftByte, List.dropWhile, noMatch]

theorem trimLeftByte_empty (byte : BitVec 8) : trimLeftByte [] byte = [] := rfl

end GoLean.Strings
