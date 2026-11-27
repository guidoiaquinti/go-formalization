import GoLean.Path.Filepath.HasPrefix
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let windows :: length :: fields := fields | throw (IO.userError "expected profile and path length")
  if length < 0 || length.toNat > fields.length then
    throw (IO.userError "invalid path length")
  let raw ← fields.mapM (wordInput 8)
  return if hasPathPrefix (windows == 1) (raw.take length.toNat) (raw.drop length.toNat)
    then "true" else "false"
