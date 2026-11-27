import GoLean.Math.Log
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let input ← wordInput 64 word
  let result := if profile < 2 then GoLean.Math.logWord (profile == 1) input
    else GoLean.Math.logGenericWord (profile == 4) (profile == 3) input
  match result with
  | some result => return toString result
  | none => return "NaN"
