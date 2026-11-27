import GoLean.Math.Inf
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let [sign] := fields | throw (IO.userError "expected one signed integer")
  return toString (GoLean.Math.inf sign).toNat
