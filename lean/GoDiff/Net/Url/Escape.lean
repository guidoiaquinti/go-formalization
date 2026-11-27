import GoLean.Net.Url.Escape
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def render (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun fields => do
  let mode :: raw := fields | throw (IO.userError "expected encoding and input")
  let mode ← wordInput 8 mode
  let bytes ← raw.mapM (wordInput 8)
  let classes := bytes.map fun b =>
    s!"{shouldEscape b mode}:{isHex b}:{(unhex b).toNat}"
  return String.intercalate "|" [String.intercalate "," classes, render (escape bytes mode),
    render (queryEscape bytes), render (pathEscape bytes)]
