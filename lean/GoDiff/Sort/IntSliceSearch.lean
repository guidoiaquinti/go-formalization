import GoLean.Sort.IntSliceSearch
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let key :: xs := fields | throw (IO.userError "expected search key and slice elements")
  return toString (GoLean.Sort.intSliceSearch xs key)
