import Init.Data.BitVec.Lemmas

namespace GoLean.Strings

/-- Go 1.27.0 src/strings/strings.go:1055–1060. The helper consumes raw
    matching bytes from the right, irrespective of whether they are ASCII. -/
def trimRightByte (input : List (BitVec 8)) (byte : BitVec 8) : List (BitVec 8) :=
  (input.reverse.dropWhile (· == byte)).reverse

theorem trimRightByte_suffix (byte : BitVec 8) (count : Nat) (leading : List (BitVec 8)) :
    trimRightByte (leading ++ List.replicate count byte) byte = trimRightByte leading byte := by
  simp [trimRightByte, List.reverse_append, List.reverse_replicate]

theorem trimRightByte_nonmatching (last byte : BitVec 8) (leading : List (BitVec 8))
    (different : last ≠ byte) :
    trimRightByte (leading ++ [last]) byte = leading ++ [last] := by
  have noMatch : (last == byte) = false := by simpa using different
  simp [trimRightByte, List.reverse_append, noMatch]

theorem trimRightByte_empty (byte : BitVec 8) : trimRightByte [] byte = [] := rfl

end GoLean.Strings
