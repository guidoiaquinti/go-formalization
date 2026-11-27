import GoLean.Math.Frexp
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [word] := fields | throw (IO.userError "expected one binary64 word")
  let bits ← wordInput 64 word
  let (fraction, exponent) := GoLean.Math.frexp bits.toNat
  return s!"{fraction} {exponent}"
