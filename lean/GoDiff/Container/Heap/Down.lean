import GoLean.Container.Heap.Down
import GoDiff.Math.Bits.Tuple

open GoLean.Container.Heap

private def encode : Call → String
  | .less child parent => s!"L:{child}:{parent}"
  | .swap parent child => s!"S:{parent}:{child}"

def main : IO Unit := runTuple fun fields => do
  let [index, length, low, high] := fields
    | throw (IO.userError "expected index, length, and two decision words")
  if index < 0 || length < 0 then
    throw (IO.userError "expected nonnegative index and length")
  let low ← wordInput 64 low
  let high ← wordInput 64 high
  let (moved, calls) := down index.toNat length.toNat (low.toNat + 2 ^ 64 * high.toNat)
  return String.intercalate " " (toString moved :: calls.map encode)
