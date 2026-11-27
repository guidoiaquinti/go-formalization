import GoLean.Path.Filepath.Dir
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let windows :: raw := fields | throw (IO.userError "expected profile and path bytes")
  let path ← raw.mapM (wordInput 8)
  let dir := dirName (windows == 1) path
  return "bytes:" ++ String.intercalate "," (dir.map (toString ∘ BitVec.toNat))
