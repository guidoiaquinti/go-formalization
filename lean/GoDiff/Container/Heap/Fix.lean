import GoLean.Container.Heap.Fix
import GoDiff.Math.Bits.Tuple

open GoLean.Container.Heap

private def encode : FixCall → String
  | .len => "len"
  | .sift (.less a b) => s!"L:{a}:{b}"
  | .sift (.swap a b) => s!"S:{a}:{b}"

def main : IO Unit := runTuple fun fields => do
  let [index, length, decisions] := fields
    | throw (IO.userError "expected index, length and decisions")
  if index < 0 || index ≥ length || decisions < 0 then
    throw (IO.userError "expected an in-range heap index and nonnegative decisions")
  return String.intercalate " " ((fix index.toNat length.toNat decisions.toNat).map encode)
