import GoLean.Path.Filepath.FromSlash
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def encode (path : List (BitVec 8)) : String :=
  String.intercalate "," (path.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected row marker and path bytes")
  let path ← raw.mapM (wordInput 8)
  return s!"{encode (fromSlash slash path)}|{encode (fromSlash 92 path)}"
