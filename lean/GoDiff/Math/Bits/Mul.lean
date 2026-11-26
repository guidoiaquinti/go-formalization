import GoLean.Math.Bits.Arithmetic
import GoDiff.Math.Bits.Tuple

def main : IO Unit :=
  runTuple fun fields => do
    let [x, y] := fields | throw (IO.userError "wrong number of inputs")
    let x ← wordInput 64 x
    let y ← wordInput 64 y
    let (first, second) := GoLean.Math.Bits.mul x y
    return s!"{first.toNat} {second.toNat}"
