import GoLean.Sort.StringSliceLess
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let i :: j :: xs := fields | throw (IO.userError "expected two indices followed by elements")
  let xs ← if xs.isEmpty then pure [] else do
    let chunks := xs.splitOn (-1)
    let [] :: strings := chunks | throw (IO.userError "expected byte strings introduced by -1")
    strings.mapM fun bytes => do
      unless bytes.all (fun x => 0 ≤ x && x < 256) do
        throw (IO.userError "invalid string byte")
      pure (bytes.map fun x => x.toNat.toUInt8)
  match GoLean.Sort.stringSliceLess xs i j with
  | .ok result => return toString result
  | .error index =>
    return if index < 0 then s!"runtime error: index out of range [{index}]"
      else s!"runtime error: index out of range [{index}] with length {xs.length}"
