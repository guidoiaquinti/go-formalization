import GoLean.Bytes.NewBuffer
import GoLean.Bytes.BufferWriteByte

namespace GoLean.Bytes

/-- Go 1.27.0 NewBufferString copies raw bytes into a nonnil slice.
    This models pinned small-object rounding through 512 bytes. -/
def newBufferString (input : List (BitVec 8)) : BufferReadState :=
  let capacity := if input.isEmpty then 0 else smallByteCapacity input.length
  newBuffer (some input) (List.replicate (capacity - input.length) 0)

theorem newBufferString_empty : newBufferString [] = ⟨⟨some [], [], 0⟩, 0⟩ := rfl

theorem newBufferString_bytes (input : List (BitVec 8)) :
    (newBufferString input).tail.stored = some input := rfl

theorem newBufferString_small :
    newBufferString [0, 255] = ⟨⟨some [0, 255], List.replicate 6 0, 0⟩, 0⟩ := rfl

end GoLean.Bytes
