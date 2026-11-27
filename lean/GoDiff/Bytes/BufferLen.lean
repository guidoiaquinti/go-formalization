import GoLean.Bytes.BufferLen
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | [1, length, offset] =>
      if offset < 0 || offset > length then
        throw (IO.userError "expected a reachable buffer offset")
      else pure (some (length.toNat, offset.toNat))
    | _ => throw (IO.userError "expected nil marker or length and offset")
  return ((GoLean.Bytes.bufferLen receiver).map toString).getD "panic"
