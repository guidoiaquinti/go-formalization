import GoLean.Sort.SortEntrypoints
import GoDiff.Math.Bits.Tuple
open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let kind :: width :: tag :: report :: raw := row |
    throw (IO.userError "expected entry point, word width, tag, Len and data")
  let state : LocalSortState := ⟨raw.toArray, []⟩
  let result := match kind with
    | 0 => sortInterface width.toNat state (if tag = 0 then some report else none)
    | 1 => intSliceSort width.toNat state
    | _ =>
      let shape := if tag = 2 then ReflectionSlice.invalid
        else if tag ≥ 3 then ReflectionSlice.nonSlice else ReflectionSlice.valid
      sortSlice width.toNat state shape (tag != 1)
  let values := String.intercalate "," (result.state.values.toList.map toString)
  let calls := if kind = 1 then "" else String.intercalate ";" (result.state.reversedCalls.reverse.map fun (swap, i, j) =>
    s!"{if swap then "S" else "L"}{i},{j}")
  let lenCalls := if kind = 1 then 0 else result.lenCalls
  return s!"{result.status}:{lenCalls}:{values}:{calls}"
