import GoLean.Container.Heap.Up
import GoDiff.Math.Bits.Tuple

open GoLean.Container.Heap

private def encode (call : Call) : String :=
  match call with
  | .less child parent => s!"L:{child}:{parent}"
  | .swap parent child => s!"S:{parent}:{child}"

def main : IO Unit := runTuple fun fields => do
  let [child, decisions] := fields | throw (IO.userError "expected index and decisions")
  if child < 0 || decisions < 0 then
    throw (IO.userError "expected nonnegative inputs")
  return String.intercalate " " ("trace" :: (up child.toNat decisions.toNat).map encode)
