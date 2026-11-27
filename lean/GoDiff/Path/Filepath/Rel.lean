import GoLean.Path.Filepath.Rel
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

def main : IO Unit := runTuple fun fields => do
  let windows :: length :: fields := fields | throw (IO.userError "expected profile and base length")
  if length < 0 || length.toNat > fields.length then
    throw (IO.userError "invalid base length")
  let raw ← fields.mapM (wordInput 8)
  let render (result : List (BitVec 8)) :=
    String.intercalate "," (result.map (toString ∘ BitVec.toNat))
  match relativePath (windows == 1) (raw.take length.toNat) (raw.drop length.toNat) with
  | none => return "diverges"
  | some (.ok output) => return "ok:" ++ render output
  | some (.error message) => return "err:" ++ render message
