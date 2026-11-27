import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

inductive BufferPeekResult where
  | nilReceiver
  | bounds (state : BufferReadState)
  | view (value : Option BufferTailView) (eof : Bool) (state : BufferReadState)

/-- Go 1.27.0 Buffer.Peek with legal initial geometry. The returned
    view starts at the unread offset and retains the original capacity.
    Neither successful observation nor bounds failure updates the buffer. -/
def bufferPeek (receiver : Option BufferReadState) (count : Int) : BufferPeekResult :=
  match receiver with
  | none => .nilReceiver
  | some state =>
    if count < 0 then .bounds state
    else
      let tail := state.tail
      let unread := (tail.stored.getD []).drop tail.offset
      let eof := unread.length < count.toNat
      let view := tail.stored.map fun _ =>
        ⟨unread.take count.toNat, unread.drop count.toNat ++ tail.spare⟩
      .view view eof state

theorem bufferPeek_nil (count : Int) : bufferPeek none count = .nilReceiver := rfl

theorem bufferPeek_negative (state : BufferReadState) (count : Int) (negative : count < 0) :
    bufferPeek (some state) count = .bounds state := by
  simp only [bufferPeek, ite_eq_left negative]

theorem bufferPeek_zeroNil (flag : Int) :
    bufferPeek (some ⟨⟨none, [], 0⟩, flag⟩) 0 = .view none false ⟨⟨none, [], 0⟩, flag⟩ := rfl

theorem bufferPeek_retainsCapacity :
    bufferPeek (some ⟨⟨some [11, 22, 33], [44], 1⟩, 2⟩) 1 =
      .view (some ⟨[22], [33, 44]⟩) false ⟨⟨some [11, 22, 33], [44], 1⟩, 2⟩ := rfl

theorem bufferPeek_eof :
    bufferPeek (some ⟨⟨some [11, 22, 33], [44], 1⟩, 2⟩) 4 =
      .view (some ⟨[22, 33], [44]⟩) true ⟨⟨some [11, 22, 33], [44], 1⟩, 2⟩ := rfl

end GoLean.Bytes
