import GoLean.Path.Filepath.Match
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let windows :: length :: raw := fields | throw (IO.userError "expected profile, pattern length and bytes")
  if length < 0 || length.toNat > raw.length then
    throw (IO.userError "invalid pattern length")
  let bytes ← raw.mapM (wordInput 8)
  let (matched, bad) := matchPath (windows == 1) (bytes.take length.toNat) (bytes.drop length.toNat)
  return s!"{matched}:{bad}"
