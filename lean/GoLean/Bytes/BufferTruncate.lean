import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/buffer.go:106–118. None denotes a nil pointer
    panic. A false result denotes a bounds panic with the last read
    already invalidated. The zero case delegates to Reset. -/
def bufferTruncate (receiver : Option BufferReadState) (count : Int) :
    Option (Bool × BufferReadState) :=
  if count = 0 then (bufferReset receiver).map (fun state => (true, state))
  else receiver.map fun state =>
    let invalidated := { state with lastRead := 0 }
    let stored := state.tail.stored.getD []
    if count < 0 ∨ count > Int.ofNat (stored.length - state.tail.offset) then
      (false, invalidated)
    else
      let endpoint := state.tail.offset + count.toNat
      (true, { invalidated with tail := { state.tail with
        stored := state.tail.stored.map (·.take endpoint)
        spare := stored.drop endpoint ++ state.tail.spare } })

theorem bufferTruncate_zero (receiver : Option BufferReadState) :
    bufferTruncate receiver 0 = (bufferReset receiver).map (fun state => (true, state)) := by
  simp [bufferTruncate]

theorem bufferTruncate_nil (count : Int) : bufferTruncate none count = none := by
  by_cases zero : count = 0 <;> simp [bufferTruncate, zero, bufferReset]

theorem bufferTruncate_bounds (state : BufferReadState) (count : Int)
    (nonzero : count ≠ 0)
    (outside : count < 0 ∨ count > Int.ofNat ((state.tail.stored.getD []).length - state.tail.offset)) :
    bufferTruncate (some state) count = some (false, { state with lastRead := 0 }) := by
  simp only [bufferTruncate, ite_eq_right nonzero, Option.map_some]
  rw [ite_eq_left outside]

theorem bufferTruncate_positive (state : BufferReadState) (count : Int)
    (positive : 0 < count)
    (fits : count ≤ Int.ofNat ((state.tail.stored.getD []).length - state.tail.offset)) :
    bufferTruncate (some state) count =
      let endpoint := state.tail.offset + count.toNat
      some (true, ⟨⟨state.tail.stored.map (·.take endpoint),
        (state.tail.stored.getD []).drop endpoint ++ state.tail.spare, state.tail.offset⟩, 0⟩) := by
  cases state
  simp_all [bufferTruncate, Int.ne_of_gt positive, Int.not_lt.mpr (Int.le_of_lt positive)]

end GoLean.Bytes
