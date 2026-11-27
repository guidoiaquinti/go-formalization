import GoLean.Sort.LocalSorts

namespace GoLean.Sort

private def pivotLess (state : LocalSortState) (i j : Nat) : Bool × LocalSortState :=
  (state.values[i]! < state.values[j]!,
    { state with reversedCalls := (false, i, j) :: state.reversedCalls })

private def pivotSwap (state : LocalSortState) (i j : Nat) : LocalSortState :=
  { values := (state.values.set! i state.values[j]!).set! j state.values[i]!
    reversedCalls := (true, i, j) :: state.reversedCalls }

private def scanLeft (a j : Nat) (equal : Bool) :
    Nat → LocalSortState → Nat → Nat × LocalSortState
  | 0, state, i => (i, state)
  | fuel + 1, state, i =>
    if j < i then (i, state)
    else
      let (less, state) := if equal then pivotLess state a i else pivotLess state i a
      if less != equal then scanLeft a j equal fuel state (i + 1) else (i, state)

private def scanRight (a i : Nat) (equal : Bool) :
    Nat → LocalSortState → Nat → Nat × LocalSortState
  | 0, state, j => (j, state)
  | fuel + 1, state, j =>
    if j < i then (j, state)
    else
      let (less, state) := if equal then pivotLess state a j else pivotLess state j a
      if less == equal then scanRight a i equal fuel state (j - 1) else (j, state)

private def splitRange (a width : Nat) (equal : Bool) :
    Nat → LocalSortState → Nat → Nat → Bool → Nat × Nat × Bool × LocalSortState
  | 0, state, i, j, untouched => (i, j, untouched, state)
  | fuel + 1, state, i, j, untouched =>
    let (i, state) := scanLeft a j equal width state i
    let (j, state) := scanRight a i equal width state j
    if j < i then (i, j, untouched, state)
    else splitRange a width equal fuel (pivotSwap state i j) (i + 1) (j - 1) false

/-- Go 1.27.0 src/sort/zsortinterface.go:135. The initial pivot stays at
the first index during the scan, then moves to the last lower element's slot. -/
def partition (state : LocalSortState) (a b pivot : Nat) : Nat × Bool × LocalSortState :=
  let width := b - a
  let (_, j, untouched, state) :=
    splitRange a width false width (pivotSwap state a pivot) (a + 1) (b - 1) true
  (j, untouched, pivotSwap state j a)

/-- Go 1.27.0 src/sort/zsortinterface.go:173. Leave the pivot at the first
index and return the first position after the equal group. -/
def partitionEqual (state : LocalSortState) (a b pivot : Nat) : Nat × LocalSortState :=
  let width := b - a
  let (i, _, _, state) :=
    splitRange a width true width (pivotSwap state a pivot) (a + 1) (b - 1) true
  (i, state)

def partitionFunc := partition
def partitionEqualFunc := partitionEqual

theorem partition_single (state : LocalSortState) :
    (partition state 0 1 0).1 = 0 ∧
      (partition state 0 1 0).2.1 = true ∧
      (partition state 0 1 0).2.2.reversedCalls =
        (true, 0, 0) :: (true, 0, 0) :: state.reversedCalls := by
  exact ⟨rfl, rfl, rfl⟩

theorem partitionEqual_single (state : LocalSortState) :
    (partitionEqual state 0 1 0).1 = 1 ∧
      (partitionEqual state 0 1 0).2.reversedCalls =
        (true, 0, 0) :: state.reversedCalls := by
  exact ⟨rfl, rfl⟩

theorem partitionFunc_same (s : LocalSortState) (a b pivot : Nat) :
    partitionFunc s a b pivot = partition s a b pivot := rfl

theorem partitionEqualFunc_same (s : LocalSortState) (a b pivot : Nat) :
    partitionEqualFunc s a b pivot = partitionEqual s a b pivot := rfl

end GoLean.Sort
