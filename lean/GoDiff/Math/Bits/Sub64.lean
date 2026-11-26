import GoLean.Math.Bits.Arithmetic
import GoDiff.Math.Bits.Tuple

def main : IO Unit :=
  runTuple fun fields => do
    let [x, y, flag] := fields | throw (IO.userError "wrong number of inputs")
    let x ← wordInput 64 x
    let y ← wordInput 64 y
    if flag != 0 && flag != 1 then throw (IO.userError "expected a one-bit flag")
    let (first, second) := GoLean.Math.Bits.sub64 x y (flag == 1)
    return s!"{first.toNat} {second.toNat}"
