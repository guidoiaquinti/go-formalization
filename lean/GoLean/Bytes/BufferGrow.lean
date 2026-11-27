import GoLean.Bytes.BufferWriteByte

namespace GoLean.Bytes

inductive BufferGrowResult where
  | negative
  | nilReceiver
  | tooLarge (state : BufferReadState)
  | success (result : BufferWriteByteResult)

private def growInPlace (state : BufferReadState) : BufferGrowResult :=
  let tail := state.tail
  .success ⟨state, !(tail.stored.getD []).isEmpty || !tail.spare.isEmpty⟩

private def prepareGrow (state : BufferReadState) : BufferReadState :=
  let tail := state.tail
  if tail.offset > 0 && tail.offset == (tail.stored.getD []).length then
    (bufferReset (some state)).getD state
  else state

/-- Go 1.27.0 Buffer.Grow. The negative-count check precedes receiver
    access. The helper reserves space, then Grow restores the content length
    without clearing an existing read flag unless the helper resets first. -/
def bufferGrow (receiver : Option BufferReadState) (count : Int) (maxInt : Nat) :
    BufferGrowResult :=
  if count < 0 then .negative
  else
    match receiver with
    | none => .nilReceiver
    | some original =>
      let state := prepareGrow original
      let n := count.toNat
      let tail := state.tail
      let stored := tail.stored.getD []
      let capacity := stored.length + tail.spare.length
      let unread := stored.drop tail.offset
      if n ≤ tail.spare.length then growInPlace state
      else if tail.stored.isNone && n ≤ 64 then
        .success ⟨⟨⟨some [], List.replicate 64 0, 0⟩, state.lastRead⟩, false⟩
      else if unread.length + n ≤ capacity / 2 then
        .success ⟨⟨⟨some unread, (stored ++ tail.spare).drop unread.length, 0⟩, state.lastRead⟩, true⟩
      else if 2 * capacity + n > maxInt then .tooLarge state
      else
        let requested := max (stored.length + n) (2 * (capacity - tail.offset))
        let spare := List.replicate (smallByteCapacity requested - unread.length) 0
        .success ⟨⟨⟨some unread, spare, 0⟩, state.lastRead⟩, false⟩

theorem bufferGrow_negative (receiver : Option BufferReadState) (count : Int) (maxInt : Nat)
    (negative : count < 0) : bufferGrow receiver count maxInt = .negative := by
  simp [bufferGrow, negative]

theorem bufferGrow_nil (count : Int) (maxInt : Nat) (nonnegative : ¬count < 0) :
    bufferGrow none count maxInt = .nilReceiver := by
  simp [bufferGrow, nonnegative]

theorem bufferGrow_exhaustedZero :
    bufferGrow (some ⟨⟨some [1, 2], [3], 2⟩, 1⟩) 0 2147483647 =
      .success ⟨⟨⟨some [], [1, 2, 3], 0⟩, 0⟩, true⟩ := rfl

theorem bufferGrow_preservesRead :
    bufferGrow (some ⟨⟨some [1, 2, 3, 4, 5, 6], [], 4⟩, 2⟩) 1 2147483647 =
      .success ⟨⟨⟨some [5, 6], [3, 4, 5, 6], 0⟩, 2⟩, true⟩ := rfl

theorem bufferGrow_tooLarge (maxInt : Nat) (large : 2 < maxInt) :
    bufferGrow (some ⟨⟨some [1, 2], [], 2⟩, 1⟩) (Int.ofNat maxInt) maxInt =
      .tooLarge ⟨⟨some [], [1, 2], 0⟩, 0⟩ := by
  have aboveOne : ¬maxInt ≤ 1 := by omega
  have aboveTwo : ¬maxInt ≤ 2 := by omega
  simp [bufferGrow, prepareGrow, bufferReset, aboveOne, aboveTwo]

end GoLean.Bytes
