import GoLean.Math.Pow
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, x, y] := fields | throw (IO.userError "expected profile and two binary64 words")
  match GoLean.Math.powWords (profile == 2) (profile == 1) (← wordInput 64 x) (← wordInput 64 y) with
  | some output => return toString output
  | none => return "NaN"
