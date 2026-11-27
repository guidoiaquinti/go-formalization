import GoLean.Math.SinPi
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, word] := fields | throw (IO.userError "expected profile and binary64 word")
  let result := GoLean.Math.sinPiWord (profile == 2) (profile == 1) (← wordInput 64 word)
  let some result := result | throw (IO.userError "expected a finite nonnegative argument")
  return toString result
