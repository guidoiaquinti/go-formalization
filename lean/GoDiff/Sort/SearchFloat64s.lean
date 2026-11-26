import GoLean.Sort.SearchFloat64s
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let values ← fields.mapM fun word => do
    return Float.ofBits (UInt64.ofBitVec (← wordInput 64 word))
  let key :: xs := values | throw (IO.userError "expected search key and slice elements")
  return toString (GoLean.Sort.searchFloat64s xs key)
