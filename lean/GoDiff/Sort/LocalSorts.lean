import GoLean.Sort.LocalSorts
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let kind :: variant :: a :: b :: c :: raw := row |
    throw (IO.userError "expected helper, variant, range and data")
  let state : LocalSortState := ⟨raw.toArray, []⟩
  let state := match kind, variant with
    | 0, 0 => insertionSort state a.toNat b.toNat
    | 0, _ => insertionSortFunc state a.toNat b.toNat
    | 1, 0 => siftDown state a.toNat b.toNat c.toNat
    | 1, _ => siftDownFunc state a.toNat b.toNat c.toNat
    | _, 0 => heapSort state a.toNat b.toNat
    | _, _ => heapSortFunc state a.toNat b.toNat
  let values := String.intercalate "," (state.values.toList.map toString)
  let calls := state.reversedCalls.reverse.map fun (swap, i, j) =>
    s!"{if swap then "S" else "L"}{i},{j}"
  return values ++ ":" ++ String.intercalate ";" calls
