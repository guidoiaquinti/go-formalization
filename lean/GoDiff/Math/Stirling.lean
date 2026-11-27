import GoLean.Math.Stirling
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let result := GoLean.Math.stirlingWords (profile == 4) (profile == 2 || profile == 3)
    (profile == 1 || profile == 3) (← wordInput 64 word)
  let some (first, second) := result | throw (IO.userError "expected Stirling domain")
  return s!"{first} {second}"
