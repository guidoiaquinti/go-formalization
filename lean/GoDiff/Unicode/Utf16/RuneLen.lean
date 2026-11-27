import GoLean.Unicode.Utf16.RuneLen
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r] := fields | throw (IO.userError "expected one rune")
  return toString (GoLean.Unicode.Utf16.runeLen r)
