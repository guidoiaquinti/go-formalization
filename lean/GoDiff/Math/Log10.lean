import GoLean.Math.Log10
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  match GoLean.Math.log10Word (profile == 1) (← wordInput 64 word) with
  | some result => return toString result
  | none => return "NaN"
