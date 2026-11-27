import GoLean.Path.Filepath.HasMeta
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected row marker and path bytes")
  let path ← raw.mapM (wordInput 8)
  return s!"{hasMeta false path}|{hasMeta true path}"
