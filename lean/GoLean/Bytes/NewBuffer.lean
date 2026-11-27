import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

/-- Go 1.27.0 NewBuffer adopts the supplied slice geometry and
    initializes the cursor and previous-read operation to zero. -/
def newBuffer (input : Option (List (BitVec 8))) (spare : List (BitVec 8)) :
    BufferReadState :=
  ⟨⟨input, spare, 0⟩, 0⟩

theorem newBuffer_nil : newBuffer none [] = ⟨⟨none, [], 0⟩, 0⟩ := rfl

theorem newBuffer_adopts (bytes spare : List (BitVec 8)) :
    newBuffer (some bytes) spare = ⟨⟨some bytes, spare, 0⟩, 0⟩ := rfl

end GoLean.Bytes
