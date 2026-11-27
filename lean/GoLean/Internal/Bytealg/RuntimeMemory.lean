import GoLean.Bytes.BufferWriteByte

namespace GoLean.Internal.Bytealg

structure ByteAllocation where
  length : Nat
  capacity : Nat
  nonNil : Bool

/-- Go 1.27.0 runtime/slice.go:432–437. Observe slice geometry after
    roundupsize for small successful requests. No uninitialized byte is modeled. -/
def makeNoZero (length : Int) (maxAllocation : Nat) : Option ByteAllocation :=
  if length < 0 || Int.ofNat maxAllocation < length then none
  else if length = 0 then some ⟨0, 0, true⟩
  else some ⟨length.toNat, GoLean.Bytes.smallByteCapacity length.toNat, true⟩

theorem makeNoZero_zero (maxAllocation : Nat) :
    makeNoZero 0 maxAllocation = some ⟨0, 0, true⟩ := by
  simp [makeNoZero]

theorem makeNoZero_negative (length : Int) (maxAllocation : Nat) (negative : length < 0) :
    makeNoZero length maxAllocation = none := by simp [makeNoZero, negative]

theorem makeNoZero_roundsNine :
    makeNoZero 9 512 = some ⟨9, 16, true⟩ := rfl

/-- The fixed-size runtime ABI compares only the indicated bytes.
    Callers establish that the windows are readable when size is positive. -/
def runtimeMemequal (size : Nat) (left right : List (BitVec 8)) : Bool :=
  decide (left.take size = right.take size)

theorem runtimeMemequal_spec (size : Nat) (left right : List (BitVec 8)) :
    runtimeMemequal size left right = true ↔ left.take size = right.take size := by
  simp [runtimeMemequal]

theorem runtimeMemequal_zero (left right : List (BitVec 8)) :
    runtimeMemequal 0 left right = true := rfl

end GoLean.Internal.Bytealg
