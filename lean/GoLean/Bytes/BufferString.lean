import GoLean.Bytes.BufferAvailableBuffer

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:68–78. Return an immutable snapshot
    of the raw unread bytes, with a literal marker for a nil receiver. -/
def bufferString (receiver : Option BufferTailState) : List (BitVec 8) :=
  match receiver with
  | none => [60, 110, 105, 108, 62]
  | some state => (state.stored.getD []).drop state.offset

theorem bufferString_nil : bufferString none = [60, 110, 105, 108, 62] := rfl

theorem bufferString_unread (stored spare : List (BitVec 8)) (offset : Nat) :
    bufferString (some ⟨some stored, spare, offset⟩) = stored.drop offset := rfl

theorem bufferString_overwrite (consumed unread replacement spare : List (BitVec 8)) :
    let before := bufferString (some ⟨some (consumed ++ unread), spare, consumed.length⟩)
    let after := bufferString (some ⟨some (consumed ++ replacement), spare, consumed.length⟩)
    (before, after) = (unread, replacement) := by
  simp [bufferString]

end GoLean.Bytes
