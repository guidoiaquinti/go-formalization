import GoLean.Math.Lgamma
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let (result, sign) := GoLean.Math.lgammaWord (profile == 2) (profile == 1) (← wordInput 64 word)
  return s!"{(result.map toString).getD "NaN"} {sign}"
