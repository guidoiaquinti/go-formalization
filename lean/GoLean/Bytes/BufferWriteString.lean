import GoLean.Bytes.BufferWrite

namespace GoLean.Bytes

/-- Go 1.27.0 Buffer.WriteString uses the same growth operations as
    Buffer.Write, with a string whose bytes remain stable during the call. -/
def bufferWriteString (receiver : Option BufferReadState) (raw : List (BitVec 8)) :
    Option (Nat × BufferWriteByteResult) :=
  bufferWrite receiver (.detached raw)

theorem bufferWriteString_detached (receiver : Option BufferReadState) (raw : List (BitVec 8)) :
    bufferWriteString receiver raw = bufferWrite receiver (.detached raw) := rfl

theorem bufferWriteString_nil (raw : List (BitVec 8)) :
    bufferWriteString none raw = none := rfl

theorem bufferWriteString_emptyNil (flag : Int) :
    bufferWriteString (some ⟨⟨none, [], 0⟩, flag⟩) [] =
      some (0, ⟨⟨⟨none, [], 0⟩, 0⟩, false⟩) := rfl

theorem bufferWriteString_appendRaw :
    bufferWriteString (some ⟨⟨some [1, 2], [3, 4, 5], 1⟩, 2⟩) [0, 255] =
      some (2, ⟨⟨⟨some [1, 2, 0, 255], [5], 1⟩, 0⟩, true⟩) := rfl

theorem bufferWriteString_compact :
    bufferWriteString (some ⟨⟨some [1, 2, 3, 4, 5, 6, 7, 8], [], 6⟩, 1⟩) [0, 255] =
      some (2, ⟨⟨⟨some [7, 8, 0, 255], [5, 6, 7, 8], 0⟩, 0⟩, true⟩) := rfl

end GoLean.Bytes
