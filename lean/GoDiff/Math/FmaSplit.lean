import GoLean.Math.FmaSplit
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [raw] := fields | throw (IO.userError "expected one uint64 word")
  let word ← wordInput 64 raw
  let (sign, exponent, significand) := GoLean.Math.fmaSplit word
  return s!"{sign} {exponent} {significand}"
