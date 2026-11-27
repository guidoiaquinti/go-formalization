import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

structure BufferReadResult where
  count : Nat
  eof : Bool
  destination : List (BitVec 8)
  state : BufferReadState

/-- Go 1.27.0 Buffer.Read with an independent destination slice. -/
def bufferRead (receiver : Option BufferReadState) (destination : List (BitVec 8)) :
    Option BufferReadResult :=
  receiver.map fun state =>
    let tail := state.tail
    let unread := (tail.stored.getD []).drop tail.offset
    match unread with
    | [] =>
      ⟨0, !destination.isEmpty, destination, (bufferReset (some state)).getD state⟩
    | byte :: rest =>
      let copied := (byte :: rest).take destination.length
      let count := copied.length
      let dst := copied ++ destination.drop count
      ⟨count, false, dst, ⟨{tail with offset := tail.offset + count}, if count == 0 then 0 else -1⟩⟩

theorem bufferRead_nil (destination : List (BitVec 8)) : bufferRead none destination = none := rfl

theorem bufferRead_emptyReset :
    bufferRead (some ⟨⟨some [11], [12], 1⟩, 2⟩) [8] =
      some ⟨0, true, [8], ⟨⟨some [], [11, 12], 0⟩, 0⟩⟩ := rfl

theorem bufferRead_zeroClears :
    bufferRead (some ⟨⟨some [11], [12], 0⟩, 2⟩) [] =
      some ⟨0, false, [], ⟨⟨some [11], [12], 0⟩, 0⟩⟩ := rfl

theorem bufferRead_partialDestination :
    bufferRead (some ⟨⟨some [11, 22], [12], 1⟩, 2⟩) [8, 9] =
      some ⟨1, false, [22, 9], ⟨⟨some [11, 22], [12], 2⟩, -1⟩⟩ := rfl

end GoLean.Bytes
