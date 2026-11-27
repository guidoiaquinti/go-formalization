import GoLean.Math.Exp2
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  return toString (GoLean.Math.exp2Word (profile == 2) (profile == 1 || profile == 3) (profile == 3) (← wordInput 64 word))
