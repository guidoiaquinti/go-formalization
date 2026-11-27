import GoLean.Bytes.ReaderLen
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let receiver ← match row with
    | [0] => pure none
    | [1, length, offset] =>
      if length < 0 || offset < 0 then throw (IO.userError "expected legal reader geometry")
      else pure (some (length.toNat, offset.toNat))
    | _ => throw (IO.userError "expected nil marker or length and offset")
  return ((GoLean.Bytes.readerLen receiver).map toString).getD "panic"
