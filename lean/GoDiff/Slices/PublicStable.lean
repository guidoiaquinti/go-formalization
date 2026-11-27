import GoLean.Slices.PublicStable
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (o : PublicSortObservation Int Int) : String :=
  let status := match o.failure with
    | none => "ok"
    | some (.sequence .nilSequence) => "nilSequence"
    | some (.sequence (.producer code)) => s!"producer:{code}"
    | some (.sort (.callback code)) => s!"compare:{code}"
    | some (.sort .nilComparator) | some (.sort .nilCounter) => "nilCompare"
  let data := match o.values with
    | none => "nil"
    | some xs => "slice:" ++ String.intercalate "," (xs.map toString)
  let calls := String.intercalate ";" (o.calls.map fun (x,y) => s!"{x},{y}")
  s!"{status}/{data}/{o.acknowledged}/{calls}"

def main : IO Unit := runTuple fun fields => do
  let tag :: producer :: mode :: stop :: mask :: cells := fields |
    throw (IO.userError "expected slice and iterator setup, comparator script, and cells")
  let input := if tag == 0 then none else some cells
  let sequence := if producer == 0 then none else some (cells,
    if producer == 2 then Except.error 7 else Except.ok ())
  let compare := if mode == 0 then none else some (fun (i : Nat) (x y : Int) =>
    if Int.ofNat i == stop then Except.error stop
    else if mode == 1 then Except.ok (if x / 256 < y / 256 then -1 else 0)
    else Except.ok (if mask.toNat.testBit (i % 12) then -1 else 0))
  return wire (sliceSortStableFunc input compare) ++ "|" ++
    wire (sliceSortedStableFunc sequence compare)
