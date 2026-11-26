import GoLean.Sort.Float64SliceLess
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let i :: j :: xs := fields | throw (IO.userError "expected two indices followed by elements")
  let xs ← xs.mapM fun word => do
    return Float.ofBits (UInt64.ofBitVec (← wordInput 64 word))
  match GoLean.Sort.float64SliceLess xs i j with
  | .ok result => return toString result
  | .error index =>
    return if index < 0 then s!"runtime error: index out of range [{index}]"
      else s!"runtime error: index out of range [{index}] with length {xs.length}"
