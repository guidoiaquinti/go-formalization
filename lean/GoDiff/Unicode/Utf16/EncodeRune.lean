import GoLean.Unicode.Utf16.EncodeRune
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r] := fields | throw (IO.userError "expected one rune")
  let (high, low) := GoLean.Unicode.Utf16.encodeRune r
  return s!"{high} {low}"
