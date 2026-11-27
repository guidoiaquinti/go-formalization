import GoLean.Sort.ComparisonEntrypoints
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

private def reply (tag : Int) : Except Unit Bool :=
  if tag = 2 then .error () else .ok (tag = 1)

private def showResult (r : Except String Bool) : String :=
  match r with
  | .ok b => toString b
  | .error s => s

def main : IO Unit := runTuple fun row => do
  let helper :: shape :: a :: b :: raw := row |
    throw (IO.userError "expected helper, shape, arguments and replies")
  if helper = 0 then
    let shape := if shape < 5 then ReflectedLength.supported
      else if shape = 5 then .invalid else .unsupported
    let (r, calls) := sliceIsSorted shape a.toNat (fun i _ => reply (raw[i - 1]?.getD 0))
    return showResult r ++ ":" ++ String.intercalate ";" (calls.map fun i => s!"{i},{i - 1}")
  else
    let callback := if shape = 0 then some (fun j i => if b = 2 then .error () else .ok (j < i)) else none
    let (r, calls) := reverseLess callback a (raw.headD 0)
    return showResult r ++ ":" ++ String.intercalate ";" (calls.map fun (i, j) => s!"{i},{j}")
