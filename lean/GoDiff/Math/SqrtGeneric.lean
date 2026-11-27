import GoLean.Math.Sqrt
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [raw] := fields | throw (IO.userError "expected a binary64 word")
  let word ← wordInput 64 raw
  return (GoLean.Math.sqrtWord word |>.map toString).getD "NaN"
