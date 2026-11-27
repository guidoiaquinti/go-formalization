import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:431–440. The Boolean records the
    errUnreadByte sentinel; None denotes a nil receiver panic. -/
def bufferUnreadByte (receiver : Option BufferReadState) : Option (Bool × BufferReadState) :=
  receiver.map fun state =>
    if state.lastRead == 0 then
      (true, state)
    else
      let tail := state.tail
      (false, ⟨{tail with offset := tail.offset - 1}, 0⟩)

theorem bufferUnreadByte_nil : bufferUnreadByte none = none := rfl

theorem bufferUnreadByte_invalid (state : BufferReadState) (invalid : state.lastRead = 0) :
    bufferUnreadByte (some state) = some (true, state) := by
  simp [bufferUnreadByte, invalid]

theorem bufferUnreadByte_rewind (state : BufferReadState) (valid : state.lastRead ≠ 0) :
    bufferUnreadByte (some state) =
      some (false, ⟨{state.tail with offset := state.tail.offset - 1}, 0⟩) := by
  simp [bufferUnreadByte, valid]

theorem bufferUnreadByte_zeroOffset (state : BufferReadState) (valid : state.lastRead ≠ 0)
    (atStart : state.tail.offset = 0) :
    bufferUnreadByte (some state) = some (false, ⟨state.tail, 0⟩) := by
  rw [bufferUnreadByte_rewind state valid]
  cases state with
  | mk tail lastRead =>
    cases tail
    simp_all

end GoLean.Bytes
