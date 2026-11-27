import GoLean.Slices.BreakPatterns
import GoLean.Slices.ChoosePivot
import GoLean.Slices.Partition
import GoLean.Slices.PartialInsertion
import GoLean.Slices.ReverseRange

namespace GoLean.Slices

private structure PdqOperations (α ε : Type) where
  insert : LocalSortState α ε → Nat → Nat → LocalSortState α ε
  heap : LocalSortState α ε → Nat → Nat → LocalSortState α ε
  choose : LocalSortState α ε → Nat → Nat → ChoosePivotObservation α ε
  partialInsert : LocalSortState α ε → Nat → Nat → Option Bool × LocalSortState α ε
  compare : LocalCompare α ε
  partition : LocalSortState α ε → Nat → Nat → Nat → Option (Nat × Bool) × LocalSortState α ε
  partitionEqual : LocalSortState α ε → Nat → Nat → Nat → Option Nat × LocalSortState α ε

private def replaceCells (s : LocalSortState α ε) (cells : Option (List α)) : LocalSortState α ε :=
  (cells.map fun values => { s with cells := values.toArray }).getD s

private def pdqLoop [Inhabited α] (width : Nat) (ops : PdqOperations α ε) :
    Nat → LocalSortState α ε → Nat → Nat → Int → Bool → Bool → LocalSortState α ε
  | 0, s, _, _, _, _, _ => s
  | fuel + 1, s, a, b, limit, balanced, partitioned =>
    let length := b - a
    if s.failure.isSome then s
    else if length ≤ 12 then ops.insert s a b
    else if limit = 0 then ops.heap s a b
    else
      let (s, limit) := if balanced then (s, limit) else
        let broken := breakPatternsOrdered width (some s.cells.toList) s.cells.size (Int.ofNat a) (Int.ofNat b)
        (replaceCells s broken.cells, (BitVec.ofInt width (limit - 1)).toInt)
      let chosen := ops.choose s a b
      let s := {s with calls := s.calls ++ chosen.calls}
      match chosen.result with
      | .error error => {s with failure := some error}
      | .ok (pivot, hint) =>
        let (pivot, hint, s) := if hint = 2 then
          ((b - 1) - (pivot - a), 1, replaceCells s (reverseRangeOrdered (some s.cells.toList) a b))
          else (pivot, hint, s)
        let (done, s) := if balanced && partitioned && hint = 1
          then ops.partialInsert s a b else (some false, s)
        if done != some false then s else
          let (duplicate, s) := if 0 < a then
            let (less, s) := lessAt ops.compare s (a - 1) pivot
            (!less, s)
            else (false, s)
          if s.failure.isSome then s
          else if duplicate then
            let (mid, s) := ops.partitionEqual s a b pivot
            match mid with
            | none => s
            | some mid => pdqLoop width ops fuel s mid b limit balanced partitioned
          else
            let (mid, s) := ops.partition s a b pivot
            match mid with
            | none => s
            | some (mid, partitioned) =>
              let left := mid - a
              let right := b - mid
              if left < right then
                let smaller := pdqLoop width ops fuel s a mid limit true true
                pdqLoop width ops fuel smaller (mid + 1) b limit (left ≥ length / 8) partitioned
              else
                let smaller := pdqLoop width ops fuel s (mid + 1) b limit true true
                pdqLoop width ops fuel smaller a mid limit (right ≥ length / 8) partitioned

/-- Go 1.27.0 src/slices/zsortordered.go:63–133. Each child and tail iteration
    removes at least one element, so the original length plus one bounds depth. -/
def pdqsortOrdered [Inhabited α] (width : Nat) (s : LocalSortState α ε) (a b : Nat) (limit : Int)
    (less : α → α → Bool) : LocalSortState α ε :=
  let ops : PdqOperations α ε := {
    insert := fun s a b => insertionSortOrdered s a b less
    heap := fun s a b => heapSortOrdered s a b less
    choose := fun s a b => choosePivotOrdered width (fun i => s.cells[i]!) a b less
    partialInsert := fun s a b => partialInsertionSortOrdered s a b less
    compare := orderedCompare less
    partition := fun s a b p => partitionOrdered s a b p less
    partitionEqual := fun s a b p => partitionEqualOrdered s a b p less }
  pdqLoop width ops (b - a + 1) s a b limit true true

/-- Go 1.27.0 src/slices/zsortanyfunc.go:61–131. The comparator script
    uses a global invocation index spanning recursion, medians, and fallbacks. -/
def pdqsortCmpFunc [Inhabited α] (width : Nat) (s : LocalSortState α ε) (a b : Nat) (limit : Int)
    (compare : Option (Nat → α → α → Except ε Int)) : LocalSortState α ε :=
  let ops : PdqOperations α ε := {
    insert := fun s a b => insertionSortCmpFunc s a b compare
    heap := fun s a b => heapSortCmpFunc s a b compare
    choose := fun s a b => choosePivotCmpFunc width (fun i => s.cells[i]!) a b
      (compare.map fun f i => f (s.calls.length + i))
    partialInsert := fun s a b => partialInsertionSortCmpFunc s a b compare
    compare := scriptedCompare compare
    partition := fun s a b p => partitionCmpFunc s a b p compare
    partitionEqual := fun s a b p => partitionEqualCmpFunc s a b p compare }
  pdqLoop width ops (b - a + 1) s a b limit true true

theorem pdqsortOrdered_short [Inhabited α] (width : Nat) (s : LocalSortState α ε)
    (a b : Nat) (limit : Int) (less : α → α → Bool)
    (healthy : s.failure = none) (short : b - a ≤ 12) :
    pdqsortOrdered width s a b limit less = insertionSortOrdered s a b less := by
  simp [pdqsortOrdered, pdqLoop, healthy, short]

theorem pdqsortCmpFunc_fallback [Inhabited α] (width : Nat) (s : LocalSortState α ε)
    (a b : Nat) (compare : Option (Nat → α → α → Except ε Int))
    (healthy : s.failure = none) (large : 12 < b - a) :
    pdqsortCmpFunc width s a b 0 compare = heapSortCmpFunc s a b compare := by
  simp [pdqsortCmpFunc, pdqLoop, healthy, Nat.not_le.mpr large]

end GoLean.Slices
