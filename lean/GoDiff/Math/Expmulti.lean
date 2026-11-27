import GoLean.Math.Expmulti
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [profile, hi, lo, k] := fields | throw (IO.userError "expected profile, two binary64 words and exponent")
  let result := GoLean.Math.expmultiWord (profile == 2) (profile == 1)
    (← wordInput 64 hi) (← wordInput 64 lo) k
  return toString result
