import GoLean.Encoding.Binary.EndianRead
import GoLean.Encoding.Binary.EndianPut
import Init.Omega

namespace GoLean.Compress.Flate

open GoLean.Encoding.Binary

def loadLE8 (backing : List (BitVec 8)) (length : Nat) (start : Int) : Option Nat :=
  if start < 0 || length ≤ start.toNat then none else backing[start.toNat]?.map BitVec.toNat

private def indexValue (width : Nat) (signed : Bool) (word : Int) : Int :=
  let value := BitVec.ofInt width word
  if signed then value.toInt else Int.ofNat value.toNat

private def loadLE (size width : Nat) (signed : Bool) (backing : List (BitVec 8))
    (start : Int) : Option Nat :=
  let stop := indexValue width signed (start + Int.ofNat size)
  if start < 0 || stop < start || backing.length < stop.toNat then none
  else readWord false size ((backing.take stop.toNat).drop start.toNat)

def loadLE32 := loadLE 4

def loadLE64 := loadLE 8

/-- Go checks visible length before the first write. -/
def storeLE64 (backing : List (BitVec 8)) (length : Nat) (word : BitVec 64) :
    Bool × List (BitVec 8) :=
  if length < 8 then (false, backing) else littlePutUint64 word backing

theorem loadLE8_short (b : List (BitVec 8)) : loadLE8 b 0 0 = none := rfl

theorem loadLE32_capacity :
    loadLE32 8 false [1, 2, 3, 4] 0 = some 67305985 := by decide

theorem loadLE64_wrappedIndex (b : List (BitVec 8)) :
    loadLE64 8 false b 252 = none := rfl

theorem storeLE64_short (b : List (BitVec 8)) (n : Nat) (word : BitVec 64) (short : n < 8) :
    storeLE64 b n word = (false, b) := by simp [storeLE64, short]

end GoLean.Compress.Flate
