import GoLean.Strings.IsSeparator
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [r] := fields | throw (IO.userError "expected one signed rune")
  return toString (GoLean.Strings.isSeparator r)
