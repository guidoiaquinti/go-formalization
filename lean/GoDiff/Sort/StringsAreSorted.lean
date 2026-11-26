import GoLean.Sort.StringsAreSorted
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: xs := fields | throw (IO.userError "expected marker followed by strings")
  let xs ← if xs.isEmpty then pure [] else do
    let chunks := xs.splitOn (-1)
    let [] :: strings := chunks | throw (IO.userError "expected byte strings introduced by -1")
    strings.mapM fun bytes => do
      unless bytes.all (fun x => 0 ≤ x && x < 256) do
        throw (IO.userError "invalid string byte")
      pure (bytes.map fun x => x.toNat.toUInt8)
  return toString (GoLean.Sort.stringsAreSorted xs)
