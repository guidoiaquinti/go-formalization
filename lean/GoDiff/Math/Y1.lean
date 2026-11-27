import GoLean.Math.Y1
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let result := GoLean.Math.y1Word (profile == 2) (profile == 1) (← wordInput 64 word)
  return (result.map toString).getD "NaN"
