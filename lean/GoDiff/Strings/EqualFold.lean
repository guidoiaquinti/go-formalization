import GoLean.Strings.EqualFold
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  match fields with
  | [0, rune] => return toString (GoLean.Strings.equalFoldSuccessor rune.toNat)
  | 1 :: count :: raw =>
    let bytes ← raw.mapM (wordInput 8)
    let bytes := bytes.map BitVec.toNat
    return toString (GoLean.Strings.equalFold (bytes.take count.toNat) (bytes.drop count.toNat))
  | _ => throw (IO.userError "expected a fold table lookup or two byte strings")
