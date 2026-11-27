import GoLean.Path.Filepath.IsAbs
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected marker and path bytes")
  let path ← raw.mapM (wordInput 8)
  return s!"{isAbs .unix path}|{isAbs .windows path}|{isAbs .plan9 path}"
