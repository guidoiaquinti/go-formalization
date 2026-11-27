import GoLean.Path.Filepath.Split
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def wire (raw : List (BitVec 8)) : String :=
  String.intercalate "," (raw.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun fields => do
  let windows :: raw := fields | throw (IO.userError "expected profile and path bytes")
  let path ← raw.mapM (wordInput 8)
  let (dir, file) := splitPath (windows == 1) path
  return wire dir ++ "|" ++ wire file
