import GoLean.Bytes.BufferWriteByte

namespace GoLean.Bytes

/-- Go 1.27.0 bytes.growSlice for successful allocations with small
    capacity requests. The input length plus count uses signed machine
    arithmetic, and extending the result reveals only zero-filled spare. -/
def growSlice (bytes : List (BitVec 8)) (cap : Nat) (count : Int) (maxInt : Nat) :
    Option (List (BitVec 8) × List (BitVec 8)) :=
  let sum := Int.ofNat bytes.length + count
  let signed := if sum > Int.ofNat maxInt then sum - 2 * (Int.ofNat maxInt + 1) else sum
  let request := max signed (Int.ofNat (2 * cap))
  if request == 0 then none
  else some (bytes, List.replicate (smallByteCapacity request.toNat - bytes.length) 0)

theorem growSlice_empty : growSlice [] 0 0 2147483647 = none := rfl

theorem growSlice_doubles : growSlice [1, 255] 6 (-4) 2147483647 =
    some ([1, 255], List.replicate 14 0) := rfl

theorem growSlice_preserves (bytes : List (BitVec 8)) (cap : Nat) (count : Int) (maxInt : Nat)
    (result : List (BitVec 8) × List (BitVec 8))
    (success : growSlice bytes cap count maxInt = some result) : result.1 = bytes := by
  dsimp only [growSlice] at success
  split at success <;> split at success <;> first | contradiction | (cases success; rfl)

end GoLean.Bytes
