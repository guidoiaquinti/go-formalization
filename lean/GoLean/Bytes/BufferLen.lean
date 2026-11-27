import Init.Data.List.Lemmas

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:18–25, 95–97. The receiver stores
    the byte slice length and read offset. Reachable states satisfy
    offset ≤ length. None denotes a nil receiver pointer panic. -/
def bufferLen (receiver : Option (Nat × Nat)) : Option Nat :=
  receiver.map fun (length, offset) => length - offset

theorem bufferLen_unread (bytes : List α) (offset : Nat) :
    bufferLen (some (bytes.length, offset)) = some (bytes.drop offset).length := by
  simp [bufferLen, List.length_drop]

theorem bufferLen_consumed (length : Nat) :
    bufferLen (some (length, length)) = some 0 := by
  simp [bufferLen]

theorem bufferLen_nil : bufferLen none = none := rfl

end GoLean.Bytes
