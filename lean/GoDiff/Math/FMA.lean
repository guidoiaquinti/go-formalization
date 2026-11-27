import GoLean.Math.FMA
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [x, y, z] := fields | throw (IO.userError "expected three binary64 words")
  match GoLean.Math.fmaWords (← wordInput 64 x) (← wordInput 64 y) (← wordInput 64 z) with
  | none => return "NaN"
  | some word => return toString word
