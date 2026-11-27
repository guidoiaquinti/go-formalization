import GoLean.Math.TrigReduce
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [word] := fields | throw (IO.userError "expected a binary64 word")
  let (octant, remainder) := GoLean.Math.trigReduceWord (← wordInput 64 word)
  return s!"{octant} {remainder}"
