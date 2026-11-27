import GoLean.Path.Filepath.MatchChunk
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def encode (value : List Nat × Bool × Bool) : String :=
  let (rest, ok, bad) := value
  s!"{ok}:{bad}:{String.intercalate "," (rest.map toString)}"

def main : IO Unit := runTuple fun fields => do
  let length :: raw := fields | throw (IO.userError "expected pattern length and input")
  if length < 0 || length.toNat > raw.length then
    throw (IO.userError "invalid pattern length")
  let bytes ← raw.mapM (wordInput 8)
  let bytes := bytes.map BitVec.toNat
  let pattern := bytes.take length.toNat
  let input := bytes.drop length.toNat
  return encode (matchChunk false pattern input) ++ "|" ++ encode (matchChunk true pattern input)
