import GoLean.Sort.IntSliceSwap
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let i :: j :: xs := fields | throw (IO.userError "expected two indices followed by elements")
  let (status, output) := match GoLean.Sort.intSliceSwap xs i j with
    | .ok result => ("ok", result)
    | .error index =>
      (if index < 0 then s!"runtime error: index out of range [{index}]"
        else s!"runtime error: index out of range [{index}] with length {xs.length}", xs)
  return s!"{status} [{String.intercalate " " (output.map toString)}]"
