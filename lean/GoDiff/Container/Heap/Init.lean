import GoLean.Container.Heap.Init
import GoDiff.Math.Bits.Tuple

open GoLean.Container.Heap

private def encode : InitCall → String
  | .len => "len"
  | .sift (.less a b) => s!"L:{a}:{b}"
  | .sift (.swap a b) => s!"S:{a}:{b}"

def main : IO Unit := runTuple fun fields => do
  let [length, decisions] := fields | throw (IO.userError "expected length and decisions")
  if length < 0 || decisions < 0 then
    throw (IO.userError "expected nonnegative length and decisions")
  return String.intercalate " " ((GoLean.Container.Heap.init length.toNat decisions.toNat).map encode)
