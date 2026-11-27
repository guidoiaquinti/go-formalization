import GoLean.Path.Filepath.Clean
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let windows :: raw := fields | throw (IO.userError "expected profile and path bytes")
  let path ← raw.mapM (wordInput 8)
  return "bytes:" ++ String.intercalate "," ((cleanPath (windows == 1) path).map (toString ∘ BitVec.toNat))
