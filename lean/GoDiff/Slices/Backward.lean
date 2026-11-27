import GoLean.Slices.Backward
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let kind :: stop :: edit :: values := fields |
    throw (IO.userError "expected mode, stop index, edit delta and elements")
  let read (cells : List Int) (i : Nat) := cells[i]?.getD 0
  let callback (i : Nat) (_ : Int) (cells : List Int) : List Int × Except Nat Bool :=
    let updated := if i = 0 then cells else cells.set (i - 1) (read cells (i - 1) + edit)
    (updated, if (i : Int) = stop then
      if kind = 2 then .error (i + 1) else .ok false
      else .ok true)
  let (trace, final, result) := sliceBackward values.length read
    (if kind = 3 then none else some callback) values
  let status := match result with
    | .ok () => "ok"
    | .error .nilYield => "nil-panic"
    | .error (.callback payload) => s!"panic {payload}"
  return String.intercalate " " (status :: final.map toString ++ ["trace"] ++
    trace.map (fun (i, v) => s!"{i} {v}"))
