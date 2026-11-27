import GoLean.Math.Asin
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let result := GoLean.Math.asinWord (profile == 2) (profile != 0) (← wordInput 64 word)
  return result.map toString |>.getD "NaN"
