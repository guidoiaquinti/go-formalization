import GoLean.Path.Filepath.VolumeName
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected marker and path bytes")
  let path ← raw.mapM (wordInput 8)
  let output := volumeName true path
  return "bytes:|bytes:" ++ String.intercalate "," (output.map (toString ∘ BitVec.toNat))
