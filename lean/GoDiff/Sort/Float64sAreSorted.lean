import GoLean.Sort.Float64sAreSorted
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: xs := fields | throw (IO.userError "expected marker and slice elements")
  let xs ← xs.mapM fun word => do
    return Float.ofBits (UInt64.ofBitVec (← wordInput 64 word))
  return toString (GoLean.Sort.float64sAreSorted xs)
