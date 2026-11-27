import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

structure BufferNextResult where
  boundsPanic : Bool
  state : BufferReadState
  view : Option BufferTailView

/-- Go 1.27.0 src/bytes/buffer.go:353–371. Next preserves the underlying
    data, returns a cursor-relative aliased view, and advances the cursor
    by the clamped count. The read flag is invalidated before bounds checks. -/
def bufferNext (receiver : Option BufferReadState) (count : Int) : Option BufferNextResult :=
  receiver.map fun state =>
    let invalidated := { state with lastRead := 0 }
    if count < 0 then ⟨true, invalidated, none⟩
    else
      let tail := state.tail
      let unread := (tail.stored.getD []).drop tail.offset
      let size := min count.toNat unread.length
      let view := tail.stored.map fun _ => ⟨unread.take size, unread.drop size ++ tail.spare⟩
      let advanced : BufferReadState :=
        ⟨{tail with offset := tail.offset + size}, if size = 0 then 0 else -1⟩
      ⟨false, advanced, view⟩

theorem bufferNext_nil (count : Int) : bufferNext none count = none := rfl

theorem bufferNext_negative (state : BufferReadState) (count : Int) (negative : count < 0) :
    bufferNext (some state) count = some ⟨true, {state with lastRead := 0}, none⟩ := by
  simp [bufferNext, negative]

theorem bufferNext_advance (state : BufferReadState) (count : Int) (nonnegative : 0 ≤ count) :
    (bufferNext (some state) count).map (·.state.tail.offset) =
      some (state.tail.offset +
        min count.toNat ((state.tail.stored.getD []).drop state.tail.offset).length) := by
  simp [bufferNext, Int.not_lt.mpr nonnegative]

theorem bufferNext_readFlag (state : BufferReadState) (count : Int) (nonnegative : 0 ≤ count) :
    (bufferNext (some state) count).map (·.state.lastRead) =
      some (if min count.toNat ((state.tail.stored.getD []).drop state.tail.offset).length = 0
            then 0 else -1) := by
  simp [bufferNext, Int.not_lt.mpr nonnegative]

end GoLean.Bytes
