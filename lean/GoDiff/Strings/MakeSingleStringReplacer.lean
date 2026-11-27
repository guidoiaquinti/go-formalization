import GoLean.Strings.MakeSingleStringReplacer
import GoDiff.Math.Bits.Tuple

private def numbers (values : List Nat) : String :=
  String.intercalate "," (values.map toString)

def main : IO Unit := runTuple fun fields => do
  let count :: fields := fields | throw (IO.userError "expected pattern length, pattern and value bytes")
  let bytes ← fields.mapM (wordInput 8)
  if count < 0 || count.toNat > bytes.length then
    throw (IO.userError "invalid pattern length")
  let replacer := GoLean.Strings.makeSingleStringReplacer (bytes.take count.toNat) (bytes.drop count.toNat)
  let finder := replacer.finder
  return numbers (finder.pattern.map BitVec.toNat) ++ ":" ++ numbers (replacer.value.map BitVec.toNat) ++
    ":" ++ numbers finder.badCharSkip ++ ":" ++ numbers finder.goodSuffixSkip
