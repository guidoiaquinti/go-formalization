import GoLean.Container.Heap.Pop
import GoDiff.Math.Bits.Tuple

open GoLean.Container.Heap

private def encode : PopCall → String
  | .len => "len"
  | .rootSwap last => s!"S:0:{last}"
  | .sift (.less a b) => s!"L:{a}:{b}"
  | .sift (.swap a b) => s!"S:{a}:{b}"
  | .pop => "pop"

def main : IO Unit := runTuple fun fields => do
  let [length, decisions, value] := fields | throw (IO.userError "expected length, decisions and return value")
  if length < 0 || decisions < 0 || value < 0 then
    throw (IO.userError "expected nonnegative inputs")
  let (result, calls) := pop length.toNat decisions.toNat value.toNat
  return String.intercalate " " (toString result :: calls.map encode)
