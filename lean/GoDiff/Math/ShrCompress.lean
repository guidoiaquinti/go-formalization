import GoLean.Math.ShrCompress
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [high, low, shift] := fields | throw (IO.userError "expected two words and a shift")
  let high ← wordInput 64 high
  let low ← wordInput 64 low
  let shift ← wordInput 64 shift
  let (outHigh, outLow) := GoLean.Math.shrcompress high low shift.toNat
  return s!"{outHigh.toNat} {outLow.toNat}"
