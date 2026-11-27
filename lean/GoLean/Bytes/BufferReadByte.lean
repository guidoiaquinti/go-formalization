import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

structure BufferReadByteResult where
  byte : BitVec 8
  eof : Bool
  state : BufferReadState

/-- Go 1.27.0 src/bytes/buffer.go:93, 374–384. Return one raw byte and
    mark an ordinary read, or reset the exhausted buffer and return EOF.
    None denotes a nil pointer panic. -/
def bufferReadByte (receiver : Option BufferReadState) : Option BufferReadByteResult := do
  let state ← receiver
  let tail := state.tail
  match (tail.stored.getD [])[tail.offset]? with
  | some byte =>
    return ⟨byte, false, ⟨{tail with offset := tail.offset + 1}, -1⟩⟩
  | none =>
    let reset ← bufferReset (some state)
    return ⟨0, true, reset⟩

theorem bufferReadByte_nil : bufferReadByte none = none := rfl

theorem bufferReadByte_live (state : BufferReadState) (byte : BitVec 8)
    (present : (state.tail.stored.getD [])[state.tail.offset]? = some byte) :
    bufferReadByte (some state) =
      some ⟨byte, false, ⟨{state.tail with offset := state.tail.offset + 1}, -1⟩⟩ := by
  simp [bufferReadByte, present]

theorem bufferReadByte_eof (state : BufferReadState)
    (exhausted : (state.tail.stored.getD [])[state.tail.offset]? = none) :
    bufferReadByte (some state) = (bufferReset (some state)).map (fun reset => ⟨0, true, reset⟩) := by
  simp [bufferReadByte, exhausted, bufferReset]

end GoLean.Bytes
