import GoLean.Sort.PivotHelpers
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let kind :: a :: b :: c :: initial :: raw := row |
    throw (IO.userError "expected helper, indices, initial swaps and Less replies")
  let script := raw.map (· != 0)
  let state : PivotState := ⟨script, [], initial⟩
  let a := a.toNat
  let b := b.toNat
  let c := c.toNat
  let (first, second, state) :=
    match kind with
    | 0 =>
      let ((x, y), state) := order2 state a b
      (x, Int.ofNat y, state)
    | 1 =>
      let (x, state) := median state a b c
      (x, -1, state)
    | 2 =>
      let (x, state) := medianAdjacent state a
      (x, -1, state)
    | _ =>
      let ((x, hint), state) := choosePivot script a b
      (x, Int.ofNat hint, {state with swaps := 0})
  let calls := state.calls.reverse.map (fun (x, y) => s!"{x},{y}")
  return s!"{first}:{second}:{state.swaps}:" ++ String.intercalate ";" calls
