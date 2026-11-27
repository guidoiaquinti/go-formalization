import GoLean.Unicode.Utf16.DecodeRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [first, second] := fields | throw (IO.userError "expected two runes")
  return toString (GoLean.Unicode.Utf16.decodeRune first second)
