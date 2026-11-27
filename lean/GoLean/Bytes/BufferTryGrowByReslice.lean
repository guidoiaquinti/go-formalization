import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

inductive BufferResliceResult where
  | nilReceiver
  | bounds (state : BufferReadState)
  | returned (index : Nat) (ok : Bool) (state : BufferReadState)

/-- Go 1.27.0 Buffer.tryGrowByReslice. Even when a negative count
    shortens the slice, the original cursor and read flag stay intact. -/
def bufferTryGrowByReslice (receiver : Option BufferReadState) (count : Int) :
    BufferResliceResult :=
  match receiver with
  | none => .nilReceiver
  | some state =>
    let tail := state.tail
    let stored := tail.stored.getD []
    let newLength := (Int.ofNat stored.length) + count
    if count > Int.ofNat tail.spare.length then .returned 0 false state
    else if newLength < 0 then .bounds state
    else
      let backing := stored ++ tail.spare
      let length := newLength.toNat
      let updated := {tail with
        stored := tail.stored.map (fun _ => backing.take length)
        spare := backing.drop length}
      .returned stored.length true {state with tail := updated}

theorem bufferTryGrowByReslice_nil (count : Int) :
    bufferTryGrowByReslice none count = .nilReceiver := rfl

theorem bufferTryGrowByReslice_reject (state : BufferReadState) (count : Int)
    (tooLarge : count > Int.ofNat state.tail.spare.length) :
    bufferTryGrowByReslice (some state) count = .returned 0 false state := by
  simp only [bufferTryGrowByReslice, ite_eq_left tooLarge]

theorem bufferTryGrowByReslice_fit :
    bufferTryGrowByReslice (some ⟨⟨some [1, 2], [3, 4], 1⟩, 3⟩) 1 =
      .returned 2 true ⟨⟨some [1, 2, 3], [4], 1⟩, 3⟩ := rfl

theorem bufferTryGrowByReslice_shrink :
    bufferTryGrowByReslice (some ⟨⟨some [1, 2], [3, 4], 2⟩, 1⟩) (-2) =
      .returned 2 true ⟨⟨some [], [1, 2, 3, 4], 2⟩, 1⟩ := rfl

theorem bufferTryGrowByReslice_bounds :
    bufferTryGrowByReslice (some ⟨⟨some [1, 2], [3, 4], 2⟩, 1⟩) (-3) =
      .bounds ⟨⟨some [1, 2], [3, 4], 2⟩, 1⟩ := rfl

end GoLean.Bytes
