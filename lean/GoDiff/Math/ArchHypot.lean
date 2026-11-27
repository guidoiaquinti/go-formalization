import GoLean.Math.ArchHypot
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, x, y] := fields | throw (IO.userError "expected profile and two binary64 words")
  match GoLean.Math.archHypotWord (profile == 0) (← wordInput 64 x) (← wordInput 64 y) with
  | .ok result => return toString result
  | .error message => return "panic:" ++ message
