import GoLean.Sort.Partition
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let kind :: variant :: a :: b :: p :: raw := row |
    throw (IO.userError "expected helper, variant, indices and data")
  let state : LocalSortState := ⟨raw.toArray, []⟩
  let (index, status, state) :=
    if kind = 0 then
      let (index, flag, state) := if variant = 0 then partition state a.toNat b.toNat p.toNat
        else partitionFunc state a.toNat b.toNat p.toNat
      (index, toString flag, state)
    else
      let (index, state) := if variant = 0 then partitionEqual state a.toNat b.toNat p.toNat
        else partitionEqualFunc state a.toNat b.toNat p.toNat
      (index, "-", state)
  let values := String.intercalate "," (state.values.toList.map toString)
  let calls := state.reversedCalls.reverse.map fun (swap, i, j) =>
    s!"{if swap then "S" else "L"}{i},{j}"
  return s!"{index}:{status}:{values}:" ++ String.intercalate ";" calls
