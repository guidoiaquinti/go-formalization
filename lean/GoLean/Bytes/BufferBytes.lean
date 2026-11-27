import GoLean.Bytes.BufferAvailableBuffer

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:55–60. The view starts at the read
    cursor, with the original spare region following the unread bytes.
    Outer None denotes a nil receiver panic, inner None a nil slice. -/
def bufferBytes (receiver : Option BufferTailState) : Option (Option BufferTailView) :=
  receiver.map fun state => state.stored.map fun stored => ⟨stored.drop state.offset, state.spare⟩

theorem bufferBytes_view (stored spare : List (BitVec 8)) (offset : Nat) :
    bufferBytes (some ⟨some stored, spare, offset⟩) =
      some (some ⟨stored.drop offset, spare⟩) := rfl

theorem bufferBytes_consumed (stored spare : List (BitVec 8)) :
    bufferBytes (some ⟨some stored, spare, stored.length⟩) = some (some ⟨[], spare⟩) := by
  have empty : stored.drop stored.length = [] := List.drop_length
  simp [bufferBytes, empty]

theorem bufferBytes_nilStorage : bufferBytes (some ⟨none, [], 0⟩) = some none := rfl

theorem bufferBytes_nil : bufferBytes none = none := rfl

end GoLean.Bytes
