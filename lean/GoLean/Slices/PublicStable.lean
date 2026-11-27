import GoLean.Slices.PublicSort
import GoLean.Slices.Stable

namespace GoLean.Slices

/-- Go 1.27.0 src/slices/sort.go:37–39. Stable sort operates on the
    caller's whole slice, retaining mutations made before a callback fault. -/
def sliceSortStableFunc [Inhabited α] (input : Option (List α))
    (compare : Option (Nat → α → α → Except ε Int)) : PublicSortObservation α ε :=
  let values := input.getD []
  if values.isEmpty then ⟨none, input, [], 0⟩
  else
    let s : LocalSortState α ε := ⟨values.toArray, none, []⟩
    let sorted := stableCmpFunc s s.cells.size compare
    ⟨sorted.failure.map PublicSortFailure.sort, input.map (fun _ => sorted.cells.toList), sorted.calls, 0⟩

/-- Go 1.27.0 src/slices/iter.go:86–90. Collect the iterator before
    calling the comparator. A panicking call does not return a result. -/
def sliceSortedStableFunc [Inhabited α] (sequence : Option (List α × Except ε Unit))
    (compare : Option (Nat → α → α → Except ε Int)) : PublicSortObservation α ε :=
  let collected := sliceCollect sequence
  match collected.result with
  | .error failure => ⟨some (.sequence failure), none, [], collected.acknowledged⟩
  | .ok input =>
    let sorted := sliceSortStableFunc input compare
    let values := if sorted.failure.isSome then none else sorted.values
    ⟨sorted.failure, values, sorted.calls, collected.acknowledged⟩

theorem sliceSortStableFunc_ties :
    sliceSortStableFunc (ε := Unit) (some ([513, 256, 512] : List Int))
      (some (fun _ x y => .ok (x / 256 - y / 256))) =
      ⟨none, some [256, 513, 512], [(256, 513), (512, 513)], 0⟩ := rfl

theorem sliceSortedStableFunc_empty [Inhabited α]
    (compare : Option (Nat → α → α → Except ε Int)) :
    sliceSortedStableFunc (some ([], .ok ())) compare = ⟨none, none, [], 0⟩ := rfl

theorem sliceSortedStableFunc_producerFailure [Inhabited α] (values : List α) (payload : ε)
    (compare : Option (Nat → α → α → Except ε Int)) :
    sliceSortedStableFunc (some (values, .error payload)) compare =
      ⟨some (.sequence (.producer payload)), none, [], values.length⟩ := rfl

end GoLean.Slices
