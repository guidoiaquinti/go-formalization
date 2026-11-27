import GoLean.Bytes.BufferCap
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | [1, capacity] =>
      if capacity < 0 then throw (IO.userError "expected nonnegative capacity")
      else pure (some capacity.toNat)
    | _ => throw (IO.userError "expected nil marker or backing capacity")
  return ((GoLean.Bytes.bufferCap receiver).map toString).getD "panic"
