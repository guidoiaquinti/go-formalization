import GoLean.Bytes.BufferReadSlice

namespace GoLean.Bytes

structure BufferReadStringResult where
  bytes : List (BitVec 8)
  eof : Bool
  state : BufferReadState

/-- Go 1.27.0 Buffer.ReadString snapshots the delimiter read as raw
    string bytes and preserves the helper's error and receiver transition. -/
def bufferReadString (receiver : Option BufferReadState) (delimiter : BitVec 8) :
    Option BufferReadStringResult :=
  (bufferReadSlice receiver delimiter).map fun result =>
    ⟨(result.view.map BufferTailView.bytes).getD [], result.eof, result.state⟩

theorem bufferReadString_nil (delimiter : BitVec 8) :
    bufferReadString none delimiter = none := rfl

theorem bufferReadString_empty (delimiter : BitVec 8) (flag : Int) :
    bufferReadString (some ⟨⟨none, [], 0⟩, flag⟩) delimiter =
      some ⟨[], true, ⟨⟨none, [], 0⟩, -1⟩⟩ := rfl

theorem bufferReadString_first :
    bufferReadString (some ⟨⟨some [0, 255, 9, 255], [7], 1⟩, 2⟩) 255 =
      some ⟨[255], false, ⟨⟨some [0, 255, 9, 255], [7], 2⟩, -1⟩⟩ := rfl

theorem bufferReadString_missing :
    bufferReadString (some ⟨⟨some [0, 255, 128], [7], 1⟩, 2⟩) 0 =
      some ⟨[255, 128], true, ⟨⟨some [0, 255, 128], [7], 3⟩, -1⟩⟩ := rfl

end GoLean.Bytes
