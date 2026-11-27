import GoLean.Math.RandV2.AbsInt32
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [word] := fields | throw (IO.userError "expected one uint32 word")
  let input ← wordInput 32 word
  return toString (GoLean.Math.RandV2.absInt32 input).toNat
