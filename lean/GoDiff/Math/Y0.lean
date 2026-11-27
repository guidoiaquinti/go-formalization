import GoLean.Math.Y0
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let result := GoLean.Math.y0Word (profile == 2) (profile == 1) (← wordInput 64 word)
  return (result.map toString).getD "NaN"
