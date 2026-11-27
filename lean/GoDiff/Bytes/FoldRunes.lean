import GoLean.Bytes.FoldRunes
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def describe (r : RunesResult) : String :=
  s!"slice:{r.capacity}:{String.intercalate "," (r.values.map toString)}"

def main : IO Unit := runTuple fun fields => do
  match fields with
  | [0, rune] => return toString (GoLean.Strings.equalFoldSuccessor rune.toNat)
  | 1 :: count :: raw =>
    let bytes ← raw.mapM (wordInput 8)
    let data := bytes.map BitVec.toNat
    let left := data.take count.toNat
    let right := data.drop count.toNat
    return s!"{equalFoldBytes left right}|{describe (runesBytes left)}|{describe (runesBytes right)}"
  | _ => throw (IO.userError "expected a fold lookup or two byte slices")
