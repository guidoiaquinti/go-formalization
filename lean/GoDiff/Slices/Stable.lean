import GoLean.Slices.Stable
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
  let kind :: a :: m :: b :: mode :: stop :: mask :: cells := fields |
    throw (IO.userError "expected entry kind, range, and comparator script")
  let s : LocalSortState Int Int := ⟨cells.toArray, none, []⟩
  let compare := if mode == 0 then none else some (fun (i : Nat) (x y : Int) =>
    if Int.ofNat i == stop then Except.error stop
    else if mode == 1 then Except.ok (if x / 256 < y / 256 then -1 else 0)
    else Except.ok (if mask.toNat.testBit (i % 12) then -1 else 0))
  let less (x y : Int) := decide (x < y)
  let (ordered, callback) := if kind == 0 then
    (symMergeOrdered s a.toNat m.toNat b.toNat less,
      symMergeCmpFunc s a.toNat m.toNat b.toNat compare)
    else (stableOrdered s b.toNat less, stableCmpFunc s b.toNat compare)
  return wire ordered ++ "|" ++ wire callback
