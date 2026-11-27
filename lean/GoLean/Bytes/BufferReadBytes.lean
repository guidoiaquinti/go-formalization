import GoLean.Bytes.BufferReadSlice
import GoLean.Bytes.BufferWriteByte

namespace GoLean.Bytes

/-- Go 1.27.0 Buffer.ReadBytes detaches the returned bytes from the
    Buffer. The allocation rounding shown here covers raw requests through 512. -/
def bufferReadBytes (receiver : Option BufferReadState) (delimiter : BitVec 8) :
    Option BufferReadSliceResult :=
  (bufferReadSlice receiver delimiter).map fun result =>
    let bytes := result.view.map BufferTailView.bytes |>.getD []
    let copied := match bytes with
      | [] => none
      | head :: rest =>
        let spare := List.replicate (smallByteCapacity bytes.length - bytes.length) 0
        some ⟨head :: rest, spare⟩
    {result with view := copied}

theorem bufferReadBytes_nil (delimiter : BitVec 8) :
    bufferReadBytes none delimiter = none := rfl

theorem bufferReadBytes_empty (delimiter : BitVec 8) :
    bufferReadBytes (some ⟨⟨some [1], [2], 1⟩, 4⟩) delimiter =
      some ⟨none, true, ⟨⟨some [1], [2], 1⟩, -1⟩⟩ := rfl

theorem bufferReadBytes_detached :
    bufferReadBytes (some ⟨⟨some [1, 255, 7], [8], 1⟩, 3⟩) 255 =
      some ⟨some ⟨[255], List.replicate 7 0⟩, false, ⟨⟨some [1, 255, 7], [8], 2⟩, -1⟩⟩ := rfl

end GoLean.Bytes
