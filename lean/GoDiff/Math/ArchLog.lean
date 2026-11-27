import GoLean.Math.ArchLog
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let input ← wordInput 64 word
  match GoLean.Math.archLogWord (profile == 0) input with
  | .ok result => return (result.map toString).getD "NaN"
  | .error message => return "panic:" ++ message
