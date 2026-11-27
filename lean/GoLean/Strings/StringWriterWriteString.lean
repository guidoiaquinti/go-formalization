import Init.Data.BitVec.Basic

namespace GoLean.Strings

structure StringWriterResult where
  count : Int
  error : Nat
  writes : List (List (BitVec 8))

/-- Go 1.27.0 src/strings/replace.go:314–320. The string's raw byte
    representation is supplied once to Write, and the response is returned
    unchanged. Error zero represents nil. -/
def stringWriterWriteString (write : List (BitVec 8) → Int × Nat)
    (input : List (BitVec 8)) : StringWriterResult :=
  let (count, error) := write input
  ⟨count, error, [input]⟩

theorem stringWriterWriteString_count (write : List (BitVec 8) → Int × Nat)
    (input : List (BitVec 8)) : (stringWriterWriteString write input).count = (write input).1 := by
  simp [stringWriterWriteString]

theorem stringWriterWriteString_error (write : List (BitVec 8) → Int × Nat)
    (input : List (BitVec 8)) : (stringWriterWriteString write input).error = (write input).2 := by
  simp [stringWriterWriteString]

theorem stringWriterWriteString_once (write : List (BitVec 8) → Int × Nat)
    (input : List (BitVec 8)) : (stringWriterWriteString write input).writes = [input] := by
  simp [stringWriterWriteString]

end GoLean.Strings
