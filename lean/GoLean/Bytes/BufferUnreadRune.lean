import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:414–423. The Boolean records a
    fresh UnreadRune error; None denotes a nil receiver panic. -/
def bufferUnreadRune (receiver : Option BufferReadState) : Option (Bool × BufferReadState) :=
  receiver.map fun state =>
    if state.lastRead ≤ 0 then
      (true, state)
    else
      let width := state.lastRead.toNat
      let tail := state.tail
      let offset := if tail.offset ≥ width then tail.offset - width else tail.offset
      (false, ⟨{tail with offset}, 0⟩)

theorem bufferUnreadRune_nil : bufferUnreadRune none = none := rfl

theorem bufferUnreadRune_reject (state : BufferReadState) (invalid : state.lastRead ≤ 0) :
    bufferUnreadRune (some state) = some (true, state) := by
  simp [bufferUnreadRune, invalid]

theorem bufferUnreadRune_rewind (state : BufferReadState) (valid : ¬state.lastRead ≤ 0)
    (fits : state.lastRead.toNat ≤ state.tail.offset) :
    bufferUnreadRune (some state) =
      some (false, ⟨{state.tail with offset := state.tail.offset - state.lastRead.toNat}, 0⟩) := by
  simp only [bufferUnreadRune, Option.map_some, ite_eq_right valid, ite_eq_left fits]

theorem bufferUnreadRune_shortOffset (state : BufferReadState) (valid : ¬state.lastRead ≤ 0)
    (short : ¬state.lastRead.toNat ≤ state.tail.offset) :
    bufferUnreadRune (some state) = some (false, ⟨state.tail, 0⟩) := by
  simp only [bufferUnreadRune, Option.map_some, ite_eq_right valid, ite_eq_right short]

end GoLean.Bytes
