import GoLean.Sort.PartialInsertion
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let variant :: a :: b :: raw := row | throw (IO.userError "expected variant, range and data")
  let state : LocalSortState := ⟨raw.toArray, []⟩
  let (sorted, state) := if variant = 0 then partialInsertionSort state a.toNat b.toNat
    else partialInsertionSortFunc state a.toNat b.toNat
  let values := String.intercalate "," (state.values.toList.map toString)
  let calls := state.reversedCalls.reverse.map fun (swap, i, j) =>
    s!"{if swap then "S" else "L"}{i},{j}"
  return s!"{sorted}:{values}:" ++ String.intercalate ";" calls
