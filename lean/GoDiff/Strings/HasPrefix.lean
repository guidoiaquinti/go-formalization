import GoLean.Strings.CutPrefix
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let length :: words := fields | throw (IO.userError "expected the prefix length and bytes")
  if length < 0 || words.length < length.toNat then
    throw (IO.userError "invalid prefix length")
  let bytes ← words.mapM (wordInput 8)
  return toString (GoLean.Strings.hasPrefix (bytes.drop length.toNat) (bytes.take length.toNat))
