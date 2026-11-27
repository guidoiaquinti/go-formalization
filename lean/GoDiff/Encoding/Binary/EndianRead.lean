import GoLean.Encoding.Binary.EndianRead
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected row marker and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let values := [bigUint16 bytes, bigUint32 bytes, bigUint64 bytes,
    littleUint16 bytes, littleUint32 bytes, littleUint64 bytes]
  return String.intercalate "|" (values.map fun value => value.map toString |>.getD "panic")
