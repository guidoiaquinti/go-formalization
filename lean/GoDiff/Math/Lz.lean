import GoLean.Math.Lz
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [high, low] := fields | throw (IO.userError "expected two uint64 words")
  let hi ← wordInput 64 high
  let lo ← wordInput 64 low
  return toString (GoLean.Math.lz hi lo)
