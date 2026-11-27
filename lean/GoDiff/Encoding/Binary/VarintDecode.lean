import GoLean.Encoding.Binary.VarintDecode
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

def main : IO Unit := runTuple fun fields => do
  let length :: raw := fields | throw (IO.userError "expected byte length")
  if length < 0 || length.toNat != raw.length then
    throw (IO.userError "wrong byte length")
  let bytes ← raw.mapM (wordInput 8)
  let (u, n) := uvarint bytes
  let (s, m) := varint bytes
  return s!"{u.toNat}:{n}|{s}:{m}"
