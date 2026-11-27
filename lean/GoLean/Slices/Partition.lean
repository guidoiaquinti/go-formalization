import GoLean.Slices.SortStateOps

namespace GoLean.Slices

private def scanLeft [Inhabited α] (compare : LocalCompare α ε)
    (a j : Nat) (equal : Bool) : Nat → Nat → LocalSortState α ε → Nat × LocalSortState α ε
  | 0, i, state => (i, state)
  | fuel + 1, i, state =>
    if j < i || state.failure.isSome then (i, state)
    else
      let (less, state) := if equal then lessAt compare state a i else lessAt compare state i a
      if state.failure.isNone && less != equal then scanLeft compare a j equal fuel (i + 1) state
      else (i, state)

private def scanRight [Inhabited α] (compare : LocalCompare α ε)
    (a i : Nat) (equal : Bool) : Nat → Nat → LocalSortState α ε → Nat × LocalSortState α ε
  | 0, j, state => (j, state)
  | fuel + 1, j, state =>
    if j < i || state.failure.isSome then (j, state)
    else
      let (less, state) := if equal then lessAt compare state a j else lessAt compare state j a
      if state.failure.isNone && less == equal then scanRight compare a i equal fuel (j - 1) state
      else (j, state)

private def splitRange [Inhabited α] (compare : LocalCompare α ε) (a width : Nat)
    (equal : Bool) : Nat → Nat → Nat → Bool → LocalSortState α ε →
      Nat × Nat × Bool × LocalSortState α ε
  | 0, i, j, untouched, state => (i, j, untouched, state)
  | fuel + 1, i, j, untouched, state =>
    let (i, state) := scanLeft compare a j equal width i state
    let (j, state) := scanRight compare a i equal width j state
    if j < i || state.failure.isSome then (i, j, untouched, state)
    else splitRange compare a width equal fuel (i + 1) (j - 1) false (exchange state i j)

private def partition [Inhabited α] (compare : LocalCompare α ε) (state : LocalSortState α ε)
    (a b pivot : Nat) : Option (Nat × Bool) × LocalSortState α ε :=
  let width := b - a
  let (_, j, untouched, state) :=
    splitRange compare a width false width (a + 1) (b - 1) true (exchange state a pivot)
  let state := exchange state j a
  (if state.failure.isSome then none else some (j, untouched), state)

private def partitionEqual [Inhabited α] (compare : LocalCompare α ε)
    (state : LocalSortState α ε) (a b pivot : Nat) : Option Nat × LocalSortState α ε :=
  let width := b - a
  let (i, _, _, state) :=
    splitRange compare a width true width (a + 1) (b - 1) true (exchange state a pivot)
  (if state.failure.isSome then none else some i, state)

/-- Go 1.27.0 src/slices/zsortordered.go:137–170. Counts only cross-swaps
    between scan stops when reporting whether the range was already partitioned. -/
def partitionOrdered [Inhabited α] (state : LocalSortState α ε) (a b pivot : Nat)
    (less : α → α → Bool) : Option (Nat × Bool) × LocalSortState α ε :=
  partition (orderedCompare less) state a b pivot

/-- Go 1.27.0 src/slices/zsortanyfunc.go:135–168. Comparator panic preserves
    the initial pivot exchange and earlier cross-swaps. -/
def partitionCmpFunc [Inhabited α] (state : LocalSortState α ε) (a b pivot : Nat)
    (compare : Option (Nat → α → α → Except ε Int)) : Option (Nat × Bool) × LocalSortState α ε :=
  partition (scriptedCompare compare) state a b pivot

/-- Go 1.27.0 src/slices/zsortordered.go:175–195. The pivot remains at the
    head and the return index is one past the group. -/
def partitionEqualOrdered [Inhabited α] (state : LocalSortState α ε) (a b pivot : Nat)
    (less : α → α → Bool) : Option Nat × LocalSortState α ε :=
  partitionEqual (orderedCompare less) state a b pivot

/-- Go 1.27.0 src/slices/zsortanyfunc.go:173–193. The pivot is always the
    left comparator operand in the equal-group scan. -/
def partitionEqualCmpFunc [Inhabited α] (state : LocalSortState α ε) (a b pivot : Nat)
    (compare : Option (Nat → α → α → Except ε Int)) : Option Nat × LocalSortState α ε :=
  partitionEqual (scriptedCompare compare) state a b pivot

theorem partitionOrdered_single :
    partitionOrdered (⟨#[8, 3, 9], none, []⟩ : LocalSortState Int Nat) 1 2 1
      (fun x y => x < y) = (some (1, true), ⟨#[8, 3, 9], none, []⟩) := rfl

theorem partitionCmpFunc_initialSwap :
    partitionCmpFunc (⟨#[7, 3, 5], none, []⟩ : LocalSortState Int Nat) 0 3 2
      (some (fun _ _ _ => .error 0)) =
      (none, ⟨#[5, 3, 7], some (.callback 0), [(3, 5)]⟩) := rfl

theorem partitionEqualOrdered_allEqual :
    partitionEqualOrdered (⟨#[8, 3, 3, 3, 9], none, []⟩ : LocalSortState Int Nat) 1 4 2
      (fun x y => x < y) = (some 4, ⟨#[8, 3, 3, 3, 9], none, []⟩) := rfl

theorem partitionEqualCmpFunc_pivotFirst :
    partitionEqualCmpFunc (⟨#[7, 5, 3], none, []⟩ : LocalSortState Int Nat) 0 3 2
      (some (fun _ _ _ => .error 0)) =
      (none, ⟨#[3, 5, 7], some (.callback 0), [(3, 5)]⟩) := rfl

end GoLean.Slices
