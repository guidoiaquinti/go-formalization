import Init.Data.BitVec.Basic

namespace GoLean.Bytes

structure BufferTailState where
  stored : Option (List (BitVec 8))
  spare : List (BitVec 8)
  offset : Nat

structure BufferTailView where
  bytes : List (BitVec 8)
  spare : List (BitVec 8)

/-- Go 1.27.0 src/bytes/buffer.go:62–66. The returned slice starts at
    the end of the stored bytes and aliases the remaining backing bytes.
    The outer None denotes a nil receiver panic, the inner None a nil view. -/
def bufferAvailableBuffer (receiver : Option BufferTailState) :
    Option (Option BufferTailView) :=
  receiver.map fun state => state.stored.map fun _ => ⟨[], state.spare⟩

theorem bufferAvailableBuffer_view (stored spare : List (BitVec 8)) (offset : Nat) :
    bufferAvailableBuffer (some ⟨some stored, spare, offset⟩) = some (some ⟨[], spare⟩) := rfl

theorem bufferAvailableBuffer_nilStorage (offset : Nat) :
    bufferAvailableBuffer (some ⟨none, [], offset⟩) = some none := rfl

theorem bufferAvailableBuffer_nil : bufferAvailableBuffer none = none := rfl

end GoLean.Bytes
