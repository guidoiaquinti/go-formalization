import GoLean.Math.Bits.Arithmetic
import GoDiff.Math.Bits.Tuple

def main : IO Unit :=
  runTuple fun fields => do
    let [x, y] := fields | throw (IO.userError "wrong number of inputs")
    let x ← wordInput 32 x
    let y ← wordInput 32 y
    let (first, second) := GoLean.Math.Bits.mul32 x y
    return s!"{first.toNat} {second.toNat}"
