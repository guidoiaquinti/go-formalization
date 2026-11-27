import GoLean.Bytes.BufferAvailable
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | [1, length, capacity] =>
      if length < 0 || capacity < length then
        throw (IO.userError "expected reachable length and capacity")
      else pure (some (length.toNat, capacity.toNat))
    | _ => throw (IO.userError "expected nil marker or length and capacity")
  return ((GoLean.Bytes.bufferAvailable receiver).map toString).getD "panic"
