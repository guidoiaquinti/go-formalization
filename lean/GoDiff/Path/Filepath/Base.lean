import GoLean.Path.Filepath.Base
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let windows :: raw := fields | throw (IO.userError "expected profile and path bytes")
  let path ← raw.mapM (wordInput 8)
  let base := baseName (windows == 1) path
  return "bytes:" ++ String.intercalate "," (base.map (toString ∘ BitVec.toNat))
