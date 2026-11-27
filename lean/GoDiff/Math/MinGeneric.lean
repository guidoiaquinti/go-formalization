import GoLean.Math.Min
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [a, b] := fields | throw (IO.userError "expected two binary64 encodings")
  let word := GoLean.Math.minFloat64 (← wordInput 64 a).toNat (← wordInput 64 b).toNat
  return toString word
