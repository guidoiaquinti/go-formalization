import GoLean.Strings.CutPrefix
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let length :: words := fields | throw (IO.userError "expected the prefix length and bytes")
  if length < 0 || words.length < length.toNat then
    throw (IO.userError "invalid prefix length")
  let bytes ← words.mapM (wordInput 8)
  let out := GoLean.Strings.trimPrefix (bytes.drop length.toNat) (bytes.take length.toNat)
  return String.intercalate " " ("bytes" :: out.map (fun byte => toString byte.toNat))
