import GoLean.Slices.Partition
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (r : Option String) (s : LocalSortState Int Int) : String :=
  let status := match s.failure with
    | none => r.getD "missing"
    | some .nilComparator | some .nilCounter => "nil"
    | some (.callback x) => s!"panic:{x}"
  let values := String.intercalate "," (s.cells.toList.map toString)
  let calls := String.intercalate ";" (s.calls.map fun (x,y) => s!"{x},{y}")
  s!"{status}/{values}/{calls}"

def main : IO Unit := runTuple fun fields => do
  let kind :: a :: b :: pivot :: mode :: stop :: mask :: cells := fields |
    throw (IO.userError "expected operation, range, pivot, comparator script and cells")
  let s : LocalSortState Int Int := ⟨cells.toArray, none, []⟩
  let compare := if mode == 0 then none else some (fun (i : Nat) (x y : Int) =>
    if Int.ofNat i == stop then Except.error stop
    else if mode == 1 then Except.ok (if x < y then -1 else 0)
    else Except.ok (if mask.toNat.testBit (i % 12) then -1 else 0))
  let less (a b : Int) := decide (a < b)
  let a := a.toNat
  let b := b.toNat
  let pivot := pivot.toNat
  if kind == 0 then
    let (i, ordered) := partitionOrdered s a b pivot less
    let (j, callback) := partitionCmpFunc s a b pivot compare
    let render (pair : Nat × Bool) := s!"ok:{pair.1},{pair.2}"
    return wire (i.map render) ordered ++ "|" ++ wire (j.map render) callback
  else
    let (i, ordered) := partitionEqualOrdered s a b pivot less
    let (j, callback) := partitionEqualCmpFunc s a b pivot compare
    return wire (i.map (fun n => s!"ok:{n}")) ordered ++ "|" ++
      wire (j.map (fun n => s!"ok:{n}")) callback
