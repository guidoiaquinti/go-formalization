import GoLean.Slices.Order2

namespace GoLean.Slices

structure LocalSortState (α ε : Type) where
  cells : Array α
  failure : Option (Order2Failure ε) := none
  calls : List (α × α) := []

private abbrev LocalCompare (α ε : Type) :=
  α → α → LocalSortState α ε → Bool × LocalSortState α ε

private def orderedCompare (less : α → α → Bool) : LocalCompare α ε :=
  fun x y state => (less x y, state)

private def scriptedCompare (compare : Option (Nat → α → α → Except ε Int)) :
    LocalCompare α ε := fun x y state =>
  match compare with
  | none => (false, { state with failure := some .nilComparator })
  | some f =>
    let calls := state.calls ++ [(x, y)]
    match f state.calls.length x y with
    | .error e => (false, { state with failure := some (.callback e), calls })
    | .ok v => (v < 0, { state with calls })

private def lessAt [Inhabited α] (compare : LocalCompare α ε)
    (state : LocalSortState α ε) (i j : Nat) : Bool × LocalSortState α ε :=
  if state.failure.isSome then (false, state)
  else compare state.cells[i]! state.cells[j]! state

private def exchange [Inhabited α] (state : LocalSortState α ε) (i j : Nat) : LocalSortState α ε :=
  if state.failure.isSome then state
  else { state with cells := (state.cells.set! i state.cells[j]!).set! j state.cells[i]! }

private def insertLeft [Inhabited α] (compare : LocalCompare α ε) (a : Nat) :
    Nat → LocalSortState α ε → LocalSortState α ε
  | 0, s => s
  | j + 1, s =>
    if j + 1 ≤ a then s
    else
      let (invert, s) := lessAt compare s (j + 1) j
      if invert then insertLeft compare a j (exchange s (j + 1) j) else s

private def insertionSort [Inhabited α] (compare : LocalCompare α ε)
    (s : LocalSortState α ε) (a b : Nat) : LocalSortState α ε :=
  (List.range (b - (a + 1))).foldl (fun acc delta => insertLeft compare a (a + 1 + delta) acc) s

private def sift [Inhabited α] (compare : LocalCompare α ε) (hi first : Nat) :
    Nat → Nat → LocalSortState α ε → LocalSortState α ε
  | 0, _, s => s
  | fuel + 1, root, s =>
    let child := 2 * root + 1
    if hi ≤ child then s
    else
      let (right, s) :=
        if child + 1 < hi then lessAt compare s (first + child) (first + child + 1)
        else (false, s)
      let greater := if right then child + 1 else child
      let (invert, s) := lessAt compare s (first + root) (first + greater)
      if invert then sift compare hi first fuel greater (exchange s (first + root) (first + greater))
      else s

private def siftDown [Inhabited α] (compare : LocalCompare α ε)
    (s : LocalSortState α ε) (lo hi first : Nat) : LocalSortState α ε :=
  sift compare hi first hi lo s

private def heapSort [Inhabited α] (compare : LocalCompare α ε)
    (s : LocalSortState α ε) (a b : Nat) : LocalSortState α ε :=
  let hi := b - a
  let built := (List.range ((hi - 1) / 2 + 1)).reverse.foldl
    (fun acc root => siftDown compare acc root hi a) s
  (List.range hi).reverse.foldl
    (fun acc stop => siftDown compare (exchange acc a (a + stop)) 0 stop a) built

/-- Go 1.27.0 src/slices/zsortordered.go:12–18. Operates on a valid subrange. -/
def insertionSortOrdered [Inhabited α] (s : LocalSortState α ε) (a b : Nat)
    (less : α → α → Bool) : LocalSortState α ε :=
  insertionSort (orderedCompare less) s a b

/-- Go 1.27.0 src/slices/zsortanyfunc.go:10–16. A panic preserves prior writes. -/
def insertionSortCmpFunc [Inhabited α] (s : LocalSortState α ε) (a b : Nat)
    (cmp : Option (Nat → α → α → Except ε Int)) : LocalSortState α ε :=
  insertionSort (scriptedCompare cmp) s a b

/-- Go 1.27.0 src/slices/zsortordered.go:22–38. Indices are relative to first. -/
def siftDownOrdered [Inhabited α] (s : LocalSortState α ε) (lo hi first : Nat)
    (less : α → α → Bool) : LocalSortState α ε :=
  siftDown (orderedCompare less) s lo hi first

/-- Go 1.27.0 src/slices/zsortanyfunc.go:20–36. Child selection precedes root comparison. -/
def siftDownCmpFunc [Inhabited α] (s : LocalSortState α ε) (lo hi first : Nat)
    (cmp : Option (Nat → α → α → Except ε Int)) : LocalSortState α ε :=
  siftDown (scriptedCompare cmp) s lo hi first

/-- Go 1.27.0 src/slices/zsortordered.go:40–55. Build right to left, then pop through index zero. -/
def heapSortOrdered [Inhabited α] (s : LocalSortState α ε) (a b : Nat)
    (less : α → α → Bool) : LocalSortState α ε :=
  heapSort (orderedCompare less) s a b

/-- Go 1.27.0 src/slices/zsortanyfunc.go:38–53. Uses the same heap build and pop loops. -/
def heapSortCmpFunc [Inhabited α] (s : LocalSortState α ε) (a b : Nat)
    (cmp : Option (Nat → α → α → Except ε Int)) : LocalSortState α ε :=
  heapSort (scriptedCompare cmp) s a b

theorem insertionSortOrdered_single [Inhabited α] (s : LocalSortState α ε) (a : Nat)
    (less : α → α → Bool) : insertionSortOrdered s a (a + 1) less = s := by
  simp [insertionSortOrdered, insertionSort]

theorem insertionSortCmpFunc_partial :
    insertionSortCmpFunc (⟨#[3, 2, 1], none, []⟩ : LocalSortState Int Nat) 0 3
      (some (fun i _ _ => if i = 1 then .error 1 else .ok (-1))) =
      ⟨#[2, 3, 1], some (.callback 1), [(2, 3), (1, 3)]⟩ := rfl

theorem siftDownOrdered_leaf [Inhabited α] (s : LocalSortState α ε) (lo hi first : Nat)
    (less : α → α → Bool) (leaf : hi ≤ 2 * lo + 1) :
    siftDownOrdered s lo hi first less = s := by
  cases hi with
  | zero => rfl
  | succ n => simp [siftDownOrdered, siftDown, sift, leaf]

theorem siftDownCmpFunc_childFirst :
    siftDownCmpFunc (⟨#[0, 1, 2], none, []⟩ : LocalSortState Int Nat) 0 3 0
      (some (fun _ _ _ => .error 0)) =
      ⟨#[0, 1, 2], some (.callback 0), [(1, 2)]⟩ := rfl

theorem heapSortOrdered_two :
    heapSortOrdered (⟨#[8, 3, 1, 9], none, []⟩ : LocalSortState Int Nat) 1 3
      (fun a b => a < b) = ⟨#[8, 1, 3, 9], none, []⟩ := rfl

theorem heapSortCmpFunc_empty [Inhabited α] (s : LocalSortState α ε) (a : Nat)
    (compare : Option (Nat → α → α → Except ε Int)) :
    heapSortCmpFunc s a a compare = s := by
  simp [heapSortCmpFunc, heapSort, siftDown, sift]

end GoLean.Slices
