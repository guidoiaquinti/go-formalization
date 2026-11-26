import GoLean.Sort.Search
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [n, cut] := fields | throw (IO.userError "expected length and threshold")
  return toString (GoLean.Sort.search n fun i => decide (cut ≤ Int.ofNat i))
