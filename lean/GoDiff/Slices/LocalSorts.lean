import GoLean.Slices.LocalSorts
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
  let kind :: a :: b :: root :: mode :: stop :: mask :: cells := fields |
    throw (IO.userError "expected operation, range, root, comparator script and cells")
  let s : LocalSortState Int Int := ⟨cells.toArray, none, []⟩
  let compare := if mode == 0 then none else some (fun (i : Nat) (x y : Int) =>
    if Int.ofNat i == stop then Except.error stop
    else if mode == 1 then Except.ok (if x < y then -1 else 0)
    else Except.ok (if mask.toNat.testBit (i % 12) then -1 else 0))
  let less (a b : Int) := decide (a < b)
  let a := a.toNat
  let b := b.toNat
  let (ordered, callback) := if kind == 0 then
    (insertionSortOrdered s a b less, insertionSortCmpFunc s a b compare)
  else if kind == 1 then
    (siftDownOrdered s root.toNat (b - a) a less, siftDownCmpFunc s root.toNat (b - a) a compare)
  else
    (heapSortOrdered s a b less, heapSortCmpFunc s a b compare)
  return wire ordered ++ "|" ++ wire callback
