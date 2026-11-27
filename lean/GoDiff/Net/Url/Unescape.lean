import GoLean.Net.Url.Unescape
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def bytes (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

private def wire : Except UnescapeError (List (BitVec 8)) → String
  | .ok output => "ok:" ++ bytes output
  | .error (.escape raw) => "escape:" ++ bytes raw
  | .error (.host b) => s!"host:{b.toNat}"

def main : IO Unit := runTuple fun fields => do
  let mode :: raw := fields | throw (IO.userError "expected mode and bytes")
  let mode ← wordInput 8 mode
  let data ← raw.mapM (wordInput 8)
  return String.intercalate "|" [wire (unescape data mode), wire (queryUnescape data), wire (pathUnescape data)]
