import GoLean.Math.Big.IntDivision
import GoDiff.Math.Bits.Tuple

open GoLean.Math.Big

private def wire (n : Option Int) : String := (n.map toString).getD "panic"
private def pairWire (pair : Option (Int × Int)) : String :=
  (pair.map fun (q, r) => s!"{q},{r}").getD "panic"

def main : IO Unit := runTuple fun fields => do
  let [x, y] := fields | throw (IO.userError "expected dividend and divisor")
  return String.intercalate "|" [wire (intQuo x y), wire (intRem x y),
    pairWire (intQuoRem x y), wire (intDiv x y), wire (intMod x y), pairWire (intDivMod x y)]
