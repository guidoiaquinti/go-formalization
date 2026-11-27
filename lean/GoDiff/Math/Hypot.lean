import GoLean.Math.Hypot
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, x, y] := fields | throw (IO.userError "expected profile and two binary64 words")
  match GoLean.Math.hypotWords (profile == 1) (← wordInput 64 x) (← wordInput 64 y) with
  | none => return "NaN"
  | some word => return toString word
