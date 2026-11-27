import GoLean.Slices.LocalSorts

namespace GoLean.Slices

abbrev LocalCompare (α ε : Type) :=
  α → α → LocalSortState α ε → Bool × LocalSortState α ε

def orderedCompare (less : α → α → Bool) : LocalCompare α ε :=
  fun x y state => (less x y, state)

def scriptedCompare (compare : Option (Nat → α → α → Except ε Int)) :
    LocalCompare α ε := fun x y state =>
  match compare with
  | none => (false, { state with failure := some .nilComparator })
  | some f =>
    let calls := state.calls ++ [(x, y)]
    match f state.calls.length x y with
    | .error e => (false, { state with failure := some (.callback e), calls })
    | .ok v => (v < 0, { state with calls })

def lessAt [Inhabited α] (compare : LocalCompare α ε)
    (state : LocalSortState α ε) (i j : Nat) : Bool × LocalSortState α ε :=
  if state.failure.isSome then (false, state)
  else compare state.cells[i]! state.cells[j]! state

def exchange [Inhabited α] (state : LocalSortState α ε) (i j : Nat) : LocalSortState α ε :=
  if state.failure.isSome then state
  else { state with cells := (state.cells.set! i state.cells[j]!).set! j state.cells[i]! }

end GoLean.Slices
