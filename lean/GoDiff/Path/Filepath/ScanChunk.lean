import GoLean.Path.Filepath.ScanChunk
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def wire (path : List (BitVec 8)) : String :=
  String.intercalate "," (path.map (toString ∘ BitVec.toNat))

private def encode (result : Bool × List (BitVec 8) × List (BitVec 8)) : String :=
  let (star, chunk, rest) := result
  s!"{star}:{wire chunk}:{wire rest}"

def main : IO Unit := runTuple fun fields => do
  let 0 :: raw := fields | throw (IO.userError "expected marker and pattern bytes")
  let pattern ← raw.mapM (wordInput 8)
  return encode (scanChunk false pattern) ++ "|" ++ encode (scanChunk true pattern)
