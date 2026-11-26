import GoLean.Math.Bits.Rotation
import GoDiff.Math.Bits.Tuple

def main : IO Unit :=
  runTuple fun fields => do
    let [x, k] := fields | throw (IO.userError "wrong number of inputs")
    let x ← wordInput 16 x
    return toString (GoLean.Math.Bits.rotateLeft16 x k).toNat
