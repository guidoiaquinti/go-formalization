import GoLean.Compress.Flate.HCode

namespace GoLean.Compress.Flate

structure HuffmanGeometry where
  length : Nat
  capacity : Nat
  zeroWords : Bool
  deriving DecidableEq

/-- The length calculation is on uint32(size-1). Successful allocations
    are modeled by their length, capacity, and zero initialization. -/
def newHuffmanEncoder (width : Nat) (size : Int) : Option HuffmanGeometry :=
  let predecessor := (BitVec.ofInt 32 (size - 1)).toNat
  let exponent := if predecessor = 0 then 0 else predecessor.log2 + 1
  let capacity := 2 ^ exponent
  if size < 0 || capacity ≥ 2 ^ (width - 1) || capacity < size.toNat then none
  else some ⟨size.toNat, capacity, true⟩

theorem newHuffmanEncoder_negative (width : Nat) (size : Int) (negative : size < 0) :
    newHuffmanEncoder width size = none := by simp [newHuffmanEncoder, negative]

theorem newHuffmanEncoder_boundary :
    newHuffmanEncoder 64 256 = some ⟨256, 256, true⟩ ∧
      newHuffmanEncoder 64 257 = some ⟨257, 512, true⟩ := by decide

/-- Both uint16 fields of the sentinel are maximum values. -/
def maxNode : BitVec 16 × BitVec 16 := (65535, 65535)

theorem maxNode_fields : maxNode.1.toNat = 65535 ∧ maxNode.2.toNat = 65535 := by decide

end GoLean.Compress.Flate
