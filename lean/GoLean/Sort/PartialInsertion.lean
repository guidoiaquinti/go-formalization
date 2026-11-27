import GoLean.Sort.LocalSorts

namespace GoLean.Sort

private def compareAdjacent (s : LocalSortState) (j : Nat) : Bool × LocalSortState :=
  (s.values[j]! < s.values[j - 1]!,
   {s with reversedCalls := (false, j, j - 1) :: s.reversedCalls})

private def swapAdjacent (s : LocalSortState) (j : Nat) : LocalSortState :=
  { values := (s.values.set! j s.values[j - 1]!).set! (j - 1) s.values[j]!
    reversedCalls := (true, j, j - 1) :: s.reversedCalls }

private def advance (b : Nat) : Nat → Nat → LocalSortState → Nat × LocalSortState
  | 0, i, s => (i, s)
  | fuel + 1, i, s =>
    if b ≤ i then (i, s)
    else
      let (inverted, s) := compareAdjacent s i
      if inverted then (i, s) else advance b fuel (i + 1) s

private def shiftLeft : Nat → LocalSortState → LocalSortState
  | 0, s => s
  | j + 1, s =>
    let (inverted, s) := compareAdjacent s (j + 1)
    if inverted then shiftLeft j (swapAdjacent s (j + 1)) else s

private def shiftRight (b : Nat) : Nat → Nat → LocalSortState → LocalSortState
  | 0, _, s => s
  | fuel + 1, j, s =>
    if b ≤ j then s
    else
      let (inverted, s) := compareAdjacent s j
      if inverted then shiftRight b fuel (j + 1) (swapAdjacent s j) else s

private def repair (a b : Nat) : Nat → Nat → LocalSortState → Bool × LocalSortState
  | 0, _, s => (false, s)
  | steps + 1, i, s =>
    let (i, s) := advance b (b - i) i s
    if i = b then (true, s)
    else if b - a < 50 then (false, s)
    else
      let s := swapAdjacent s i
      let s := if 2 ≤ i - a then shiftLeft (i - 1) s else s
      let s := if 2 ≤ b - i then shiftRight b (b - i) (i + 1) s else s
      repair a b steps i s

/-- Go 1.27.0 src/sort/zsortinterface.go:195. At most five adjacent inversions
are repaired. The left shift can visit indices below a, stopping at zero. -/
def partialInsertionSort (s : LocalSortState) (a b : Nat) : Bool × LocalSortState :=
  repair a b 5 (a + 1) s

def partialInsertionSortFunc := partialInsertionSort

theorem partialInsertionSort_empty (s : LocalSortState) (a : Nat) :
    partialInsertionSort s a a = (false, s) := by simp [partialInsertionSort, repair, advance, Nat.sub_eq_zero_of_le]

theorem partialInsertionSort_single (s : LocalSortState) (a : Nat) :
    partialInsertionSort s a (a + 1) = (true, s) := by simp [partialInsertionSort, repair, advance]

theorem partialInsertionSortFunc_same (s : LocalSortState) (a b : Nat) :
    partialInsertionSortFunc s a b = partialInsertionSort s a b := rfl

end GoLean.Sort
