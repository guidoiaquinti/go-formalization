import GoLean.Math.Sincos
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  match GoLean.Math.sincosWords (profile == 2) (profile == 1) (← wordInput 64 word) with
  | some (sine, cosine) => return s!"{sine} {cosine}"
  | none => return "NaN NaN"
