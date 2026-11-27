import GoLean.Slices.Pdqsort
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (s : LocalSortState Int Int) : String :=
  let status := match s.failure with
    | none => "ok"
    | some .nilComparator | some .nilCounter => "nil"
    | some (.callback x) => s!"panic:{x}"
  let values := String.intercalate "," (s.cells.toList.map toString)
  let calls := String.intercalate ";" (s.calls.map fun (x,y) => s!"{x},{y}")
  s!"{status}/{values}/{calls}"

def main : IO Unit := runTuple fun fields => do
  let width :: a :: b :: limit :: mode :: stop :: mask :: cells := fields |
    throw (IO.userError "expected width, range, limit, comparator script and cells")
  let s : LocalSortState Int Int := ⟨cells.toArray, none, []⟩
  let compare := if mode == 0 then none else some (fun (i : Nat) (x y : Int) =>
    if Int.ofNat i == stop then Except.error stop
    else if mode == 1 then Except.ok (if x < y then -1 else 0)
    else Except.ok (if mask.toNat.testBit (i % 12) then -1 else 0))
  let less (a b : Int) := decide (a < b)
  return wire (pdqsortOrdered width.toNat s a.toNat b.toNat limit less) ++ "|" ++
    wire (pdqsortCmpFunc width.toNat s a.toNat b.toNat limit compare)
