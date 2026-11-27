import GoLean.Unicode.Utf16.Decode
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: words := fields | throw (IO.userError "expected zero marker followed by UTF-16 words")
  let input ← words.mapM fun word => do
    return (← wordInput 16 word).toNat
  return String.intercalate " " ("runes" :: (GoLean.Unicode.Utf16.decode input).map toString)
