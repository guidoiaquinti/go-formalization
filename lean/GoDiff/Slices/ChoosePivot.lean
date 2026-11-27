import GoLean.Slices.ChoosePivot
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def wire (result : ChoosePivotObservation Int Int) : String :=
  let report := match result.result with
    | .ok (i, hint) => s!"ok:{i},{hint}"
    | .error .nilCounter => "nil"
    | .error .nilComparator => "nil"
    | .error (.callback payload) => s!"panic:{payload}"
  let calls := String.intercalate ";" (result.calls.map fun (a,b) => s!"{a},{b}")
  s!"{report}/{calls}"

def main : IO Unit := runTuple fun fields => do
  let width :: a :: b :: mode :: failAt :: replyMask :: cells := fields |
    throw (IO.userError "expected native width, range, callback script and int64 cells")
  let get := fun (i : Nat) => cells[i]?.getD 0
  let compare := if mode == 0 then none else some (fun (i : Nat) (x y : Int) =>
    if Int.ofNat i == failAt then Except.error failAt
    else if mode == 1 then Except.ok (if x < y then -1 else if y < x then 1 else 0)
    else Except.ok (if replyMask.toNat.testBit i then -1 else 0))
  return wire (choosePivotOrdered width.toNat get a.toNat b.toNat (fun x y => x < y)) ++ "|" ++
    wire (choosePivotCmpFunc width.toNat get a.toNat b.toNat compare)
