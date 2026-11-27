import GoLean.Slices.Values
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let kind :: stop :: edit :: values := fields |
    throw (IO.userError "expected mode, stop index, edit delta and elements")
  let read (cells : List Int × Nat) (i : Nat) := cells.1[i]?.getD 0
  let callback (_ : Int) (cells : List Int × Nat) : (List Int × Nat) × Except Nat Bool :=
    let i := cells.2
    let updated := cells.1.set (i + 1) (read cells (i + 1) + edit)
    ((updated, i + 1), if (i : Int) = stop then
      if kind = 2 then .error (i + 1) else .ok false
      else .ok true)
  let (trace, (final, _), result) := sliceValues values.length read
    (if kind = 3 then none else some callback) (values, 0)
  let (replay, _, _) := sliceValues (ε := Nat) values.length read
    (some fun _ s => (s, .ok true)) (final, 0)
  let status := match result with
    | .ok () => "ok"
    | .error .nilYield => "nil-panic"
    | .error (.callback payload) => s!"panic {payload}"
  return String.intercalate " " (status :: final.map toString ++ ["trace"] ++ trace.map toString ++
    ["replay"] ++ replay.map toString)
