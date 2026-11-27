import GoLean.Strings.SingleStringReplacerReplace
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let patternLength :: valueLength :: fields := fields
    | throw (IO.userError "expected pattern length, value length and bytes")
  let bytes ← fields.mapM (wordInput 8)
  if patternLength ≤ 0 || valueLength < 0 || patternLength.toNat + valueLength.toNat > bytes.length then
    throw (IO.userError "invalid pattern or value length")
  let pattern := bytes.take patternLength.toNat
  let afterPattern := bytes.drop patternLength.toNat
  let value := afterPattern.take valueLength.toNat
  let input := afterPattern.drop valueLength.toNat
  let result := GoLean.Strings.singleStringReplacerReplace pattern value input
  return "bytes:" ++ String.intercalate "," (result.map (fun byte => toString byte.toNat))
