import GoLean.Path.Filepath.CleanGlobPathWindows
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def wire (result : Nat × List (BitVec 8)) : String :=
  s!"{result.1}:{String.intercalate "," (result.2.map (toString ∘ BitVec.toNat))}"

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected marker and path bytes")
  let path ← raw.mapM (wordInput 8)
  return wire (cleanGlobPathWindows false path) ++ "|" ++ wire (cleanGlobPathWindows true path)
