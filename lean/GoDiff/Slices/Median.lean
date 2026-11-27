import GoLean.Slices.Median
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (result : MedianObservation Int Int) : String :=
  let report := match result.result with
    | .ok i => s!"ok:{i}"
    | .error .nilCounter => "nil"
    | .error .nilComparator => "nil"
    | .error (.callback payload) => s!"panic:{payload}"
  let calls := String.intercalate ";" (result.calls.map fun (a,b) => s!"{a},{b}")
  s!"{report}/{(result.counter.map toString).getD "nil"}/{calls}"

def main : IO Unit := runTuple fun fields => do
  let width :: a :: b :: c :: center :: present :: counter :: mode :: failAt :: r0 :: r1 :: r2 :: cells := fields |
    throw (IO.userError "expected native width, indices, callback script and int64 cells")
  let get := fun (i : Nat) => cells[i]?.getD 0
  let count := if present == 0 then none else some counter
  let less (x y : Int) := decide (x < y)
  let compare := if mode == 0 then none else some (fun i x y : Int =>
    if i == failAt then Except.error i
    else if mode == 1 then Except.ok (if x < y then -1 else if y < x then 1 else 0)
    else Except.ok ([r0, r1, r2][i.toNat]?.getD 0))
  let callback := compare.map fun cmp => fun (i : Nat) x y => cmp (Int.ofNat i) x y
  return String.intercalate "|" [
    wire (medianOrdered width.toNat get a.toNat b.toNat c.toNat count less),
    wire (medianCmpFunc width.toNat get a.toNat b.toNat c.toNat count callback),
    wire (medianAdjacentOrdered width.toNat get center.toNat count less),
    wire (medianAdjacentCmpFunc width.toNat get center.toNat count callback)]
