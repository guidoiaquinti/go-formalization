import GoLean.Sort.Float64SliceSwap
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let i :: j :: xs := fields | throw (IO.userError "expected two indices followed by elements")
  let xs ← xs.mapM fun word => do
    return UInt64.ofBitVec (← wordInput 64 word)
  let (status, output) := match GoLean.Sort.float64SliceSwap xs i j with
    | .ok result => ("ok", result)
    | .error index =>
      (if index < 0 then s!"runtime error: index out of range [{index}]"
        else s!"runtime error: index out of range [{index}] with length {xs.length}", xs)
  return s!"{status} [{String.intercalate " " (output.map toString)}]"
