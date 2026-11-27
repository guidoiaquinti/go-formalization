import GoLean.Container.Heap.Remove
import GoDiff.Math.Bits.Tuple

open GoLean.Container.Heap

private def encode : RemoveCall → String
  | .len => "len"
  | .initialSwap a b => s!"S:{a}:{b}"
  | .repair (.sift (.less a b)) => s!"L:{a}:{b}"
  | .repair (.sift (.swap a b)) => s!"S:{a}:{b}"
  | .repair .len => "unexpected Len"
  | .pop => "pop"

def main : IO Unit := runTuple fun fields => do
  let [index, length, decisions, value] := fields
    | throw (IO.userError "expected index, length, decisions and return value")
  if index < 0 || index ≥ length || decisions < 0 || value < 0 then
    throw (IO.userError "expected an in-range index and nonnegative decisions and value")
  let (result, calls) := remove index.toNat length.toNat decisions.toNat value.toNat
  return String.intercalate " " (toString result :: calls.map encode)
