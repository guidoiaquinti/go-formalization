import GoLean.Path.Filepath.Join
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def readParts : Nat → List Int → IO (List (List (BitVec 8)))
  | 0, [] => pure []
  | 0, _ => throw (IO.userError "trailing bytes")
  | _ + 1, [] => throw (IO.userError "missing element length")
  | count + 1, length :: fields => do
    if length < 0 || length.toNat > fields.length then
      throw (IO.userError "invalid element length")
    let part ← (fields.take length.toNat).mapM (wordInput 8)
    return part :: (← readParts count (fields.drop length.toNat))

def main : IO Unit := runTuple fun fields => do
  let windows :: count :: raw := fields | throw (IO.userError "expected profile and element count")
  if count < 0 then throw (IO.userError "invalid element count")
  let parts ← readParts count.toNat raw
  return "bytes:" ++ String.intercalate "," ((joinPath (windows == 1) parts).map (toString ∘ BitVec.toNat))
