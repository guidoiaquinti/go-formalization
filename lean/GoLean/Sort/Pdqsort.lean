import GoLean.Sort.BreakPatterns
import GoLean.Sort.RangeMoves
import GoLean.Sort.PartialInsertion
import GoLean.Sort.Partition

namespace GoLean.Sort

private def pdqLess (s : LocalSortState) (i j : Nat) : Bool × LocalSortState :=
  (s.values[i]! < s.values[j]!,
   {s with reversedCalls := (false, i, j) :: s.reversedCalls})

private def applySwaps (s : LocalSortState) (pairs : List (Nat × Nat)) : LocalSortState :=
  pairs.foldl (fun s (i, j) =>
    { values := (s.values.set! i s.values[j]!).set! j s.values[i]!
      reversedCalls := (true, i, j) :: s.reversedCalls }) s

private def orderPair (s : LocalSortState) (inversions i j : Nat) :
    Nat × Nat × Nat × LocalSortState :=
  let (reversed, s) := pdqLess s j i
  if reversed then (j, i, inversions + 1, s) else (i, j, inversions, s)

private def chooseMedian (s : LocalSortState) (inversions i j k : Nat) :
    Nat × Nat × LocalSortState :=
  let (i, j, inversions, s) := orderPair s inversions i j
  let (j, _, inversions, s) := orderPair s inversions j k
  let (_, j, inversions, s) := orderPair s inversions i j
  (j, inversions, s)

private def pivot (s : LocalSortState) (a b : Nat) : Nat × Nat × LocalSortState :=
  let quarter := (b - a) / 4
  let i := a + quarter
  let j := a + 2 * quarter
  let k := a + 3 * quarter
  if b - a < 50 then chooseMedian s 0 i j k
  else
    let (i, count, s) := chooseMedian s 0 (i - 1) i (i + 1)
    let (j, count, s) := chooseMedian s count (j - 1) j (j + 1)
    let (k, count, s) := chooseMedian s count (k - 1) k (k + 1)
    chooseMedian s count i j k

private def sortSteps (wordWidth : Nat) :
    Nat → LocalSortState → Nat → Nat → Nat → Bool → Bool → LocalSortState
  | 0, s, _, _, _, _, _ => s
  | fuel + 1, s, a, b, limit, balanced, partitioned =>
    let length := b - a
    if length ≤ 12 then insertionSort s a b
    else if limit = 0 then heapSort s a b
    else
      let (s, limit) := if balanced then (s, limit)
        else (applySwaps s (breakPatterns wordWidth a b), limit - 1)
      let (p, inversions, s) := pivot s a b
      let (p, increasing, s) := if inversions = 12 then
          (b - 1 - (p - a), true, applySwaps s (reverseRange a b))
        else (p, inversions == 0, s)
      let (done, s) := if balanced && partitioned && increasing then partialInsertionSort s a b
        else (false, s)
      if done then s
      else
        let (duplicate, s) := if 0 < a then
            let (less, s) := pdqLess s (a - 1) p
            (!less, s)
          else (false, s)
        if duplicate then
          let (mid, s) := partitionEqual s a b p
          sortSteps wordWidth fuel s mid b limit balanced partitioned
        else
          let (mid, partitioned, s) := partition s a b p
          let left := mid - a
          let right := b - mid
          if left < right then
            let s := sortSteps wordWidth fuel s a mid limit true true
            sortSteps wordWidth fuel s (mid + 1) b limit (left ≥ length / 8) partitioned
          else
            let s := sortSteps wordWidth fuel s (mid + 1) b limit true true
            sortSteps wordWidth fuel s a mid limit (right ≥ length / 8) partitioned

/-- Go 1.27.0 src/sort/zsortinterface.go:61. Each recursive child and the next
loop interval exclude the current pivot, or at least one equal element. -/
def pdqsort (width : Nat) (s : LocalSortState) (a b limit : Nat) : LocalSortState :=
  sortSteps width (b - a + 1) s a b limit true true

def pdqsortFunc := pdqsort

theorem pdqsort_small (width : Nat) (s : LocalSortState) (a b limit : Nat) (h : b - a ≤ 12) :
    pdqsort width s a b limit = insertionSort s a b := by simp [pdqsort, sortSteps, h]

theorem pdqsort_fallback (width : Nat) (s : LocalSortState) (a b : Nat) (h : 12 < b - a) :
    pdqsort width s a b 0 = heapSort s a b := by
  simp [pdqsort, sortSteps, Nat.not_le.mpr h]

theorem pdqsortFunc_same (w : Nat) (s : LocalSortState) (a b limit : Nat) :
    pdqsortFunc w s a b limit = pdqsort w s a b limit := rfl

end GoLean.Sort
