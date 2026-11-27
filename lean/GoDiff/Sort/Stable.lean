import GoLean.Sort.Stable
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let kind :: tag :: n :: raw := row | throw (IO.userError "expected entrypoint, shape, Len and keys")
  let data : MergeState := ⟨raw.zipIdx.toArray, []⟩
  let result : StableResult := match kind with
    | 0 => ⟨"ok", 0, stableSort data n⟩
    | 1 => ⟨"ok", 0, stableSortFunc data n⟩
    | 2 => stableInterface data (if tag = 0 then some n else none)
    | _ =>
      let shape := if tag = 2 then ReflectionSlice.invalid
        else if tag = 3 then ReflectionSlice.nonSlice else ReflectionSlice.valid
      stableSlice data shape (tag != 1)
  let ids := String.intercalate "," (result.state.values.toList.map (toString ∘ Prod.snd))
  let calls := result.state.reversedCalls.reverse.filterMap fun (swap, i, j) =>
    if swap && kind = 3 then none else some s!"{if swap then "S" else "L"}{i},{j}"
  return s!"{result.status}:{result.lenCalls}:{ids}:" ++ String.intercalate ";" calls
