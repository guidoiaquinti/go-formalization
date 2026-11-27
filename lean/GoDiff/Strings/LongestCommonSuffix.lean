import GoLean.Strings.LongestCommonSuffix
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let count :: bytes := fields | throw (IO.userError "expected first string length and bytes")
  let raw ← bytes.mapM (wordInput 8)
  return toString (GoLean.Strings.longestCommonSuffix (raw.take count.toNat) (raw.drop count.toNat))
