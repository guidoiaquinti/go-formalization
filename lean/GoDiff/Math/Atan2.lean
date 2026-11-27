import GoLean.Math.Atan2
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, y, x] := fields | throw (IO.userError "expected profile and two binary64 words")
  let result := GoLean.Math.atan2Word (profile == 1) (← wordInput 64 y) (← wordInput 64 x)
  return result.map toString |>.getD "NaN"
