import GoLean.Slices.SortStateOps

namespace GoLean.Slices

private def compareAdjacent [Inhabited α] (compare : LocalCompare α ε)
    (s : LocalSortState α ε) (j : Nat) : Bool × LocalSortState α ε :=
  lessAt compare s j (j - 1)

private def swapAdjacent [Inhabited α] (s : LocalSortState α ε) (j : Nat) : LocalSortState α ε :=
  exchange s j (j - 1)

private def advance [Inhabited α] (compare : LocalCompare α ε) :
    Nat → Nat → LocalSortState α ε → Nat × LocalSortState α ε
  | 0, i, s => (i, s)
  | fuel + 1, i, s =>
    let (inverted, s) := compareAdjacent compare s i
    if inverted || s.failure.isSome then (i, s)
    else advance compare fuel (i + 1) s

private def shiftLeft [Inhabited α] (compare : LocalCompare α ε) :
    Nat → LocalSortState α ε → LocalSortState α ε
  | 0, s => s
  | j + 1, s =>
    let (inverted, s) := compareAdjacent compare s (j + 1)
    if inverted then shiftLeft compare j (swapAdjacent s (j + 1)) else s

private def shiftRight [Inhabited α] (compare : LocalCompare α ε) :
    Nat → Nat → LocalSortState α ε → LocalSortState α ε
  | 0, _, s => s
  | fuel + 1, j, s =>
    let (inverted, s) := compareAdjacent compare s j
    if inverted then shiftRight compare fuel (j + 1) (swapAdjacent s j) else s

private def repair [Inhabited α] (compare : LocalCompare α ε) (a b : Nat) :
    Nat → Nat → LocalSortState α ε → Option Bool × LocalSortState α ε
  | 0, _, s => (some false, s)
  | steps + 1, i, s =>
    let (i, s) := advance compare (b - i) i s
    if s.failure.isSome then (none, s)
    else if i = b then (some true, s)
    else if b - a < 50 then (some false, s)
    else
      let s := swapAdjacent s i
      let s := if 2 ≤ i - a then shiftLeft compare (i - 1) s else s
      let s := if 2 ≤ b - i then shiftRight compare (b - (i + 1)) (i + 1) s else s
      if s.failure.isSome then (none, s)
      else repair compare a b steps i s

/-- Go 1.27.0 src/slices/zsortordered.go:197–237. After at most five
    repairs, the helper returns false. The left shift stops at index zero. -/
def partialInsertionSortOrdered [Inhabited α] (s : LocalSortState α ε) (a b : Nat)
    (less : α → α → Bool) : Option Bool × LocalSortState α ε :=
  repair (orderedCompare less) a b 5 (a + 1) s

/-- Go 1.27.0 src/slices/zsortanyfunc.go:195–235. A callback panic stops
    the current scan and preserves earlier exchanges. -/
def partialInsertionSortCmpFunc [Inhabited α] (s : LocalSortState α ε) (a b : Nat)
    (compare : Option (Nat → α → α → Except ε Int)) : Option Bool × LocalSortState α ε :=
  repair (scriptedCompare compare) a b 5 (a + 1) s

theorem partialInsertionSortOrdered_empty (s : LocalSortState Int Nat) (a : Nat)
    (less : Int → Int → Bool) (healthy : s.failure = none) :
    partialInsertionSortOrdered s a a less = (some false, s) := by
  simp [partialInsertionSortOrdered, repair, advance, healthy, Nat.sub_eq_zero_of_le]

theorem partialInsertionSortOrdered_movesPrefix :
    let state : LocalSortState Int Nat :=
      ⟨(#[900, 901, 10, 11, 0] ++ ((List.range 47).map (fun n => Int.ofNat n + 13)).toArray), none, []⟩
    let (result, state) := partialInsertionSortOrdered state 2 52 (fun x y => x < y)
    result = some true ∧ state.cells[0]? = some 0 := by decide

theorem partialInsertionSortCmpFunc_panic :
    let state : LocalSortState Int Nat := ⟨#[2, 1] ++ (List.replicate 48 3).toArray, none, []⟩
    let (result, state) := partialInsertionSortCmpFunc state 0 50
      (some (fun call _ _ => if call = 1 then .error 1 else .ok (-1)))
    result = none ∧ state.cells[0]? = some 1 ∧
      state.failure = some (.callback 1) ∧ state.calls = [(1, 2), (3, 2)] := by
  exact ⟨rfl, rfl, rfl, rfl⟩

end GoLean.Slices
