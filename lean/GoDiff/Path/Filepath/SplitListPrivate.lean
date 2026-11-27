import GoLean.Path.Filepath.SplitList
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def encodePart (part : List (BitVec 8)) : String :=
  String.intercalate "," (part.map fun byte => toString byte.toNat)

private def encode (parts : List (List (BitVec 8))) : String :=
  s!"slice:{parts.length}:{String.intercalate ";" (parts.map encodePart)}"

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected row marker and path-list bytes")
  let path ← raw.mapM (wordInput 8)
  return s!"{encode (splitList 58 path)}|{encode (splitList 59 path)}|{encode (splitList 0 path)}"
