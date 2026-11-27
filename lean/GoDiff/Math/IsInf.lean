import GoLean.Math.IsInf
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [bits, sign] := fields | throw (IO.userError "expected binary64 bits and a signed integer")
  let repr ← wordInput 64 bits
  return toString (GoLean.Math.isInf repr sign)
