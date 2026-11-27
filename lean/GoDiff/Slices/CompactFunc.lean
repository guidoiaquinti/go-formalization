import GoLean.Slices.CompactFunc
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let mode :: fault :: xs := fields | throw (IO.userError "expected mode, fault and slice")
  let comparePair (i : Nat) (current previous : Int) : Except Nat Bool :=
    if (i : Int) = fault then .error (i + 1)
    else .ok (if mode = 0 then current = previous
      else if mode = 1 then current = previous + 1 else i % 3 = 0)
  let (out, backing, trace) := sliceCompactFunc
    (if mode = 3 then none else some comparePair) 0 xs
  let tag := match out with
    | .ok kept => s!"ok {kept.length}"
    | .error payload => s!"panic {payload}"
  return String.intercalate " " (tag :: backing.map toString ++
    ["trace"] ++ trace.map (fun (a, b) => s!"{a} {b}"))
