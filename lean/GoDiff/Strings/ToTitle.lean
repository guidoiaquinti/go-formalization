import GoLean.Strings.ToTitle
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  match fields with
  | [0, rune] => return toString (GoLean.Strings.titleRune rune.toNat)
  | 1 :: raw =>
    let bytes ← raw.mapM (wordInput 8)
    let output := GoLean.Strings.toTitle bytes
    return "bytes:" ++ String.intercalate "," (output.map toString)
  | _ => throw (IO.userError "expected a titlecase table lookup or a byte string")
