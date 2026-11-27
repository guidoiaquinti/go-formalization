import GoLean.Bytes.SimpleCase
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def render (r : CaseBytes) : String :=
  let capacity := (r.exactCapacity.map toString).getD "growth"
  s!"slice:{capacity}:{String.intercalate "," (r.bytes.map toString)}"

def main : IO Unit := runTuple fun fields => do
  match fields with
  | [0, r] =>
    let rune := r.toNat
    return s!"{GoLean.Strings.upperRune rune},{GoLean.Strings.lowerRune rune},{GoLean.Strings.titleRune rune}"
  | 1 :: raw =>
    let input ← raw.mapM (wordInput 8)
    return String.intercalate "|" [render (toUpperBytes input), render (toLowerBytes input),
      render (toTitleBytes input)]
  | _ => throw (IO.userError "expected a case lookup or byte slice")
