import GoLean.Math.ArchExp2
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let input ← wordInput 64 word
  match GoLean.Math.archExp2Word (profile == 2) input with
  | .ok result => return toString result
  | .error message => return "panic:" ++ message
