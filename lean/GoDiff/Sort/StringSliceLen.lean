import GoLean.Sort.StringSliceLen
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [size, lo, hi] := fields | throw (IO.userError "expected size lo hi")
  unless 0 ≤ lo && lo ≤ hi && hi ≤ size do
    throw (IO.userError "invalid slice bounds")
  let backing := (List.range size.toNat).map fun i => (toString i).toUTF8
  let xs := (backing.drop lo.toNat).take (hi - lo).toNat
  return toString (GoLean.Sort.stringSliceLen xs)
