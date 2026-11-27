import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

structure BufferReadSliceResult where
  view : Option BufferTailView
  eof : Bool
  state : BufferReadState

/-- Go 1.27.0 Buffer.readSlice includes the first raw delimiter and
    sets the ordinary read flag even if no bytes are returned. -/
def bufferReadSlice (receiver : Option BufferReadState) (delimiter : BitVec 8) :
    Option BufferReadSliceResult :=
  receiver.map fun state =>
    let tail := state.tail
    let unread := (tail.stored.getD []).drop tail.offset
    let index := unread.idxOf delimiter
    let eof := index == unread.length
    let count := min (index + 1) unread.length
    let view := tail.stored.map fun _ =>
      ⟨unread.take count, unread.drop count ++ tail.spare⟩
    ⟨view, eof, ⟨{tail with offset := tail.offset + count}, -1⟩⟩

theorem bufferReadSlice_nil (delimiter : BitVec 8) :
    bufferReadSlice none delimiter = none := rfl

theorem bufferReadSlice_empty (delimiter : BitVec 8) (flag : Int) :
    bufferReadSlice (some ⟨⟨none, [], 0⟩, flag⟩) delimiter =
      some ⟨none, true, ⟨⟨none, [], 0⟩, -1⟩⟩ := rfl

theorem bufferReadSlice_first :
    bufferReadSlice (some ⟨⟨some [0, 255, 9, 255], [7], 1⟩, 2⟩) 255 =
      some ⟨some ⟨[255], [9, 255, 7]⟩, false, ⟨⟨some [0, 255, 9, 255], [7], 2⟩, -1⟩⟩ := rfl

theorem bufferReadSlice_missing :
    bufferReadSlice (some ⟨⟨some [0, 22, 33], [7], 1⟩, 2⟩) 0 =
      some ⟨some ⟨[22, 33], [7]⟩, true, ⟨⟨some [0, 22, 33], [7], 3⟩, -1⟩⟩ := rfl

end GoLean.Bytes
