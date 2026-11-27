import GoLean.Math.Exp
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let input ← wordInput 64 word
  let result := if profile < 3 then GoLean.Math.expWord (profile == 2) (profile == 1) input
    else GoLean.Math.expGenericWord (profile == 5) (profile == 4) input
  return toString result
