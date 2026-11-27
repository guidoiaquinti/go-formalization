import GoLean.Unicode.Utf16.Encode
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: runes := fields | throw (IO.userError "expected a zero marker followed by runes")
  let (words, capacity) := GoLean.Unicode.Utf16.encode runes
  return String.intercalate " " ("present" :: toString capacity :: words.map toString)
