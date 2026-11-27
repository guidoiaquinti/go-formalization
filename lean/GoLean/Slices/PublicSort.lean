import GoLean.Slices.AppendCollect
import GoLean.Slices.Pdqsort

namespace GoLean.Slices

inductive PublicSortFailure (ε : Type) where
  | sequence (failure : AppendSeqError ε)
  | sort (failure : Order2Failure ε)

structure PublicSortObservation (α ε : Type) where
  failure : Option (PublicSortFailure ε)
  values : Option (List α)
  calls : List (α × α)
  acknowledged : Nat

private def initialLimit (width length : Nat) : Int :=
  let n := (BitVec.ofNat width length).toNat
  if n = 0 then 0 else Int.ofNat (n.log2 + 1)

private def sortState (input : Option (List α)) : LocalSortState α ε :=
  ⟨(input.getD []).toArray, none, []⟩

private def observeSort (input : Option (List α)) (state : LocalSortState α ε) :
    PublicSortObservation α ε :=
  ⟨state.failure.map PublicSortFailure.sort, input.map (fun _ => state.cells.toList), state.calls, 0⟩

/-- Go 1.27.0 src/slices/sort.go:16–19. Start pdqsort with the native bit length. -/
def sliceSort [Inhabited α] (width : Nat) (input : Option (List α)) (less : α → α → Bool) :
    PublicSortObservation α ε :=
  let s := sortState input
  observeSort input (pdqsortOrdered width s 0 s.cells.size (initialLimit width s.cells.size) less)

/-- Go 1.27.0 src/slices/sort.go:30–33. Keep exact comparator order and partial mutations. -/
def sliceSortFunc [Inhabited α] (width : Nat) (input : Option (List α))
    (compare : Option (Nat → α → α → Except ε Int)) : PublicSortObservation α ε :=
  let s := sortState input
  observeSort input (pdqsortCmpFunc width s 0 s.cells.size (initialLimit width s.cells.size) compare)

private def collectThen (sequence : Option (List α × Except ε Unit))
    (sort : Option (List α) → PublicSortObservation α ε) : PublicSortObservation α ε :=
  let collected := sliceCollect sequence
  match collected.result with
  | .error failure => ⟨some (.sequence failure), none, [], collected.acknowledged⟩
  | .ok output =>
    let sorted := sort output
    let values := if sorted.failure.isSome then none else sorted.values
    ⟨sorted.failure, values, sorted.calls, collected.acknowledged⟩

/-- Go 1.27.0 src/slices/iter.go:66–70. Collect first, then sort. -/
def sliceSorted [Inhabited α] (width : Nat) (sequence : Option (List α × Except ε Unit))
    (less : α → α → Bool) : PublicSortObservation α ε :=
  collectThen sequence (fun input => sliceSort width input less)

/-- Go 1.27.0 src/slices/iter.go:75–79. A producer panic prevents comparator calls. -/
def sliceSortedFunc [Inhabited α] (width : Nat) (sequence : Option (List α × Except ε Unit))
    (compare : Option (Nat → α → α → Except ε Int)) : PublicSortObservation α ε :=
  collectThen sequence (fun input => sliceSortFunc width input compare)

theorem sliceSort_empty [Inhabited α] (width : Nat) (less : α → α → Bool) :
    sliceSort (ε := ε) width (some []) less = ⟨none, some [], [], 0⟩ := rfl

theorem sliceSort_three :
    sliceSort (ε := Unit) 64 (some ([3, 2, 1] : List Int)) (fun x y => x < y) =
      ⟨none, some [1, 2, 3], [], 0⟩ := rfl

theorem sliceSortFunc_partial :
    sliceSortFunc 64 (some ([3, 2, 1] : List Int))
      (some (fun i _ _ => if i = 1 then .error (1 : Int) else .ok (-1))) =
      ⟨some (.sort (.callback 1)), some [2, 3, 1], [(2, 3), (1, 3)], 0⟩ := rfl

theorem sliceSortFunc_nilOnPair :
    sliceSortFunc (ε := Unit) 64 (some ([2, 1] : List Int)) none =
      ⟨some (.sort .nilComparator), some [2, 1], [], 0⟩ := rfl

theorem sliceSorted_empty [Inhabited α] (width : Nat) (less : α → α → Bool) :
    sliceSorted (ε := ε) width (some ([], .ok ())) less = ⟨none, none, [], 0⟩ := rfl

theorem sliceSorted_three :
    sliceSorted (ε := Unit) 64 (some (([3, 2, 1] : List Int), .ok ())) (fun x y => x < y) =
      ⟨none, some [1, 2, 3], [], 3⟩ := rfl

theorem sliceSortedFunc_producerFailure [Inhabited α] (width : Nat) (values : List α)
    (payload : ε) (compare : Option (Nat → α → α → Except ε Int)) :
    sliceSortedFunc width (some (values, .error payload)) compare =
      ⟨some (.sequence (.producer payload)), none, [], values.length⟩ := rfl

end GoLean.Slices
