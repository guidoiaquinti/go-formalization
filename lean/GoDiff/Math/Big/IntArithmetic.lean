import GoLean.Math.Big.IntArithmetic
import GoDiff.Math.Bits.Tuple

open GoLean.Math.Big

def main : IO Unit := runTuple fun fields => do
  let [x, y] := fields | throw (IO.userError "expected two signed integers")
  return s!"{intSet x} {intAbs x} {intNeg x} {intAdd x y} {intSub x y} {intMul x y} {intCmp x y} {intCmpAbs x y}"
