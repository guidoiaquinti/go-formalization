import GoLean.Math.Ilogb
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [word] := fields | throw (IO.userError "expected one binary64 word")
  return toString (GoLean.Math.ilogb (← wordInput 64 word).toNat)
