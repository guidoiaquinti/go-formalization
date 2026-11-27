import GoLean.Sort.RangeMoves
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let [kind, variant, a, b, n] := row |
    throw (IO.userError "expected helper, variant and three range arguments")
  let a := a.toNat
  let b := b.toNat
  let n := n.toNat
  let trace := match kind, variant with
    | 0, 0 => reverseRange a b
    | 0, _ => reverseRangeFunc a b
    | 1, 0 => swapRange a b n
    | 1, _ => swapRangeFunc a b n
    | _, 0 => rotate a b n
    | _, _ => rotateFunc a b n
  return "swaps:" ++ String.intercalate ";" (trace.map fun (x, y) => s!"{x},{y}")
