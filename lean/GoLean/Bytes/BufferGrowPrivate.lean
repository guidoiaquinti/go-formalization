import GoLean.Bytes.BufferGrow
import GoLean.Bytes.BufferTryGrowByReslice

namespace GoLean.Bytes

inductive BufferGrowPrivateResult where
  | nilReceiver
  | bounds (state : BufferReadState)
  | tooLarge (state : BufferReadState)
  | returned (index : Nat) (result : BufferWriteByteResult)

private def privateGrowPrep (state : BufferReadState) : BufferReadState :=
  if state.tail.offset > 0 && state.tail.offset == (state.tail.stored.getD []).length then
    (bufferReset (some state)).getD state
  else state

private def revealReservation (result : BufferWriteByteResult) (count : Nat) :
    BufferGrowPrivateResult :=
  let state := result.state
  let tail := state.tail
  let stored := tail.stored.getD []
  let next := {tail with
    stored := tail.stored.map (fun _ => stored ++ tail.spare.take count)
    spare := tail.spare.drop count}
  .returned stored.length {result with state := {state with tail := next}}

/-- Go 1.27.0 Buffer.grow, including the signed direct reslice path.
    For a slow successful grow, the public Grow model gives the same backing
    after hiding the reserved bytes; reveal those bytes in the stored slice. -/
def bufferGrowPrivate (receiver : Option BufferReadState) (count : Int) (maxInt : Nat) :
    BufferGrowPrivateResult :=
  match receiver with
  | none => .nilReceiver
  | some original =>
    let state := privateGrowPrep original
    match bufferTryGrowByReslice (some state) count with
    | .nilReceiver => .nilReceiver
    | .bounds prior => .bounds prior
    | .returned index true resized =>
      let reused := !(state.tail.stored.getD []).isEmpty || !state.tail.spare.isEmpty
      .returned index ⟨resized, reused⟩
    | .returned _ false _ =>
      match bufferGrow (some state) count maxInt with
      | .success result => revealReservation result count.toNat
      | .tooLarge prior => .tooLarge prior
      | .negative => .bounds state
      | .nilReceiver => .nilReceiver

theorem bufferGrowPrivate_nil (count : Int) (maxInt : Nat) :
    bufferGrowPrivate none count maxInt = .nilReceiver := rfl

theorem bufferGrowPrivate_reslice :
    bufferGrowPrivate (some ⟨⟨some [1, 2], [3, 4], 1⟩, 3⟩) 1 2147483647 =
      .returned 2 ⟨⟨⟨some [1, 2, 3], [4], 1⟩, 3⟩, true⟩ := rfl

theorem bufferGrowPrivate_compact :
    bufferGrowPrivate (some ⟨⟨some [1, 2, 3, 4, 5, 6], [], 4⟩, 2⟩) 1 2147483647 =
      .returned 2 ⟨⟨⟨some [5, 6, 3], [4, 5, 6], 0⟩, 2⟩, true⟩ := rfl

theorem bufferGrowPrivate_resetThenBounds :
    bufferGrowPrivate (some ⟨⟨some [1, 2], [3], 2⟩, 1⟩) (-1) 2147483647 =
      .bounds ⟨⟨some [], [1, 2, 3], 0⟩, 0⟩ := rfl

end GoLean.Bytes
