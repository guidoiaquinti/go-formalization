import GoLean.Path.Filepath.GetEsc
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def encode (value : Nat × List Nat × Bool) : String :=
  let (rune, rest, bad) := value
  s!"{rune}:{bad}:{String.intercalate "," (rest.map toString)}"

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected marker and chunk bytes")
  let chunk ← raw.mapM (wordInput 8)
  let input := chunk.map BitVec.toNat
  return encode (getEsc false input) ++ "|" ++ encode (getEsc true input)
