import GoLean.Path.Filepath.CleanGlobPath
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let separator :: raw := fields | throw (IO.userError "expected separator and raw path")
  let separator ← wordInput 8 separator
  let path ← raw.mapM (wordInput 8)
  return "bytes:" ++ String.intercalate "," ((cleanGlobPath separator path).map (toString ∘ BitVec.toNat))
