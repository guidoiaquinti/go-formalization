import GoLean.Math.Dim
import GoDiff.Math.Bits.Tuple

def main : IO Unit :=
  runTuple fun fields => do
    let [x, y] := fields | throw (IO.userError "expected two binary64 words")
    let x := Float.ofBits (UInt64.ofBitVec (← wordInput 64 x))
    let y := Float.ofBits (UInt64.ofBitVec (← wordInput 64 y))
    let result := GoLean.Math.dim x y
    return if result.isNaN then "NaN" else toString result.toBits
