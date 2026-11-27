import GoLean.Math.Remainder
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [x, y] := fields | throw (IO.userError "expected two binary64 words")
  match GoLean.Math.remainderWords (← wordInput 64 x) (← wordInput 64 y) with
  | none => return "NaN"
  | some word => return toString word
