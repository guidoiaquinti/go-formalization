import GoLean.Bytes.BufferEmpty
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | [1, length, offset] => pure (some (length.toNat, offset))
    | _ => throw (IO.userError "expected receiver marker, length and offset")
  return ((GoLean.Bytes.bufferEmpty receiver).map toString).getD "panic"
