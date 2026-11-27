import GoLean.Math.Pow10
import GoDiff.Math.Bits.Tuple

def main : IO Unit :=
  runTuple fun fields => do
    let [n] := fields | throw (IO.userError "expected one signed exponent")
    return toString (GoLean.Math.pow10 n).toBits
