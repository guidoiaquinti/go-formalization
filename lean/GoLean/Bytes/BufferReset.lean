import GoLean.Bytes.BufferAvailableBuffer

namespace GoLean.Bytes

structure BufferReadState where
  tail : BufferTailState
  lastRead : Int

/-- Go 1.27.0 src/bytes/buffer.go:121–128. Reset empties the visible
    slice, moves the cursor to the beginning and marks the last operation
    invalid for unread. All original backing bytes are retained as spare. -/
def bufferReset (receiver : Option BufferReadState) : Option BufferReadState :=
  receiver.map fun state =>
    let stored := state.tail.stored
    ⟨⟨stored.map (fun _ => []), stored.getD [] ++ state.tail.spare, 0⟩, 0⟩

theorem bufferReset_state (stored spare : List (BitVec 8)) (offset : Nat) (lastRead : Int) :
    bufferReset (some ⟨⟨some stored, spare, offset⟩, lastRead⟩) =
      some ⟨⟨some [], stored ++ spare, 0⟩, 0⟩ := rfl

theorem bufferReset_nilStorage (lastRead : Int) :
    bufferReset (some ⟨⟨none, [], 0⟩, lastRead⟩) =
      some ⟨⟨none, [], 0⟩, 0⟩ := rfl

theorem bufferReset_nil : bufferReset none = none := rfl

end GoLean.Bytes
