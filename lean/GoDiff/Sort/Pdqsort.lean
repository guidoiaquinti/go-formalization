import GoLean.Sort.Pdqsort
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let variant :: width :: a :: b :: limit :: raw := row |
    throw (IO.userError "expected variant, word width, range, limit and data")
  let state : LocalSortState := ⟨raw.toArray, []⟩
  let state := if variant = 0 then pdqsort width.toNat state a.toNat b.toNat limit.toNat
    else pdqsortFunc width.toNat state a.toNat b.toNat limit.toNat
  let values := String.intercalate "," (state.values.toList.map toString)
  let calls := state.reversedCalls.reverse.map fun (swap, i, j) =>
    s!"{if swap then "S" else "L"}{i},{j}"
  return values ++ ":" ++ String.intercalate ";" calls
