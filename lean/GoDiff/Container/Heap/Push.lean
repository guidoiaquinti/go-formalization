import GoLean.Container.Heap.Push
import GoDiff.Math.Bits.Tuple

open GoLean.Container.Heap

private def encode : PushCall → String
  | .push value => s!"P:{value}"
  | .len => "len"
  | .lift (.less child parent) => s!"L:{child}:{parent}"
  | .lift (.swap parent child) => s!"S:{parent}:{child}"

def main : IO Unit := runTuple fun fields => do
  let [value, length, decisions] := fields
    | throw (IO.userError "expected value, post-push length and decisions")
  if value < 0 || length ≤ 0 || decisions < 0 then
    throw (IO.userError "expected nonnegative value and decisions, and positive length")
  return String.intercalate " " ((push value.toNat length.toNat decisions.toNat).map encode)
